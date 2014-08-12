package fit

import "flag"

const (
	DefaultDatabase  string = "fitness"
	DefaultMongoHost string = "127.0.0.1"
	DatabaseFlag     string = "db"
	MongoHostFlag    string = "dbHost"
)

type MongoOptions struct {
	Database string
	Host     string
}

func ParseMongoOptions() *MongoOptions {
	opts := &MongoOptions{}
	flag.StringVar(&opts.Database, DatabaseFlag, DefaultDatabase, "Database name")
	flag.StringVar(&opts.Host, MongoHostFlag, DefaultMongoHost, "Mongo DB hostname")
	flag.Parse()
	return opts
}
