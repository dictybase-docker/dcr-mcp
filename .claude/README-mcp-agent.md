# MCP Subagents for Go Development

This directory contains specialized AI subagents for Claude Code, designed to assist with MCP (Model Context Protocol) development using Go and the `mark3labs/mcp-go` package. These Claude Code subagents guide you through the complete MCP server development lifecycle, from architecture design to production deployment.

## 🤖 About Claude Code Subagents

These are specialized AI subagents for Claude Code that provide expert assistance for specific MCP development tasks. Each subagent operates with its own context window, specialized tools, and deep expertise in particular areas of MCP development using Go and `mark3labs/mcp-go`.

### How to Use Subagents

**Method 1: Automatic Delegation** (Recommended)
Claude Code will automatically select the most appropriate subagent based on your task description:
```
"Help me implement a new MCP tool for file validation using Go best practices"
"Create comprehensive tests for my MCP server with benchmarks"
"Deploy my Go MCP server to Kubernetes with monitoring"
```

**Method 2: Explicit Invocation**
Use explicit syntax when you want a specific subagent:
```
> Use the mcp-go-specialist subagent to implement a file validation tool
> Use the mcp-testing-engineer subagent to create performance benchmarks
> Use the mcp-context7-integrator subagent to fetch Go package documentation
```

**Interactive Management**
Use the `/agents` command in Claude Code for an interactive subagent management interface.

**Note**: These subagents follow your project's conventions from `CLAUDE.md` and are designed to work seamlessly with your existing Go MCP development workflow.

## 📚 Table of Contents

