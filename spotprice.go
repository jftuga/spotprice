/*

spotprice.go
-John Taylor
Mar 26 2020

Get AWS spot instance pricing

*/

package main

import (
	"bytes"
	"encoding/gob"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/olekukonko/tablewriter"
)

const version = "0.0.2"

type spotPriceItem struct {
	Region             string
	AvailabilityZone   string
	InstanceType       string
	ProductDescription string
	SpotPrice          string
}

func outputTable(allItems []spotPriceItem) {
	if 0 == len(allItems) {
		fmt.Fprintf(os.Stderr, "\n\nError: No spot instances found.\n")
		os.Exit(1)
	}
	var data [][]string
	for _, i := range allItems {
		item := []string{i.Region, i.AvailabilityZone, i.InstanceType, i.ProductDescription, i.SpotPrice}
		data = append(data, item)
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Region", "AZ", "Instance", "Desc", "Spot Price"})

	for _, v := range data {
		table.Append(v)
	}
	table.Render() // Send output
}

func inspectRegion(region string, instanceTypes string, describeCh chan []spotPriceItem) {
	var vmTypes []string
	vmTypes = strings.Split(instanceTypes, ",")
	var spotResults ec2.DescribeSpotPriceHistoryOutput
	spotResults = getSpotPriceHistory(region, vmTypes)
	_, items := createSpotInfoArray(region, spotResults)
	//fmt.Println("items:", len(items))
	describeCh <- items
}

func main() {
	argsVersion := flag.Bool("v", false, "show program version and then exit")
	argsDebug := flag.Bool("d", false, "run in debug mode")
	argsRegion := flag.String("reg", "", "A comma-separated list of regular-expressions to match regions (eg: us-*)")
	//argsAZ := flag.String("az", "", "a regular-expression to match AZs")
	argsInst := flag.String("inst", "", "A comma-separated list of exact Instance Type names (eg: t2.small,t3a.micro,c5.large")
	//argsProd := flag.String("prod", "", "A comma-separated list of exact Instance Type names (eg: Windows,Linux/UNIX,SUSE Linux,Red Hat Enterprise Linux")
	//argsOutput := flag.String("out", "", "Set output to 'csv' or 'json'")

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

	if 0 == len(*argsRegion) && 0 == len(*argsInst) {
		flag.Usage()
		os.Exit(1)
	}

	if 0 == len(*argsInst) {
		fmt.Fprintf(os.Stderr, "\nThe -inst option is required.\nThis limitation will be removed in a future release.\n\n")
		flag.Usage()
		os.Exit(1)
	}

	timeStart := time.Now()
	//*argsDebug = true

	var allSpotsAllRegions []spotPriceItem

	allRegions := getDesiredRegions(*argsRegion)
	//allRegions = allRegions[:3]
	//fmt.Printf("Regions: %v\n\n", allRegions)
	if 0 == len(allRegions) {
		fmt.Fprintf(os.Stderr, "\n\nError: No matching AWS regions.\n")
		os.Exit(1)
	}

	if !*argsDebug {
		describeCh := make(chan []spotPriceItem)
		for _, region := range allRegions {
			go inspectRegion(region, *argsInst, describeCh)
		}

		for range allRegions {
			table := <-describeCh
			if len(table) > 0 {
				for _, t := range table {
					allSpotsAllRegions = append(allSpotsAllRegions, t)
				}
			}
		}
	}

	if !*argsDebug {
		var bufOut bytes.Buffer
		enc := gob.NewEncoder(&bufOut)
		err := enc.Encode(allSpotsAllRegions)
		if err != nil {
			log.Fatal("encode error:", err)
			return
		}

		err = ioutil.WriteFile("current.dat", bufOut.Bytes(), 0600)
		if err != nil {
			log.Fatal("WriteFile error:", err)
			return
		}
	} else {
		bufIn, err := ioutil.ReadFile("current.dat")
		if err != nil {
			log.Fatal("ReadFile error:", err)
			return
		}
		//fmt.Println("bufIn:", len(bufIn))

		z := bytes.NewBuffer(bufIn)
		dec := gob.NewDecoder(z)
		err = dec.Decode(&allSpotsAllRegions)
		if err != nil {
			log.Fatal("decode error 1:", err)
		}
	}

	//fmt.Println("ok")

	//fmt.Printf("%v\n", allSpotsAllRegions)
	//sortRegion(allSpotsAllRegions, true)

	sortAvailabilityZone(allSpotsAllRegions, false)
	sortSpotPrice(allSpotsAllRegions, true)

	outputTable(allSpotsAllRegions)

	elapsed := time.Since(timeStart)
	fmt.Println()
	fmt.Printf("elapsed time : %v\n", elapsed)
}
