package main

import (
	"flag"
	"github.com/gorilla/mux"
	"github.com/tomharrison/fitness/api"
	"gopkg.in/mgo.v2"
	"net/http"
)

const (
	DefaultApiPrefix string = "/api"
	DefaultDatabase  string = "fitness"
	DefaultMongoHost string = "127.0.0.1"
	DefaultPort      string = "3000"
	DefaultPubDir    string = "./pub"
	ApiPrefixFlag    string = "apiPrefix"
	DatabaseFlag     string = "db"
	MongoHostFlag    string = "dbHost"
	PubDirFlag       string = "pubDir"
	PortFlag         string = "port"
)

type ServerOptions struct {
	ApiPrefix string
	Database  string
	MongoHost string
	PubDir    string
	Port      string
}

func main() {
	options := GetOptions()
	CreateServer(options)
}

func GetOptions() *ServerOptions {
	opts := &ServerOptions{}
	flag.StringVar(&opts.ApiPrefix, ApiPrefixFlag, DefaultApiPrefix, "API path prefix")
	flag.StringVar(&opts.Database, DatabaseFlag, DefaultDatabase, "Database name")
	flag.StringVar(&opts.MongoHost, MongoHostFlag, DefaultMongoHost, "Mongo DB hostname")
	flag.StringVar(&opts.PubDir, PubDirFlag, DefaultPubDir, "Path to the publicly accessible web root")
	flag.StringVar(&opts.Port, PortFlag, DefaultPort, "Port on which to listen for connections")
	flag.Parse()
	return opts
}

func CreateServer(opts *ServerOptions) {
	session, err := mgo.Dial(opts.MongoHost)
	if err != nil {
		panic(err)
	}
	db := session.DB(opts.Database)

	router := mux.NewRouter()
	router.Handle("/", http.FileServer(http.Dir(opts.PubDir)))

	api.NewFitnessApi(db, router.PathPrefix(opts.ApiPrefix).Subrouter())

	http.Handle("/", router)
	http.ListenAndServe(":"+opts.Port, nil)
}
