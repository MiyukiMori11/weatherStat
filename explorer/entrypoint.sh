#!/usr/bin/env bash

echoerr() { echo "$@" 1>&2; }

export CNTDAYARCHIVE="'${CNTDAYARCHIVE} days'"


envsubst < /usr/src/app/migrations/000006_add_ttl_trigger.up.sql > /usr/src/app/migrations/000006_add_ttl_trigger.up.sql

app migrate

if [ $? -ne 0 ]; then
      echoerr "Can't put interval into script";
      exit 1;
fi
