# NMOS controller that tests IS-04 and IS-05 - for testig purposes

## Flow

This script will activate NMOS node sender and nother NMOS node receiver, fetch data from the specified NMOS sender and receiver nodes, update the necessary JSON configuration files, and send PATCH requests to receiver to establish the connection between two nodes.

## Usage

Provide the appropriate IP addresses and ports for the sender NMOS node A and receiver NMOS node B.

```bash
python3 threaded-nmos-controller05.py --receiver_ip <ip-address-nmos-node-b-receiver> --sender_ip <ip-address-nmos-node-a-sender> --receiver_port <port-nmos-node-b-receiver> --sender_port <port-nmos-node-a-sender>
```

Files `sender.json` and `receiver.json` are send appropraitely to sender and reciever. Currently, no need to edit.

# Docker

```bash
# Build the Docker image
docker build -t nmos-is05-controller:latest .

# In case of issues with proxy, try:
docker build --build-arg HTTP_PROXY=<proxy> --b
uild-arg HTTPS_PROXY=<proxy> -t nmos-is
05-controller:latest .

# Run the Docker container with overwritten CMD
docker run -e RECEIVER_IP=localhost -e SENDER_IP=localhost -e RECEIVER_PORT=90 -e SENDER_PORT=95 nmos-is05-controller:latest
```

> `-e RECEIVER_IP=localhost -e SENDER_IP=localhost -e RECEIVER_PORT=90 -e SENDER_PORT=95` are the addresses of NMOS node that one acts as sender and the other as receiver.
