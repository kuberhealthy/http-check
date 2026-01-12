# Kuberhealthy HTTP Check

This check performs HTTP requests to a configured URL and verifies that a minimum percentage of requests return the expected status code.

## What It Does

1. Parses configuration from environment variables.
2. Sends a number of HTTP requests to the target URL.
3. Sleeps between requests when configured.
4. Fails if the percentage of successful responses is below the threshold.

## Configuration

All configuration is controlled via environment variables.

- `CHECK_URL`: URL to query (required).
- `COUNT`: Number of requests to perform. Default `0`.
- `SECONDS`: Delay between requests. Default `0`.
- `PASSING_PERCENT`: Percentage of successful requests required. Default `100`.
- `REQUEST_TYPE`: HTTP method to use. Default `GET`.
- `REQUEST_BODY`: Body for non-GET requests. Default `{}`.
- `EXPECTED_STATUS_CODE`: Expected status code. Default `200`.

Kuberhealthy injects these variables automatically into the check pod:

- `KH_REPORTING_URL`
- `KH_RUN_UUID`
- `KH_CHECK_RUN_DEADLINE`

## Build

Use the `Justfile` to build or test the check:

```bash
just build
just test
```

## Example HealthCheck

```yaml
apiVersion: kuberhealthy.github.io/v2
kind: HealthCheck
metadata:
  name: http-check
  namespace: kuberhealthy
spec:
  runInterval: 5m
  timeout: 5m
  podSpec:
    spec:
      containers:
        - name: http-check
          image: kuberhealthy/http-check:sha-<short-sha>
          imagePullPolicy: IfNotPresent
          env:
            - name: CHECK_URL
              value: "https://example.com"
            - name: COUNT
              value: "10"
            - name: SECONDS
              value: "1"
            - name: PASSING_PERCENT
              value: "100"
            - name: REQUEST_TYPE
              value: "GET"
            - name: EXPECTED_STATUS_CODE
              value: "200"
      restartPolicy: Never
```

A full install bundle is available in `healthcheck.yaml`.

## Image Tags

- `sha-<short-sha>` tags are published on every push to `main`.
- `vX.Y.Z` tags are published when a matching Git tag is pushed.
