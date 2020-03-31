package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"
)

// parseFloat
func pf(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		fmt.Fprintf(os.Stderr, "\nError converting '%s' to float\n\n", s)
		return 0.0000
	}
	return f
}

func sortSpotPrice(entry []spotPriceItem, ascending bool) {
	sort.SliceStable(entry, func(i, j int) bool {
		if ascending {
			return pf(entry[i].SpotPrice) < pf(entry[j].SpotPrice)
		}
		return pf(entry[i].SpotPrice) > pf(entry[j].SpotPrice)
	})
}

func sortAvailabilityZone(entry []spotPriceItem, ascending bool) {
	sort.SliceStable(entry, func(i, j int) bool {
		if ascending {
			return entry[i].AvailabilityZone < entry[j].AvailabilityZone
		}
		return entry[i].AvailabilityZone > entry[j].AvailabilityZone
	})
}

func sortInstanceType(entry []spotPriceItem, ascending bool) {
	sort.SliceStable(entry, func(i, j int) bool {
		if ascending {
			return entry[i].InstanceType < entry[j].InstanceType
		}
		return entry[i].InstanceType > entry[j].InstanceType
	})
}

func sortRegion(entry []spotPriceItem, ascending bool) {
	sort.SliceStable(entry, func(i, j int) bool {
		if ascending {
			return entry[i].Region < entry[j].Region
		}
		return entry[i].Region > entry[j].Region
	})
}

func sortProductDescription(entry []spotPriceItem, ascending bool) {
	sort.SliceStable(entry, func(i, j int) bool {
		if ascending {
			return entry[i].ProductDescription < entry[j].ProductDescription
		}
		return entry[i].ProductDescription > entry[j].ProductDescription
	})
}
