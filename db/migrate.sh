#!/usr/bin/env bash

usage()
{
	cat << EOF
usage: $0 options

This script manages the MongoDB database by executing all scripts
found at ./migrations/*.js. Note: scripts are run in alphabetical
order, and should be idempotent.

OPTIONS:
  -d    Name of the database to create/use. Defaults to fitness.
  -s    Server to connect to. Defaults to localhost.
  -p    Port to connect on. Defaults to 27017.
  -x    Drop the database. Skipped by default.
EOF
}

DB="fitness"
SERVER="localhost"
PORT="27017"
DROP=0

while getopts "hd:s:p:x" opt; do
	case "$opt" in
		h)
			usage
			exit 1
			;;
		d)
			DB="$OPTARG"
			;;
		s)
			SERVER="$OPTARG"
			;;
		p)
			PORT="$OPTARG"
			;;
		x)
			DROP=1
			;;
	esac
done

DBADDRESS="$SERVER:$PORT/$DB"

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
