# NMOS controller that tests IS-04 and IS-05 - for testig purposes

## Flow

This script will activate NMOS node sender and nother NMOS node receiver, fetch data from the specified NMOS sender and receiver nodes, update the necessary JSON configuration files, and send PATCH requests to receiver to establish the connection between two nodes.

## Usage

Provide the appropriate IP addresses and ports for the sender NMOS node A and receiver NMOS nodeB.

```bash
python3 threaded-nmos-controller05.py --receiver_ip <ip-address-nmos-node-b-receiver> --sender_ip <ip-address-nmos-node-a-sender> --receiver_port <port-nmos-node-b-receiver> --sender_port <port-nmos-node-a-sender>
```

Files `sender.json` and `receiver.json` are send appropraitely to sender and reciever. Currently, no need to edit.
