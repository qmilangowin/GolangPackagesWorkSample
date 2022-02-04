#!/bin/sh
#Just a dummy docker-entrypoint file, since original contains propriatary. 
ls -la
# echo "Sleeping"
# sleep 10s
# echo "Done Sleeping"
for i in 1 2 3 4 5
do 
    echo "Looping ... number $i"
    sleep 1s
    echo "Nothing to see here ...."
done
echo ${BUCKET}
echo {}
ls -la
