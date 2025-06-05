# Run BCS launcher Kubernetes Mode controller tests

The envtest environment requires binaries like etcd, kube-apiserver, and kubectl. You can install these binaries using the setup-envtest tool provided by controller-runtime.

## Run the following command to install the required binaries

```bash
cd <repo>/launcher/cmd
go install sigs.k8s.io/controller-runtime/tools/setup-envtest@latest
```

## Download the Kubernetes Binaries

Make sure to replace `<repo>` with your actual repository path.

```bash
cd <repo>/launcher/controller
setup-envtest use 1.29.0 --bin-dir=./bin/k8s
```

## Run the Controller Tests

You can run the controller tests using the following command:
```bash
cd <repo>/launcher/controller
# Ensure you have the necessary dependencies installed
go mod tidy
# Run the tests
go test ./... -v
```