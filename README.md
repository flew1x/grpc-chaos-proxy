# grpc-chaos-proxy

Chaos engineering for a gRPC zoo, all-in-one and without vendor lock-in.

## Overview

**grpc-chaos-proxy** is a tool for introducing chaos into your gRPC-based systems. It acts as a proxy between your gRPC clients and servers, allowing you to inject failures, delays, aborts, spamming, and other network anomalies to test the resilience of your microservices.

## Features

- Inject artificial delays and errors into gRPC traffic
- Abort requests with configurable error codes and percentage
- Simulate network latency (randomized min/max delay)
- Spam requests to backend (for load/chaos testing)
- Compose complex chaos actions (randomly pick from a set)
- Flexible configuration via YAML files (see `configs/dev.yaml`)
- Hot-reload config support
- No vendor lock-in, fully open source

## Installation

### Prerequisites

- Go 1.20+ installed
- (Optional) Docker for containerized usage

### Build from Source

```bash
git clone https://github.com/flew1x/grpc-chaos-proxy
cd grpc-chaos-proxy
make build
```

The binary will be available at `bin/proxy`.

## Quick Start

### 1. Prepare a Configuration

Example (`configs/dev.yaml`):

```yaml
listener:
  address: "localhost:5050"

backend:
  address: "localhost:5010"

rules:
  - name: "spammer-test"
    match:
      service: "protoann.Service"
      method_regex: "^GetByIinOrBin$"
    action:
      spammer:
        count: 5
        delay: { min_ms: 1000, max_ms: 5000 }

  - name: "chaos-test"
    match:
      service: "protoann.Service"
      method_regex: "^GetByIinOrBin$"
    action:
      chaos:
        actions:
          - delay: { min_ms: 100, max_ms: 300 }
          - abort:
              code: "UNAVAILABLE"
              percentage: 50

  - name: "abort-test"
    match:
      service: "protoann.Service"
      method_regex: "^GetByIinOrBin$"
    action:
      abort:
        code: "UNAVAILABLE"
        percentage: 50

  - name: "delay-test"
    match:
      service: "protoann.Service"
      method_regex: "^GetByIinOrBin$"
    action:
      delay: { min_ms: 100, max_ms: 300}

  - name: "header-inject-test"
    match:
      service: "protoann.Service"
      method_regex: "^GetByIinOrBin$"
    action:
      header:
        headers:
          x-custom:
            prefix: "pre-"
            suffix: "-suf"
            values: ["val1", "val2"]
          x-another: "static-value"
        allowlist: ["x-custom", "x-another"]
        direction: "inbound" # or "outbound", or "both"

  - name: "ratelimit-test"
    match:
      service: "protoann.Service"
      method_regex: "^GetByIinOrBin$"
    action:
      ratelimit:
        rate_limit: 5    # allowed requests per second
        burst_size: 2    # additional burst capacity
```

### 2. Start the Proxy

```bash
./bin/proxy --config configs/dev.yaml
```

### 3. Point Your gRPC Client

Change your gRPC client to connect to the proxy (`localhost:5050`), which will forward requests to your real server (`localhost:5010`) and inject chaos as configured.

## Rule Types & Examples

### Delay
Injects a random delay before forwarding the request.
```yaml
- name: "delay-test"
  match:
    service: "protoann.Service"
    method_regex: "^GetByIinOrBin$"
  action:
    delay: { min_ms: 100, max_ms: 300 }
```

### Abort
Aborts requests with a given gRPC code in a percentage of cases.
```yaml
- name: "abort-test"
  match:
    service: "protoann.Service"
    method_regex: "^GetByIinOrBin$"
  action:
    abort:
      code: "UNAVAILABLE"
      percentage: 50
```

### Spammer
Sends multiple requests to the backend, optionally with delay between them.
```yaml
- name: "spammer-test"
  match:
    service: "protoann.Service"
    method_regex: "^GetByIinOrBin$"
  action:
    spammer:
      count: 5
      delay: { min_ms: 1000, max_ms: 5000 }
```

### Chaos (Composite)
Randomly applies one of the listed actions.
```yaml
- name: "chaos-test"
  match:
    service: "protoann.Service"
    method_regex: "^GetByIinOrBin$"
  action:
    chaos:
      actions:
        - delay: { min_ms: 100, max_ms: 300 }
        - abort:
            code: "UNAVAILABLE"
            percentage: 50
```

### Network
Simulates network failures: packet loss (loss) and artificial delay (throttle).

```yaml
- name: "network-test"
  match:
    service: "protoann.Service"
    method_regex: "^GetByIinOrBin$"
  action:
    network:
      loss_percentage: 20      # probability to drop the request, %
      throttle_ms: 200         # artificial delay in milliseconds
```

- `loss_percentage`: probability that the request will be "lost" (not forwarded to the backend)
- `throttle_ms`: delay (in ms) before forwarding the request to the backend

### Header
Injects, modifies, or removes gRPC metadata headers. Supports:
- Adding/modifying headers with prefix/suffix and multiple values
- Allowlist: keep only specified headers, remove all others
- Direction: apply only on inbound, outbound, or both traffic

Example:
```yaml
- name: "header-test"
  match:
    service: "protoann.Service"
    method_regex: "^GetByIinOrBin$"
  action:
    header:
      headers:
        x-custom:
          prefix: "pre-"
          suffix: "-suf"
          values: ["val1", "val2"]
        x-another: "static-value"
      allowlist: ["x-custom", "x-another"]
      direction: "inbound" # or "outbound", or "both"
```
- `headers`: map of header names to modification rules. Each rule can have `prefix`, `suffix`, and a list of `values` (or a single string value).
- `allowlist`: if set, only these headers will be kept, all others will be removed.
- `direction`: controls when the injection is applied: `inbound`, `outbound`, or `both` (default).

