package main

import (
	"errors"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
	"github.com/stretchr/testify/assert"
)

var failGetCallerIdentity bool = true

type mockSTSClient struct {
	stsiface.STSAPI
}

func TestInit(t *testing.T) {
	var err error

	mockSTS := &mockSTSClient{}

	svc := Service{
		stsClient: mockSTS,
	}

	err = svc.Init()
	assert.NotEqual(t, err, nil)

	failGetCallerIdentity = false

	err = svc.Init()
	assert.Equal(t, err, nil)

	err = svc.RunEbsReport()
	assert.NotEqual(t, err, nil)

	err = svc.RunElbReport()
	assert.NotEqual(t, err, nil)

	os.Setenv("AWS_REGION", "fail")

	err = svc.RunEbsReport()
	assert.NotEqual(t, err, nil)

	err = svc.RunElbReport()
	assert.NotEqual(t, err, nil)

	NewService()
}

func (mock *mockSTSClient) GetCallerIdentity(input *sts.GetCallerIdentityInput) (*sts.GetCallerIdentityOutput, error) {
	if failGetCallerIdentity {
		return nil, errors.New("TEST: GetCallerIdentity")
	}

	return &sts.GetCallerIdentityOutput{
		Account: aws.String("abc123"),
	}, nil
}
