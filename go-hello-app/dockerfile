# should be pinned by digest not by tag
FROM golang:1.26-bookworm AS builder 

WORKDIR /app

COPY app/go.* ./
RUN go mod download
COPY app/ .

RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/hello-wiremind .

FROM gcr.io/distroless/static-debian13:nonroot

COPY --from=builder /out/hello-wiremind /app/hello-wiremind

USER nonroot:nonroot
EXPOSE 8080
ENTRYPOINT ["/app/hello-wiremind"]