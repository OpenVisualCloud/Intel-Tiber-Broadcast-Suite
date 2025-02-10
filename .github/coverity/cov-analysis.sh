#!/bin/bash 

ROOT_DIR="$(git rev-parse --show-toplevel)"
. ${ROOT_DIR}/.github/coverity/enviroment.sh

cd ${ROOT_DIR}

tar -czvf ${COVERITY_PROJECT}.tgz cov-int

curl \
  --form token="${COVERITY_TOKEN}" \
  --form email="${COVERITY_EMAIL}" \
  --form file=@${COVERITY_PROJECT}.tgz \
  --form version="${VERSION}" \
  --form description="${DESCRIPTION}" \
  "https://scan.coverity.com/builds?project=${COVERITY_PROJECT}"
