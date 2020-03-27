# spotprice
Get AWS spot instance pricing



### Usage

```
spotprice: Get AWS spot instance pricing
usage: spotprice [options]
       (required EC2 IAM Permissions: DescribeRegions, DescribeAvailabilityZones, DescribeSpotPriceHistory)

  -d	run in debug mode
  -inst string
    	A comma-separated list of exact Instance Type names (eg: t2.small,t3a.micro,c5.large
  -reg string
    	A comma-separated list of regular-expressions to match regions (eg: us-*)
  -v	show program version and then exit
```

### EC2 IAM Permissions
* DescribeAvailabilityZones
* DescribeRegions
* DescribeSpotPriceHistory

### Example output

```
+------------+----------+------------+------------+
|     AZ     | INSTANCE |    DESC    | SPOT PRICE |
+------------+----------+------------+------------+
| us-east-1f | t2.small | Linux/UNIX |   0.007900 |
| us-east-1a | t3a.nano | Linux/UNIX |   0.002200 |
| us-east-1f | t3a.nano | Linux/UNIX |   0.001800 |
| us-east-1c | t3a.nano | Linux/UNIX |   0.001800 |
| us-east-1a | t2.small | Linux/UNIX |   0.006900 |
| us-east-1d | t2.small | Linux/UNIX |   0.006900 |
| us-east-1e | t2.small | Linux/UNIX |   0.006900 |
| us-east-1c | t2.small | Linux/UNIX |   0.006900 |
| us-east-1b | t3a.nano | Linux/UNIX |   0.001500 |
| us-east-1d | t3a.nano | Linux/UNIX |   0.002100 |
| us-east-1a | t2.micro | Linux/UNIX |   0.003500 |
| us-east-1f | t2.micro | Linux/UNIX |   0.003500 |
| us-east-1d | t2.micro | Linux/UNIX |   0.003500 |
| us-east-1e | t2.micro | Linux/UNIX |   0.003500 |
| us-east-1c | t2.micro | Linux/UNIX |   0.003500 |
| us-east-1b | t2.micro | Linux/UNIX |   0.003500 |
| us-east-1b | t2.small | Linux/UNIX |   0.007700 |
+------------+----------+------------+------------+
```
