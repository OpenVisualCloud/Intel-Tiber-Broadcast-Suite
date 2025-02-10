#!/bin/bash -e
REPO=Intel-Tiber-Broadcast-Suite

COV_CPP_VER="$(ls /opt/coverity/cxx)"
COV_OTHER_VER="$(ls /opt/coverity/other)"
COVERITY_CPP_BIN_DIR="/opt/coverity/cxx/$COV_CPP_VER/bin"
COVERITY_OTHER_BIN_DIR="/opt/coverity/other/$COV_OTHER_VER/bin"
COV_CPP_JAVA_HOM="${COVERITY_CPP_BIN_DIR}/jre/bin/java"
COV_OTHER_JAVA_HOM="${COVERITY_OTHER_BIN_DIR}/jre/bin/java"

function export_cpp_java(){ 
  export JAVA_HOME="${COV_CPP_JAVA_HOM}"
}
function export_other_java(){ 
  export JAVA_HOME="${COV_OTHER_JAVA_HOM}"
}
function log_info(){ 
   echo "[INFO]: ${1}" 
}
function log_warning(){ 
  echo "[WARNING]: $1 "
}