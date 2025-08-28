---
allowed-tools: Bash(templ:*), Bash(go run:*), Bash(kill:*), Bash(ps:*), Bash(lsof:*), Bash(pkill:*), TodoWrite
argument-hint: [--watch] [--port=PORT] [--proxy-port=PORT]
description: Watch and rebuild templ templates with live reload for Go web applications
model: claude-3-5-haiku@20241022
---

I'll help you manage templ template development with automatic rebuilding and
live reload for your Go web application.

**Prerequisites:**
- templ CLI tool must be installed (`go install github.com/a-h/templ/cmd/templ@latest`)
- Project uses .templ files for HTML templating
- Go application can be run with `go run .`

**Development Workflow:**
1. Check for existing templ processes and clean up if needed
2. Start templ in watch mode with proxy configuration
3. Monitor template changes and regenerate Go code automatically
4. Provide live browser reload capabilities
5. Handle graceful shutdown on interruption

**Command Options:**
- **Default**: `templ generate --watch --proxy="http://localhost:8888" --cmd="go run ."`
- **Custom Port**: `templ generate --watch --proxy="http://localhost:PORT" --cmd="go run ."`
- **Watch Only**: `templ generate --watch` (no proxy/server)
- **One-time Generate**: `templ generate` (no watching)

**Process Management:**
- Automatically detects and terminates existing templ processes
- Manages both templ watcher and Go server processes
- Provides clean shutdown handling
- Reports process status and port usage

**Live Reload Features:**
- Browser automatically refreshes on template changes
- Proxy handles both template regeneration and server restart
- Fast incremental builds for rapid development
- Error reporting for template compilation issues

**Validated Command Patterns:**
- `templ generate --watch --proxy="http://localhost:8888" --cmd="go run ."` - Full development mode
- `templ generate --watch --proxy="http://localhost:3000" --cmd="go run . --port 3000"` - Custom port
- `templ generate --watch` - Watch templates only, no server
- `templ generate` - One-time generation
- `kill $(ps aux | grep '[t]empl' | awk '{print $2}')` - Process cleanup

**Error Handling:**
- Template syntax error reporting with file/line references
- Port conflict detection and resolution
- Process cleanup on unexpected termination
- Go compilation error forwarding

**Integration with Go Web Applications:**
The command integrates with your web application's development workflow:
- Templates in project directory (layout.templ, etc.)
- Automatic Go code generation to *.templ.go files
- Server restart on template changes
- Preserves application command-line arguments

**Usage Examples:**

```bash
# Start full development mode (default port 8080 or app-specific)
/templ-rebuild

# Start with custom port
/templ-rebuild --port=3000

# Watch templates only, no server
/templ-rebuild --watch

# One-time template generation
/templ-rebuild --no-watch
```

**Process Monitoring:**
The command monitors:
- templ watcher process status
- Go server process health
- Port availability and conflicts
- Template file changes and compilation status
- Browser proxy connection status

**Cleanup and Shutdown:**
- Graceful process termination on Ctrl+C
- Automatic cleanup of orphaned processes
- Preservation of generated .templ.go files
- Clean exit status reporting

**Development Tips:**
- Use with your existing `go mod download` and `go build` workflow
- Compatible with various CLI frameworks (urfave/cli, cobra, etc.)
- Works with common web features (content processing, CSS frameworks, etc.)
- Supports concurrent development of multiple template files

**Performance Optimization:**
- Fast incremental template compilation
- Efficient file watching with minimal CPU usage
- Quick server restart cycles
- Optimized proxy configuration for minimal latency

This command streamlines the development workflow for templ-based Go web applications, providing automated template watching and live reload capabilities essential for efficient development.
