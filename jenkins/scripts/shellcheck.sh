#! /bin/bash

DIRS_TO_SCAN=( "./" 
               "./doc"
               "./patches"
               "./pipelines"
)
OPTIONS="-s bash"
OUT_FILE="shellcheck_output.log"
rm -rf shellcheck_logs
mkdir -p shellcheck_logs


for _DIR in "${DIRS_TO_SCAN[@]}"
do  
    echo "*********** SCANNING ${_DIR} DIRECTORY... ***********" | tee -a shellcheck_logs/${OUT_FILE}
    BASH_FILES=$(find ${_DIR} -maxdepth 1  -name "*.sh")
    if [ -n "${BASH_FILES}" ]; then shellcheck ${OPTIONS} ${BASH_FILES} >>  shellcheck_logs/${OUT_FILE}; fi
    
done
echo "*********** SCANNING DONE ***********" >> shellcheck_logs/${OUT_FILE}
