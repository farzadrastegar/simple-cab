#!/bin/bash
serviceName="localhost"

#remove existing containers
>&2 printf "removing dangling images..."
cmdOut=`docker images -q -f dangling=true | awk 'END{print NR}'`
if [[ ${cmdOut} -ne 0 ]]; then
   `nohup docker rmi $(docker images -q -f dangling=true) >/dev/null 2>&1 &`
   >&2 echo "done" 
else
   >&2 echo "no dangling image found" 
fi

