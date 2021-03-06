package controllers

var preStopScript = `
#!/bin/bash
# This works for worker only. Coordinator doesn't support graceful shutdown.
# This script will block until the server has actually shutdown.
set -x

http_port="$(cat /usr/lib/presto/etc/config.properties | grep 'http-server.http.port' | sed 's/^.*=\(.*\)$/\1/')"
https_port="$(cat /usr/lib/presto/etc/config.properties | grep 'http-server.https.port' | sed 's/^.*=\(.*\)$/\1/')"

if [ -n "$http_port" ] ; then
    res=$(curl -s -o /dev/null -w "%{http_code}"  -XPUT --data '"SHUTTING_DOWN"' -H "Content-type: application/json" http://localhost:${http_port}/v1/info/state)
fi

if [ -z "$res" -o "$res" != "200" ] && [ -n "$https_port" ]; then
    res=$(curl -k -s -o /dev/null -w "%{http_code}"  -XPUT --data '"SHUTTING_DOWN"' -H "Content-type: application/json" https://localhost:${https_port}/v1/info/state)
fi

if [ -z "$res" -o "$res" != "200" ] ; then
  # Failed to send the shutdown request.
  exit -1
else
  # Server is shutting down. Block until the server is actually down.
  while curl http://localhost:${http_port}/v1/info/state; do
    sleep 1
  done
fi
`
