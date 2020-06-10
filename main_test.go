package main

import "testing"

func Test_search(t *testing.T) {
	result, err := search("aws s3", false, 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%#+v", result)
}

func Test_searchAll(t *testing.T) {
	results := searchAll("aws s3", false)
	for result := range results {
		t.Logf("%#+v", result)
	}
}
