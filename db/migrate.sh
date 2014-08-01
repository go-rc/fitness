#!/usr/bin/env bash

usage()
{
	cat << EOF
usage: $0 options

This script manages the MongoDB database by executing all scripts
found at ./migrations/*.js. Note: scripts are run in alphabetical
order, and should be idempotent.

OPTIONS:
  -d    Name of the database to create/use.
  -h    Hostname to connect to.
  -p    Port to connect on.
  -x    Drop the database.
EOF
}

DB="fitness"
HOST="localhost"
PORT="27017"
DROP=0

while getopts "hd:h:p:x" opt; do
	case "$opt" in
		h)
			usage
			exit 1
			;;
		d)
			DB="$OPTARG"
			;;
		h)
			HOST="$OPTARG"
			;;
		p)
			PORT="$OPTARG"
			;;
		x)
			DROP=1
			;;
	esac
done

DBADDRESS="$HOST:$PORT/$DB"

if [ $DROP -eq 1 ]
then
	echo "Dropping database: $DB"
	mongo $DBADDRESS --eval "db.dropDatabase()"
	exit 0
fi

SCRIPTS=./migrations/*.js
for s in $SCRIPTS
do
	echo "Executing $s"
	mongo $DBADDRESS $s
done
