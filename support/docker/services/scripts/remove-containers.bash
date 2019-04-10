#!/bin/bash
serviceName="localhost"

#remove existing containers
>&2 printf "removing all containers..."
cmdOut=`docker ps -a -q | awk 'END{print NR}'`
if [[ ${cmdOut} -ne 0 ]]; then
   `nohup docker rm -f $(docker ps -a -q) >/dev/null 2>&1 &`
   >&2 echo "done" 
else
   >&2 echo "no container found" 
fi

