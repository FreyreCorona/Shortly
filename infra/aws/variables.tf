variable "azs" {
  description = "The availability zones to use for the AWS resources."
  type        = list(string)
  default     = ["us-east-1a", "us-east-1b", "us-east-1c"]
}

