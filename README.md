# altemista-billing

[![Build Status](https://travis-ci.org/Altemista/altemista-billing.svg?branch=master)](https://travis-ci.org/Altemista/altemista-billing)

An AWS billing service for Altemista Cloud.

Please read through [Authentication and Config](#authconfig) if you're running the application for the first time.

<!-- Please read through [CONTRIBUTING.md](/CONTRIBUTING.md) before making any contributions. -->

## Command line interface <a name="cli"></a>
```
Usage:
  altemista-billing [command]

Available Commands:
  invoice     Create invoices for a specified month and uploads them to S3
  help        Help about any command
  serve       Serve http requests, exposing an API similar to that of cost

Flags:
  -h, --help   help for altemista-billing
```
Run `altemista-billing help <sub-command>` for flags and detailed information for each subcommand


## HTTP interface <a name="httpinterface"></a>
enabled if you invoke the serve subcommand.
```sh
altemista-billing serve -p 8080
```
Then send a get request to [localhost:8080/cost/?month=current](localhost:8080/invoice/?month=YYYY-MM).

The http API supports the same parameters as the cost CLI subcommand, in the form of get parameters. Run the following command for detailed information.

```sh
altemista-billing help cost
```

## Authentication and Config <a name="authconfig"></a>

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

Before you start to build your project, please follow the steps in [Authentication and Config](#authconfig).


### A. Build with Docker <a name="buildsdocker"></a>

Please make sure you have correctly set up authentication as described in `1. Authenticating with AWS`

Now build and run the program using docker
```shell
docker build -t $(basename $PWD) .
docker run -v /path/to/your/credentials:/home/runner/.aws/credentials $(basename $PWD)
```

### B. Build manually <a name="buildsmanual"></a>

Please make sure you have correctly set up authentication as described in `1. Authenticating with AWS`

Requirements:
- Go 1.12+
- ca-certificates (already present on most non-minimal linux distributions)

Now you can compile the program using Go 1.12+
```zsh
export AWS_SDK_LOAD_CONFIG=1 #only do this once
go run . createBill --month current --bucket <yourS3bucket>
```

## Testing <a name="testing"></a>

Automated testing is handled by Travis CI and can be configured in .travis.yml.

If you want to run tests locally, run

```zsh
export AWS_SDK_LOAD_CONFIG=1 #only do this once
go test ./...
```

## Debugging with VS Code <a name="debugging"></a>

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
