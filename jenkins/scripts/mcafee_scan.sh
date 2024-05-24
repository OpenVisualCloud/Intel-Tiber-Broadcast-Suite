#!/bin/bash -e
mkdir -p Malware
touch Malware/malware_clean_reports
MALWARE_PASS=true
for tool in mcafee; do
  TOOL_PASS=true
  abi virus_scan --target Dockerfile --logtarget Malware --tool $tool || TOOL_PASS=false
  if [ "$TOOL_PASS" = "true" ]; then 
    if [ "$tool" = "mcafee" ]; then
      echo "Malware/images_McAfee.html" >> Malware/malware_clean_reports
      echo "Malware/images_McAfee.log" >> Malware/malware_clean_reports
    elif [ "$tool" = "clamav" ]; then
      echo "Malware/images_ClamAV.log" >> Malware/malware_clean_reports
    fi
  else
    MALWARE_PASS=false
  fi
done
chmod -R a+rw Malware
echo $MALWARE_PASS>Malware/malware_pass
