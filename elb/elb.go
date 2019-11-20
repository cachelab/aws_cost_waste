package elb

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/elb/elbiface"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/aws/aws-sdk-go/service/elbv2/elbv2iface"
	"github.com/olekukonko/tablewriter"
)

const (
	ALB_PRICING = 0.0225
	NLB_PRICING = 0.0225
	CLB_PRICING = 0.025
)

type ElbService struct {
	elbClient   elbiface.ELBAPI
	elbV2Client elbv2iface.ELBV2API
	account     string
	region      string
}

type ElbWaste struct {
	LoadBalancerName *string
	Scheme           *string
	CreatedTime      string
	Type             *string
	Cost             float64
	LoadBalancerArn  *string
	TargetGroups     []*string
}

func NewElbService(account string, region string) *ElbService {
	return &ElbService{
		elbClient:   elb.New(session.New()),
		elbV2Client: elbv2.New(session.New()),
		account:     account,
		region:      region,
	}
}

func (svc *ElbService) RunReport() error {
	fmt.Println("\nELB Waste Report")

	var total float64

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Account", "Region", "Name", "Scheme", "Created", "Type", "Wasting"})

	elbWaste, err := svc.findV2LoadBalancers()
	if err != nil {
		return err
	}

	err = svc.findTargetGroups(elbWaste)
	if err != nil {
		return err
	}

	actualWaste, err := svc.findTargetHealth(elbWaste)
	if err != nil {
		return err
	}

	elbWaste, err = svc.findClassicLoadBalancers()
	if err != nil {
		return err
	}

	elbWaste, err = svc.findInstanceHealth(elbWaste)
	if err != nil {
		return err
	}

	for _, waste := range elbWaste {
		actualWaste = append(actualWaste, waste)
	}

	for _, loadBalancer := range actualWaste {
		total = loadBalancer.Cost + total

		table.Append([]string{
			svc.account,
			svc.region,
			*loadBalancer.LoadBalancerName,
			*loadBalancer.Scheme,
			loadBalancer.CreatedTime,
			*loadBalancer.Type,
			fmt.Sprintf("$%.2f", loadBalancer.Cost),
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

func (svc *ElbService) findV2LoadBalancers() ([]*ElbWaste, error) {
	input := &elbv2.DescribeLoadBalancersInput{}

	result, err := svc.elbV2Client.DescribeLoadBalancers(input)
	if err != nil {
		return nil, err
	}

	var elbWaste []*ElbWaste

	for _, loadBalancer := range result.LoadBalancers {
		var cost float64

		switch *loadBalancer.Type {
		case "application":
			cost = 24 * 30 * ALB_PRICING
			break
		case "network":
			cost = 24 * 30 * NLB_PRICING
			break
		}

		elbWaste = append(elbWaste, &ElbWaste{
			LoadBalancerName: loadBalancer.LoadBalancerName,
			Scheme:           loadBalancer.Scheme,
			CreatedTime:      loadBalancer.CreatedTime.Format("2006-01-02"),
			Type:             loadBalancer.Type,
			Cost:             cost,
			LoadBalancerArn:  loadBalancer.LoadBalancerArn,
		})
	}

	return elbWaste, nil
}

func (svc *ElbService) findTargetGroups(elbWaste []*ElbWaste) error {
	for _, waste := range elbWaste {
		var targetGroups []*string

		input := &elbv2.DescribeTargetGroupsInput{
			LoadBalancerArn: waste.LoadBalancerArn,
		}

		result, err := svc.elbV2Client.DescribeTargetGroups(input)
		if err != nil {
			return err
		}

		for _, targetGroup := range result.TargetGroups {
			targetGroups = append(targetGroups, targetGroup.TargetGroupArn)
		}

		waste.TargetGroups = targetGroups
	}

	return nil
}

func (svc *ElbService) findTargetHealth(elbWaste []*ElbWaste) ([]*ElbWaste, error) {
	var actualWaste []*ElbWaste

	for _, waste := range elbWaste {
		wasted := true

		for _, target := range waste.TargetGroups {
			input := &elbv2.DescribeTargetHealthInput{
				TargetGroupArn: target,
			}

			result, err := svc.elbV2Client.DescribeTargetHealth(input)
			if err != nil {
				return nil, err
			}

			if len(result.TargetHealthDescriptions) != 0 {
				wasted = false
			}
		}

		if wasted {
			actualWaste = append(actualWaste, waste)
		}
	}

	return actualWaste, nil
}

func (svc *ElbService) findClassicLoadBalancers() ([]*ElbWaste, error) {
	input := &elb.DescribeLoadBalancersInput{}

	result, err := svc.elbClient.DescribeLoadBalancers(input)
	if err != nil {
		return nil, err
	}

	var elbWaste []*ElbWaste

	for _, loadBalancer := range result.LoadBalancerDescriptions {
		elbWaste = append(elbWaste, &ElbWaste{
			LoadBalancerName: loadBalancer.LoadBalancerName,
			Scheme:           loadBalancer.Scheme,
			CreatedTime:      loadBalancer.CreatedTime.Format("2006-01-02"),
			Type:             aws.String("classic"),
			Cost:             24 * 30 * CLB_PRICING,
			LoadBalancerArn:  aws.String(""),
		})
	}

	return elbWaste, nil
}

func (svc *ElbService) findInstanceHealth(elbWaste []*ElbWaste) ([]*ElbWaste, error) {
	var actualWaste []*ElbWaste

	for _, waste := range elbWaste {
		input := &elb.DescribeInstanceHealthInput{
			LoadBalancerName: waste.LoadBalancerName,
		}

		result, err := svc.elbClient.DescribeInstanceHealth(input)
		if err != nil {
			return nil, err
		}

		if len(result.InstanceStates) == 0 {
			actualWaste = append(actualWaste, waste)
		}
	}

	return actualWaste, nil
}
