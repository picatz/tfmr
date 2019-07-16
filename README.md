# tfmr

Search CLI for the [Terrarform Module Registry](https://registry.terraform.io/)

## Download

```console
$ go get -u github.com/picatz/tfmr
```

## Usage

```console
$ tfmr aws s3
ID                                                           Verified   Downloads  Published
javilac/s3/aws/0.0.1                                         false      19         2018-01-22
devops-workflow/s3-buckets/aws/0.3.0                         false      2479       2018-06-21
thrashr888/s3-website/aws/0.4.0                              false      49         2018-08-02
youngfeldt/backend-s3/aws/1.3.3                              false      580        2019-02-14
tmknom/s3-cloudtrail/aws/1.0.0                               false      10         2018-10-26
appzen-oss/s3-buckets/aws/0.3.2                              false      605        2019-02-12
...
```

To search only verified modules:

```console
$ tfmr -verified consul
ID                                                           Verified   Downloads  Published
hashicorp/consul/azurerm/0.0.5                               true       160        2019-02-14
hashicorp/consul/google/0.4.0                                true       3726       2019-06-26
hashicorp/consul/aws/0.7.1                                   true       16300      2019-07-11
...
```
