---
allowed-tools: Bash(go run:*), Bash(kill:*), Bash(lsof:*), Bash(ps:*), Bash(curl:*), Bash(pkill:*), Bash(netstat:*), TodoWrite
argument-hint: [--port=PORT] [--dir=DIR] [--health-check] [--stop]
description: Start Go web application development server with hot reload, health checks, and process management
model: claude-3-5-haiku@20241022
---

I'll help you manage your Go web application development server with advanced process management, health monitoring, and development-focused features.

**Prerequisites:**
- Go application can be run with `go run .`
- Application supports common command-line options (--port, --dir, etc.)
- Development environment with curl available for health checks

**Server Management:**
1. Check for existing MDViewer server processes
2. Handle port conflicts and cleanup orphaned processes
3. Start server with specified configuration
4. Perform health checks and connectivity verification
5. Provide process monitoring and status reporting

**Command Options:**
- **Default**: Start server on port 8888 serving current directory
- **Custom Port**: `--port=3000` - Start on specified port
- **Custom Directory**: `--dir=/path/to/docs` - Serve specific markdown directory
- **Health Check**: `--health-check` - Verify server is responding
- **Stop Server**: `--stop` - Gracefully stop running application processes

**Development Features:**
- Automatic port conflict detection and resolution
- Health endpoint monitoring (`/` with markdown content check)
- Process lifecycle management (start/stop/restart)
- Development-friendly error reporting
- Integration with file watchers and live reload tools

**Validated Command Patterns:**
- `go run . --port 8888 --dir .` - Default development server
- `go run . --port 3000 --dir /docs` - Custom port and directory
- `lsof -ti:8888` - Port conflict detection
- `curl -s http://localhost:8888/` - Health check
- `kill $(ps aux | grep '[g]o run' | awk '{print $2}')` - Process cleanup

**Health Monitoring:**
- HTTP connectivity verification
- Response time measurement
- Content type validation (text/html expected)
- Error detection and reporting
- Automatic retry logic for startup delays

**Port Management:**
- Automatic detection of available ports
- Conflict resolution with existing services
- Port range suggestion (8888, 8889, 8890, etc.)
- Network interface binding verification
- IPv4/IPv6 compatibility checking

**Process Lifecycle:**
```bash
# Start default development server
/serve-dev

# Start with custom configuration
/serve-dev --port=3000 --dir=./docs

# Check if server is healthy
/serve-dev --health-check

# Stop all application processes
/serve-dev --stop

# Restart server (stop + start)
/serve-dev --restart
```

**Integration with Go Web Applications:**
- Supports common command-line options and flags
- Works with various content serving patterns
- Compatible with CSS frameworks and styling systems
- Handles dynamic content rendering
- Supports common web application features

**Development Workflow Integration:**
- Compatible with `/templ-rebuild` for template development
- Works alongside file watchers and live reload
- Integrates with testing and linting commands
- Supports multiple concurrent development sessions

**Error Handling and Diagnostics:**
- Detailed startup failure analysis
- Port binding error detection
- File permission issue reporting
- Network connectivity troubleshooting
- Process conflict resolution

**Performance Monitoring:**
- Server startup time measurement
- Response time benchmarking
- Memory usage tracking (via ps)
- Request throughput estimation
- Resource utilization reporting

**Security Considerations:**
- Development-only configuration warnings
- Local network binding restrictions
- File access permission validation
- Directory traversal protection verification
- Safe process cleanup procedures

**Status Reporting:**
The command provides comprehensive status information:
- Server process ID and status
- Port binding and network interface
- Served directory and file count
- Health check results and response times
- Resource usage and performance metrics

**Cleanup and Maintenance:**
- Graceful shutdown signal handling
- Orphaned process detection and cleanup
- Temporary file management
- Log file rotation (if applicable)
- Development artifact cleanup

**Usage Examples:**

```bash
# Quick development server start
/serve-dev

# Production-like testing with specific config
/serve-dev --port=80 --dir=/var/www/docs

# Health monitoring during development
/serve-dev --health-check

# Clean restart during debugging
/serve-dev --stop && /serve-dev
```

**Development Best Practices:**
- Use consistent port numbering across team
- Monitor server health during intensive development
- Clean up processes between development sessions
- Verify content file accessibility before serving
- Test with both relative and absolute directory paths

This command enhances your Go web application development workflow by providing robust server management, health monitoring, and integration capabilities essential for efficient development.