/*

spotprice.go
-John Taylor
Mar 26 2020

Get AWS spot instance pricing

*/

package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/olekukonko/tablewriter"
)

const version = "0.1.1"

type spotPriceItem struct {
	Region             string
	AvailabilityZone   string
	InstanceType       string
	ProductDescription string
	SpotPrice          string
}

func outputTable(allItems []spotPriceItem, maxPrice float64) {
	if 0 == len(allItems) {
		fmt.Fprintf(os.Stderr, "\n\nError: No spot instances found.\n")
		os.Exit(1)
	}

	var data [][]string
	for _, i := range allItems {
		if pf(i.SpotPrice) <= maxPrice {
			item := []string{i.Region, i.AvailabilityZone, i.InstanceType, i.ProductDescription, i.SpotPrice}
			data = append(data, item)
		}
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Region", "AZ", "Instance", "Desc", "Spot Price"})

	for _, v := range data {
		table.Append(v)
	}
	table.Render() // Send output
}

func describeRegion(region string, regionSPH ec2.DescribeSpotPriceHistoryInput, describeCh chan []spotPriceItem) {
	var spotResults ec2.DescribeSpotPriceHistoryOutput
	spotResults = getSpotPriceHistory(region, regionSPH)
	//fmt.Println("describeRegion() spotResults:", spotResults)
	_, items := createSpotInfoArray(region, spotResults)
	describeCh <- items
}

func filterAvailabilityZones(allSpotsAllRegions []spotPriceItem, rawAZs string) []spotPriceItem {
	filteredAZ := strings.Split(rawAZs, ",")
	var filterSpotRegions []spotPriceItem
	for _, zoneRE := range filteredAZ {
		re, err := regexp.Compile(zoneRE)
		if err != nil {
			fmt.Fprintf(os.Stderr, "\n\nregexp: %s\n%s\n", zoneRE, err)
			os.Exit(1)
		}
		for _, entry := range allSpotsAllRegions {
			if re.MatchString(entry.AvailabilityZone) {
				//fmt.Printf("%s matches %s\n", entry, zoneRE)
				filterSpotRegions = append(filterSpotRegions, entry)
			}
		}
	}
	return filterSpotRegions
}

func main() {
	argsVersion := flag.Bool("v", false, "show program version and then exit")
	argsRegion := flag.String("reg", "", "A comma-separated list of regular-expressions to match regions (eg: us-.*2b)")
	argsAZ := flag.String("az", "", "A comma-separated list of regular-expressions to match AZs (eg: us-*1a)")
	argsInst := flag.String("inst", "", "A comma-separated list of exact Instance Type names (eg: t2.small,t3a.micro,c5.large)")
	argsProd := flag.String("prod", "", "A comma-separated list of exact, case-sensitive Product Names (eg: Windows,Linux/UNIX,SUSE Linux,Red Hat Enterprise Linux)")
	argsLessThan := flag.Float64("less", 0.00, "Only output if price is less than or equal to given amount")
	//argsOutput := flag.String("out", "", "Set output to 'csv' or 'json'") // to do
	// maybe add option for AMI to generate a CF for spot instances

	flag.Usage = func() {
		pgmName := os.Args[0]
		if strings.HasPrefix(os.Args[0], "./") {
			pgmName = os.Args[0][2:]
		}
		fmt.Fprintf(os.Stderr, "\n%s: Get AWS spot instance pricing\n", pgmName)
		fmt.Fprintf(os.Stderr, "usage: %s [options]\n", pgmName)
		fmt.Fprintf(os.Stderr, "       (required EC2 IAM Permissions: DescribeRegions, DescribeAvailabilityZones, DescribeSpotPriceHistory)\n\n")
		flag.PrintDefaults()
	}

	flag.Parse()
	if *argsVersion {
		fmt.Fprintf(os.Stderr, "version %s\n", version)
		os.Exit(1)
	}

	timeStart := time.Now()
	var allSpotsAllRegions []spotPriceItem

	allRegions := getDesiredRegions(*argsRegion)
	if 0 == len(allRegions) {
		fmt.Fprintf(os.Stderr, "\n\nError: No matching AWS regions.\n")
		os.Exit(1)
	}

	describeCh := make(chan []spotPriceItem)
	for _, region := range allRegions {
		//fmt.Println("main() region:", region)
		regionSPH := new(ec2.DescribeSpotPriceHistoryInput)
		if len(*argsInst) > 0 {
			regionSPH.SetInstanceTypes(createAWSStringTypes(*argsInst))
		}
		if len(*argsProd) > 0 {
			regionSPH.SetProductDescriptions(createAWSStringTypes(*argsProd))
		}
		go describeRegion(region, *regionSPH, describeCh)
	}

	for range allRegions {
		table := <-describeCh
		if len(table) > 0 {
			for _, t := range table {
				allSpotsAllRegions = append(allSpotsAllRegions, t)
			}
		}
	}

	if len(*argsAZ) > 0 {
		allSpotsAllRegions = filterAvailabilityZones(allSpotsAllRegions, *argsAZ)
	}

	sortAvailabilityZone(allSpotsAllRegions, false)
	sortSpotPrice(allSpotsAllRegions, true)

	outputTable(allSpotsAllRegions, *argsLessThan)

	elapsed := time.Since(timeStart)
	fmt.Println()
	fmt.Printf("elapsed time : %v\n", elapsed)
}
