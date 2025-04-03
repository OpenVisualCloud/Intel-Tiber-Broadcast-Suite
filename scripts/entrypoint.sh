#!/bin/bash

# Check if a configuration file is provided as an argument
if [ "$#" -ne 1 ]; then
    echo "Usage: $0 <config.json>"
    exit 1
fi

CONFIG_FILE="$1"

# Check if the provided file exists
if [ ! -f "$CONFIG_FILE" ]; then
    echo "Error: Configuration file '$CONFIG_FILE' does not exist."
    exit 1
fi

# Start the mDNSResponder service
echo -e "\nStarting mDNSResponder service"
/etc/init.d/mdns start

# Run the bcs-nmos-node with the provided configuration file
echo -e "\nRunning bcs-nmos-node with configuration: $CONFIG_FILE"
cat $CONFIG_FILE
./bcs-nmos-node $CONFIG_FILE
