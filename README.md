# altemista-billing

[![Build Status](https://travis-ci.org/Altemista/altemista-billing.svg?branch=master)](https://travis-ci.org/Altemista/altemista-billing)

AWS billing service for Altemista Cloud

Please read through [CONTRIBUTING.md](/CONTRIBUTING.md) before making any contributions.

## Local builds
You can choose to either build a docker image or compile the program manually.

Docker is recommended for ease of use.

Before you decide on a build process, make sure your application can authenticate with your AWS subscription.

### 1. Authentication with AWS

The Application will look for `~/.aws/config` and `~/.aws/credentials` on your machine.

You can copy the default config included in this repository, and adapt it to your needs.
```shell
cp -r ./.aws/config ~/.aws/config
```


You will have to provide your own credentials file or use IAM Roles / Environment variables. You can read more about it here:
[AWS - Configuring sdk](https://docs.aws.amazon.com/de_de/sdk-for-go/v1/developer-guide/configuring-sdk.html)
and [AWS - Configure Files](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-files.html)

The easiest way is creating a credentials file with an AWS access token with the following contents at `~/.aws/credentials`:
```
[default]
aws_access_key_id=your_key_id_here
aws_secret_access_key=your_secret_here
```



### 2A. Build with Docker

Please make sure you have correctly set up authentication as described in `1. Authenticating with AWS`

Now build and run the program using docker
```shell
docker build -t $(basename $PWD) .
docker run -p 8080:8080 -v /path/to/your/credentials:/home/runner/.aws/credentials $(basename $PWD)
```
Then call via curl
```shell
curl -X get http://localhost:8080/costs\?start\=2019-03-29\&end\=2019-04-02
```
Or in your webbrowser
[http://localhost:8080/costs?start=2019-03-29&end=2019-04-02](http://localhost:8080/costs?start=2019-03-29&end=2019-04-02)


### 2B. Build manually

Please make sure you have correctly set up authentication as described in `1. Authenticating with AWS`

Requirements:
- Go 1.12+
- ca-certificates (already present on most non-minimal linux distributions)

Now you can compile the program using Go 1.12+
```zsh
export AWS_SDK_LOAD_CONFIG=1 #only do this once
go run .
```
Then call via curl
```shell
curl -X get http://localhost:8080/costs\?start\=2019-03-29\&end\=2019-04-02
```
Or in your webbrowser
[http://localhost:8080/costs?start=2019-03-29&end=2019-04-02](http://localhost:8080/costs?start=2019-03-29&end=2019-04-02)

## Testing

Automated testing is handled by Travis CI and can be configured in .travis.yml.

If you want to run tests locally, run

```zsh
export AWS_SDK_LOAD_CONFIG=1 #only do this once
go test ./...
```

## Debugging with VS Code

1. Make sure you have the official Microsoft Go extension installed and enabled in VS Code and have gone through all the required steps to compile the application manually

2. To avoid having to export AWS_SDK_LOAD_CONFIG=1 from the parent process of vscode, open VSCode, hit `Ctrl + Shift + P` to open the commands palette, enter `Debug: Open launch.json` and tweak the env key accordingly:
```json
{
    ...
    "env": {"AWS_SDK_LOAD_CONFIG":"1"},
    ...
}
```

3. Save the file and click the debug button in VS Code