### RateLimit
Limits the number of requests per second (token bucket algorithm). Useful for simulating backend rate limiting or throttling.

```yaml
- name: "ratelimit-test"
  match:
    service: "protoann.Service"
    method_regex: "^GetByIinOrBin$"
  action:
    ratelimit:
      rate_limit: 5    # allowed requests per second
      burst_size: 2    # additional burst capacity
```
- `rate_limit`: maximum number of requests per second.
- `burst_size`: how many extra requests can be handled in a burst (optional, default 0).

If the limit is exceeded, the request will be rejected with a rate limit error.

### Disconnect
Simulates random connection drops (disconnects) by returning a gRPC error with a specified probability. Useful for testing client resilience to network failures.

Example:
```yaml
- name: "disconnect-test"
  match:
    service: "protoann.Service"
    method_regex: "^GetByIinOrBin$"
  action:
    disconnect:
      percentage: 20  # probability (0-100) to simulate a disconnect
```
- `percentage`: probability (0-100) that the request will be forcibly disconnected (default: 0).

If triggered, the proxy returns a gRPC error with code `UNAVAILABLE` and message `chaos disconnect injected`.

### Code
Injects a custom gRPC error code with advanced options. Useful for simulating specific error scenarios, custom error messages, and more.

Example:
```yaml
- name: "code-test"
  match:
    service: "protoann.Service"
    method_regex: "^GetByIinOrBin$"
  action:
    code:
      code: "UNAVAILABLE"         # gRPC code to return (see codes.go for all options)
      message: "custom error"     # custom error message (optional)
      percentage: 30              # probability (0-100) to inject error (optional)
      delay_ms: 100               # delay before returning error (ms, optional)
      metadata:
        x-debug: "true"           # custom metadata to add to response (optional)
      only_on_methods: ["GetByIinOrBin"] # apply only to these methods (optional)
      repeat_count: 2             # how many times to repeat error (optional)
```
- `code`: gRPC error code to return (e.g., UNAVAILABLE, INTERNAL, NOT_FOUND, etc.)
- `message`: custom error message (optional)
- `percentage`: probability (0-100) to inject error (optional, default: 100)
- `delay_ms`: delay before returning error in milliseconds (optional)
- `metadata`: map of metadata keys/values to add to the response (optional)
- `only_on_methods`: list of method names to apply the rule to (optional)
- `repeat_count`: how many times to repeat the error for the same request (optional)

If triggered, the proxy returns a gRPC error with the specified code and message, and can add custom metadata or delay the response.

### Script
Executes a custom shell script (sh/bash) as part of the chaos action. Useful for dynamic, programmable chaos scenarios, integration with external systems, or advanced request/response mutation.

Example:
```yaml
- name: "script-test"
  match:
    service: "protoann.Service"
    method_regex: "^GetByIinOrBin$"
  action:
    script:
      language: sh
      source: |
        if [ "$1" = "fail" ]; then
          echo "X-CHAOS-ERROR: custom script error"
          exit 1
        fi
        echo "X-CHAOS-HEADER: x-script=ok"
      args: ["fail"]
      timeout_ms: 500
      env:
        FOO: "bar"
```
- `language`: script language (currently supports `sh` or `bash`)
- `source`: script source code (string or multiline)
- `args`: arguments to pass to the script (optional)
- `timeout_ms`: script execution timeout in milliseconds (optional)
- `env`: environment variables for the script (optional)

**Special output handling:**
- If the script outputs a line starting with `X-CHAOS-ERROR:`, the proxy will treat it as an error and return it to the client.
- If the script outputs a line starting with `X-CHAOS-HEADER: key=value`, the proxy will add this header to the response metadata.

This allows you to implement custom, programmable chaos logic directly in your configuration.

## Configuration Reference

- `listener.address`: Address to listen for incoming gRPC requests (e.g., `localhost:5050`)
- `backend.address`: Address of the real gRPC server (e.g., `localhost:5010`)
- `rules`: List of rules to apply
  - `name`: Rule name (for logging)
  - `match.service`: Service name to match (exact, case-insensitive)
  - `match.method_regex`: Regex for method name
  - `action`: One of:
    - `delay`: `{ min_ms: int, max_ms: int }`
    - `abort`: `{ code: string, percentage: int }`
    - `spammer`: `{ count: int, delay: { min_ms: int, max_ms: int } }`
    - `chaos`: `{ actions: [ ... ] }` (list of actions)
    - `network`: `{ loss_percentage: int, throttle_ms: int }`
    - `header`: `{ headers: { ... }, allowlist: [ ... ], direction: "inbound|outbound|both" }`
    - `ratelimit`: `{ rate_limit: int, burst_size: int }`
    - `disconnect`: `{ percentage: int }`
    - `code`: `{ code: string, message: string, percentage: int, delay_ms: int, metadata: { ... }, only_on_methods: [ ... ], repeat_count: int }`
    - `script`: `{ language: string, source: string, args: [ ... ], timeout_ms: int, env: { ... } }`

## Hot Reload

The proxy supports hot-reloading of the config file. Update the YAML file and the proxy will reload rules automatically.

## Contributing

Contributions are welcome! Please open issues or pull requests.

## License

MIT License. See [LICENSE](LICENSE) for details.
