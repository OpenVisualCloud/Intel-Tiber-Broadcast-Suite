# Overview
The `FFmpeg_wrapper_service.cc` file implements a gRPC server that executes FFmpeg commands received from clients. The server handles requests asynchronously using a completion queue and processes each request in a separate instance of the CallData class.

### Key Components : 

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
