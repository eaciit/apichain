#!/bin/bash

# Start db
mongod &

# wait 3 sec
sleep 3

# check for db exist
COLCOUNT=`mongo scbapichain --quiet --eval "printjson(db.getCollectionNames().length)"`

# restore on empty collection
if [ "$COLCOUNT" == "0" ]; then
  echo "Restore default database"
  mongorestore /srv/dump
fi

# Start app
pushd /srv/app/src/eaciit/apichain/webapp
./webapp &
popd

# Naive check runs checks once a minute to see if either of the processes exited.
# This illustrates part of the heavy lifting you need to do if you want to run
# more than one service in a container. The container will exit with an error
# if it detects that either of the processes has exited.
# Otherwise it will loop forever, waking up every 60 seconds

while /bin/true; do
  ps aux | grep mongod | grep -q -v grep
  PROCESS_1_STATUS=$?
  ps aux | grep webapp | grep -q -v grep
  PROCESS_2_STATUS=$?
  # If the greps above find anything, they will exit with 0 status
  # If they are not both 0, then something is wrong
  if [ $PROCESS_1_STATUS -ne 0 -o $PROCESS_2_STATUS -ne 0 ]; then
    echo "Process exit, shutting down..."
    exit -1
  fi
  sleep 5
done

