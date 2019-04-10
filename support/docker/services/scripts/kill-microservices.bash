#!/bin/bash
homeDir=`pwd`

#kill gateway, zombie_driver, driver_location
>&2 printf "killing microservices..."
cmdOut=`cat ${homeDir}/run.pid | awk 'END{print NR}'`
if [[ ${cmdOut} -ne 0 ]]; then
   `nohup pkill -9 -s $(cat ${homeDir}/run.pid | head -n1) >/dev/null 2>&1 &`
   >&2 echo "done" 
else
   >&2 echo "no microservice found" 
fi

