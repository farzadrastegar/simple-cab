#!/bin/bash
serviceName="localhost"

>&2 printf "waiting for kibana to become available..."

#expect huge output from the curl command
cmdOut=`curl -s http://localhost:5601/app/kibana | awk 'END{print NR}'`
while [[ ${cmdOut} -lt 10 ]]; do
   >&2 printf "."
   sleep 1
   cmdOut=`curl -s http://localhost:5601/app/kibana | awk 'END{print NR}'`
done

>&2 echo
>&2 echo "kibana is up"
