# gRPC Overview

## FFmpeg_wrapper_service.cc File

The `FFmpeg_wrapper_service.cc` file implements a gRPC server that executes FFmpeg commands received from clients. The server handles requests asynchronously using a completion queue and processes each request in a separate instance of the CallData class.

### Requirements
To be able to build and run the service :
- Requirement to build grp first: https://grpc.io/docs/protoc-installation/
- Then run the compilation script `./compile.sh` and using optional argument `--unit_testing` to build and run the unit tests. i.e `./compile.sh --unit_testing`.

### Run the service

gRPC service needs two arguments, interface and port. These arguments are passed in command line style.

- When using the wrapper sevice as a command line utility : ./FFmpeg_wrapper_service <interface/ip> <port> i.e `./FFmpeg_wrapper_service 10.10.10.10 5555`
- When spawning the service as a docker image entry point : docker run <name of the docker image> <interface/ip> <port> i.e `docker run video_production_image 10.10.10.10 5555`

### Key Components

**1. Includes and Using Declarations**:
- Includes necessary headers for gRPC, standard I/O, and signal handling.
- Uses gRPC classes and functions for server and request handling.

**2. Constants**:

- Defines constants for various FFmpeg command statuses and messages.

**3. CmdPassImpl Class**:

- Manages the server lifecycle and handles incoming RPC calls.

**4. CallData Class**:

-  Manages the lifecycle of an individual RPC call.
- Processes the FFmpeg command and sends the response back to the client.

**5. Signal Handling**:

- Sets up signal handlers to gracefully shut down the server on receiving SIGTERM or SIGINT.

**6. Main Function**:

- Parses command-line arguments to get the server address and port.
- Initializes and runs the server.

### Detailed Explanation

**1. CmdPassImpl Class**:

Run Method:

- Initializes the server with the provided address and port.
- Registers the service and adds a completion queue.
- Starts the server and prints a message indicating that the server is running.
- Creates a thread to wait for the stop flag and shut down the server gracefully.
- Calls HandleRpcs to start processing incoming RPC calls.

Shutdown Method:

- Sets the stop flag to true and notifies the waiting thread to shut down the server.

**2. CallData Class**:

Constructor:

- Initializes the service, completion queue, responder, and status.
- Calls Proceed to start processing the request.
- Proceed Method:
- Handles the different stages of the RPC call lifecycle (CREATE, PROCESS, FINISH).
- In the CREATE stage, requests the next FFmpeg command execution.
- In the PROCESS stage, constructs the FFmpeg command from the request and executes it using popen.
- Captures the output of the FFmpeg command and prints it.
- Sends the appropriate response based on the success or failure of the command execution.
- In the FINISH stage, deletes the CallData instance.

**3. HandleRpcs Method**:

- Creates a new CallData instance to handle incoming requests.
- Continuously processes incoming requests from the completion queue until the server is shut down.

**4. Signal Handling**:

- Sets up signal handlers to call Shutdown on receiving SIGTERM or SIGINT.

**5. Main Function** :

- Parses command-line arguments to get the server address and port.
- Initializes a CmdPassImpl instance and sets the global pointer.
- Sets up signal handlers.
- Calls Run to start the server.

--------------------------------------------------------------------------------------------------------------------

## FFmpeg_wrapper_client.cc File

The `FFmpeg_wrapper_client.cc` file implements a gRPC client that sends FFmpeg commands to a gRPC server and handles the responses asynchronously. The client is designed to manage multiple requests concurrently and ensures that all requests are completed before shutting down.

### Key Components
**1. Includes and Using Declarations**:

- Includes necessary headers for gRPC, standard I/O, and string manipulation.
- Includes the generated protobuf headers for FFmpeg command wrapping.

**2. CmdPassClient Class**:

- Manages the client lifecycle and handles sending FFmpeg commands to the server.
- Contains methods for initiating and processing asynchronous RPC calls.

**3. Constructor**:

- Initializes the gRPC channel and stub for communication with the server.
- Starts a thread to process the completion queue for handling asynchronous responses.

**4. Destructor**:

- Shuts down the completion queue and joins the processing thread to ensure a clean shutdown.

**5. FFmpegCmdExec Method**:

- Constructs the request object from a vector of command pairs.
- Initiates an asynchronous RPC call to the server.
- Increments the count of pending requests.

**6. AsyncCompleteRpc Method**:

- Continuously processes responses from the completion queue.
- Handles successful and failed RPC calls by printing appropriate messages.
- Decrements the count of pending requests and notifies if all requests are completed.

**7. WaitForAllRequests Method**:

- Waits for all pending requests to be completed before proceeding.

### Detailed Explanation

**1. CmdPassClient Class**:

- Takes the server interface and port as arguments.
- Creates a gRPC channel and stub for communication with the server.
- Starts a thread to process the completion queue for handling asynchronous responses.
- Destructor:
- Shuts down the completion queue to stop processing responses.
- Joins the processing thread to ensure it has completed before the client is destroyed.

**2. FFmpegCmdExec Method**:

- Takes a vector of command pairs as input.
- Constructs a ReqCmds request object and populates it with the command pairs.
- Creates a new AsyncClientCall object to manage the asynchronous RPC call.
- Increments the count of pending requests.
- Initiates the asynchronous RPC call and sets up the response reader to handle the response.

**3. AsyncCompleteRpc Method**:

- Runs in a separate thread to process responses from the completion queue.
- Continuously waits for responses and processes them as they arrive.
- Handles successful responses by printing a success message.
- Handles failed responses by printing error details.
- Decrements the count of pending requests and notifies if all requests are completed.

**4. WaitForAllRequests Method**:

- Uses a condition variable to wait until all pending requests are completed.
- Ensures that the client does not shut down until all requests have been processed.
