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
      service: "companyinfov1.CompanyInfoService"
      method_regex: "^GetCompanyInfoByIinOrBin$"
    action:
      spammer:
        count: 5
        delay: { min_ms: 1000, max_ms: 5000 }

  - name: "chaos-test"
    match:
      service: "companyinfov1.CompanyInfoService"
      method_regex: "^GetCompanyInfoByIinOrBin$"
    action:
      chaos:
        actions:
          - delay: { min_ms: 100, max_ms: 300 }
          - abort:
              code: "UNAVAILABLE"
              percentage: 50

  - name: "abort-test"
    match:
      service: "companyinfov1.CompanyInfoService"
      method_regex: "^GetCompanyInfoByIinOrBin$"
    action:
      abort:
        code: "UNAVAILABLE"
        percentage: 50

  - name: "delay-test"
    match:
      service: "companyinfov1.CompanyInfoService"
      method_regex: "^GetCompanyInfoByIinOrBin$"
    action:
      delay: { min_ms: 100, max_ms: 300}
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
    service: "companyinfov1.CompanyInfoService"
    method_regex: "^GetCompanyInfoByIinOrBin$"
  action:
    delay: { min_ms: 100, max_ms: 300 }
```

### Abort
Aborts requests with a given gRPC code in a percentage of cases.
```yaml
- name: "abort-test"
  match:
    service: "companyinfov1.CompanyInfoService"
    method_regex: "^GetCompanyInfoByIinOrBin$"
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
    service: "companyinfov1.CompanyInfoService"
    method_regex: "^GetCompanyInfoByIinOrBin$"
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
    service: "companyinfov1.CompanyInfoService"
    method_regex: "^GetCompanyInfoByIinOrBin$"
  action:
    chaos:
      actions:
        - delay: { min_ms: 100, max_ms: 300 }
        - abort:
            code: "UNAVAILABLE"
            percentage: 50
```

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
    - `chaos`: `{ actions: [ ... ] }` (list of actions as above)
```
```

## Hot Reload

The proxy supports hot-reloading of the config file. Just update the YAML file and the proxy will reload rules automatically.

## Contributing

Contributions are welcome! Please open issues or pull requests.

## License

MIT License. See [LICENSE](LICENSE) for details.
