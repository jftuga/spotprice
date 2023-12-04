/*
spotprice.go
-John Taylor
2023-12-02

version 2
=========
* Now uses the AWS GO SDK v2
* Output all instance types as well as all regions, use -l
* You can now use regular expressions to match multiple instance types, use -I
* Shortcuts for products, use -prod
* * lin => Linux/UNIX
* * red => Red Hat Enterprise Linux
* * suse => SUSE Linux
* * win => Windows

This code is feature completed, but still needs a few more things such as:
* variable name refactoring for more consistency
* better function documentation

Examples
========
go run .\spotprice.go -I "c\d.xlarge" -max 0.2 -prod "lin" -reg "us-.*-1" -az "[bf]$"

go build -ldflags="-s -w" .\spotprice.go
.\spotprice.exe -I "c\d.xlarge" -max 0.38 -prod "lin,red" -reg "us-.*-1,eu-west.*" -az "[bdf]$" -prof someone

*/

package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/olekukonko/tablewriter"
	"os"
	"regexp"
	"slices"
	"sort"
	"strings"
	"sync"
	"time"
)

const pgmName string = "spotprice"
const pgmVersion string = "2.0.0"
const pgmUrl = "https://github.com/jftuga/spotprice"
const pgmDescription = "Quickly get AWS spot instance pricing across multiple regions"

type pgmArgs struct {
	Profile  string
	Region   string
	AZ       string
	Inst     string
	InstRE   string
	Prod     string
	MaxPrice float64
}

type spotInfo struct {
	Region   string
	AZ       string
	Instance string
	Product  string
	Price    string
}

// when given a string slice (hackstack),
// only return items that match a regular expression (needle)
func matcher(haystack []string, needles string) []string {
	allRegExprs := strings.Split(needles, ",")
	var found []string
	for _, a := range haystack {
		for _, regExpr := range allRegExprs {
			re, err := regexp.Compile(regExpr)
			if err != nil {
				fmt.Fprintf(os.Stderr, "\n\nregexp: %s\n%s\n", regExpr, err)
				os.Exit(1)
			}

			if re.MatchString(a) {
				found = append(found, a)
			}
		}
	}
	return found
}

func getTime() *time.Time {
	utc, _ := time.LoadLocation("UTC")
	now := time.Now().In(utc)
	return &now
}

func processArgs() (pgmArgs, bool) {
	var usedArgs pgmArgs
	argsVersion := flag.Bool("v", false, "show program version and then exit")
	argsList := flag.Bool("l", false, "List regions & instance types, then exit")
	flag.StringVar(&usedArgs.Profile, "prof", "default", "AWS profile to use")
	flag.StringVar(&usedArgs.Region, "reg", "", "A comma-separated list of regular-expressions to match regions, eg: us-.*-2,ap-.*east-\\d")
	flag.StringVar(&usedArgs.AZ, "az", "", "A comma-separated list of regular-expressions to match AZs (eg: us-*1a)")
	flag.StringVar(&usedArgs.Inst, "inst", "", "A comma-separated list of exact Instance Type names, eg: t2.small,t3a.micro,c5.large")
	flag.StringVar(&usedArgs.InstRE, "I", "", "A comma-separated list of regular-expressions to match Instance Type names, eg: t2.*,c5(a\\.|n\\.|\\.)4xlarge")
	flag.StringVar(&usedArgs.Prod, "prod", "", "A comma-separated list of exact, case-sensitive Product Names (eg: Windows,win,Linux/UNIX,lin,SUSE Linux,Red Hat Enterprise Linux)")
	flag.Float64Var(&usedArgs.MaxPrice, "max", 0.00, "Only output if spot price is less than or equal to given amount")

	flag.Usage = func() {
		pgmName := os.Args[0]
		if strings.HasPrefix(os.Args[0], "./") {
			pgmName = os.Args[0][2:]
		}
		fmt.Fprintf(os.Stderr, "\n\n%s: %s\n\n", pgmName, pgmDescription)
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
		fmt.Printf("%v v%v\n%v\n", pgmName, pgmVersion, pgmUrl)
		os.Exit(0)
	}
	return usedArgs, *argsList
}

func getRegionList(region, profile string, listOnly bool) []string {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region), config.WithSharedConfigProfile(profile))
	if err != nil {
		fmt.Fprintf(os.Stderr, "\n%v\n", err)
		os.Exit(255)
	}

	ctx := context.TODO()
	ec2Client := ec2.NewFromConfig(cfg)
	regionInfo, err := ec2Client.DescribeRegions(ctx, &ec2.DescribeRegionsInput{AllRegions: aws.Bool(false)})
	if err != nil {
		fmt.Fprintf(os.Stderr, "\n%v\n", err)
		os.Exit(255)
	}

	allRegions := make([]string, len(regionInfo.Regions))
	for i, regionName := range regionInfo.Regions {
		allRegions[i] = *regionName.RegionName
	}
	sort.Strings(allRegions)

	if listOnly {
		for _, region := range allRegions {
			fmt.Printf("%v\n", region)
		}
	}
	return allRegions

}

