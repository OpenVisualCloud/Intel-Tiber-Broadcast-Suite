COV_VER="$(ls /opt/coverity/bin/)"
COVERITY_BIN_DIR="/opt/coverity/bin/$COV_VER/bin"
export COV_BUILD="${COVERITY_BIN_DIR}/cov-build"
export JAVA_HOME="/opt/coverity/bin/$COV_VER/jre/bin/java"