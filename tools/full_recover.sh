#!/bin/bash
set -e

#恢复数据
bk_dir='/redis/recover/'
bkfile=`ls "$bk_dir" | awk '{print $(NF)}' | grep 'tar.gz' | sort -r | head -1`
if [ ! -z "$bkfile" ]; then
bkdir_name=$(echo $bkfile | awk -F'.' '{print $1}')
prepare_dir="$bk_dir/$bkdir_name"
mkdir $prepare_dir
tar -izxvf $bk_dir/$bkfile -C $prepare_dir
\cp -rf $prepare_dir/* /
rm -rf $bk_dir/*
fi
