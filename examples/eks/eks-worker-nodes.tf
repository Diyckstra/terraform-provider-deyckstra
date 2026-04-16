resource "aws_eks_node_group" "example" {
  cluster_name    = aws_eks_cluster.example.name
  instance_types  = ["c5.large"]
  node_group_name = "terraform-eks-example"
  subnet_ids      = aws_subnet.example[*].id

  scaling_config {
    desired_size = 1
    max_size     = 1
    min_size     = 1
  }

  tags = {
    Name = "terraform-eks-example"
  }
}
