package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"strings"
	"time"
)

var flagVerified bool

func init() {
	flag.BoolVar(&flagVerified, "verified", false, "return only verified modules")

	flag.Parse()
}

type module struct {
	ID          string    `json:"id"`
	Owner       string    `json:"owner"`
	Namespace   string    `json:"namespace"`
	Name        string    `json:"name"`
	Version     string    `json:"version"`
	Provider    string    `json:"provider"`
	Description string    `json:"description"`
	Source      string    `json:"source"`
	Tag         string    `json:"tag"`
	PublishedAt time.Time `json:"published_at"`
	Downloads   int       `json:"downloads"`
	Verified    bool      `json:"verified"`
}

func (m *module) String() string {
	return fmt.Sprintf("%-60v %-10v %-10v %.10v", m.ID, m.Verified, m.Downloads, m.PublishedAt)
}

type searchResult struct {
	Meta struct {
		Limit         int    `json:"limit"`
		CurrentOffset int    `json:"current_offset"`
		NextOffset    int    `json:"next_offset"`
		NextURL       string `json:"next_url"`
	} `json:"meta"`
	Modules []module `json:"modules"`
}

func search(query string, verified bool, limit, offset int) (*searchResult, error) {
	req, err := http.NewRequest("GET", "https://registry.terraform.io/v1/modules/search", nil)
	if err != nil {
		return nil, err
	}

	// add required url parameters
	q := req.URL.Query()
	q.Add("q", query)
	q.Add("limit", fmt.Sprintf("%v", limit))
	q.Add("offset", fmt.Sprintf("%v", offset))
	q.Add("verified", fmt.Sprintf("%v", verified))
	req.URL.RawQuery = q.Encode()

	// add some headers
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("User-Agent", "Terraform Module Registry CLI github.com/picatz/tfmr")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Requested-With", "Picat's ♥")

	// make request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// decode json body
	decoder := json.NewDecoder(resp.Body)
	sr := searchResult{}
	err = decoder.Decode(&sr)
	if err != nil {
		panic(err)
	}

	return &sr, nil
}

func searchAll(query string, verified bool) <-chan *searchResult {
	results := make(chan (*searchResult))
	go func() {
		defer close(results)

		limit := 9
		offset := 0

		for {
			nextResult, err := search(query, verified, limit, offset)
			if err != nil {
				return
			}
			results <- nextResult
			// def not shotty next offset logic
			if nextResult.Meta.NextOffset != 0 && nextResult.Meta.NextOffset >= offset {
				offset += nextResult.Meta.NextOffset
			} else {
				return
			}
		}
	}()
	return results
}

func main() {
	searchQuery := strings.Join(flag.Args(), " ")

	// if flag is out of order (not provided first)
	if strings.Contains(searchQuery, "-verified") {
		flagVerified = true
		searchQuery = strings.Replace(searchQuery, "-verified", "", -1)
	}

	results := searchAll(searchQuery, flagVerified)

	topLine := fmt.Sprintf("%-60v %-10v %-10v %v", "ID", "Verified", "Downloads", "Published")

	fmt.Println(topLine)

	for result := range results {
		for _, module := range result.Modules {
			fmt.Println(module.String())
		}
	}
}
