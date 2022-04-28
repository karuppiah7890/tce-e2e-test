# tce-e2e-test

This repo contains [TCE](https://github.com/vmware-tanzu/community-edition) E2E tests

## Running the tests

<!-- TODO: Add instructions on how to run all the tests together - AWS, Azure, vSphere, Docker -->
<!-- TODO: Add instructions on how to run all the tests separately and independently - AWS, Azure, vSphere, Docker -->

Ensure you have latest golang `1.17.x` installed

```bash
git clone https://github.com/karuppiah7890/tce-e2e-test

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
