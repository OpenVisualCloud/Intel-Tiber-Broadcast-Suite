!/bin/bash -e





function update_and_build_dockerfile() {
      echo "update dockerfile with patched ffmpeg and onevpl instalation steps"
      INJECTION_LINE=$(awk '/Tiber Suite final-stage/{print NR-2}' docker/app/Dockerfile )
      awk 'NR=='${INJECTION_LINE}' {
            while(getline line < "bdba_command.txt") 
                  print line} 1' docker/app/Dockerfile > docker/app/Dockerfile.new

      mv docker/app/Dockerfile.new docker/app/Dockerfile

      echo " updating dockerfile with non patched ffmpeg and onevpl instalation steps"
      INJECTION_LINE=$(awk '/Tiber Suite final-stage/{print NR-2}' docker/app/Dockerfile )
      awk 'NR=='${INJECTION_LINE}' {
                  while(getline line < "bdba_command.txt") 
                        print line} 1' docker/app/Dockerfile > docker/app/Dockerfile.new

      mv docker/app/Dockerfile.new docker/app/Dockerfile

      echo "build modified dockerfile with patched and unpatched code"

      echo "export build stage as docker image"
      sed -i 's/--target final-stage/--target build-stage/' build.sh
      export IMAGE_TAG="bdba-build"
      export BUILD_TYPE=CI
      ./build.sh

      # revoke changes
      git checkout .

      echo "open another container and run:"
      echo "(sudo) docker run -it tiber-broadcast-suite:bdba-build bash"
      echo "then get back to first terminal and run:"
      echo " (sudo) ./BDBA.sh -f"
}
function get_docker_container_binaries(){
      echo "get container ID"
      CONTAINER_ID="$(docker ps | grep tiber-broadcast-suite:bdba-build | awk '{print $1}')"

      echo "prepare directories for binaries"
      BDBA_BIN_FOLDER="bdba_bins"
      rm -rf ${BDBA_BIN_FOLDER}
      mkdir -p ${BDBA_BIN_FOLDER}/ffmpeg-patched
      mkdir -p ${BDBA_BIN_FOLDER}/ffmpeg-unpatched
      mkdir -p ${BDBA_BIN_FOLDER}/onevpl-patched/bin
      mkdir -p ${BDBA_BIN_FOLDER}/onevpl-patched/lib
      mkdir -p ${BDBA_BIN_FOLDER}/onevpl-unpatched/bin
      mkdir -p ${BDBA_BIN_FOLDER}/onevpl-unpatched/lib

      echo "copy binaries from running container ${CONTAINER_ID} to host"
      sudo docker cp ${CONTAINER_ID}:/tmp/ffmpeg_pure/*_g ./${BDBA_BIN_FOLDER}/ffmpeg-unpatched
      sudo docker cp ${CONTAINER_ID}:/tmp/ffmpeg_pure/ffmpeg ./${BDBA_BIN_FOLDER}/ffmpeg-unpatched
      sudo docker cp ${CONTAINER_ID}:/tmp/ffmpeg_pure/ffplay ./${BDBA_BIN_FOLDER}/ffmpeg-unpatched
      sudo docker cp ${CONTAINER_ID}:/tmp/ffmpeg_pure/ffprobe ./${BDBA_BIN_FOLDER}/ffmpeg-unpatched

      # patched ffmpeg is eventually moved to /usr/bin 
      sudo docker cp ${CONTAINER_ID}:/buildout/usr/bin ./${BDBA_BIN_FOLDER}/ffmpeg-patched

      sudo docker cp ${CONTAINER_ID}:/tmp/onevpl_pure/build/__bin/release/  ./${BDBA_BIN_FOLDER}/onevpl-unpatched/bin
      sudo docker cp ${CONTAINER_ID}:/tmp/onevpl_pure/build/__lib/release/ ./${BDBA_BIN_FOLDER}/onevpl-unpatched/lib

      sudo docker cp ${CONTAINER_ID}:/tmp/onevpl/build/__bin/release/ ./${BDBA_BIN_FOLDER}/onevpl-patched/bin
      sudo docker cp ${CONTAINER_ID}:/tmp/onevpl/build/__lib/release/ ./${BDBA_BIN_FOLDER}/onevpl-patched/lib

      echo "binaries copied to ${BDBA_BIN_FOLDER} folder"

      echo "use: "
      echo "scp -r gta@SED-Val-2:<abs path>/${BDBA_BIN_FOLDER} ./BCS_BDBA"
      echo "to copy the binaries to your local machine"

      echo "cleanup"
      echo "stop container"
      docker stop ${CONTAINER_ID}
      echo "remove container"
      docker rm ${CONTAINER_ID}
}


function usage() {
    echo "Usage: $0 [-h] [-p] [-u]"
    echo "  -h  Show this help message"
    echo "  -b  Build modified image"
    echo "  -f  fetch binaries from running container"
}
# Parse command line arguments
while getopts ":hbf" opt; do
    case $opt in
        h)
            usage
            exit 0
            ;;
        b)
            update_and_build_dockerfile
            exit 0
            ;;
        f)
            get_docker_container_binaries
            exit 0
            ;;
        \?)
            echo "Invalid option: -$OPTARG" >&2
            usage
            exit 1
            ;;
    esac
done
