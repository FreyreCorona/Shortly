resource "aws_security_group" "postgres" {
  name_prefix = "shortly-postgres-"
  vpc_id      = module.vpc.vpc_id
}

resource "aws_security_group_rule" "postgres_ingress" {
  type                     = "ingress"
  from_port                = 5432
  to_port                  = 5432
  protocol                 = "tcp"
  security_group_id        = aws_security_group.postgres.id
  source_security_group_id = module.eks.node_security_group_id
}

module "rds" {
  source  = "terraform-aws-modules/rds/aws"
  version = "6.10.0"

  identifier = "shortly-postgres"

  engine               = "postgres"
  engine_version       = "17"
  family               = "postgres17"
  major_engine_version = "17"
  instance_class       = "db.t3.micro"
  allocated_storage    = 20
  max_allocated_storage = 100
  storage_encrypted    = true

  db_name  = "shortly"
  username = "postgres"
  password = var.rds_password
  port     = 5432

  vpc_security_group_ids = [aws_security_group.postgres.id]
  db_subnet_group_name   = module.vpc.database_subnet_group_name
  create_db_subnet_group = false

  multi_az            = false
  skip_final_snapshot = true
  apply_immediately   = true
  backup_retention_period = 7
}