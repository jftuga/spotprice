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

const version = "1.0.0"

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

	var returnAll = true
	if maxPrice > 0.00000 {
		returnAll = false
	}

	var data [][]string
	for _, i := range allItems {
		item := []string{i.Region, i.AvailabilityZone, i.InstanceType, i.ProductDescription, i.SpotPrice}
		if (returnAll) || pf(i.SpotPrice) <= maxPrice {
			data = append(data, item)
		}
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Region", "AZ", "Instance", "Desc", "Spot Price"})

	for _, v := range data {
		table.Append(v)
	}
	if 0 == len(data) {
		fmt.Fprintf(os.Stderr, "\n\nError: No spot instances at or below -max price of %f\n", maxPrice)
		os.Exit(1)
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
	argsList := flag.Bool("list", false, "List regions, services and then exit")
	argsRegion := flag.String("reg", "", "A comma-separated list of regular-expressions to match regions (eg: us-.*2b)")
	argsAZ := flag.String("az", "", "A comma-separated list of regular-expressions to match AZs (eg: us-*1a)")
	argsInst := flag.String("inst", "", "A comma-separated list of exact Instance Type names (eg: t2.small,t3a.micro,c5.large)")
	argsProd := flag.String("prod", "", "A comma-separated list of exact, case-sensitive Product Names (eg: Windows,win,Linux/UNIX,lin,SUSE Linux,Red Hat Enterprise Linux)")
	argsMaxPrice := flag.Float64("max", 0.00, "Only output if spot price is less than or equal to given amount")
	//argsOutput := flag.String("out", "", "Set output to 'csv' or 'json'") // to do

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
	if 1 == len(os.Args) {
		flag.Usage()
		os.Exit(1)
	}

	if *argsVersion {
		fmt.Fprintf(os.Stderr, "version %s\n", version)
		os.Exit(0)
	}

	if *argsList {
		ListAllInfo()
		os.Exit(0)
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
			if *argsProd == "lin" {
				*argsProd = "Linux/UNIX"
			} else if *argsProd == "win" {
				*argsProd = "Windows"
			}
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

	outputTable(allSpotsAllRegions, *argsMaxPrice)

	elapsed := time.Since(timeStart)
	fmt.Println()
	fmt.Printf("elapsed time : %v\n", elapsed)
}
