# ElastiCache Redis Cluster
# Equivalent to Azure Redis Cache (azurerm_redis_cache)

# Generate random password for Redis auth token
# ElastiCache auth_token only allows alphanumeric characters (no @, ", /, or other special chars)
resource "random_password" "redis_auth_token" {
  count            = var.use_redis_cache ? 1 : 0
  length           = 32
  special          = false # ElastiCache doesn't allow special characters in auth_token
  override_special = ""    # No special characters allowed
}

# Create ServiceLinkedRole for ElastiCache (required by AWS)
# This role allows ElastiCache to manage resources on your behalf
resource "aws_iam_service_linked_role" "elasticache" {
  count            = var.use_redis_cache ? 1 : 0
  aws_service_name = "elasticache.amazonaws.com"
  description      = "Service-linked role for ElastiCache"
}

# ElastiCache Subnet Group
resource "aws_elasticache_subnet_group" "redis" {
  count      = var.use_redis_cache ? 1 : 0
  name       = "${local.name_prefix}-redis-subnet-group"
  subnet_ids = aws_subnet.private[*].id

  depends_on = [aws_iam_service_linked_role.elasticache]

  tags = merge(
    local.common_tags,
    {
      Name = "${local.name_prefix}-redis-subnet-group"
    }
  )
}

# Security group for ElastiCache
resource "aws_security_group" "elasticache" {
  name        = "${local.name_prefix}-elasticache-sg"
  description = "Security group for ElastiCache Redis"
  vpc_id      = aws_vpc.main.id

  ingress {
    from_port       = 6379
    to_port         = 6379
    protocol        = "tcp"
    security_groups = [aws_security_group.lambda.id]
    description     = "Redis access from Lambda functions"
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = merge(
    local.common_tags,
    {
      Name = "${local.name_prefix}-elasticache-sg"
    }
  )
}

# ElastiCache Redis Cluster
# Equivalent to azurerm_redis_cache with Basic SKU
resource "aws_elasticache_replication_group" "redis" {
  count = var.use_redis_cache ? 1 : 0

  replication_group_id = "${local.name_prefix}-redis"
  description          = "Redis cluster for ${var.project_name}"

  # Engine configuration
  engine               = "redis"
  engine_version       = "7.0"
  node_type            = "cache.t3.micro" # Equivalent to Basic C1 in Azure
  port                 = 6379
  parameter_group_name = "default.redis7"

  # Cluster configuration
  num_cache_clusters = 1 # Single node for Basic tier (equivalent to Azure Basic)

  # Network configuration
  subnet_group_name  = aws_elasticache_subnet_group.redis[0].name
  security_group_ids = [aws_security_group.elasticache.id]

  # Authentication - use the generated random password
  auth_token                 = random_password.redis_auth_token[0].result
  transit_encryption_enabled = true
  at_rest_encryption_enabled = true

  # Automatic failover (disabled for single node)
  automatic_failover_enabled = false

  # Snapshot configuration
  snapshot_retention_limit = 7 # days
  snapshot_window          = "03:00-05:00"

  # Maintenance window
  maintenance_window = "mon:05:00-mon:07:00"

  tags = merge(
    local.common_tags,
    {
      Name = "${local.name_prefix}-redis"
    }
  )

  depends_on = [
    aws_elasticache_subnet_group.redis[0],
    aws_security_group.elasticache,
    random_password.redis_auth_token[0],
    aws_iam_service_linked_role.elasticache[0]
  ]
}
