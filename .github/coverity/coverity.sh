#!/bin/bash -e

ROOT_DIR="$(git rev-parse --show-toplevel)"
. ${ROOT_DIR}/.github/coverity/enviroment.sh


log_info "Installing Coverity Scan for CPP "
export LANGUAGE="cxx"
.install_coverity.sh  
log_info "Installing Coverity Scan for GO "
export LANGUAGE="other"
.install_coverity.sh


log_info "Cloning the repository"
git clone https://github.com/OpenVisualCloud/${REPO}.git
cd ${REPO}
git submodule update --init --recursive




log_info "Building the project"
. cov_build.sh all