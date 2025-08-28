---
allowed-tools: Bash(go test:*), Bash(wrk:*), Bash(curl:*), Bash(ab:*), Bash(hey:*), Bash(time:*), Bash(ps:*), Bash(top:*), TodoWrite
argument-hint: [--server=URL] [--duration=30s] [--connections=10] [--tool=wrk|ab|hey] [--endpoint=PATH]
description: Run performance benchmarks on Go web server endpoints and content processing performance
model: claude-3-5-sonnet@20241022
---

I'll help you benchmark your Go web server performance, including endpoint response times, content processing throughput, and resource utilization analysis.

**Prerequisites:**
- Running Go web server (or ability to start one)
- One of: wrk, ab (Apache Bench), or hey benchmarking tools
- Go test framework for internal benchmarks
- Sample markdown files for realistic testing

**Benchmark Categories:**
1. **HTTP Endpoint Performance**: Response times, throughput, concurrency handling
2. **Content Processing**: Content processing speed, template execution
3. **Static Asset Delivery**: CSS, JavaScript, and other static resources
4. **Memory and CPU Usage**: Resource consumption under load
5. **Concurrency Testing**: Multi-user scenario simulation

**Benchmarking Tools:**

**wrk (Preferred):**
- Modern HTTP benchmarking tool
- Lua scripting capabilities for complex scenarios
- Detailed latency distribution reporting
- High concurrency support

**Apache Bench (ab):**
- Widely available HTTP benchmarking
- Simple command-line interface
- Basic concurrency and duration controls
- Standard output formatting

**hey:**
- Go-based HTTP load generator
- Good performance and detailed metrics
- HTTP/2 support
- JSON output options

**Command Options:**
- **Default**: Test localhost:8888 with moderate load (10 connections, 30 seconds)
- **Custom Server**: `--server=http://localhost:3000` - Test specific server
- **Duration**: `--duration=60s` - Benchmark duration
- **Connections**: `--connections=50` - Concurrent connection count
- **Tool Selection**: `--tool=wrk` - Specify benchmarking tool
- **Endpoint Testing**: `--endpoint=/README` - Test specific markdown file

**Test Scenarios:**

**Basic Endpoint Tests:**
- Root endpoint (`/`) response time
- Markdown file rendering (`/README`, `/docs/guide`)
- Non-existent file handling (404 responses)
- Static asset delivery (CSS, JavaScript)
- Health check endpoint performance

**Load Testing Scenarios:**
- Gradual load increase (1, 5, 10, 25, 50 connections)
- Sustained load testing (constant traffic over time)
- Spike testing (sudden traffic bursts)
- Mixed workload (different file types and sizes)
- Cache effectiveness testing

**Markdown-Specific Benchmarks:**
- Simple markdown files vs complex documents
- Mermaid diagram rendering performance
- Code syntax highlighting overhead
- Table rendering performance
- Large file processing times

**Validated Command Patterns:**
- `wrk -t4 -c10 -d30s http://localhost:8888/` - Basic load test
- `ab -n 1000 -c 10 http://localhost:8888/README` - Apache Bench test
- `hey -n 1000 -c 10 http://localhost:8888/` - hey load test
- `go test -bench=. -benchmem ./...` - Go internal benchmarks
- `curl -w "%{time_total}\n" -o /dev/null -s http://localhost:8888/` - Single request timing

**Performance Metrics:**
- **Latency**: Min, max, average, percentiles (50th, 90th, 95th, 99th)
- **Throughput**: Requests per second, bytes per second
- **Error Rate**: Failed requests, timeout percentage
- **Resource Usage**: CPU utilization, memory consumption
- **Concurrency**: Maximum sustainable concurrent users

**Report Sections:**
1. **Server Information**: Version, configuration, system specs
2. **Test Configuration**: Tool used, duration, connection count
3. **Response Time Analysis**: Latency distribution and percentiles
4. **Throughput Metrics**: RPS, bandwidth utilization
5. **Error Analysis**: Failed requests and error types
6. **Resource Utilization**: CPU, memory, and system load
7. **Recommendations**: Performance optimization suggestions

**Usage Examples:**

```bash
# Quick performance check
/benchmark-server

# Intensive load testing
/benchmark-server --connections=100 --duration=120s

# Test specific markdown file
/benchmark-server --endpoint=/docs/large-document

# Compare different tools
/benchmark-server --tool=wrk --duration=30s
/benchmark-server --tool=ab --duration=30s

# Test custom server instance
/benchmark-server --server=http://production-server:8080
```

**Integration with Go Web Applications:**
- Tests file serving and routing performance
- Validates content processing efficiency
- Measures CSS framework delivery performance
- Evaluates dynamic content processing overhead
- Assesses template rendering speed

**Advanced Testing Features:**
- **Realistic Workloads**: Mixed file sizes and types
- **Cache Testing**: Repeated requests to measure caching effectiveness
- **Progressive Load**: Gradually increasing concurrent users
- **Endurance Testing**: Long-duration stability testing
- **Memory Leak Detection**: Resource usage monitoring over time

**Performance Analysis:**
```
🚀 MDViewer Performance Benchmark Report
═══════════════════════════════════════

🖥️  Server: http://localhost:8080
⏱️  Test Duration: 30 seconds
🔗 Concurrent Connections: 10
🛠️  Tool: wrk v4.2.0

📊 Response Time Analysis:
   Average: 15.2ms
   50th percentile: 12.1ms
   90th percentile: 28.4ms
   95th percentile: 35.7ms
   99th percentile: 52.3ms
   Maximum: 67.8ms

📈 Throughput Metrics:
   Requests/sec: 657.32
   Transfer/sec: 2.1MB
   Total Requests: 19,720
   Total Transfer: 63.2MB

✅ Reliability:
   Success Rate: 100% (0 errors)
   Timeouts: 0
   Socket Errors: 0

🧠 Resource Usage:
   Peak CPU: 45%
   Peak Memory: 89MB
   System Load: 1.2

📝 Endpoint Breakdown:
   /: 8,945 requests (avg: 12.1ms)
   /README: 5,432 requests (avg: 18.7ms)
   /docs/guide: 3,221 requests (avg: 22.4ms)
   /static/pico.css: 2,122 requests (avg: 3.2ms)

🎯 Performance Rating: Excellent
   • Response times well under 100ms
   • High throughput with stable performance
   • No errors or timeouts detected
   • Resource usage within acceptable limits

💡 Recommendations:
   • Consider response caching for popular endpoints
   • Monitor performance with larger markdown files
   • Test with higher concurrency for production readiness
```

**Performance Optimization Insights:**
- Identifies bottlenecks in markdown processing
- Reveals template rendering performance characteristics
- Measures static asset delivery efficiency
- Evaluates server resource utilization patterns
- Provides baseline metrics for optimization efforts

**CI/CD Integration:**
- Exit codes for pass/fail criteria
- JSON output for automated processing
- Performance regression detection
- Threshold-based alerting
- Historical performance tracking

This command provides comprehensive performance analysis for your Go web application, helping ensure optimal user experience and identify optimization opportunities before deployment.