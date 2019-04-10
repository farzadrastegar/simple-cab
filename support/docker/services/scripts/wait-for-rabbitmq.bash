#!/bin/bash
serviceName="localhost"

#expect a 'server' header
cmdOut=`curl -s -i ${serviceName}:15672 | grep "server" | head -n1 | awk 'END{print NR}'`
>&2 printf "waiting for rabbitmq to become available..."
while [[ ${cmdOut} -ne 1 ]]; do
   >&2 printf "."
   sleep 1
   cmdOut=`curl -s -i ${serviceName}:15672 | grep "server" | head -n1 | awk 'END{print NR}'`
done

>&2 echo
>&2 echo "rabbitmq is up"
