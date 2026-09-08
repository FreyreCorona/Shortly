resource "aws_security_group" "redis" {
  name_prefix = "shortly-redis-"
  vpc_id      = module.vpc.vpc_id
}

resource "aws_security_group_rule" "redis_ingress" {
  type                     = "ingress"
  from_port                = 6379
  to_port                  = 6379
  protocol                 = "tcp"
  security_group_id        = aws_security_group.redis.id
  source_security_group_id = module.eks.node_security_group_id
}

resource "aws_elasticache_subnet_group" "shortly" {
  name       = "shortly-redis"
  subnet_ids = module.vpc.elasticache_subnets
}

resource "aws_elasticache_cluster" "shortly" {
  cluster_id           = "shortly-redis"
  engine               = "redis"
  engine_version       = "7.1"
  node_type            = "cache.t3.micro"
  num_cache_nodes      = 1
  parameter_group_name = "default.redis7"
  port                 = 6379
  subnet_group_name    = aws_elasticache_subnet_group.shortly.name
  security_group_ids   = [aws_security_group.redis.id]
  apply_immediately    = true

  tags = { Project = "Shortly" }
}