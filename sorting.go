package main

import "sort"

func sortSpotPrice(entry []spotPriceItem, ascending bool) {
	sort.SliceStable(entry, func(i, j int) bool {
		if ascending {
			return entry[i].SpotPrice < entry[j].SpotPrice
		}
		return entry[i].SpotPrice > entry[j].SpotPrice
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
