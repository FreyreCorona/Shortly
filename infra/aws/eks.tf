module "eks" {
  source  = "terraform-aws-modules/eks/aws"
  version = "21.25.0"
  name = "shortly"
  kubernetes_version = "1.35"
  
  vpc_id = module.vpc.vpc_id
  subnet_ids = module.vpc.private_subnets
  control_plane_subnet_ids = module.vpc.private_subnets

  endpoint_public_access = true
  endpoint_public_access_cidrs = ["181.224.50.187/32"]
  endpoint_private_access = true
  enable_cluster_creator_admin_permissions = true

  addons = {
    coredns = { }
    eks-pod-identity-agent = {before_compute = true}
    kube-proxy = { }
    vpc-cni = { before_compute = true }
  }
  eks_managed_node_groups = {
    main = {
        desired_size = 1
        min_size = 1
        max_size = 2
        instance_types = ["t3.small"]
        ami_type = "AL2023_x86_64_STANDARD"
    }
  }
}

