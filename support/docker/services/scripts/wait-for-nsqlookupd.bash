#!/bin/bash
serviceName="localhost"

#expect a 'status_txt' in output json
cmdOut=`curl -s ${serviceName}:4161 | awk 'END{print NR}'`
>&2 printf "waiting for nsqlookupd to become available..."
while [[ ${cmdOut} -eq 0 ]]; do
   >&2 printf "."
   sleep 1
   cmdOut=`curl -s ${serviceName}:4161 | awk 'END{print NR}'`
done

#expect a 'status_txt' in output json
cmdOut=`curl -s ${serviceName}:4161/lookup?topic=* | python -c "import sys, json; print json.load(sys.stdin)['status_txt']" | awk 'END{print NR}'`
while [[ ${cmdOut} -ne 1 ]]; do
   >&2 printf "."
   sleep 1
   cmdOut=`curl -s ${serviceName}:4161/lookup?topic=* | python -c "import sys, json; print json.load(sys.stdin)['status_txt']" | awk 'END{print NR}'`
done

>&2 echo
>&2 echo "nsqlookupd is up"

# ececute everything in the arguments
exec $@
