#!/bin/bash

backupdir="/opt/hotso/mongobackup"
outdir="$backupdir/dump"
archivedir="$backupdir/archive"

if [ ! -d "$outdir" ]; then
    mkdir -p "$outdir"
fi
if [ ! -d "$archivedir" ]; then
    mkdir -p "$archivedir"
fi

date=`date +%Y%m%d`
days=7 #删除七天前的备份

/usr/bin/mongodump -d hotso -o $outdir/dump_$date

tar -zcvP -f $archivedir/mongo_bak_$date.tar.gz $outdir/dump_$date