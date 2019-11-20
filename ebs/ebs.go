package ebs

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/olekukonko/tablewriter"
)

const (
	EBS_GP2_PRICING              = 0.1
	EBS_IO1_PRICING              = 0.125
	EBS_SC1_PRICING              = 0.025
	EBS_ST1_PRICING              = 0.045
	EBS_STANDARD_PRICING         = 0.05
	EBS_PROVISIONED_IOPS_PRICING = 0.065
)

type EbsService struct {
	ec2Client ec2iface.EC2API
	account   string
	region    string
}

type EbsWaste struct {
	VolumeId   *string
	VolumeType *string
	CreateTime string
	Size       *int64
	Cost       float64
}

func NewEbsService(account string, region string) *EbsService {
	return &EbsService{
		ec2Client: ec2.New(session.New()),
		account:   account,
		region:    region,
	}
}

func (svc *EbsService) RunReport() error {
	fmt.Println("\nEBS Waste Report")

	var total float64

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Account", "Region", "ID", "Volume Type", "Created", "Size", "Wasting"})

	volumes, err := svc.findAvailableVolumes()
	if err != nil {
		return err
	}

	for _, volume := range volumes {
		total = volume.Cost + total

		table.Append([]string{
			svc.account,
			svc.region,
			*volume.VolumeId,
			*volume.VolumeType,
			volume.CreateTime,
			fmt.Sprintf("%d GB", *volume.Size),
			fmt.Sprintf("$%.2f", volume.Cost),
		})
	}

	table.SetFooter([]string{"", "", "", "", "", "Monthly Total", fmt.Sprintf("$%.2f", total)})
	table.SetFooterColor(
		tablewriter.Colors{},
		tablewriter.Colors{},
		tablewriter.Colors{},
		tablewriter.Colors{},
		tablewriter.Colors{},
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.FgHiRedColor},
	)
	table.Render()

	return nil
}

// Private

func (svc *EbsService) findAvailableVolumes() ([]*EbsWaste, error) {
	input := &ec2.DescribeVolumesInput{Filters: []*ec2.Filter{
		&ec2.Filter{
			Name:   aws.String("status"),
			Values: []*string{aws.String("available")},
		},
	}}

	result, err := svc.ec2Client.DescribeVolumes(input)
	if err != nil {
		return nil, err
	}

	ebsWaste := []*EbsWaste{}

	for _, volume := range result.Volumes {
		var cost float64

		switch *volume.VolumeType {
		case "gp2":
			cost = float64(*volume.Size) * EBS_GP2_PRICING
			break
		case "io1":
			cost = (float64(*volume.Size) * EBS_IO1_PRICING) + (float64(*volume.Iops) * EBS_PROVISIONED_IOPS_PRICING)
			break
		case "sc1":
			cost = float64(*volume.Size) * EBS_SC1_PRICING
			break
		case "st1":
			cost = float64(*volume.Size) * EBS_ST1_PRICING
			break
		case "standard":
			cost = float64(*volume.Size) * EBS_STANDARD_PRICING
			break
		}

		ebsWaste = append(ebsWaste, &EbsWaste{
			VolumeId:   volume.VolumeId,
			VolumeType: volume.VolumeType,
			CreateTime: volume.CreateTime.Format("2006-01-02"),
			Size:       volume.Size,
			Cost:       cost,
		})
	}

	return ebsWaste, nil
}
