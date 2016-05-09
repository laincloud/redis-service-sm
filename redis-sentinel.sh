#!/bin/bash
set -xe
mkdir -p /redis/sentinel
sleep 5

echo 'redis sentinel start init'
/usr/bin/supervisord -c /lain/app/supervisord_sentinel.conf
