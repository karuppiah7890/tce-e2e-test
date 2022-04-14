# tce-e2e-test

This repo contains [TCE](https://github.com/vmware-tanzu/community-edition) E2E tests

## Running the tests

## On \*nix systems (Linux, MacOS)

Ensure you have latest golang `1.17.x` installed

```bash
# clone the repo
git clone https://github.com/karuppiah7890/tce-e2e-test

# change directory to the project
cd tce-e2e-test

# copy sameple .env file
cp .env.sample.azure .env.azure.prod

# modify the .env file with your secrets
vi .env.azure.prod

# source it to export the environment variables
source .env.azure.prod

# To clean up any golang test result cache data
go clean -testcache

go test -v ./... -timeout 2h
```

## On Windows

```bat
:: clone the repo
git clone https://github.com/karuppiah7890/tce-e2e-test

:: change directory to the project
cd tce-e2e-test

:: copy sameple .env file
copy .env.sample.azure.cmd .env.my-azure.cmd

:: modify the .env file with your secrets using a text editor
vi .env.my-azure.cmd

:: run it to export the environment variables
.env.my-azure.cmd

:: To clean up any golang test result cache data
go clean -testcache

:: To run the golang tests with high timeout so that test has enough time to conmplete
:: and Golang does not end them abruptly
go test -v ./... -timeout 2h
```
