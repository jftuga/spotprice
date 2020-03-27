# spotprice
Get AWS spot instance pricing

The [Releases Page](https://github.com/jftuga/spotprice/releases) contains binaries for Windows, MacOS, Linux and FreeBSD.

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

spotprice.exe -reg us -inst t2.micro,t2.small

+-----------+------------+----------+------------+------------+
|  REGION   |     AZ     | INSTANCE |    DESC    | SPOT PRICE |
+-----------+------------+----------+------------+------------+
| us-west-2 | us-west-2c | t2.micro | Linux/UNIX |   0.003500 |
| us-west-2 | us-west-2b | t2.micro | Linux/UNIX |   0.003500 |
| us-west-2 | us-west-2a | t2.micro | Linux/UNIX |   0.003500 |
| us-east-2 | us-east-2c | t2.micro | Linux/UNIX |   0.003500 |
| us-east-2 | us-east-2b | t2.micro | Linux/UNIX |   0.003500 |
| us-east-2 | us-east-2a | t2.micro | Linux/UNIX |   0.003500 |
| us-east-1 | us-east-1f | t2.micro | Linux/UNIX |   0.003500 |
| us-east-1 | us-east-1e | t2.micro | Linux/UNIX |   0.003500 |
| us-east-1 | us-east-1d | t2.micro | Linux/UNIX |   0.003500 |
| us-east-1 | us-east-1c | t2.micro | Linux/UNIX |   0.003500 |
| us-east-1 | us-east-1b | t2.micro | Linux/UNIX |   0.003500 |
| us-east-1 | us-east-1a | t2.micro | Linux/UNIX |   0.003500 |
| us-west-1 | us-west-1b | t2.micro | Linux/UNIX |   0.004100 |
| us-west-1 | us-west-1a | t2.micro | Linux/UNIX |   0.004100 |
| us-west-2 | us-west-2c | t2.small | Linux/UNIX |   0.006900 |
| us-west-2 | us-west-2b | t2.small | Linux/UNIX |   0.006900 |
| us-west-2 | us-west-2a | t2.small | Linux/UNIX |   0.006900 |
| us-east-2 | us-east-2c | t2.small | Linux/UNIX |   0.006900 |
| us-east-2 | us-east-2b | t2.small | Linux/UNIX |   0.006900 |
| us-east-2 | us-east-2a | t2.small | Linux/UNIX |   0.006900 |
| us-east-1 | us-east-1e | t2.small | Linux/UNIX |   0.006900 |
| us-east-1 | us-east-1d | t2.small | Linux/UNIX |   0.006900 |
| us-east-1 | us-east-1c | t2.small | Linux/UNIX |   0.006900 |
| us-east-1 | us-east-1a | t2.small | Linux/UNIX |   0.006900 |
| us-east-1 | us-east-1b | t2.small | Linux/UNIX |   0.007700 |
| us-east-1 | us-east-1f | t2.small | Linux/UNIX |   0.008200 |
| us-west-1 | us-west-1b | t2.small | Linux/UNIX |   0.008300 |
| us-west-1 | us-west-1a | t2.small | Linux/UNIX |   0.008300 |
+-----------+------------+----------+------------+------------+
```
