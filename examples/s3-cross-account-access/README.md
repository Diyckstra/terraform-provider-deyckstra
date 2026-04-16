# S3 bucket cross-account access example

This example demonstrates how to create an S3 bucket in one account and grant access to that bucket to the user from another account using the bucket ACL.
To access infrastructure in different accounts, two providers are configured with different aliases.

The example uses several variables.
You can copy `terraform.template.tfvars` to `terraform.tfvars` and specify their values in this file.

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
	-var="second_secret_key=your-second-secret-key"
```

Destroying the example:

```
$ terraform destroy

# or
$ terraform destroy \
	-var="first_access_key=your-first-access-key" \
	-var="first_secret_key=your-first-secret-key" \
	-var="second_access_key=your-second-access-key" \
	-var="second_secret_key=your-second-secret-key"
```
