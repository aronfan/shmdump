#!/bin/bash

SHMKEY=$1
SHMCFG=sample.xml
SHMFILE=sample.shm
SHMFILE2=sample2.shm
VERBOSE=1

echo "shmdump begin to delete shm $SHMKEY"
./shmdump -e "shmkey=$SHMKEY&op=del" 2>error.log
RC=$?
if [ $RC != 0 ]; then
    echo "shmdump failed to del shm $SHMKEY: $RC"
    exit 1
fi
echo "shmdump delete shm $SHMKEY ok"

echo "shmdump begin to load shm $SHMKEY from '$FILE'"
./shmdump -e "shmkey=$SHMKEY&op=load&cfg=sample.xml&file=$SHMFILE"
RC=$?
if [ $RC != 0 ]; then
    echo "shmdump failed to load shm $SHMKEY: $RC"
    ./shmdump -e "shmkey=$SHMKEY&op=del" 2>error.log
    exit 1
fi
echo "shmdump load shm $SHMKEY from '$FILE' ok"

echo "shmdump begin to save shm $SHMKEY"
FILE=`./shmdump -e "shmkey=$SHMKEY&op=save&file=$SHMFILE2&verbose=$VERBOSE" 2>error.log`
RC=$?
if [ $RC != 0 ]; then
    echo "shmdump failed to save shm $SHMKEY: $RC"
    exit 1
fi
echo "shmdump save shm to '$FILE' ok"

echo "compare two files: $SHMFILE vs $SHMFILE2"
cmp $SHMFILE $SHMFILE2
RC=$?
if [ $RC != 0 ]; then
    echo "shmdump load & save procedure was possible wrong!!!"
    exit 1
fi
echo "file: $SHMFILE and $SHMFILE2 are the same!"
