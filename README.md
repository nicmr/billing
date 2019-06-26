# altemista-billing

[![Build Status](https://travis-ci.org/Altemista/altemista-billing.svg?branch=master)](https://travis-ci.org/Altemista/altemista-billing)

An AWS billing service for Altemista Cloud.

Please read through [AWS Authentication and Config](#awsauthconfig) if you're running the application for the first time.

<!-- Please read through [CONTRIBUTING.md](/CONTRIBUTING.md) before making any contributions. -->

## Command line interface <a name="cli"></a>
```
Usage:
  altemista-billing [command]

Available Commands:
  help        Help about any command
  invoice     Analyzes costs and creates billing documents for a single month

Flags:
      --config string   supplies the application with the specified config file
  -h, --help            help for altemista-billing

Use "altemista-billing [command] --help" for more information about a command.
```
Run `altemista-billing help <sub-command>` for flags and detailed information for each subcommand

## Config file
The configuration of the application can be changed with a config file.
The application will search for a `config.yaml` in the working directory.

All possible command line flags are valid keys:
```
bucket: mybucket
month: current
provider: aws
```
The desired subcommand can currently not be specified, and has to be passed on the command line.

## AWS Auth and Config <a name="awsauthconfig"></a>

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


## Local builds <a name="builds"></a>
You can choose to either build a docker image or compile the program manually.

Before you start to build your project, please follow the steps in [AWS Auth and Config](#awsauthconfig).


### A. Build with Docker <a name="buildsdocker"></a>

Please make sure you have correctly set up authentication as described in `1. Authenticating with AWS`

Now build and run the program using docker
```shell
docker build -t $(basename $PWD) .
docker run -v /path/to/your/credentials:/home/runner/.aws/credentials $(basename $PWD)
```

### B. Build manually <a name="buildsmanual"></a>

Requirements:
- Go 1.12+
- ca-certificates (already present on most non-minimal linux distributions)

Now you can compile the program using Go 1.12+
```zsh
go run . createBill --month current --bucket <yourS3bucket>
```

## Testing <a name="testing"></a>

Automated testing is handled by Travis CI and can be configured in .travis.yml.

If you want to run tests locally, run

```zsh
go test ./...
```
