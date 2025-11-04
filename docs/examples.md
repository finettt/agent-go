# Examples and Best Practices

This document provides practical examples and best practices for using Agent-Go effectively in various scenarios.

## Table of Contents

- [Basic Usage Examples](#basic-usage-examples)
- [Advanced Workflows](#advanced-workflows)
- [Integration Scenarios](#integration-scenarios)
- [Performance Optimization](#performance-optimization)
- [Security Best Practices](#security-best-practices)
- [Common Use Cases](#common-use-cases)
- [Troubleshooting Examples](#troubleshooting-examples)

## Basic Usage Examples

### 1. File Operations

**Create and Edit Files**
```
> Create a Python script that calculates factorial
Created factorial.py with the following content:
```python
def factorial(n):
    if n == 0:
        return 1
    else:
        return n * factorial(n-1)

# Test the function
print(factorial(5))
```

> Run the factorial script
Factorial of 5 is 120
```

**Batch File Operations**
```
> Create a directory structure for a new project and initialize it
Created directory 'my-project'
Created directory 'my-project/src'
Created directory 'my-project/docs'
Created directory 'my-project/tests'
Initialized git repository in my-project
Created requirements.txt with basic dependencies
```

### 2. System Administration

**System Information**
```
> Show system information and disk usage
System Information:
- OS: Linux Ubuntu 22.04
- Kernel: 5.15.0-88-generic
- CPU: 4 cores, 2.4 GHz
- Memory: 8GB total, 4GB available

Disk Usage:
- /: 45GB used, 55GB available
- /home: 25GB used, 75GB available
```

**Network Operations**
```
> Check network connectivity and open ports
Network connectivity:
- Internet: ✅ Connected
- DNS: ✅ Working (8.8.8.8)
- Latency: 15ms

Open ports:
- 22/tcp: SSH (OpenSSH)
- 80/tcp: HTTP (nginx)
- 443/tcp: HTTPS (nginx)
```

### 3. Development Tasks

**Code Generation**
```
> Generate a React component for a user profile
Created UserProfile.jsx with the following content:
```jsx
import React from 'react';

const UserProfile = ({ user }) => {
  return (
    <div className="user-profile">
      <img src={user.avatar} alt={user.name} />
      <h2>{user.name}</h2>
      <p>{user.email}</p>
      <p>Role: {user.role}</p>
    </div>
  );
};

export default UserProfile;
```
```

**Database Operations**
```
> Create a SQL schema for a blog application
Created schema.sql with the following content:
```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE posts (
    id SERIAL PRIMARY KEY,
    title VARCHAR(200) NOT NULL,
    content TEXT NOT NULL,
    user_id INTEGER REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```
```

## Advanced Workflows

### 1. Multi-Step Project Setup

**Complete Web Application Setup**
```
> Set up a complete web application with React backend and PostgreSQL database
1. Created project directory 'web-app'
2. Initialized React app with Vite
3. Set up Express.js backend
4. Configured PostgreSQL database connection
5. Created API routes for CRUD operations
6. Set up environment variables
7. Added Docker configuration
8. Created basic CI/CD pipeline

Project structure created successfully:
web-app/
├── frontend/          # React frontend
├── backend/           # Node.js backend
├── database/          # Database scripts
├── docker-compose.yml # Docker configuration
└── .github/           # CI/CD configuration
```

### 2. Data Processing Pipeline

**Automated Data Analysis**
```
> Create a data processing pipeline that analyzes sales data and generates reports
1. Created data processing script with pandas
2. Set up database connection for data extraction
3. Implemented data cleaning and transformation
4. Created visualization charts
5. Generated PDF report
6. Set up automated scheduling

Data processing pipeline created:
- Extracts data from PostgreSQL database
- Cleans and transforms sales data
- Generates monthly sales reports
- Creates visualizations and charts
- Saves reports to shared drive
```

### 3. DevOps Automation

**Infrastructure Setup**
```
> Set up complete infrastructure on AWS using Terraform
1. Created VPC with public and private subnets
2. Configured security groups and NACLs
3. Set up EC2 instances with auto-scaling
4. Configured load balancer
5. Set up RDS database
6. Created S3 buckets for storage
7. Configured monitoring and logging

Infrastructure deployed successfully:
- VPC: 10.0.0.0/16
- Public Subnets: 2
- Private Subnets: 2
- EC2 Instances: 3 (with auto-scaling)
- RDS Instance: PostgreSQL 14
- S3 Buckets: 3
- CloudWatch Alarms: 5
```

## Integration Scenarios

### 1. CI/CD Pipeline Integration

**GitHub Actions Example**
```yaml
# .github/workflows/agent-go.yml
name: Agent-Go CI/CD

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main ]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3
    
    - name: Set up Go
      uses: actions/setup-go@v3
      with:
        go-version: 1.22
    
    - name: Run tests
      run: |
        go test -v ./...
        go vet ./...
        golangci-lint run
    
    - name: Build application
      run: |
        make build
    
    - name: Run integration tests
      run: |
        export OPENAI_KEY=${{ secrets.OPENAI_API_KEY }}
        ./agent-go --test-mode
```

### 2. Docker Integration

**Dockerfile Example**
```dockerfile
# Dockerfile
FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o agent-go ./src

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/agent-go .
COPY --from=builder /app/config ./config

CMD ["./agent-go"]
```

**docker-compose.yml Example**
```yaml
version: '3.8'
services:
  agent-go:
    build: .
    environment:
      - OPENAI_KEY=${OPENAI_KEY}
      - OPENAI_MODEL=gpt-4-turbo
      - RAG_ENABLED=1
      - RAG_PATH=/documents
    volumes:
      - ./documents:/documents
      - ./config:/root/.config/agent-go
    ports:
      - "8080:8080"
```

### 3. Kubernetes Integration

**Kubernetes Deployment**
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: agent-go
spec:
  replicas: 3
  selector:
    matchLabels:
      app: agent-go
  template:
    metadata:
      labels:
        app: agent-go
    spec:
      containers:
      - name: agent-go
        image: agent-go:latest
        ports:
        - containerPort: 8080
        env:
        - name: OPENAI_KEY
          valueFrom:
            secretKeyRef:
              name: agent-go-secrets
              key: openai-key
        - name: OPENAI_MODEL
          value: "gpt-4-turbo"
        volumeMounts:
        - name: config-volume
          mountPath: /root/.config/agent-go
        - name: documents
          mountPath: /documents
      volumes:
      - name: config-volume
        configMap:
          name: agent-go-config
      - name: documents
        persistentVolumeClaim:
          claimName: agent-go-documents
```

## Performance Optimization

### 1. Configuration Optimization

**High-Performance Configuration**
```json
{
  "api_url": "https://api.openai.com",
  "model": "gpt-4-turbo",
  "temperature": 0.0,
  "max_tokens": 2000,
  "rag_enabled": true,
  "rag_path": "/optimized/documents",
  "rag_snippets": 10,
  "message_history_limit": 15,
  "request_timeout": 30,
  "max_retries": 3
}
```

**Environment Variables for Performance**
```bash
# Optimized environment setup
export OPENAI_KEY="your-api-key"
export OPENAI_MODEL="gpt-4-turbo"
export RAG_ENABLED=1
export RAG_PATH="/fast/documents"
export RAG_SNIPPETS=8
export MESSAGE_HISTORY_LIMIT=15
export REQUEST_TIMEOUT=25
```

### 2. RAG Optimization

**Document Organization**
```
# Optimal directory structure for RAG
documents/
├── technical/          # Technical documentation
│   ├── api/
│   ├── architecture/
│   └── deployment/
├── business/          # Business requirements
│   ├── requirements/
│   └── specifications/
├── code/              # Code documentation
│   ├── src/
│   └── tests/
└── external/          # External references
    ├── standards/
    └── best-practices/
```

**Document Indexing Strategy**
```bash
# Pre-process documents for better RAG performance
find documents -type f -name "*.md" -exec grep -l "important" {} \; > important_docs.txt
find documents -type f -name "*.txt" -exec wc -l {} \; | sort -n > document_sizes.txt
```

### 3. Memory Management

**Optimized Configuration**
```json
{
  "message_history_limit": 10,
  "cache_enabled": true,
  "cache_ttl": 300,
  "memory_limit": "512MB",
  "gc_percent": 100
}
```

## Security Best Practices

### 1. Secure Configuration

**Production Configuration**
```json
{
  "api_url": "https://secure-api-provider.com",
  "model": "gpt-4",
  "temperature": 0.0,
  "rag_enabled": false,
  "max_tokens": 1000,
  "message_history_limit": 5,
  "enable_logging": false,
  "sanitize_inputs": true
}
```

**Environment Variables for Security**
```bash
# Secure environment setup
export OPENAI_KEY="${OPENAI_API_KEY}"
export OPENAI_BASE="https://secure-api.com"
export OPENAI_MODEL="gpt-4"
export RAG_ENABLED=0
export LOG_LEVEL="error"
export ENABLE_METRICS=0
```

### 2. Input Validation

**Safe Command Execution**
```
# Example of safe commands
> Create a backup of /etc/nginx/nginx.conf
Created backup: /etc/nginx/nginx.conf.bak

> List files in /var/log/nginx
Files in /var/log/nginx:
- access.log
- error.log
- access.log.1
- error.log.1

> Check disk space on /var partition
Disk space on /var:
- Total: 50GB
- Used: 35GB
- Available: 15GB
- Usage: 70%
```

### 3. Network Security

**Secure Network Configuration**
```bash
# Firewall rules for Agent-Go
sudo ufw allow 22/tcp    # SSH
sudo ufw allow 8080/tcp  # Agent-Go web interface (if enabled)
sudo ufw deny out 443    # Block direct API access
sudo ufw enable
```

## Common Use Cases

### 1. Development Assistant

**Code Review and Optimization**
```
> Review the following code and suggest improvements
Code review completed:
1. Added error handling for file operations
2. Implemented proper logging
3. Added input validation
4. Optimized database queries
5. Added security headers

Improved code:
```go
func processFile(filename string) error {
    if filename == "" {
        return fmt.Errorf("filename cannot be empty")
    }
    
    file, err := os.Open(filename)
    if err != nil {
        log.Printf("Failed to open file: %v", err)
        return fmt.Errorf("file access error: %w", err)
    }
    defer file.Close()
    
    // Process file content
    return nil
}
```
```

### 2. System Administration

**Automated Maintenance**
```
> Perform system maintenance: update packages, clean logs, check disk space
System maintenance completed:
1. Updated 15 packages
2. Cleaned log files (freed 2.3GB)
3. Disk space checked and optimized
4. Services restarted if needed
5. Backup verification passed

System status:
- CPU Usage: 15%
- Memory Usage: 45%
- Disk Usage: 65%
- Services: All running
```

### 3. Data Analysis

**Business Intelligence**
```
> Analyze sales data from Q1 and generate insights
Sales analysis completed:
- Total revenue: $2.3M (↑ 15% from Q4)
- Top product: Product A (35% of sales)
- Best region: North America (40% of revenue)
- Customer retention: 85%
- Average order value: $125

Recommendations:
1. Increase marketing in Europe
2. Optimize inventory for Product A
3. Improve customer retention strategies
```

## Troubleshooting Examples

### 1. API Connection Issues

**Diagnosing API Problems**
```
> Check API connectivity and diagnose issues
API connectivity check:
- Status: ❌ Failed
- Error: Connection timeout
- URL: https://api.openai.com
- Port: 443
- DNS Resolution: ✅ Working
- Network Reachability: ✅ Working

Troubleshooting steps:
1. Check internet connection
2. Verify API key is valid
3. Check firewall settings
4. Try alternative API endpoint
```

### 2. Performance Issues

**Performance Analysis**
```
> Analyze application performance and identify bottlenecks
Performance analysis completed:
- Response time: 2.3s (target: <1s)
- Memory usage: 512MB (high)
- CPU usage: 45% (peak)
- API calls: 15/min
- Cache hit rate: 60%

Optimization recommendations:
1. Increase message history limit
2. Enable response caching
3. Optimize RAG search parameters
4. Add request batching
```

### 3. Configuration Issues

**Configuration Validation**
```
> Validate configuration and identify issues
Configuration validation:
- API Key: ✅ Valid
- API URL: ✅ Reachable
- Model: ✅ Available
- RAG Path: ❌ Not accessible
- Permissions: ❌ Insufficient

Issues found:
1. RAG directory does not exist
2. Missing read permissions for config file
3. Environment variables not set

Fixes applied:
1. Created RAG directory
2. Set proper file permissions
3. Updated environment variables
```

## Best Practices Summary

### 1. Configuration Management
- Use environment variables for sensitive data
- Keep configuration files in version control
- Document all configuration options
- Test configuration changes in development first

### 2. Security
- Never hardcode API keys
- Use proper file permissions
- Validate all user inputs
- Implement proper error handling

### 3. Performance
- Monitor response times
- Use caching where appropriate
- Optimize RAG search parameters
- Set reasonable limits on message history

### 4. Development
- Write comprehensive tests
- Follow Go coding standards
- Use proper error handling
- Document complex logic

### 5. Deployment
- Use containerization for consistency
- Implement proper logging
- Set up monitoring and alerting
- Use CI/CD pipelines for automation

---

These examples and best practices should help you get the most out of Agent-Go in various scenarios. Remember to adapt them to your specific needs and requirements.