#!/bin/bash -xe


ROOT_DIR="$(git rev-parse --show-toplevel)"
. ${ROOT_DIR}/.github/coverity/enviroment.sh

COVERITY_URL=https://scan.coverity.com/download/${LANGUAGE}/linux64
COVERITY_TOOL_FILE=coverity_tool_${LANGUAGE}.tgz

if [ -e "${COVERITY_TOOL_FILE}" ]; then
 echo coverity already installed, skipping installation
 exit 0
fi 
# Check if required environment variables are set
if [ -z "$COVERITY_PROJECT" ] || [ -z "$COVERITY_TOKEN" ] || [ -z "$LANGUAGE" ]; then
  echo "Error: COVERITY_PROJECT, COVERITY_TOKEN and LANGUAGE environment variables must be set."
  exit 1
fi

# Download Coverity Scan
echo "Downloading Coverity Scan..."
curl -L ${COVERITY_URL} \
  --output ${COVERITY_TOOL_FILE} \
  --data "token=${COVERITY_TOKEN}" \
  --data "project=${COVERITY_PROJECT}"

# Extract Coverity Scan
mkdir -p /opt/coverity/${LANGUAGE}
echo "Extracting Coverity Scan..."
tar -xzf ${COVERITY_TOOL_FILE} -C /opt/coverity/${LANGUAGE}

echo "coverity installation completed successfully."
echo "binary installed in /opt/coverity/${LANGUAGE}/bin/ folder"
