# spotprice
Get AWS spot instance pricing

This program is similar to using `aws ec2 describe-spot-price-history` but is faster and has a few more options.

The [Releases Page](https://github.com/jftuga/spotprice/releases) contains binaries for Windows, MacOS, Linux and FreeBSD.

### Usage

```

spotprice: Get AWS spot instance pricing
usage: spotprice [options]
       (required EC2 IAM Permissions: DescribeRegions, DescribeAvailabilityZones, DescribeSpotPriceHistory)

  -az string
    	A comma-separated list of regular-expressions to match AZs (eg: us-*1a)
  -inst string
    	A comma-separated list of exact Instance Type names (eg: t2.small,t3a.micro,c5.large)
  -max float
    	Only output if spot price is less than or equal to given amount
  -prod string
    	A comma-separated list of exact, case-sensitive Product Names (eg: Windows,Linux/UNIX,SUSE Linux,Red Hat Enterprise Linux)
  -reg string
    	A comma-separated list of regular-expressions to match regions (eg: us-.*2b)
  -v	show program version and then exit

```

### EC2 IAM Permissions
* DescribeAvailabilityZones
* DescribeRegions
* DescribeSpotPriceHistory

### Examples

#### Only return pricing for US and Canada regions; Windows OS, (these 4 instance types); less than or equal to $2.00; in AZs that end in either an a, b, or d (such as us-east-2b)

* spotprice.exe -reg us-,ca- -prod Windows -inst r5.8xlarge,x1.32xlarge,x1e.32xlarge,c4.8xlarge -max 2.00 -az "(a|b|d)$"

```
+--------------+---------------+------------+---------+------------+
|    REGION    |      AZ       |  INSTANCE  |  DESC   | SPOT PRICE |
+--------------+---------------+------------+---------+------------+
| us-east-2    | us-east-2b    | c4.8xlarge | Windows |   1.789600 |
| us-east-2    | us-east-2a    | c4.8xlarge | Windows |   1.789600 |
| us-east-2    | us-east-2b    | r5.8xlarge | Windows |   1.806500 |
| us-east-2    | us-east-2a    | r5.8xlarge | Windows |   1.806500 |
| ca-central-1 | ca-central-1d | r5.8xlarge | Windows |   1.980800 |
| ca-central-1 | ca-central-1b | r5.8xlarge | Windows |   1.980800 |
| ca-central-1 | ca-central-1a | r5.8xlarge | Windows |   1.980800 |
| us-west-2    | us-west-2b    | c4.8xlarge | Windows |   1.990600 |
| us-west-2    | us-west-2a    | c4.8xlarge | Windows |   1.990600 |
| us-east-1    | us-east-1d    | c4.8xlarge | Windows |   1.993500 |
| us-east-1    | us-east-1b    | c4.8xlarge | Windows |   1.993500 |
| us-east-1    | us-east-1a    | c4.8xlarge | Windows |   1.993500 |
+--------------+---------------+------------+---------+------------+
```

#### Only return pricing for all US regions with instance types of either t2.micro or t2.small

* spotprice.exe -reg us -inst t2.micro,t2.small

```
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