- [🎯 Subagent Overview](#-subagent-overview)
- [🚀 Quick Start Guide](#-quick-start-guide)
- [📋 Usage Patterns](#-usage-patterns)
  - [🏗️ Design Phase](#️-design-phase)
  - [💻 Implementation Phase](#-implementation-phase)
  - [🧪 Testing Phase](#-testing-phase)
  - [🚀 Deployment Phase](#-deployment-phase)
- [Task-Specific Workflows](#task-specific-workflows)
- [🎛️ Multi-Agent Workflows](#️-multi-agent-workflows)
- [📊 Agent Priority Matrix](#-agent-priority-matrix)
- [🛠️ Project Integration](#️-project-integration)
- [💡 Best Practices](#-best-practices)
- [🔗 Related Resources](#-related-resources)

## 🎯 Subagent Overview

### Core Subagents
- **`@mcp-server-architect`** - MCP server architecture and design patterns using Go
- **`@mcp-developer`** - Go MCP implementation with mark3labs/mcp-go SDK
- **`@mcp-testing-engineer`** - Go testing strategies, benchmarking, and quality assurance
- **`@mcp-deployment-orchestrator`** - Go binary deployment, containerization, and Kubernetes

### Specialized Subagents
- **`@mcp-go-specialist`** - Expert in mark3labs/mcp-go package specifics and Go idioms
- **`@mcp-context7-integrator`** - Real-time documentation access via Context7 MCP server

### Supporting Subagents
- **`@mcp-protocol-specialist`** - Protocol compliance and specification guidance
- **`@mcp-integration-engineer`** - General system integration patterns
- **`@mcp-registry-navigator`** - MCP registry and discovery
- **`@mcp-security-auditor`** - Security and compliance validation

## 🚀 Quick Start Guide

### For New MCP Projects
```bash
# 1. Architecture & Design
"Design an MCP server architecture for [your use case] using Go best practices"
"Set up Go project structure following CLAUDE.md conventions with mark3labs/mcp-go"

# 2. Implementation
"Implement core MCP server with tools and resources using Go idioms"
"Fetch documentation for mark3labs/mcp-go package and server configuration"

# 3. Testing & Validation
"Create comprehensive test suite with benchmarks for my Go MCP server"

# 4. Deployment
"Create production-ready Dockerfile and Kubernetes manifests for Go MCP server"
```

### For Adding New Features
```bash
# Tool Development
"Implement a new MCP tool for [functionality] using Go validation patterns"
"Create comprehensive tests for my new MCP tool with error cases"

# Resource Development  
"Add dynamic resource endpoints for [data type] with proper Go error handling"
"Get documentation for Go packages related to [specific functionality]"
```

### Explicit Subagent Usage
When you need a specific subagent's expertise:
```bash
> Use the mcp-go-specialist subagent to implement middleware for my MCP server
> Use the mcp-testing-engineer subagent to debug performance bottlenecks
> Use the mcp-context7-integrator subagent to fetch real-time API documentation
> Use the mcp-deployment-orchestrator subagent to optimize container builds
```

## 📋 Usage Patterns

### Phase-Based Development

#### 🏗️ Design Phase
**Primary Subagents**: `@mcp-server-architect` → `@mcp-go-specialist`

Use when:
- Planning server architecture
- Choosing MCP capabilities
- Designing tool interfaces
- Setting up project structure

**Example**:
```
"I need to design an MCP server architecture that handles file operations, 
integrates with external APIs, and supports real-time data streaming using Go"
```

#### 💻 Implementation Phase
**Primary Subagents**: `mcp-go-specialist` → `mcp-developer` → `mcp-context7-integrator`

Use when:
- Writing Go code
- Implementing tools/resources/prompts
- Following project conventions
- Accessing documentation

**Example**:
```
"Help me implement a file processing tool that validates input using go-playground/validator 
and follows the project's error handling patterns from CLAUDE.md"
```

#### 🧪 Testing Phase
**Primary Subagent**: `mcp-testing-engineer`

Use when:
- Writing unit and integration tests
- Performance benchmarking
- Protocol compliance testing
- Debugging and optimization

**Example**:
```
"Create comprehensive tests for my file processing tool including table-driven tests, 
benchmarks, and error cases using testify"
```

#### 🚀 Deployment Phase
**Primary Subagent**: `mcp-deployment-orchestrator`

Use when:
- Creating Docker containers
- Kubernetes deployment
- Setting up monitoring
- Production optimization

**Example**:
```
"Create a production-ready deployment for my Go MCP server with health checks, 
metrics, and auto-scaling using Kubernetes"
```

### Task-Specific Workflows

#### 📚 Documentation & Learning
**Automatic delegation** to `mcp-context7-integrator`:
```bash
# Get real-time package documentation
"Fetch documentation for mark3labs/mcp-go server configuration and usage examples"

# Learn about specific patterns
"Get examples for Go context handling in web servers and best practices"
```

#### 🔧 Go Implementation
**Automatic delegation** to `mcp-go-specialist`:
```bash
# Follow Go best practices
"Show me the proper way to implement middleware using the options pattern from CLAUDE.md"

# Use mark3labs/mcp-go specifics
"How do I add custom hooks to my MCP server using mark3labs/mcp-go?"
```

#### 🐛 Troubleshooting
**Explicit subagent selection** when you need specific expertise:
```bash
# Protocol issues
> Use the mcp-protocol-specialist subagent to fix JSON-RPC batching issues

# Performance problems
> Use the mcp-testing-engineer subagent to debug memory leaks in my MCP server

# Integration issues
> Use the mcp-integration-engineer subagent to resolve authentication problems
```

## 🎛️ Multi-Agent Workflows

### Complete Feature Development
```bash
1. @mcp-server-architect: "Design architecture for user management feature"
2. @mcp-go-specialist: "Implement user management tools using mark3labs/mcp-go"
3. @mcp-context7-integrator: "Get documentation for authentication libraries"
4. @mcp-testing-engineer: "Create comprehensive tests with security testing"
5. @mcp-deployment-orchestrator: "Update deployment with new security configurations"
```

### Performance Optimization
```bash
1. @mcp-testing-engineer: "Identify performance bottlenecks using Go profiling"
2. @mcp-go-specialist: "Optimize code based on profiling data"
3. @mcp-context7-integrator: "Get documentation for performance optimization techniques"
4. @mcp-deployment-orchestrator: "Update resource limits and scaling policies"
```

### Security Review
```bash
1. @mcp-security-auditor: "Conduct security assessment of MCP server"
2. @mcp-go-specialist: "Implement security recommendations using Go patterns"
3. @mcp-testing-engineer: "Add security tests and validation"
4. @mcp-deployment-orchestrator: "Harden deployment configuration"
```

## 📊 Agent Priority Matrix

| Task Type | Primary Agent | Secondary Agent | Tertiary Agent |
|-----------|---------------|-----------------|----------------|
| **New Tool Development** | `mcp-go-specialist` | `mcp-developer` | `mcp-context7-integrator` |
| **Server Architecture** | `mcp-server-architect` | `mcp-go-specialist` | `mcp-deployment-orchestrator` |
| **Testing & QA** | `mcp-testing-engineer` | `mcp-go-specialist` | - |
| **Deployment** | `mcp-deployment-orchestrator` | `mcp-testing-engineer` | - |
| **Documentation** | `mcp-context7-integrator` | `mcp-go-specialist` | - |
| **Protocol Issues** | `mcp-protocol-specialist` | `mcp-developer` | `mcp-integration-engineer` |
| **Security** | `mcp-security-auditor` | `mcp-go-specialist` | `mcp-deployment-orchestrator` |
| **Performance** | `mcp-testing-engineer` | `mcp-go-specialist` | `mcp-deployment-orchestrator` |

## 🛠️ Project Integration

All agents are designed to work with your project's specific conventions:

- **Follow `CLAUDE.md`**: All Go code follows project conventions
- **Use `mark3labs/mcp-go`**: Primary MCP implementation package
- **Integrate Context7**: Real-time documentation access
- **Support Testing**: `testify`, `gotestsum`, benchmarking
- **Enable Deployment**: Docker, Kubernetes, monitoring

## 💡 Best Practices

### Subagent Selection
1. **Prefer Automatic Delegation**: Let Claude Code choose the best subagent based on your task description
2. **Be Specific**: Provide clear, detailed task descriptions for better subagent matching
3. **Use Explicit Invocation**: When you need a specific subagent's expertise or context
4. **Chain Workflows**: Use output from one subagent to inform subsequent tasks

### Communication Patterns
```bash
# ✅ Good: Specific, clear request (automatic delegation)
"Implement a JSON validation tool using go-playground/validator with proper error handling following CLAUDE.md conventions"

# ✅ Good: Explicit subagent when needed
> Use the mcp-go-specialist subagent to optimize my tool's performance

# ❌ Avoid: Vague request
"Help me with validation"
```

### Workflow Integration
- Leverage automatic delegation for most development tasks
- Use explicit invocation for specialized expertise
- Build on previous subagent outputs in conversation
- Maintain consistency with project conventions
- Use `/agents` command for subagent management

## 🔗 Related Resources

- [mark3labs/mcp-go Documentation](https://github.com/mark3labs/mcp-go)
- [Project Conventions (CLAUDE.md)](../CLAUDE.md)
- [Context7 MCP Server](https://github.com/upstash/context7)
- [Go Testing Best Practices](https://github.com/stretchr/testify)

---

**Pro Tip**: Start with `@mcp-go-specialist` for most Go MCP development tasks, then use supporting agents for specific aspects like testing, deployment, or documentation.