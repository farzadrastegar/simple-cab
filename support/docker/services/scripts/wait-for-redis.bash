#!/bin/bash
serviceName="redis"

#expect output from the command below
cmdOut=`docker ps | grep ${serviceName} | tail -n1 | cut -d' ' -f1 | awk '{print "docker exec -it " $1 " redis-cli --version"}' | awk 'END{print NR}'`
>&2 printf "waiting for redis to become available..."
while [[ ${cmdOut} -ne 1 ]]; do
   >&2 printf "."
   sleep 1
   cmdOut=`docker ps | grep ${serviceName} | tail -n1 | cut -d' ' -f1 | awk '{print "docker exec -it " $1 " redis-cli --version"}' | awk 'END{print NR}'`
done

#expect redis version starting with 'redis-cli'
cmd=`docker ps | grep ${serviceName} | tail -n1 | cut -d' ' -f1 | awk '{print "docker exec -it " $1 " redis-cli --version"}'`
cmdOut=`${cmd} | grep ^redis-cli | wc -l`
while [[ ${cmdOut} -ne 1 ]]; do
   >&2 printf "."
   sleep 1
   cmdOut=`${cmd} | grep ^redis-cli`
done

>&2 echo
>&2 echo "redis is up"
