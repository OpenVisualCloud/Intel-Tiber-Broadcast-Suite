# Test pipeline on baremetal

Open 5 terminal   windows:

- 1. terminal  
Run TX (NMOS node with sender)
```bash
cd <repo>/nmos/nmos-cpp/Development/build
./nmos-cpp-node ../nmos-cpp-node/node1_baremetal.json
```

- 2. terminal  
Run RX (NMOS node with receiver)
```bash
cd <repo>/nmos/nmos-cpp/Development/build
./nmos-cpp-node ../nmos-cpp-node/node2_baremetal.json
```

- 3. terminal  
Run ffmpeg for TX. Port 50051 is defined in `../nmos-cpp-node/node1_baremetal.json` as `ffmpeg_grpc_server_port`
```bash
cd <repo>/gRPC/build/
./FFmpeg_wrapper_service localhost 50051
```

- 4. terminal  
Run ffmpeg for RX. Port 50052 is defined in `../nmos-cpp-node/node2_baremetal.json` as `ffmpeg_grpc_server_port`
```bash
cd <repo>/gRPC/build/
./FFmpeg_wrapper_service localhost 50052
```

- 5. terminal  
Connect sender and receiver using NMOS IS-04 and IS-05. `--receiver_port 95 --sender_port 90` are defined in `../nmos-cpp-node/node2_baremetal.json` as `http_port`. On this port NMOS node is exposed.
```bash
cd <repo>/nmos/nmos-is05-controller
python3 threaded-nmos-controller05.py --receiver_ip localhost --sender_ip localhost --receiver_port 95 --sender_port 90
```