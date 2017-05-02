#!/usr/bin/env bash

URL='https://hooks.slack.com/services/my-hook-url'
TO="#monitoring"
EMOJI=':frowning:'
USERNAME='InfoCX Reporter'

MESSAGE="We have a problem with infoCX feeders! exit code is x"
PAYLOAD="payload={\"channel\": \"${TO}\", \"username\": \"${USERNAME}\", \"text\": \"${MESSAGE}\", \"icon_emoji\": \"${EMOJI}\"}"


ICX_SUBSCRIBER_OPTS="-Dlogback.configurationFile=./logback.xml"
export ICX_SUBSCRIBER_OPTS

STATE="file.state"

./bin/icx-subscriber $@
EXITCODE=$?
#echo ${EXITCODE}

if [[ ${EXITCODE} != 0 ]]; then
  if [ -e -f } ]; then
    echo "We have a problem with infoCX feeders! exit code is ${EXITCODE}"
  else
    MESSAGE="We have a problem with infoCX feeders! exit code is ${EXITCODE}"
    PAYLOAD="payload={\"channel\": \"${TO}\", \"username\": \"${USERNAME}\", \"text\": \"${MESSAGE}\", \"icon_emoji\": \"${EMOJI}\"}"
    curl -m 5 --data-urlencode "${PAYLOAD}" $URL
    touch ${STATE}
else
  if [ -e ${STATE} ]; then
    rm -f ${STATE}
  fi
fi
