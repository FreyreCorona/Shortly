resource "aws_security_group" "rabbitmq" {
  name_prefix = "shortly-rabbitmq-"
  vpc_id      = module.vpc.vpc_id
}

resource "aws_security_group_rule" "rabbitmq_ingress" {
  type                     = "ingress"
  from_port                = 5671   # amqps/TLS — recordarás la gotcha #1 del código Go
  to_port                  = 5671
  protocol                 = "tcp"
  security_group_id        = aws_security_group.rabbitmq.id
  source_security_group_id = module.eks.node_security_group_id
}

resource "aws_mq_broker" "shortly" {
  broker_name        = "shortly-rabbitmq"
  engine_type        = "RabbitMQ"
  engine_version     = "3.13"
  host_instance_type = "mq.m7g.medium"
  auto_minor_version_upgrade = true
  user {
    username = "shortly"
    password = var.rabbit_password
  }

  security_groups     = [aws_security_group.rabbitmq.id]
  subnet_ids          = [module.vpc.private_subnets[0]]  # single instance → 1 subnet
  publicly_accessible = false
  deployment_mode     = "SINGLE_INSTANCE"

  tags = { Project = "Shortly" }
}