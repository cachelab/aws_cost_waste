# AWS Cost Waste

Reviews EBS volumes that are available as well as ELBs that do not have
instances attached.

[![CircleCI](https://circleci.com/gh/cachelab/aws_cost_waste.svg?style=svg)](https://circleci.com/gh/cachelab/aws_cost_waste)

## Example

![alt text](/images/screenshot.png)

## Contributing

* `make run` - runs aws_cost_waste in a docker container
* `make build` - builds your aws_cost_waste docker container
* `make vet` - go fmt and vet code
* `make test` - run unit tests

Before you submit a pull request please update the semantic version inside of
`main.go` with what you feel is appropriate and then edit the `CHANGELOG.md`
with your changes and follow a similar structure to what is there.
