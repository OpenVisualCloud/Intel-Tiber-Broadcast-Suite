# Run BCS launcher Kubernetes Mode controller tests

The envtest environment requires binaries like etcd, kube-apiserver, and kubectl. You can install these binaries using the setup-envtest tool provided by controller-runtime.

## Run the following command to install the required binaries

```bash
cd <repo>/launcher/cmd
go install sigs.k8s.io/controller-runtime/tools/setup-envtest@latest
export PATH=$PATH:$(go env GOPATH)/bin
```

If you are using go version Golang 1.23, it is recommended by official [publisher](https://pkg.go.dev/sigs.k8s.io/controller-runtime/tools/setup-envtest@v0.0.0-20250604165838-d6126d850224#section-readme), to download branch: 
```bash
cd <repo>/launcher/cmd
go install sigs.k8s.io/controller-runtime/tools/setup-envtest@release-0.20
export PATH=$PATH:$(go env GOPATH)/bin
```

## Download the Kubernetes Binaries

Make sure to replace `<repo>` with your actual repository path.

```bash
cd <repo>/launcher/internal/controller
setup-envtest use 1.29.0 --bin-dir=./bin/k8s
```

You should have directory `bin` with content created in `<repo>/launcher/internal/controller`

## Run the Controller Tests

You can run the K8s controller tests using the following command:
```bash
cd <repo>/launcher/internal/controller
# Ensure you have the necessary dependencies installed
go mod tidy
# Run the tests
go test ./... -v
```
