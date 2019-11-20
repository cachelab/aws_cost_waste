package elb

import (
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/elb/elbiface"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/aws/aws-sdk-go/service/elbv2/elbv2iface"
	"github.com/stretchr/testify/assert"
)

var failFindV2LoadBalancers bool = true
var failFindTargetGroups bool = true
var failFindTargetHealth bool = true
var failFindClassicLoadBalancers bool = true
var findInstanceHealth bool = true

type mockELBClient struct {
	elbiface.ELBAPI
}

type mockELBV2Client struct {
	elbv2iface.ELBV2API
}

func TestInit(t *testing.T) {
	var err error

	mockELB := &mockELBClient{}
	mockELBV2 := &mockELBV2Client{}

	elbService := &ElbService{
		elbClient:   mockELB,
		elbV2Client: mockELBV2,
		account:     "123456789",
		region:      "us-east-1",
	}

	err = elbService.RunReport()
	assert.NotEqual(t, err, nil)

	failFindV2LoadBalancers = false

	err = elbService.RunReport()
	assert.NotEqual(t, err, nil)

	failFindTargetGroups = false

	err = elbService.RunReport()
	assert.NotEqual(t, err, nil)

	failFindTargetHealth = false

	err = elbService.RunReport()
	assert.NotEqual(t, err, nil)

	failFindClassicLoadBalancers = false

	err = elbService.RunReport()
	assert.NotEqual(t, err, nil)

	findInstanceHealth = false

	err = elbService.RunReport()
	assert.Equal(t, err, nil)

	NewElbService("123456789", "us-east-1")
}

func (mock *mockELBV2Client) DescribeLoadBalancers(input *elbv2.DescribeLoadBalancersInput) (*elbv2.DescribeLoadBalancersOutput, error) {
	if failFindV2LoadBalancers {
		return nil, errors.New("TEST: DescribeLoadBalancers")
	}

	t := time.Now()

	return &elbv2.DescribeLoadBalancersOutput{
		LoadBalancers: []*elbv2.LoadBalancer{
			&elbv2.LoadBalancer{
				LoadBalancerName: aws.String("abc123"),
				Scheme:           aws.String("internet-facing"),
				CreatedTime:      &t,
				Type:             aws.String("application"),
				LoadBalancerArn:  aws.String("abc123"),
			},
			&elbv2.LoadBalancer{
				LoadBalancerName: aws.String("abc123"),
				Scheme:           aws.String("internet-facing"),
				CreatedTime:      &t,
				Type:             aws.String("network"),
				LoadBalancerArn:  aws.String("abc123"),
			},
		},
	}, nil
}

func (mock *mockELBV2Client) DescribeTargetGroups(input *elbv2.DescribeTargetGroupsInput) (*elbv2.DescribeTargetGroupsOutput, error) {
	if failFindTargetGroups {
		return nil, errors.New("TEST: DescribeTargetGroups")
	}

	return &elbv2.DescribeTargetGroupsOutput{
		TargetGroups: []*elbv2.TargetGroup{
			&elbv2.TargetGroup{
				TargetGroupArn: aws.String("abc123"),
			},
		},
	}, nil
}

func (mock *mockELBV2Client) DescribeTargetHealth(input *elbv2.DescribeTargetHealthInput) (*elbv2.DescribeTargetHealthOutput, error) {
	if failFindTargetHealth {
		return nil, errors.New("TEST: DescribeTargetHealth")
	}

	return &elbv2.DescribeTargetHealthOutput{}, nil
}

func (mock *mockELBClient) DescribeLoadBalancers(input *elb.DescribeLoadBalancersInput) (*elb.DescribeLoadBalancersOutput, error) {
	if failFindClassicLoadBalancers {
		return nil, errors.New("TEST: DescribeLoadBalancers")
	}

	t := time.Now()

	return &elb.DescribeLoadBalancersOutput{
		LoadBalancerDescriptions: []*elb.LoadBalancerDescription{
			&elb.LoadBalancerDescription{
				LoadBalancerName: aws.String("abc123"),
				Scheme:           aws.String("internet-facing"),
				CreatedTime:      &t,
			},
		},
	}, nil
}

func (mock *mockELBClient) DescribeInstanceHealth(input *elb.DescribeInstanceHealthInput) (*elb.DescribeInstanceHealthOutput, error) {
	if findInstanceHealth {
		return nil, errors.New("TEST: DescribeInstanceHealth")
	}

	return &elb.DescribeInstanceHealthOutput{}, nil
}
