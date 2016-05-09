#!/bin/bash
set -e

/usr/bin/supervisord -c /lain/app/supervisord_ctl.conf
