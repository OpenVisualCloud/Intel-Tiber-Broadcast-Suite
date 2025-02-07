#!/bin/bash 

ROOT_DIR="$(git rev-parse --show-toplevel)"
. ${ROOT_DIR}/.github/coverity/enviroment.sh

cd ${ROOT_DIR}

tar -czvf ${COVERITY_PROJECT}.tgz cov-int

curl \
  --form token="${COVERITY_TOKEN}" \
  --form email="eee.ddd@iii.com" \
  --form file=@${COVERITY_PROJECT}.tgz \
  --form version="2024.6.1" \
  --form description="total" \
  "https://scan.coverity.com/builds?project=${COVERITY_PROJECT}"