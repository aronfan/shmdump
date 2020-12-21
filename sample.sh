#!/bin/bash

SHMKEY=$1
SHMFILE=sample.shm
VERBOSE=1

echo "shmdump begin to save shm $SHMKEY"
FILE=`./shmdump -e "shmkey=$SHMKEY&op=save&file=$SHMFILE&verbose=$VERBOSE" 2>/dev/null`
RC=$?
if [ $RC != 0 ]; then
    echo "shmdump failed to save shm $SHMKEY: $RC"
    exit 1
fi
echo "shmdump save shm to '$FILE' ok"

echo "shmdump begin to delete shm $SHMKEY"
./shmdump -e "shmkey=$SHMKEY&op=del" 2>/dev/null
RC=$?
if [ $RC != 0 ]; then
    echo "shmdump failed to del shm $SHMKEY: $RC"
    exit 1
fi
echo "shmdump delete shm $SHMKEY ok"

echo "shmdump begin to load shm $SHMKEY from '$FILE'"
echo "shmdump load shm $SHMKEY from '$FILE' ok"

