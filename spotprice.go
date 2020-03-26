/*

spotprice.go
-John Taylor
Mar 26 2020

Get AWS spot instance pricing

*/

package main

import (
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/olekukonko/tablewriter"
)

type sportPriceHistory struct {
	AvailabilityZone   string
	InstanceType       string
	ProductDescription string
	SpotPrice          string
	Timestamp          string
}

func printRegions() {
	p := endpoints.AwsPartition()
	fmt.Printf("partition: %v\n\n", p)

	fmt.Println("Regions for", p.ID())
	for id := range p.Regions() {
		fmt.Println("*", id)
	}

	fmt.Println()
	fmt.Println("Services for", p.ID())
	fmt.Println()
	for id := range p.Services() {
		fmt.Println("*", id)
	}
}

func getRegions() []string {
	p := endpoints.AwsPartition()
	var regions []string
	for id := range p.Regions() {
		regions = append(regions, id)
	}
	return regions
}

// https://docs.aws.amazon.com/sdk-for-go/api/service/ec2/#EC2.DescribeSpotPriceHistory

const longForm = "Jan 2, 2006 at 3:04pm (MST)"

func getSpotPriceHistoryOLD() string {
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})
	svc := ec2.New(sess)
	startTime, _ := time.Parse(longForm, "Mar 26, 2020 at 7:38am (EDT)")
	endTime, _ := time.Parse(longForm, "Mar 26, 2020 at 8:38am (EDT)")
	input := &ec2.DescribeSpotPriceHistoryInput{
		EndTime: &endTime,
		InstanceTypes: []*string{
			aws.String("t2.nano"), aws.String("t2.micro"), aws.String("t2.small"), aws.String("t3a.nano"),
			//aws.String("t3a.micro"), aws.String("t3a.small"), aws.String("t3.nano"), aws.String("t3.micro"),
			//aws.String("t3.small"), aws.String("t1.micro"),
		},
		ProductDescriptions: []*string{
			aws.String("Linux/UNIX (Amazon VPC)"),
		},
		StartTime: &startTime,
	}

	result, err := svc.DescribeSpotPriceHistory(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println("Error 1:", aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println("Error 2:", err.Error())
		}
		return ""
	}

	createSpotInfoArray(*result)
	return fmt.Sprintf("%s", result)
}

func getSpotPriceHistory(region string) ec2.DescribeSpotPriceHistoryOutput {
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	svc := ec2.New(sess)
	startTime, _ := time.Parse(longForm, "Mar 26, 2020 at 7:38am (EDT)")
	endTime, _ := time.Parse(longForm, "Mar 26, 2020 at 8:38am (EDT)")
	input := &ec2.DescribeSpotPriceHistoryInput{
		EndTime: &endTime,
		InstanceTypes: []*string{
			aws.String("t2.nano"), aws.String("t2.micro"), aws.String("t2.small"), aws.String("t3a.nano"),
			//aws.String("t3a.micro"), aws.String("t3a.small"), aws.String("t3.nano"), aws.String("t3.micro"),
			//aws.String("t3.small"), aws.String("t1.micro"),
		},
		ProductDescriptions: []*string{
			aws.String("Linux/UNIX (Amazon VPC)"),
		},
		StartTime: &startTime,
	}

	result, err := svc.DescribeSpotPriceHistory(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println("Error 1:", aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println("Error 2:", err.Error())
		}
		item := new(ec2.DescribeSpotPriceHistoryOutput)
		return *item
	}

	//createSpotInfoArray(*result)
	return *result
}

func createSpotInfoArray(data ec2.DescribeSpotPriceHistoryOutput) [][]string {
	//fmt.Printf("length: %d\n", len(data.SpotPriceHistory))
	var table [][]string
	for i := range data.SpotPriceHistory {
		table = append(table, []string{*data.SpotPriceHistory[i].AvailabilityZone, *data.SpotPriceHistory[i].InstanceType, *data.SpotPriceHistory[i].ProductDescription, *data.SpotPriceHistory[i].SpotPrice})
		/*
			fmt.Printf("AvalabilityZone : %s\n", *data.SpotPriceHistory[i].AvailabilityZone)
			fmt.Printf("InstanceType    : %s\n", *data.SpotPriceHistory[i].InstanceType)
			fmt.Printf("ProductDesc     : %s\n", *data.SpotPriceHistory[i].ProductDescription)
			fmt.Printf("SpotPrice       : %s\n", *data.SpotPriceHistory[i].SpotPrice)
			fmt.Println("==========================================================")
		*/
	}
	return table
}

func outputTable(data [][]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"AZ", "Instance", "Desc", "Spot Price"})

	for _, v := range data {
		table.Append(v)
	}
	table.Render() // Send output
}

func inspectRegion(region string, describeCh chan [][]string) {
	var spotResults ec2.DescribeSpotPriceHistoryOutput
	spotResults = getSpotPriceHistory(region)
	table := createSpotInfoArray(spotResults)
	describeCh <- table
}

func main() {
	allRegions := getRegions()
	fmt.Printf("%v\n\n", allRegions)

	describeCh := make(chan [][]string)

	timeStart := time.Now()
	for _, region := range allRegions {
		go inspectRegion(region, describeCh)
	}

	for range allRegions {
		table := <-describeCh
		if len(table) > 0 {
			outputTable(table)
		}
	}

	elapsed := time.Since(timeStart)
	fmt.Println()
	fmt.Printf("elapsed time : %v\n", elapsed)
}
