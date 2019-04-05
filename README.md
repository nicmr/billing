# altemista-billing

[![Build Status](https://travis-ci.org/Altemista/altemista-billing.svg?branch=master)](https://travis-ci.org/Altemista/altemista-billing)

AWS billing service for Altemista Cloud

Please read through [CONTRIBUTING.md](/CONTRIBUTING.md) before making any contributions.

## Local builds
You can choose to either build a docker image or compile the program manually.

Docker is recommended for ease of use.

### 1. Docker
Create a credentials file with an AWS access tokens at `altemista-billing/aws/credentials` with the following contents:
```
[default]
aws_access_key_id=your_key_id_here
aws_secret_access_key=your_secret_here
```

Now build and run the program using docker
```shell
docker build -t $(basename $PWD) .
docker run -p 8080:8080 $(basename $PWD)
```
Then call via curl
```shell
curl -X get http://localhost:8080/costs\?start\=2019-03-29\&end\=2019-04-02
```
Or in your webbrowser
[http://localhost:8080/costs?start=2019-03-29&end=2019-04-02](http://localhost:8080/costs?start=2019-03-29&end=2019-04-02)


### 2. Compile manually

Requirements:
- Go 1.12+
- ca-certificates (already present on most non-minimal linux distributions)

The AWS SDK will look for `~/.aws/config` and `~/.aws/credentials` on your machine.

You can copy the default config included in this repository, and adapt it to your needs.
```shell
cp -r ./.aws ~/.aws
```

You will have to provide your own credentials file or use IAM Roles / Environment variables. You can read more about it here:
https://docs.aws.amazon.com/de_de/sdk-for-go/v1/developer-guide/configuring-sdk.html

Now you can compile the program using Go 1.12+
```shell
go build .
AWS_SDK_LOAD_CONFIG=1 ./altemista-billing
```
Then call via curl
```shell
curl -X get http://localhost:8080/costs\?start\=2019-03-29\&end\=2019-04-02
```
Or in your webbrowser
[http://localhost:8080/costs?start=2019-03-29&end=2019-04-02](http://localhost:8080/costs?start=2019-03-29&end=2019-04-02)
