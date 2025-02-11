#!/bin/bash

# Config file processing procedures
sed_escape() {
  sed -e 's/[]\/$*.^[]/\\&/g'
}

if [ $# -gt 0 ]; then
   if [ -f "$1" ]; then
       echo "Reading configuration from $1"
       export node_json="$1"
   else
       exec "$@"
   fi
fi

  echo -e "\nStarting NMOS Node with following config"
  cat $node_json
  /home/nmos-cpp-node $node_json
  ret=$?

exit $ret  # Make sure we really exit
