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

Running the example:

```
$ export AWS_ACCESS_KEY_ID="your-access-key"
$ export AWS_SECRET_ACCESS_KEY="your-secret-key"
$ terraform init
$ terraform apply
```

Destroying the example:

```
$ terraform destroy
```
