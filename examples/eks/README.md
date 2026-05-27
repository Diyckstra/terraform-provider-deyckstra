[Terraform http provider]: https://www.terraform.io/docs/providers/http/index.html

# EKS Example

The EKS (Elastic Kubernetes Service) example launches an EKS cluster and EKS node group with necessary network resources in all availability zones in the region.

Running the example:

```
$ export AWS_ACCESS_KEY_ID="your-access-key"
$ export AWS_SECRET_ACCESS_KEY="your-secret-key"
$ terraform init
$ terraform apply
```

Get kubeconfig for the created cluster:

```
$ terraform output -raw kubeconfig > ~/.kube/config
```

Destroying the example:

```
$ terraform destroy
```

This example uses [Terraform http provider] to send a request to https://icanhazip.com/ to determine the local workstation external IP for the security group configuration.
This request is optional and can be replaced by specifying the IP address manually in the [workstation-external-ip.tf](workstation-external-ip.tf).