// return a sorted list of all EC2 Instance Types
// this make at least 6 AWS API calls, for pagination
func getAllInstanceTypes(region, profile string, listOnly bool) []string {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region), config.WithSharedConfigProfile(profile))
	if err != nil {
		fmt.Fprintf(os.Stderr, "\n%v\n", err)
		os.Exit(255)
	}

	ctx := context.TODO()
	ec2Client := ec2.NewFromConfig(cfg)
	var allInstanceNames []string
	paginator := ec2.NewDescribeInstanceTypesPaginator(ec2Client, &ec2.DescribeInstanceTypesInput{})
	for paginator.HasMorePages() {
		instanceTypeInfo, err := paginator.NextPage(ctx)
		//fmt.Printf("PAGE: %v\n", len(instanceTypeInfo.InstanceTypes))
		if err != nil {
			fmt.Fprintf(os.Stderr, "\n%v\n", err)
			os.Exit(255)
		}
		for _, i := range instanceTypeInfo.InstanceTypes {
			allInstanceNames = append(allInstanceNames, string(i.InstanceType))
		}
	}
	sort.Strings(allInstanceNames)

	if listOnly {
		for _, name := range allInstanceNames {
			fmt.Printf("%v\n", name)
		}
	}
	return allInstanceNames
}

func getRegionalSpotInfo(region, profile, commaDelimitedInstances, commaDelimitedProducts, commaDelimitedAZ, maxPrice string, all map[string][]spotInfo) {
	allInstances := strings.Split(commaDelimitedInstances, ",")
	allProducts := strings.Split(commaDelimitedProducts, ",")
	for i := range allProducts {
		switch allProducts[i] {
		case "lin":
			allProducts[i] = "Linux/UNIX"
		case "red":
			allProducts[i] = "Red Hat Enterprise Linux"
		case "suse":
			allProducts[i] = "SUSE Linux"
		case "win":
			allProducts[i] = "Windows"
		}
	}
	var allInstanceTypes []types.InstanceType
	for _, instance := range allInstances {
		allInstanceTypes = append(allInstanceTypes, types.InstanceType(instance))
	}
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region), config.WithSharedConfigProfile(profile))
	if err != nil {
		fmt.Fprintf(os.Stderr, "\n%v\n", err)
		os.Exit(255)
	}

	ctx := context.TODO()
	ec2Client := ec2.NewFromConfig(cfg)
	input := ec2.DescribeSpotPriceHistoryInput{InstanceTypes: allInstanceTypes, ProductDescriptions: allProducts, StartTime: getTime()}
	v, err := ec2Client.DescribeSpotPriceHistory(ctx, &input)

	for _, info := range v.SpotPriceHistory {
		if *info.SpotPrice > maxPrice && maxPrice != "-1" {
			continue
		}
		var matchedAZ []string
		if len(commaDelimitedAZ) > 0 {
			matchedAZ = matcher([]string{*info.AvailabilityZone}, commaDelimitedAZ)
		}
		if len(commaDelimitedAZ) == 0 || (len(commaDelimitedAZ) > 0 && slices.Contains(matchedAZ, *info.AvailabilityZone)) {
			all[region] = append(all[region], spotInfo{region, *info.AvailabilityZone, string(info.InstanceType), string(info.ProductDescription), *info.SpotPrice})
		}
	}
}

func sortResults(allRegionSpotInfo map[string][]spotInfo) []spotInfo {
	var mergedSpotInfo []spotInfo
	for region := range allRegionSpotInfo {
		for _, info := range allRegionSpotInfo[region] {
			mergedSpotInfo = append(mergedSpotInfo, info)
		}
	}

	sort.SliceStable(mergedSpotInfo, func(i, j int) bool {
		if mergedSpotInfo[i].AZ < mergedSpotInfo[j].AZ {
			return true
		} else {
			return false
		}
	})

	sort.SliceStable(mergedSpotInfo, func(i, j int) bool {
		if mergedSpotInfo[i].Price < mergedSpotInfo[j].Price {
			return true
		} else {
			return false
		}
	})

	return mergedSpotInfo
}

func main() {
	args, listOnly := processArgs()
	if listOnly {
		getRegionList(args.Region, args.Profile, listOnly)
		fmt.Println("---")
		getAllInstanceTypes(args.Region, args.Profile, true)
		return
	}

	var matchedRegions []string
	if len(args.Region) > 0 {
		matchedRegions = matcher(getRegionList("us-east-1", args.Profile, false), args.Region)
	} else {
		matchedRegions = getRegionList("us-east-1", args.Profile, false)
	}

	var matchedInstanceNames string
	if len(args.InstRE) > 0 {
		matchedInstanceNames = strings.Join(matcher(getAllInstanceTypes("us-east-1", args.Profile, false), args.InstRE), ",")
	} else {
		matchedInstanceNames = args.Inst
	}

	var maxPrice string
	if args.MaxPrice > 0 {
		maxPrice = fmt.Sprintf("%f", args.MaxPrice)
	} else {
		maxPrice = "-1"
	}

	allRegionSpotInfo := make(map[string][]spotInfo)
	var wg sync.WaitGroup
	for _, region := range matchedRegions {
		wg.Add(1)
		allRegionSpotInfo[region] = []spotInfo{}
		go func(region, profile, inst, allProducts, filteredAZ, maxPrice string, allRegionSpotInfo map[string][]spotInfo) {
			defer wg.Done()
			getRegionalSpotInfo(region, profile, inst, allProducts, filteredAZ, maxPrice, allRegionSpotInfo)
		}(region, args.Profile, matchedInstanceNames, args.Prod, args.AZ, maxPrice, allRegionSpotInfo)
	}
	wg.Wait()

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Region", "AZ", "Instance", "Product", "Spot Price"})

	sortedSportInfoResults := sortResults(allRegionSpotInfo)
	for _, info := range sortedSportInfoResults {
		table.Append([]string{string(info.Region), string(info.AZ), string(info.Instance), string(info.Product), string(info.Price)})
	}

	table.Render()
}
