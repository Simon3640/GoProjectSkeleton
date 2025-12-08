# RDS PostgreSQL Database
# Equivalent to Azure PostgreSQL Flexible Server (azurerm_postgresql_flexible_server)

# DB Subnet Group for RDS (required for RDS in VPC)
resource "aws_db_subnet_group" "main" {
  name       = "${local.name_prefix}-db-subnet-group"
  subnet_ids = aws_subnet.private[*].id

  tags = merge(
    local.common_tags,
    {
      Name = "${local.name_prefix}-db-subnet-group"
    }
  )
}

# Security group for RDS
resource "aws_security_group" "rds" {
  name        = "${local.name_prefix}-rds-sg"
  description = "Security group for RDS PostgreSQL"
  vpc_id      = aws_vpc.main.id

  ingress {
    from_port       = 5432
    to_port         = 5432
    protocol        = "tcp"
    security_groups = [aws_security_group.lambda.id]
    description     = "PostgreSQL access from Lambda functions"
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
      Name = "${local.name_prefix}-rds-sg"
    }
  )
}

# RDS PostgreSQL Instance
# Equivalent to azurerm_postgresql_flexible_server
resource "aws_db_instance" "postgres" {
  identifier = "${local.name_prefix}-postgres"

  # Engine configuration
  engine         = "postgres"
  engine_version = "15"
  instance_class = "db.t3.micro" # Equivalent to B_Standard_B1ms in Azure

  # Storage configuration
  allocated_storage     = 20  # GB (minimum for PostgreSQL)
  max_allocated_storage = 100 # Auto-scaling up to 100GB
  storage_type          = "gp3"
  storage_encrypted     = true

  # Database configuration
  db_name  = var.db_name
  username = var.db_user
  password = aws_secretsmanager_secret_version.db_password.secret_string

  # Network configuration
  db_subnet_group_name   = aws_db_subnet_group.main.name
  vpc_security_group_ids = [aws_security_group.rds.id]
  publicly_accessible    = false # Keep private for security

  # Use parameter group
  parameter_group_name = aws_db_parameter_group.postgres.name

  # Backup configuration
  # Free tier allows max 1 day retention, production should use 7+
  backup_retention_period = 1 # days (free tier max, increase to 7+ for production)
  backup_window           = "03:00-04:00"
  maintenance_window      = "mon:04:00-mon:05:00"

  # Performance and availability
  multi_az                        = false # Set to true for production
  performance_insights_enabled    = true
  enabled_cloudwatch_logs_exports = ["postgresql", "upgrade"]

  # Deletion protection
  deletion_protection = false # Set to true for production
  skip_final_snapshot = true  # Set to false for production

  tags = merge(
    local.common_tags,
    {
      Name = "${local.name_prefix}-postgres"
    }
  )

  depends_on = [
    aws_db_subnet_group.main,
    aws_security_group.rds,
    aws_secretsmanager_secret_version.db_password
  ]
}

# RDS Parameter Group for PostgreSQL configuration
resource "aws_db_parameter_group" "postgres" {
  name   = "${local.name_prefix}-postgres-params"
  family = "postgres15"

  # max_connections is a static parameter, requires pending-reboot
  # For dynamic parameters, use apply_method = "immediate"
  parameter {
    name         = "max_connections"
    value        = "100"
    apply_method = "pending-reboot" # Required for static parameters
  }

  tags = merge(
    local.common_tags,
    {
      Name = "${local.name_prefix}-postgres-params"
    }
  )
}
