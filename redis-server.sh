#!/bin/bash
set -e
mkdir -p /redis/{data,run,log}
bkExist=`ls /redis/recover/`
#当redis/data 数据文件夹为空并且备份文件夹不为空时，触发数据恢复操作
if ([ ! -e /redis/data/*.aof ] && [ ! -e /redis/data/*.rdb ]) && [ -n -a "$bkExist" ]; then
    cd tools && ./full_recover.sh
    cd ..
fi
echo 'redis server start init'
/usr/bin/supervisord -c /lain/app/supervisord_redis.conf
