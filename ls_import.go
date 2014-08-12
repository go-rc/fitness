// ls_import is a utility which downloads diet data from one's Livestrong diary
// and stores it in MongoDB. Livestrong unfortunately doesn't have much of a
// public API, but we can still use the net/http library to automate logging
// into one's account and using the "Export CSV" feature of the Diary page.
package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"github.com/tomharrison/fitness/fit"
	"gopkg.in/mgo.v2"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strconv"
	"time"
)

// Names of the arguments which may be given as command line parameters.
const (
	EndDateFlag      string = "end"
	PasswordFlag     string = "passwd"
	StartDateFlag    string = "start"
	UsernameFlag     string = "uname"
)

// Default values for all of the command line parameters.
const (
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
	cookieJar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: cookieJar}
	mongo := getMongo(dbOpts)
	entries := mongo.DB(dbOpts.Database).C("entries")
	repo := fit.NewEntryRepository(entries)

	// Download and store data.
	login(opts.Username, opts.Password, client)
	dietResp := requestDietData(&start, &end, client)
	storeDietData(dietResp, repo)
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

// Log into Livestrong.com with the given username and password, via the given
// http client.
//
// Precondition: the given http client has a cookie jar.
func login(username string, password string, c *http.Client) {
	params := make(url.Values)
	params.Set("login_user", username)
	params.Set("login_password", password)

	resp, err := c.PostForm("https://www.livestrong.com/login/", params)
	defer resp.Body.Close()

	if err != nil {
		panic(err)
	}
}

// Request diet data logged between the given start and end times, via the given
// http client. The response's payload will be CSV values each containing a date,
// calorie goal, gross calories consumed, total calories burned, and a net calorie
// count (gross minus burned).
func requestDietData(start *time.Time, end *time.Time, c *http.Client) *http.Response {
	params := make(url.Values)
	params.Set("start_Month", start.Format("01"))
	params.Set("start_Day", start.Format("02"))
	params.Set("start_Year", start.Format("2006"))
	params.Set("end_Month", end.Format("01"))
	params.Set("end_Day", end.Format("02"))
	params.Set("end_Year", end.Format("2006"))
	params.Set("ftype", "overview")
	params.Set("fltype", "csv")

	resp, err := c.PostForm("http://www.livestrong.com/thedailyplate/diary/csv/tomharrison/", params)
	if err != nil {
		panic(err)
	}
	return resp
}

// storeDietData parses CSV data from the given http response's body into 
// Entry records and stores them in MongoDB via the given repository.
func storeDietData(r *http.Response, repo *fit.EntryRepository) {
	defer r.Body.Close()
	reader := csv.NewReader(r.Body)
	lineCount := 0
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error: ", err)
			return
		}

		if lineCount == 0 {
			lineCount += 1
			continue
		}

		entry := newEntryFromRecord(record)
		repo.Upsert(entry)
		lineCount += 1
	}
}

// newEntryFromRecord is a factory function which intializes the Entry type
// using an array of strings which were parsed from a CSV record.
func newEntryFromRecord(r []string) *fit.Entry {
	e := &fit.Entry{}
	reg := regexp.MustCompile("[a-z]{2},")
	r[0] = reg.ReplaceAllString(r[0], ",")
	e.Date, _ = time.Parse("January _2, 2006", r[0])
	e.CalorieGoal, _ = strconv.ParseFloat(r[1], 64)
	e.CaloriesConsumed, _ = strconv.ParseFloat(r[2], 64)
	e.CaloriesBurned, _ = strconv.ParseFloat(r[3], 64)
	e.NetCalories, _ = strconv.ParseFloat(r[4], 64)
	return e
}
