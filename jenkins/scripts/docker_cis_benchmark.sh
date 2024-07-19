#!/bin/bash -x


TARGET_IMAGE=${1}

TOOL_VER=1.6.1
TOOL_IMAGE=docker-bench-security:${TOOL_VER}

CURRENT_DIR=$(pwd)

git clone https://github.com/docker/docker-bench-security.git /tmp/docker-bench-security
cd /tmp/docker-bench-security
git checkout "v${TOOL_VER}"
docker build -t $TOOL_IMAGE --build-arg http_proxy --build-arg https_proxy .

cd ${CURRENT_DIR}

TARGET_IMAGE=${TARGET_IMAGE/.tar.gz/}
IMAGE_NAME="${TARGET_IMAGE##*/}"


HOST_OUTPUT_DIR=cisdockerbench_results
rm -rf ${HOST_OUTPUT_DIR}
mkdir -p ${HOST_OUTPUT_DIR}
DOCKER_OUTPUT_DIR="/root/results"
OUTPUT_NAME=docker_cis_report.txt
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
OPTIONS="-c docker_bench_security,docker_security_operations,container_images,container_runtime,community_checks"

DOCKER_BENCH_ARGS="$OPTIONS -i $IMAGE_NAME -l $OUTPUT_FILE"


echo "stopping and remove running containers"
docker stop $(docker ps -q)
docker container prune -f

# Run container in detached mode to trigger Container Runtime section
docker run -td --entrypoint sleep --name $IMAGE_NAME $IMAGE_NAME 10

docker run --net host --pid host --userns host \
--cap-add audit_control \
-e DOCKER_CONTENT_TRUST \
-e TARGET_IMAGE=$TARGET_IMAGE \
-e http_proxy \
-e https_proxy \
$VOLUMES \
--label docker_bench_security \
${TOOL_IMAGE} $DOCKER_BENCH_ARGS | tee ${HOST_OUTPUT_FILE}


NUM_OF_WARN_MSG=$(cat ${HOST_OUTPUT_FILE} | grep "WARN" -c)

if [[ $NUM_OF_WARN_MSG!=0 ]]; then
    echo "::warning::Review ${NUM_OF_WARN_MSG} [WARN] messages in artifact report"
fi
