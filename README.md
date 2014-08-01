Fitness Tracker in Go
=====================

I'm learning Go by building this system to track fitness statistics.

Requirements
------------

* Go 1.3
* Mongo DB

Building/Running the Tracker
----------------------------

    $ go get github.com/tomharrison/fitness
    $ cd $GOPATH/src/github.com/tomharrison/fitness
    $ make db
    $ make serve

Clear out the Database
----------------------

    $ make dropdb

Advanced Use of the Database Migration Tool
-------------------------------------------

    $ cd db
    $ ./migrate.sh -h
    usage: ./migrate.sh options

    This script manages the MongoDB database by executing all scripts
    found at ./migrations/*.js. Note: scripts are run in alphabetical
    order, and should be idempotent.

    OPTIONS:
      -d    Name of the database to create/use. Defaults to fitness.
      -s    Server to connect to. Defaults to localhost.
      -p    Port to connect on. Defaults to 27017.
      -x    Drop the database. Skipped by default.

    $ ./migrate.sh -h localhost -p 27017 -d fitness

