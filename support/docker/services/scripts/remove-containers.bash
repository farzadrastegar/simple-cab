#!/bin/bash
serviceName="localhost"

#remove existing containers
>&2 printf "removing all containers..."
cmdOut=`docker ps -a -q | awk 'END{print NR}'`
if [[ ${cmdOut} -ne 0 ]]; then
   `docker rm -f $(docker ps -a -q)`
   >&2 echo "done" 
else
   >&2 echo "no container found" 
fi

