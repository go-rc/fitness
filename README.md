Fitness Tracker in Go
=====================

I'm learning Go by building this system to track fitness statistics.

Requirements
------------

* Go 1.3
* Mongo DB

Setting up the Database
-----------------------

    $ cd db
    $ ./migrate.sh -h
    usage: ./migrate.sh options

    This script manages the MongoDB database by executing all scripts
    found at ./migrations/*.js. Note: scripts are run in alphabetical
    order, and should be idempotent.

    OPTIONS:
      -d    Name of the database to create/use.
      -h    Hostname to connect to.
      -p    Port to connect on.
      -x    Drop the database.

    $ ./migrate.sh -h localhost -p 27017 -d fitness

Building/Running the Tracker
----------------------------

    $ go get github.com/tomharrison/fitness
    $ cd $GOPATH/src/github.com/tomharrison/fitness
    $ make serve
