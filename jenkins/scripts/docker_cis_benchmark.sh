#!/bin/bash 


TARGET_IMAGE=${1}

git clone https://github.com/docker/docker-bench-security.git
TOOL_IMAGE=docker-bench-security:CI
cd docker-bench-security
docker build --no-cache -t $TOOL_IMAGE --build-arg http_proxy --build-arg https_proxy .

TARGET_IMAGE=${TARGET_IMAGE}
IMAGE_NAME="${TARGET_IMAGE##*/}"


HOST_OUTPUT_DIR=cisdockerbench_results
mkdir -p ${HOST_OUTPUT_DIR}
DOCKER_OUTPUT_DIR="/root/results"
OUTPUT_NAME=ubuntu22_${IMAGE_NAME}.txt
OUTPUT_FILE="$DOCKER_OUTPUT_DIR/$OUTPUT_NAME"
HOST_OUTPUT_FILE="$HOST_OUTPUT_DIR/$OUTPUT_NAME"
echo "Sending output to: $HOST_OUTPUT_DIR/$OUTPUT_NAME"

# Including volumes per documentation
VOLUMES="-v $(pwd)/$HOST_OUTPUT_DIR:$DOCKER_OUTPUT_DIR:rw "
STD_VOLUMES=("/etc" 
            "/lib/systemd/system"
            "/usr/bin/containerd"
            "/usr/bin/runc"
            "/usr/lib/systemd"
            "/var/lib"
            "/var/run/docker.sock")
for DIR in "${STD_VOLUMES[@]}"; do VOLUMES+="-v $DIR:$DIR:ro "; done

# Sections of report applicable to pipeline framework
OPTIONS="-c container_images,\
            container_runtime,\
            docker_security_operations,\
            community_checks"

DOCKER_BENCH_ARGS="$OPTIONS -i $IMAGE_NAME -l $OUTPUT_FILE"

# Run container in detached mode to trigger Container Runtime section
docker run -td --name $IMAGE_NAME $TARGET_IMAGE
sleep 3

docker run --net host --pid host --userns host \
--cap-add audit_control \
-e DOCKER_CONTENT_TRUST \
-e TARGET_IMAGE=$TARGET_IMAGE \
-e http_proxy \
-e https_proxy \
$VOLUMES \
--label docker_bench_security $TOOL_IMAGE $DOCKER_BENCH_ARGS | tee ${HOST_OUTPUT_FILE}


NUM_OF_WARN_MSG=$(cat ${HOST_OUTPUT_FILE} | grep "WARN" -c)

if [[ $NUM_OF_WARN_MSG!=0 ]]; then
    echo "::warning::Review ${NUM_OF_WARN_MSG} [WARN] messages in artifact report"
    exit 1
fi
