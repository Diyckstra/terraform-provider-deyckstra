# Networking Example

This example creates various network resources:

* VPCs in each of the two regions;
* one or two subnets in each VPC, depending on the number of availability zones in region;
* security groups for different networks.

This example also demonstrates the use of modules to create several copies of the same resource set with different arguments.
The child modules in this directory are:

* `region`: a module for all the network resources within a region. This module is instantiated once per region;
* `subnet`: a module for all the subnet resources within the given availability zone.
  This module is instantiated once or twice per region, depending on the number of availability zones.

The example takes credentials from a `c2rc.sh` file.
Get the file for your project and place it in this directory before running the example.

Running the example:

```shell
$ terraform init
$ terraform apply
```

Destroying the example:

```shell
$ terraform destroy
```
