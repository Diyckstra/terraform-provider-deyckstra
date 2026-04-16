# EC2 transit gateway cross-account VPC attachment example

This example demonstrates how to create a transit gateway in one account and share it with a second account to attach it to a VPC.
To access infrastructure in different accounts, two providers are configured with different aliases.

The example uses several variables. You can copy `terraform.template.tfvars` to `terraform.tfvars` and specify their values in this file.

```
$ cp terraform.template.tfvars terraform.tfvars 
```

Running the example:

```
$ terraform init
$ terraform apply

# or
$ terraform apply \
	-var="first_access_key=your-first-access-key" \
	-var="first_secret_key=your-first-secret-key" \
	-var="second_access_key=your-second-access-key" \
	-var="second_secret_key=your-second-secret-key" \
	-var="second_account_id=your-second-account-id"
```

Destroying the example:

```
$ terraform destroy

# or
$ terraform destroy \
	-var="first_access_key=your-first-access-key" \
	-var="first_secret_key=your-first-secret-key" \
	-var="second_access_key=your-second-access-key" \
	-var="second_secret_key=your-second-secret-key" \
	-var="second_account_id=your-second-account-id"
```
