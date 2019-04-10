#!/bin/bash
serviceName="localhost"
port=$1

#expect the input port to be listening
cmdOut=`lsof -i -P | grep -E ":${port} .*LISTEN" | head -n1 | awk 'END{print NR}'`
>&2 printf "waiting for logstash to become available..."
while [[ ${cmdOut} -ne 1 ]]; do
   >&2 printf "."
   sleep 1
   cmdOut=`lsof -i -P | grep -E ":${port} .*LISTEN" | head -n1 | awk 'END{print NR}'`
done

>&2 echo
>&2 echo "logstash is up"
