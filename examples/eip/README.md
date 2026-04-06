# Elastic IP Example

The EIP example launches the nginx web server and creates a security group.

This example assumes that an SSH key pair was created in the cloud console.

Running the example:

```
$ export AWS_ACCESS_KEY_ID="your-access-key"
$ export AWS_SECRET_ACCESS_KEY="your-secret-key"
$ terraform init
$ terraform apply -var="key_name=your-key-name"
```

In a few minutes, you can reach the nginx welcome page at the Elastic IP address.

Destroying the example:

```
$ terraform destroy -var="key_name=your-key-name"
```

Instead of using `-var`, you can copy `terraform.template.tfvars` to `terraform.tfvars` and use it to specify variable values.
