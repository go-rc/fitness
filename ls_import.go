// ls_import is a utility which downloads diet data from one's Livestrong diary
// and stores it in MongoDB. Livestrong unfortunately doesn't have much of a
// public API, but we can still use the net/http library to automate logging
// into one's account and using the "Export CSV" feature of the Diary page.
package main

import (
	"flag"
	"github.com/tomharrison/fitness/fit"
	"github.com/tomharrison/go-livestrong"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

// Command line parameter names and default values.
const (
	EndDateFlag      string = "end"
	PasswordFlag     string = "passwd"
	StartDateFlag    string = "start"
	UsernameFlag     string = "uname"
	DefaultEndDate   string = "2012-11-01"
	DefaultPassword  string = ""
	DefaultStartDate string = "2012-11-30"
	DefaultUsername  string = "tomharrison"
)

// Options is a struct which encapsulates the settings which may be given
// as command line arguments.
type Options struct {
	EndDate   string
	Password  string
	StartDate string
	Username  string
}

func main() {
	// Parse command line options.
	opts := getOptions()
	dbOpts := fit.ParseMongoOptions()
	start := parseDateOption(opts.StartDate)
	end := parseDateOption(opts.EndDate)

	// Initialize dependencies.
	livestrong, lsErr := livestrong.NewSiteClient(opts.Username, opts.Password)
	if lsErr != nil {
		panic(lsErr)
	}

	mongo := getMongo(dbOpts)
	entries := mongo.DB(dbOpts.Database).C("entries")

	// Download and store data.
	importDietData(&start, &end, livestrong, entries)
	importWeightData(&start, &end, livestrong, entries)
}

// Connect to MongoDB.
func getMongo(opts *fit.MongoOptions) *mgo.Session {
	session, err := mgo.Dial(opts.Host)
	if err != nil {
		panic(err)
	}
	return session
}

// getOptions parses the command line parameters and returns their values
// inside an Options struct. Any argument which is not given will be
// initialized to a default value, which is defined above as a constant.
func getOptions() *Options {
	opts := &Options{}

	flag.StringVar(&opts.EndDate, EndDateFlag, DefaultEndDate, "End date")
	flag.StringVar(&opts.Password, PasswordFlag, DefaultPassword, "Password")
	flag.StringVar(&opts.StartDate, StartDateFlag, DefaultStartDate, "Start date")
	flag.StringVar(&opts.Username, UsernameFlag, DefaultUsername, "Username")

	flag.Parse()
	return opts
}

// Parse a given date option from a string into an instance of time.Time. Dates
// must be formatted like 2006-01-02.
func parseDateOption(dateStr string) time.Time {
	d, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		panic(err)
	}
	return d
}

// storeDietData parses CSV data from the given http response's body into
// Entry records and stores them in MongoDB via the given repository.
func importDietData(start *time.Time, end *time.Time, ls *livestrong.SiteClient, entries *mgo.Collection) {
	dietRecords, err := ls.Diet.Query(start, end)
	if err != nil {
		panic(err)
	}

	for _, entry := range dietRecords {
		criteria := bson.M{"date": entry.Date}

		update := bson.M{"$set": bson.M{
			"date":              entry.Date,
			"calorie_goal":      entry.CalorieGoal,
			"calories_consumed": entry.CaloriesConsumed,
			"calories_burned":   entry.CaloriesBurned,
			"net_calories":      entry.NetCalories,
		}}

		entries.Upsert(criteria, update)
	}
}

// Parses weight JSON data out of the given response body, and stores it in
// the given MongoDB collection.
func importWeightData(start *time.Time, end *time.Time, ls *livestrong.SiteClient, entries *mgo.Collection) {
	data, err := ls.Weight.Query(start, end)
	if err != nil {
		panic(err)
	}

	for _, record := range data {
		entryDate, dateErr := time.Parse("2006-01-02", record.DisplayDate)
		if dateErr != nil {
			continue
		}

		criteria := bson.M{"date": entryDate}
		update := bson.M{"$set": bson.M{"date": entryDate, "weight": record.Weight}}
		entries.Upsert(criteria, update)
	}
}
