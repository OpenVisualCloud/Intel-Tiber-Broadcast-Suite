#!/bin/bash -e

<<<<<<< HEAD

COVERITY_URL=https://scan.coverity.com/download/${LANGUAGE}/linux64

# Check if required environment variables are set
if [ -z "$COVERITY_PROJECT" ] || [ -z "$COVERITY_TOKEN" ] || [ -z "$COVERITY_EMAIL" ] || [ -z "$LANGUAGE" ]; then
  echo "Error: COVERITY_PROJECT, COVERITY_TOKEN, COVERITY_EMAIL and LANGUAGE environment variables must be set."
=======
COVERITY_URL=https://scan.coverity.com/download/${LANGUAGE}/linux64
COVERITY_TOOL_FILE=coverity_tool_${LANGUAGE}.tgz
 
# Check if required environment variables are set
if [ -z "$COVERITY_PROJECT" ] || [ -z "$COVERITY_TOKEN" ] || [ -z "$LANGUAGE" ]; then
  echo "Error: COVERITY_PROJECT, COVERITY_TOKEN and LANGUAGE environment variables must be set."
>>>>>>> 3e5c1ac (change coverity build machine)
  exit 1
fi

# Download Coverity Scan
<<<<<<< HEAD
echo "Downloading Coverity Scan..."
curl -L ${COVERITY_URL} \
  --output coverity_tool.tgz \
=======
echo "Downloading ${LANGUAGE} Coverity Scan..."
curl -L ${COVERITY_URL} \
  --output ${COVERITY_TOOL_FILE} \
>>>>>>> 3e5c1ac (change coverity build machine)
  --data "token=${COVERITY_TOKEN}" \
  --data "project=${COVERITY_PROJECT}"

# Extract Coverity Scan
<<<<<<< HEAD
mkdir -p /opt/coverity
mkdir -p /opt/coverity/bin
echo "Extracting Coverity Scan..."
tar -xzf coverity_tool.tgz -C /opt/coverity/bin

echo "coverity installation completed successfully."
echo "run binary from /opt/coverity/bin/ folder"
=======
sudo mkdir -p /opt/coverity/${LANGUAGE}
echo "Extracting ${LANGUAGE} Coverity Scan..."
sudo tar -xzf ${COVERITY_TOOL_FILE} -C /opt/coverity/${LANGUAGE}

echo "coverity installation completed successfully."
echo "binary installed in /opt/coverity/${LANGUAGE}/bin/ folder"
>>>>>>> 3e5c1ac (change coverity build machine)
