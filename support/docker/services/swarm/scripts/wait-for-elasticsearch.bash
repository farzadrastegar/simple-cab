#!/bin/bash
serviceName="localhost"

#expect 100% green
cmdOut=`curl -s http://${serviceName}:9200/_cat/health | grep -E "green.*100.0%" | head -n1 | awk 'END{print NR}'`
>&2 printf "waiting for elasticsearch to become available..."
while [[ ${cmdOut} -ne 1 ]]; do
   >&2 printf "."
   sleep 1
   cmdOut=`curl -s http://${serviceName}:9200/_cat/health | grep -E "green.*100.0%" | head -n1 | awk 'END{print NR}'`
done

>&2 echo
>&2 echo "elasticsearch is up"
