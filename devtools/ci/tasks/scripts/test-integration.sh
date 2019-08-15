#!/bin/sh
export GOALERT_DB_URL="$DB_URL/$(cat db/NAME)"
export CYPRESS_DB_URL="$GOALERT_DB_URL"
set -x

export PATH=$PATH:$(pwd)/bin
mkdir -p logs

COMMIT=$(cat COMMIT)

trap "tar czf ../debug/debug-$COMMIT.tgz -C .. goalert" EXIT

mockslack \
  -client-id=000000000000.000000000000 \
  -client-secret=00000000000000000000000000000000 \
  -access-token=xoxp-000000000000-000000000000-000000000000-00000000000000000000000000000000 \
  -prefix=/slack \
  -single-user=bob \
  -addr=localhost:3046 >logs/mockslack.log 2>&1 &
simpleproxy -addr=localhost:3030 /slack/=http://127.0.0.1:3046 http://127.0.0.1:3042 >logs/simpleproxy.log 2>&1 &

goalert migrate
goalert --api-only --listen=:3042 --slack-base-url=http://127.0.0.1:3046/slack >logs/goalert.log 2>&1 &

if [ "$MOBILE" = "1" ]
then
  CYPRESS_viewportWidth=375 CYPRESS_viewportHeight=667 cypress run
else
  CYPRESS_viewportWidth=1440 CYPRESS_viewportHeight=900 cypress run
fi
