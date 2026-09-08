terraform {
  required_version = ">= 1.16.1"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 6.63.0"
    }
  }
  backend "s3" {
    bucket         = "shortly-terraform-state-522779858672"
    key            = "aws/terraform.tfstate"
    region         = "us-east-1"
    use_lockfile   = true
    encrypt        = true
  }
}

provider "aws" {
  region = "us-east-1"
  default_tags {
    tags = {
      Project = "Shortly"
      Owner   = "FreyreCorona"
    }
  }
}