module "vpc" {
    name = "shortly-vpc"
    source = "terraform-aws-modules/vpc/aws"
    version = "6.7.2"
    azs = var.azs
    cidr = "10.0.0.0/16"
    enable_nat_gateway = true
    single_nat_gateway = true
    enable_dns_hostnames = true

    private_subnets      = ["10.0.1.0/24", "10.0.2.0/24"]
    public_subnets       = ["10.0.101.0/24", "10.0.102.0/24"]
    database_subnets     = ["10.0.21.0/24", "10.0.22.0/24"]
    elasticache_subnets  = ["10.0.31.0/24", "10.0.32.0/24"]
}
