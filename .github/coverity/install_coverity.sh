#!/bin/bash -xe



COVERITY_URL=https://scan.coverity.com/download/${LANGUAGE}/linux64
COVERITY_TOOL_FILE=coverity_tool.tgz

if [ -e ${COVERITY_TOOL_FILE} ]; then
 echo coverity already installed, skipping installation
 exit(0)
fi 
# Check if required environment variables are set
if [ -z "$COVERITY_PROJECT" ] || [ -z "$COVERITY_TOKEN" ] || [ -z "$COVERITY_EMAIL" ] || [ -z "$LANGUAGE" ]; then
  echo "Error: COVERITY_PROJECT, COVERITY_TOKEN, COVERITY_EMAIL and LANGUAGE environment variables must be set."
  exit 1
fi

# Download Coverity Scan
echo "Downloading Coverity Scan..."
curl -L ${COVERITY_URL} \
  --output ${COVERITY_TOOL_FILE} \
  --data "token=${COVERITY_TOKEN}" \
  --data "project=${COVERITY_PROJECT}"

# Extract Coverity Scan
mkdir -p /opt/coverity
mkdir -p /opt/coverity/bin
echo "Extracting Coverity Scan..."
tar -xzf coverity_tool.tgz -C /opt/coverity/bin

echo "coverity installation completed successfully."
echo "run binary from /opt/coverity/bin/ folder"
