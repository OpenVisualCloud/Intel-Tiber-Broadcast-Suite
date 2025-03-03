#!/bin/bash -e

COVERITY_URL="https://scan.coverity.com/download/${LANGUAGE}/linux64"
COVERITY_TOOL_FILE="coverity_tool_${LANGUAGE}.tgz"
 
# Check if required environment variables are set
if [ -z "$COVERITY_PROJECT" ] || [ -z "$COVERITY_TOKEN" ] || [ -z "$LANGUAGE" ]; then
  echo "Error: COVERITY_PROJECT, COVERITY_TOKEN and LANGUAGE environment variables must be set."
  exit 1
fi

# Download Coverity Scan
echo "Downloading ${LANGUAGE} Coverity Scan..."
curl -L "${COVERITY_URL}" \
  --output "${COVERITY_TOOL_FILE}" \
  --data "token=${COVERITY_TOKEN}" \
  --data "project=${COVERITY_PROJECT}"

# Extract Coverity Scan
sudo mkdir -p "/opt/coverity/${LANGUAGE}"
echo "Extracting ${LANGUAGE} Coverity Scan..."
sudo tar -xzf "${COVERITY_TOOL_FILE}" -C "/opt/coverity/${LANGUAGE}"

echo "coverity installation completed successfully."
echo "binary installed in /opt/coverity/${LANGUAGE}/bin/ folder"
