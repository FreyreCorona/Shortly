output "cluster_name" { value= module.eks.cluster_name }
output "rds_endpoint" { value= module.rds.db_instance_address }
output "rds_secret_arn" { value= module.rds.db_instance_master_user_secret_arn }

output "redis_endpoint" { value = aws_elasticache_cluster.shortly.cache_nodes[0].address }
output "mq_endpoint"    { value = aws_mq_broker.shortly.instances[0].endpoints[0] }