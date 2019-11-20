package ebs

import (
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/stretchr/testify/assert"
)

var fail bool = true

type mockEC2Client struct {
	ec2iface.EC2API
}

func TestInit(t *testing.T) {
	var err error

	mockEC2 := &mockEC2Client{}

	ebsService := &EbsService{
		ec2Client: mockEC2,
		account:   "123456789",
		region:    "us-east-1",
	}

	err = ebsService.RunReport()
	assert.NotEqual(t, err, nil)

	fail = false

	err = ebsService.RunReport()
	assert.Equal(t, err, nil)

	NewEbsService("123456789", "us-east-1")
}

func (mock *mockEC2Client) DescribeVolumes(input *ec2.DescribeVolumesInput) (*ec2.DescribeVolumesOutput, error) {
	if fail {
		return nil, errors.New("TEST: DescribeVolumes")
	}

	t := time.Now()
	i := int64(10)

	return &ec2.DescribeVolumesOutput{
		Volumes: []*ec2.Volume{
			&ec2.Volume{
				VolumeId:   aws.String("abc123"),
				VolumeType: aws.String("gp2"),
				Size:       &i,
				CreateTime: &t,
			},
			&ec2.Volume{
				VolumeId:   aws.String("abc123"),
				VolumeType: aws.String("io1"),
				Size:       &i,
				Iops:       &i,
				CreateTime: &t,
			},
			&ec2.Volume{
				VolumeId:   aws.String("abc123"),
				VolumeType: aws.String("sc1"),
				Size:       &i,
				CreateTime: &t,
			},
			&ec2.Volume{
				VolumeId:   aws.String("abc123"),
				VolumeType: aws.String("st1"),
				Size:       &i,
				CreateTime: &t,
			},
			&ec2.Volume{
				VolumeId:   aws.String("abc123"),
				VolumeType: aws.String("standard"),
				Size:       &i,
				CreateTime: &t,
			},
		},
	}, nil
}
