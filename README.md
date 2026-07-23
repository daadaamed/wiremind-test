# hello-wiremind

A minimal Go web server, containerized and deployed.

## Endpoints

- `GET /hello?firstname=X&lastname=Y`
  - Returns `Hello X Y`.
  - Responds with `400 Bad Request` if either parameter is missing.

- `GET /healthz`
  - Liveness/readiness check.

## Run locally
 
```bash
go run main.go
curl "http://localhost:8080/hello?firstname=Jane&lastname=Doe"
```
Port defaults to `8080`, override with `PORT`.

## Build & push image
 
```bash
docker build -t <registry>/<repo>/hello-wiremind:<tag> .
docker push <registry>/<repo>/hello-wiremind:<tag>
```

## Deploy to Kubernetes
 
```bash
kubectl apply -f k8s/ -n <namespace>
```

## Key choices

**Application**
- Go: chosen for its static compilation, no runtime dependencies and a small binary.
- Graceful shutdown on `SIGTERM`, structured JSON logs via `slog`.
 
**Image**
- Multi-stage build: `golang:1.26-bookworm` to compile a static binary (`CGO_ENABLED=0`), `gcr.io/distroless/static-debian13:nonroot` to run it. Result: ~3.4 MB, no shell, no package manager.
- Distroless over `scratch` specifically to keep CA certificates available.
- `-trimpath -ldflags="-s -w"` for a reproducible, stripped binary.

- Runs as nonroot (`runAsNonRoot: true`, `runAsUser: 65532`, matching distroless's built-in user)

**Kubernetes setup**
- 2 replicas for zero-downtime rolling updates, even though the cluster is single-node.
- Sized for a single-node cluster: No anti-affinity or `PodDisruptionBudget`

## What I'd rather handle in CI
 
- Building and pushing the image (manual here, via `docker build`/`docker push`).
- Tagging by git SHA instead of a static tag.
- Vulnerability scanning and push to Artifact Registry.

## What I'd rather handle at runtime
 
- Any environment-specific config (port, log level, feature flags), not hardcoded at build time, so the same image runs unchanged across environments.