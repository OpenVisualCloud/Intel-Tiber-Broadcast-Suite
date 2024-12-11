# 
# SPDX-FileCopyrightText: Copyright (c) 2024 Intel Corporation
# 
# SPDX-License-Identifier: BSD-3-Clause
# 

# Build the manager binary
FROM golang:1.23 AS builder
ARG TARGETOS
ARG TARGETARCH

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

# Copy the go source
COPY cmd/main.go cmd/main.go
COPY api/ api/
COPY internal/ internal/
COPY resources_library/ resources_library/

# Build
# the GOARCH has not a default value to allow the binary be built according to the host where the command
# was called. 
RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} go build -a -o manager cmd/main.go

# # Use distroless as minimal base image to package the manager binary
# # Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /workspace/manager .
USER 65532:65532

HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 CMD [ "curl", "-f", "http://localhost:8081/healthz" ] || exit 1

ENTRYPOINT ["/manager"]
