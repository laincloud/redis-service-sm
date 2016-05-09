#!/bin/bash
set -xe

sleep 5

/usr/bin/supervisord -c /lain/app/supervisord_proxy.conf
