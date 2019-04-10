#!/bin/bash
serviceName="localhost"
appName=$1

#expect output from the command below
cmdOut=`curl -s http://${serviceName}:8888/${appName}/dev/master | awk 'END{print NR}'`
>&2 printf "waiting for configserver to become available..."
while [[ ${cmdOut} -eq 0 ]]; do
   >&2 printf "."
   sleep 1
   cmdOut=`curl -s http://${serviceName}:8888/${appName}/dev/master | awk 'END{print NR}'`
done

#expect a 'name' key
cmdOut=`curl -s http://${serviceName}:8888/${appName}/dev/master | python -c "import sys, json; print json.load(sys.stdin)['name']" | grep ^${appName}$| awk 'END{print NR}'`
while [[ ${cmdOut} -ne 1 ]]; do
   >&2 printf "."
   sleep 1
   cmdOut=`curl -s http://${serviceName}:8888/${appName}/dev/master | python -c "import sys, json; print json.load(sys.stdin)['name']" | grep ^${appName}$| awk 'END{print NR}'`
done

>&2 echo
>&2 echo "configserver for ${appName} is up"
