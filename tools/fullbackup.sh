#!/bin/sh
set -e
back_path="/redis/data_fullbackup/"
if [ ! -d "$back_path" ]; then 
	mkdir $back_path
fi
backup_time=`date +'%Y-%m-%d-%H-%M-%S'`
tar cvzf "/redis/data_fullbackup/full_$backup_time"_bak.tar.gz /redis/data --exclude "*.aof"



