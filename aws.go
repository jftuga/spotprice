package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

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

// remove any duplicate items in the given string slice
func removeDuplicatesUnordered(elements []string) []string {
	encountered := map[string]bool{}

	// Create a map of all unique elements.
	for v := range elements {
		encountered[elements[v]] = true
	}

	// Place all keys from the map into a slice.
	result := []string{}
	for key := range encountered {
		result = append(result, key)
	}
	return result
}

// when given a string slice (hackstack),
// only return items that match a regular expression (needle)
func match(haystack []string, needle string) []string {
	//fmt.Println("haystack  :", haystack)
	//fmt.Println("needle    :", needle)
	var found []string
	for _, a := range haystack {
		//fmt.Println("     a:", a)
		re, err := regexp.Compile(needle)
		if err != nil {
			fmt.Fprintf(os.Stderr, "\n\nregexp: %s\n%s\n", needle, err)
			os.Exit(1)
		}

		if re.MatchString(a) {
			found = append(found, a)
		}
	}
	return found
}

// return all AWS region names
func getAllRegions(ep map[string]endpoints.Region) []string {
	var allRegions []string
	for name := range ep {
		allRegions = append(allRegions, name)
	}
	return allRegions
}

// requested is a slice of regular expressions used to match the desired AWS regions
// return all Regions that match the regexp listed within each slice, removing any duplicate region names
func getDesiredRegions(requested string) []string {
	desiredRegions := strings.Split(requested, ",")
	//fmt.Println("desired  :", desiredRegions)
	for i, r := range desiredRegions {
		region := strings.ToLower(strings.TrimSpace(r))
		desiredRegions[i] = region
	}
	p := endpoints.AwsPartition()
	allRegions := getAllRegions(p.Regions())
	//fmt.Println("allRegions:", allRegions)
	var regions []string

	for _, f := range desiredRegions {
		m := match(allRegions, f)
		if len(m) > 0 {
			regions = append(regions, m...)
		}
	}

	return removeDuplicatesUnordered(regions)
}

func createAWSStringTypes(commaList string) []*string {
	var stringTypes []*string
	for _, item := range strings.Split(commaList, ",") {
		item := strings.TrimSpace(item) // also strings.ToLower ??
		stringTypes = append(stringTypes, aws.String(item))
	}
	return stringTypes
}

// https://docs.aws.amazon.com/sdk-for-go/api/service/ec2/#EC2.DescribeSpotPriceHistory
func getSpotPriceHistory(region string, input ec2.DescribeSpotPriceHistoryInput) ec2.DescribeSpotPriceHistoryOutput {

	utc, _ := time.LoadLocation("UTC")
	endTime := time.Now().In(utc)
	startTime := endTime.Add(-1 * time.Minute) // 1 minute ago
	input.SetStartTime(startTime)
	input.SetEndTime(endTime)

	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})

	svc := ec2.New(sess)
	result, err := svc.DescribeSpotPriceHistory(&input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				if !strings.Contains(aerr.Error(), "AuthFailure:") {
					fmt.Fprintf(os.Stderr, "Error 1: region: %s, %s", region, aerr.Error())
				}
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Fprintf(os.Stderr, "Error 2: region: %s, %s", region, err.Error())
		}
		item := new(ec2.DescribeSpotPriceHistoryOutput)
		return *item
	}
	//fmt.Println("getSpotPriceHistory() result:", *result)

	return *result
}

func createSpotInfoArray(region string, data ec2.DescribeSpotPriceHistoryOutput) ([][]string, []spotPriceItem) {
	//fmt.Printf("length: %d\n", len(data.SpotPriceHistory))
	var table [][]string
	var items []spotPriceItem
	for i := range data.SpotPriceHistory {
		item := new(spotPriceItem)
		item.Region = region
		item.AvailabilityZone = *data.SpotPriceHistory[i].AvailabilityZone
		item.InstanceType = *data.SpotPriceHistory[i].InstanceType
		item.ProductDescription = *data.SpotPriceHistory[i].ProductDescription
		item.SpotPrice = *data.SpotPriceHistory[i].SpotPrice

		items = append(items, *item)
	}

	for i := range data.SpotPriceHistory {
		table = append(table, []string{*data.SpotPriceHistory[i].AvailabilityZone, *data.SpotPriceHistory[i].InstanceType, *data.SpotPriceHistory[i].ProductDescription, *data.SpotPriceHistory[i].SpotPrice})
		//fmt.Printf("AvalabilityZone : %s\n", *data.SpotPriceHistory[i].AvailabilityZone)
		//fmt.Printf("InstanceType    : %s\n", *data.SpotPriceHistory[i].InstanceType)
		//fmt.Printf("ProductDesc     : %s\n", *data.SpotPriceHistory[i].ProductDescription)
		//fmt.Printf("SpotPrice       : %s\n", *data.SpotPriceHistory[i].SpotPrice)
		//fmt.Println("==========================================================")
	}
	return table, items
}
