/*

spotprice.go
-John Taylor
Mar 26 2020

Get AWS spot instance pricing

*/

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
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

func getSpotPriceHistory() string {
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
			aws.String("t3a.micro"), aws.String("t3a.small"), aws.String("t3.nano"), aws.String("t3.micro"),
			aws.String("t3.small"), aws.String("t1.micro"),
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

	extractSpotInfo(*result)
	return fmt.Sprintf("%s", result)
}

func extractSpotInfo(data ec2.DescribeSpotPriceHistoryOutput) {
	fmt.Printf("extract: %d\n", len(data.SpotPriceHistory))
	for i := range data.SpotPriceHistory {
		fmt.Println("%V\n", i)
	}
}

func spotPriceToJSON(raw string) {
	reader := strings.NewReader(raw)
	dec := json.NewDecoder(reader)
	for {
		// Read one JSON object and store it in a map.
		var m map[string]interface{}
		if err := dec.Decode(&m); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		for k := range m {
			fmt.Println(k)
		}
	}
}

func main() {
	//allRegions := getRegions()
	//fmt.Printf("%s\n", allRegions)
	getSpotPriceHistory()
	//spotPriceToJSON(raw)
}
