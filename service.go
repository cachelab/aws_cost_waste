package main

import (
	"os"

	"aws_cost_waste/ebs"
	"aws_cost_waste/elb"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
)

type Service struct {
	stsClient  stsiface.STSAPI
	awsDetails *AwsDetails
}

type AwsDetails struct {
	account string
	region  string
}

func NewService() *Service {
	return &Service{
		stsClient: sts.New(session.New()),
	}
}

func (svc *Service) Init() error {
	input := &sts.GetCallerIdentityInput{}
	result, err := svc.stsClient.GetCallerIdentity(input)
	if err != nil {
		return err
	}

	svc.awsDetails = &AwsDetails{
		account: *result.Account,
		region:  os.Getenv("AWS_REGION"),
	}

	return nil
}

func (svc *Service) RunEbsReport() error {
	ebsService := ebs.NewEbsService(svc.awsDetails.account, svc.awsDetails.region)

	err := ebsService.RunReport()
	if err != nil {
		return err
	}

	return nil
}

func (svc *Service) RunElbReport() error {
	elbService := elb.NewElbService(svc.awsDetails.account, svc.awsDetails.region)

	err := elbService.RunReport()
	if err != nil {
		return err
	}

	return nil
}
