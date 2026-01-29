# Examples and Best Practices

This document provides practical examples and best practices for using Agent-Go effectively in various scenarios.

## Table of Contents

- [Agent Studio Examples](#agent-studio-examples)
- [Session Management Examples](#session-management-examples)
- [Background Command Execution](#background-command-execution)
- [Basic Usage Examples](#basic-usage-examples)
- [MCP Integration Examples](#mcp-integration-examples)
- [Todo List Management Examples](#todo-list-management-examples)
- [Notes Management Examples](#notes-management-examples)
- [Advanced Workflows](#advanced-workflows)
- [Integration Scenarios](#integration-scenarios)
- [Performance Optimization](#performance-optimization)
- [Security Best Practices](#security-best-practices)
- [Common Use Cases](#common-use-cases)
- [Troubleshooting Examples](#troubleshooting-examples)

## Agent Studio Examples

### 1. Creating Specialized Agents

**Code Review Agent**
```
> /agent studio

User: I want to create a code review agent that analyzes pull requests

Agent Studio: Welcome to Agent Studio! Please describe your agent:

1. Primary Goal: Analyze code changes for quality, security, and best practices
2. Constraints: Should not approve changes that break existing functionality
3. Workflow: Review diffs, run linting, check for security issues, provide feedback
4. Tools: execute_command, spawn_agent, get_todo_list, create_todo, update_todo

Specification saved successfully.
Agent 'code-reviewer' created with ID: agent-20251218-123456
```

**Documentation Writer Agent**
```
> /agent studio

User: Create an agent that generates technical documentation

Agent Studio: Please specify:
1. Primary Goal: Generate clear, comprehensive technical documentation
2. Constraints: Should follow company documentation standards
3. Workflow: Analyze code, extract key information, structure documentation
4. Tools: execute_command, spawn_agent, get_todo_list, create_todo, update_todo

Agent 'documentation-writer' created successfully.
```

### 2. Managing Agent Lifecycle

**Listing and Viewing Agents**
```
> /agent list
Available Agents:

1. default (Built-in)
   - Purpose: General-purpose AI assistant
   - Model: gpt-4-turbo
   - Temperature: 0.1
   - Status: Protected (cannot be deleted)

2. code-reviewer
   - Purpose: Analyze code changes for quality and security
   - Model: gpt-4-turbo
   - Temperature: 0.2
   - Created: 2025-12-18 10:30:00

> /agent view code-reviewer
=== Agent: code-reviewer ===
Purpose: Analyze code changes for quality and security
Model: gpt-4-turbo
Temperature: 0.2
Max Tokens: 2000
Created: 2025-12-18 10:30:00
Updated: 2025-12-18 10:30:00

System Prompt: [Full system prompt content]

Available Tools: execute_command, spawn_agent, get_todo_list, create_todo, update_todo
```

**Using Specialized Agents**
```
> /agent use code-reviewer
Agent 'code-reviewer' activated.
System prompt rebuilt with agent configuration.
Context cleared for focused agent operation.

> Review this pull request for potential issues
[Agent analyzes code using specialized code review tools and workflows]
```

### 3. Agent Collaboration

**Main Agent with Subagent Delegation**
```
> /agent use project-manager

User: I need to review this complex pull request

Project Manager Agent: This requires specialized analysis. Let me spawn a code review subagent.
[Spawns code-reviewer subagent with specialized tools]

Code Review Subagent: Analyzing code changes...
- Security Analysis: ✅ No vulnerabilities detected
- Performance Review: ⚠️ Potential performance issue in function X
- Code Quality: ✅ Follows best practices

Main Agent: Based on the specialized analysis, here are my recommendations...
```

## Session Management Examples

### 1. Project Context Switching

**Creating Project Sessions**
```
> /session new
Session 'project-alpha-20251218' saved successfully.

> Working on feature A implementation
[Agent provides context-aware assistance for feature A]

> /session new
Session 'feature-b-development' saved successfully.

> /session list
Saved Sessions:
1. project-alpha-20251218 (Created: 2025-12-18 10:30:00, Messages: 25, Tokens: 15,432)
2. feature-b-development (Created: 2025-12-18 10:45:00, Messages: 12, Tokens: 8,215)

Current Session: feature-b-development

> /session restore project-alpha-20251218
Session 'project-alpha-20251218' restored successfully.
Loaded 25 messages with 15,432 tokens.
```

### 2. Task-Based Session Management

**Separate Sessions for Different Tasks**
```
> /session new
Session 'database-migration' saved successfully.

> Design and implement the database migration script
[Agent assists with database migration task]

> /session new
Session 'api-integration' saved successfully.

> /session list
Saved Sessions:
1. database-migration (Created: 2025-12-18 14:00:00, Messages: 45, Tokens: 28,150)
2. api-integration (Created: 2025-12-18 14:30:00, Messages: 22, Tokens: 12,840)

> /session restore database-migration
Session 'database-migration' restored successfully.
[Back to database migration context]
```

### 3. Long-term Project Management

**Session Management Across Days**
```
> /session new
Session 'sprint-3-development' saved successfully.

[Work continues over several days...]

> /session list
Saved Sessions:
1. sprint-3-development (Created: 2025-12-15 09:00:00, Last Accessed: 2025-12-18 10:30:00, Messages: 156, Tokens: 89,215)
2. bug-fixes (Created: 2025-12-16 14:20:00, Last Accessed: 2025-12-17 16:45:00, Messages: 34, Tokens: 18,720)

> /session restore sprint-3-development
Session 'sprint-3-development' restored successfully.
[Resume sprint 3 development with full context]
```

## Background Command Execution

### 1. Long-Running Tasks

**Building and Deployment**
```
> Execute the production build in the background
$ npm run build &
[Background command started with PID: 12345]

> Monitor build progress
$ tail -f build.log
[Real-time build output...]

> List running background commands
> /list_background_commands
Running Background Commands:
1. PID: 12345, Command: npm run build, Status: Running
2. PID: 12346, Command: npm install, Status: Completed
```

### 2. Parallel Task Execution

**Multiple Background Operations**
```
> Start multiple background tasks
$ npm run build --background
$ npm run test --background
$ npm run lint --background

> Monitor all background processes
$ list_background_commands
Background Processes:
- Build: Running (75% complete)
- Tests: Running (50% complete)
- Lint: Running (90% complete)

> Get detailed logs for specific process
$ get_background_logs 12345
[Detailed build logs...]
```

### 3. Resource Management

**Managing System Resources**
```
> Check system resource usage
$ top -b | head -20
[Resource usage information...]

> Kill resource-intensive process if needed
$ kill_background_command 12345
Process 12345 terminated successfully.

> Verify process termination
$ list_background_commands
Running Background Commands:
1. PID: 12346, Command: npm run test, Status: Running
2. PID: 12347, Command: npm run lint, Status: Running

> Exit Agent-Go (prevented while background tasks running)
Cannot exit: Background tasks are still running.
Use /list_background_commands to check status.
```

## Notes Management Examples

### 1. Project Documentation Storage

**API Endpoint Management**

> Create a note for the production API endpoint
Created note: api_endpoint

> Add the staging API endpoint to notes
Created note: staging_api_endpoint

> View the production API endpoint
> /notes view api_endpoint
=== api_endpoint ===
Created: 2025-11-25 08:30:00
Updated: 2025-11-25 08:30:00

The production API endpoint is https://api.example.com/v1

**Database Configuration**

> Create notes for database connection details
Created note: db_production
Created note: db_staging

> Update database configuration notes
Updated note: db_production
Updated note: db_staging

### 2. Development Workflow Integration

**Project Setup Notes**

> Create comprehensive project setup notes
Created note: project_setup
Created note: deployment_guide
Created note: troubleshooting_common_issues

> Reference project setup during development
[Using stored notes: project_setup, deployment_guide]

Based on your project setup notes, I can see you're working with a microservices architecture. Let me help you set up the development environment...

**Code Conventions**

> Store coding standards and conventions
Created note: coding_standards
Created note: naming_conventions
Created note: review_checklist

> Apply stored conventions during code generation
[Using stored notes: coding_standards, naming_conventions]

Generated code following your established coding standards and naming conventions.

### 3. Knowledge Management

**Personal Preferences**

> Store frequently used commands and shortcuts
Created note:常用命令
Created note: debugging_commands
Created note: deployment_commands

> Store troubleshooting patterns
Created note: common_errors_solutions
Created note: performance_optimization_tips

**Project-Specific Knowledge**

> Store project-specific information
Created note: architecture_decisions
Created note: tech_stack_details
Created note: third_party_integrations

> Reference stored knowledge during development
[Using stored notes: architecture_decisions, tech_stack_details]

Based on your architecture decisions document, I can see you're using event-driven microservices. Let me help you implement the event publisher...

### 4. Cross-Session Persistence

**Session Continuity**

> Create notes at the end of a session
Created note: session_summary_20251125
Created note: next_session_tasks

> Continue work in next session
[Using stored notes: session_summary_20251125, next_session_tasks]

Welcome back! Based on your previous session notes, you were working on implementing the user authentication service. Here's what we accomplished and what's next...

**Documentation Updates**

> Update notes as project evolves
Updated note: architecture_decisions
Updated note: deployment_guide

> Maintain living documentation
Created note: api_changes_log
Created note: dependency_versions

### 5. Notes Management Best Practices

**Organized Note Structure**

```
# Optimal note organization for projects:
project_setup/          # Project initialization
├── api_endpoints       # API URLs and documentation
├── database_config     # Database connection details
├── deployment          # Deployment procedures
└── troubleshooting     # Common issues and solutions

development/            # Development workflow
├── coding_standards    # Code style and conventions
├── testing             # Testing procedures
├── debugging           # Debugging techniques
└── tools              # Development tools and setup

knowledge/              # Personal and team knowledge
├── best_practices      # Industry best practices
├── patterns            # Design patterns
├── performance         # Optimization techniques
└── security            # Security considerations
```

**Note Content Guidelines**

- **Be specific**: Use clear, descriptive names for notes
- **Include metadata**: Add timestamps and context
- **Keep updated**: Regularly review and update notes
- **Use consistent format**: Standardize note structure across projects
- **Separate concerns**: Different types of notes in different categories
- **Include examples**: Provide practical examples where helpful

**Integration with Other Features**

```
> Create a note and use it with RAG
Created note: project_context

Based on your project context note and the RAG documents, I can see you're working with a React application using TypeScript. Let me help you create a component...

> Create a note and reference it in todos
Created note: research_topics

Created todo: [ID: 1] Research new authentication methods (pending)
Created todo: [ID: 2] Update security documentation (pending)

[Using stored notes: research_topics]

Based on your research topics note, let me start with investigating OAuth 2.0 implementation...
```

## Basic Usage Examples

### 1. File Operations

**Create and Edit Files**


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



**Batch File Operations**


> Create a directory structure for a new project and initialize it
Created directory 'my-project'
Created directory 'my-project/src'
Created directory 'my-project/docs'
Created directory 'my-project/tests'
Initialized git repository in my-project
Created requirements.txt with basic dependencies



### 2. System Administration

**System Information**


> Show system information and disk usage
System Information:

- OS: Linux Ubuntu 22.04
- Kernel: 5.15.0-88-generic
- CPU: 4 cores, 2.4 GHz
- Memory: 8GB total, 4GB available

Disk Usage:

- /: 45GB used, 55GB available
- /home: 25GB used, 75GB available



**Network Operations**


> Check network connectivity and open ports
Network connectivity:

- Internet: ✅ Connected
- DNS: ✅ Working (8.8.8.8)
- Latency: 15ms

Open ports:

- 22/tcp: SSH (OpenSSH)
- 80/tcp: HTTP (nginx)
- 443/tcp: HTTPS (nginx)



### 3. Development Tasks

**Code Generation**


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



**Database Operations**


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


## Advanced Workflows

### 1. Multi-Step Project Setup

**Complete Web Application Setup**

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


### 2. Data Processing Pipeline

**Automated Data Analysis**

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


### 3. DevOps Automation

**Infrastructure Setup**

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

## Integration Scenarios

### 1. Docker Integration

The Agent-Go Dockerfile creates a `/workspace` volume for mounting your current directory, allowing seamless file access within the container.

**Building and Running with Docker**

```bash
# Build the Docker image
docker build -t agent-go .

# Run with current directory mounted
docker run -it -v $(pwd):/workspace agent-go

# Run with environment variables
docker run -it \
  -v $(pwd):/workspace \
  -e OPENAI_KEY="your-api-key" \
  -e OPENAI_MODEL="gpt-4-turbo" \
  agent-go

# Run with RAG documents
docker run -it \
  -v $(pwd):/workspace \
  -v /path/to/docs:/documents \
  -e OPENAI_KEY="your-api-key" \
  -e RAG_ENABLED=1 \
  -e RAG_PATH="/documents" \
  agent-go

# Run with persistent configuration
docker run -it \
  -v $(pwd):/workspace \
  -v ~/.config/agent-go:/home/appuser/.config/agent-go \
  -e OPENAI_KEY="your-api-key" \
  agent-go
```

**Docker Volume Configuration:**
- `/workspace` - Your current directory is mounted here for file operations
- `/home/appuser/.config/agent-go` - Configuration directory (persists across runs when mounted)
- Additional volumes for RAG documents or other data
- Container runs as non-root user `appuser` for security

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
      - AUTO_COMPRESS=1
    volumes:
      - .:/workspace
      - ./documents:/documents
      - ~/.config/agent-go:/home/appuser/.config/agent-go
    stdin_open: true
    tty: true
```

**Usage with docker-compose:**

```bash
# Start the service
docker-compose run agent-go

# With custom command
docker-compose run agent-go /bin/sh
```

### 2. Kubernetes Integration

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

## Advanced Features

### 1. Token Usage Monitoring

Agent-Go tracks and displays token usage in real-time throughout your session:

```
> Create a Python script that calculates factorial
[Tokens: 156]

> Run the factorial script with input 10
[Tokens: 342]

> /compress
Context compressed. Starting new chat with compressed summary as system message.
[Tokens: 0]  # Token counter resets after compression
```

**Understanding Token Usage:**
- Tokens are cumulative throughout the session
- Each API call adds to the total token count
- Auto-compression triggers at 75% of `model_context_length`
- Manual compression with `/compress` resets the token counter

### 2. Tool Loop Detection

Agent-Go automatically detects when the AI gets stuck in infinite loops:

```
> Try to execute a command that doesn't exist
[Tool Loop Detected: execute_command called 3 times with identical arguments]
STOP: You appear to be stuck in a loop, repeatedly calling the same tool. Please step back, analyze what you've tried so far, and try a completely different approach to solve this task.

> Analyze the situation and try a different approach
[AI now tries a different strategy, such as checking if the command exists or suggesting alternatives]
```

**How Tool Loop Detection Works:**
- Monitors tool calls within each API response
- Tracks repeated calls across multiple iterations
- Triggers after 3 identical tool calls
- Injects a stop message to guide the model
- Works in all execution modes (CLI, task, pipeline, editor)

### 3. Empty Response Retry Logic

Agent-Go handles empty model responses gracefully:

```
> Ask a question
[Warning: received an empty response from the API (attempt 1/3), retrying...]
[Warning: model returned an empty message (attempt 2/3), retrying...]
[AI provides response after successful retry]
```

**Retry Behavior:**
- Automatically retries up to 3 times
- Warns user about empty responses
- Handles transient API issues gracefully
- Works in all execution modes

### 4. Time Context Injection

Agent-Go provides temporal awareness to the AI:

```
> What time is it?
Current Time: 2026-01-29T14:30:00Z (UTC) | Local: Thu, 29 Jan 2026 17:30:00 MSK | Timezone: MSK (UTC+03:00)

> Schedule a meeting for tomorrow at 10 AM
[AI uses current time context to calculate correct time for meeting]
```

**Time Context Benefits:**
- Enables time-aware operations
- Improves context for scheduling tasks
- Helps with time-sensitive operations
- Automatically injected in all API requests

### 2. MCP-Enhanced Development Workflow

```
> /mcp add context7 npx -y @upstash/context7-mcp
> /mcp add time uvx mcp-server-time

> Create a React component that shows the current time in multiple timezones
[Using MCP: Getting React documentation]
[Using MCP: Getting timezone utilities]

Created TimeZones.jsx:
```jsx
import React, { useState, useEffect } from 'react';

const TimeZones = () => {
  const [times, setTimes] = useState({});
  
  useEffect(() => {
    // Component implementation using latest React patterns
    ...
  });
  
  return (
    <div className="timezones">
      <h2>World Clocks</h2>
      {/* Display times for multiple zones */}
    </div>
  );
};

export default TimeZones;
```

Based on latest React documentation and timezone utilities.
```

### 3. Todo-Driven Development

```
> Create a development plan for a new API endpoint
Created todo list:
- [ID: 1] Design API schema (pending)
- [ID: 2] Write route handler (pending)
- [ID: 3] Add validation middleware (pending)
- [ID: 4] Write unit tests (pending)
- [ID: 5] Update API documentation (pending)

> Let's work through these step by step. Start with the schema design.
Updated: [ID: 1] Design API schema (in-progress)

Created schema.ts:
```typescript
interface UserRequest {
  email: string;
  name: string;
  role: 'admin' | 'user';
}
```

Updated: [ID: 1] Design API schema (completed)

> Now the route handler
Updated: [ID: 2] Write route handler (in-progress)
...
```

### 4. Command-Line Task Execution

Execute tasks directly without interactive mode:

```bash
# Single task execution
./agent-go "Create a new directory called 'test-project' and initialize git"

# The agent will execute the task and exit automatically
```

**Benefits:**
- Ideal for scripting and automation
- No interactive prompts
- Returns exit code 0 on success, non-zero on failure
- Perfect for CI/CD pipelines

### 3. Shell Mode

Direct shell command execution for interactive sessions:

```
> /shell
Entered shell mode. Type 'exit' to return.
shell> ls -la
total 8
drwxr-xr-x 2 user user 4096 Oct 27 10:00 .
drwxr-xr-x 5 user user 4096 Oct 27 10:00 ..

shell> exit
Exited shell mode.
```

**Shell Mode Features:**
- Platform-aware (CMD on Windows, sh on Unix-like systems)
- All output is captured and displayed
- Type 'exit' to return to AI assistant mode
- Slash commands are not available in shell mode

### 4. Error Handling and Debugging

Agent-Go provides comprehensive error handling with helpful feedback:

**Common Error Messages:**
```
Error: OpenAI API key is not set
Solution: Set OPENAI_KEY environment variable or run interactive setup

Error: could not connect to API
Solution: Verify OPENAI_BASE URL is correct and accessible

Error: cannot access RAG path
Solution: Ensure RAG_PATH exists and is readable

Error: no messages to compress
Solution: Start a conversation first before using /compress
```

**Debug Mode:**
```bash
# Enable detailed logging for troubleshooting
export DEBUG=1
./agent-go
```

## Common Use Cases

### 1. Development Assistant

**Code Review and Optimization**


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


### 2. System Administration

**Automated Maintenance**


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


### 3. Data Analysis

**Business Intelligence**

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


## Troubleshooting Examples

### 1. API Connection Issues

**Diagnosing API Problems**

> Check API connectivity and diagnose issues
API connectivity check:

- Status: ❌ Failed
- Error: Connection timeout
- URL: <https://api.openai.com>
- Port: 443
- DNS Resolution: ✅ Working
- Network Reachability: ✅ Working

Troubleshooting steps:

1. Check internet connection
2. Verify API key is valid
3. Check firewall settings
4. Try alternative API endpoint


### 2. Performance Issues

**Performance Analysis**


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


### 3. MCP Server Issues

**Problem: MCP server not connecting**

```
> /mcp list
Configured MCP servers:
- custom-server: npx -y @custom/server (connection failed: command not found)

# Solution: Check if the package is available
> /shell
shell> npm search @custom/server
shell> npx -y @custom/server --version
shell> exit

# Update the MCP server command
> /mcp remove custom-server
> /mcp add custom-server "node /path/to/custom-server/index.js"
MCP server 'custom-server' added.
```

**Problem: MCP tool not working**

```
# Enable debug mode to see detailed MCP communication
> Ask the agent to use a specific MCP tool
[MCP Debug: Connecting to server 'time']
[MCP Debug: Calling tool 'get_current_time']
[MCP Debug: Error - invalid arguments]

# Check the tool schema
> /mcp list
# Review the tool documentation
```

### 4. Todo List Issues

**Problem: Todos not persisting**

```
# Check todo file location
> /shell
shell> ls -la ~/.config/agent-go/todos/
shell> cat ~/.config/agent-go/todos/main.json
shell> exit

# Verify permissions
shell> chmod 644 ~/.config/agent-go/todos/*.json
```

**Problem: Wrong agent todo list**

```
# Each agent has its own todo list
# Main agent: ~/.config/agent-go/todos/main.json
# Sub-agents: ~/.config/agent-go/todos/{uuid}.json

# To view main agent todos only:
> /todo
Current Todo List:
- [ID: 1] Main task (pending)
...
```

### 5. Notes Management Issues

**Problem: Notes not persisting**

```
# Check notes directory location
> /shell
shell> ls -la .agent-go/notes/
shell> cat .agent-go/notes/api_endpoint.json
shell> exit

# Verify permissions
shell> chmod 644 .agent-go/notes/*.json
```

**Problem: Notes not appearing in system prompt**

```
# Check if notes are being loaded properly
> Create a test note
Created note: test_note

> Ask about the test note
[Using stored notes: test_note]

Based on your test note, I can see...
```

If notes aren't appearing, check:
1. Notes directory exists and is readable
2. Notes are valid JSON format
3. System prompt is being built correctly
4. No permission issues preventing file access

**Problem: Notes command not working**

```
# Check if notes commands are available
> /help
# Look for /notes commands in the help output

# Try basic notes command
> /notes list
# Should show all notes or error if directory doesn't exist
```

### 6. Configuration Issues

**Configuration Validation**

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

### 5. MCP Integration

**Best Practice: Organize MCP servers by purpose**

```json
{
  "mcp_servers": {
    "context7": {
      "name": "context7",
      "command": "npx -y @upstash/context7-mcp"
    },
    "time": {
      "name": "time",
      "command": "uvx mcp-server-time"
    },
    "fs-projects": {
      "name": "fs-projects",
      "command": "npx -y @modelcontextprotocol/server-filesystem /home/user/projects"
    },
    "fs-docs": {
      "name": "fs-docs",
      "command": "npx -y @modelcontextprotocol/server-filesystem /home/user/docs"
    }
  }
}
```

**Best Practice: Use MCP for external integrations**

- Use MCP servers for accessing external APIs
- Keep sensitive credentials in MCP server configuration
- Use filesystem MCP servers for controlled file access
- Leverage existing MCP servers before building custom tools

### 6. Todo Management

**Best Practice: Use todos for complex multi-step tasks**

- Create todos at the start of complex projects
- Update status as you progress
- Review completed todos for project documentation
- Use separate todo lists for different agents/contexts

**Best Practice: Structured todo descriptions**

```
> Create detailed todos for the deployment process
- [ID: 1] Run tests and ensure all pass (pending)
- [ID: 2] Build production bundle (pending)
- [ID: 3] Run database migrations on staging (pending)
- [ID: 4] Deploy to staging and verify (pending)
- [ID: 5] Run smoke tests on staging (pending)
- [ID: 6] Deploy to production (pending)
- [ID: 7] Monitor production metrics (pending)
```

### 7. Notes Management
- **Organize by purpose**: Create separate categories for different types of information
- **Use consistent naming**: Follow a naming convention for easy discovery
- **Keep updated**: Regularly review and update stored information
- **Include context**: Add relevant metadata and examples
- **Cross-reference**: Link related notes where appropriate
- **Backup important notes**: Consider version control for critical documentation

### 8. Deployment
- Use containerization for consistency
- Implement proper logging
- Set up monitoring and alerting
- Use CI/CD pipelines for automation

---

These examples and best practices should help you get the most out of Agent-Go in various scenarios. Remember to adapt them to your specific needs and requirements.
