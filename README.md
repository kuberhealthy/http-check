# Http Check

Kuberhealthy's HTTP check

## What it is
This repository builds the container image used by Kuberhealthy to run the http-check check.

## Image
- `docker.io/kuberhealthy/http-check`
- Tags: short git SHA for `main` pushes and `vX.Y.Z` for releases.

## Quick start
- Apply the example manifest: `kubectl apply -f healthcheck.yaml`
- Edit the manifest to set any required inputs for your environment.

## Build locally
- `docker build -f ./Containerfile -t kuberhealthy/http-check:dev .`

## Contributing
Issues and PRs are welcome. Please keep changes focused and add a short README update when behavior changes.

## License
See `LICENSE`.
