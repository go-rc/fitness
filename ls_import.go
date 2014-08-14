// ls_import is a utility which downloads diet data from one's Livestrong diary
// and stores it in MongoDB. Livestrong unfortunately doesn't have much of a
// public API, but we can still use the net/http library to automate logging
// into one's account and using the "Export CSV" feature of the Diary page.
package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/tomharrison/fitness/fit"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strconv"
	"time"
)

// Names of the arguments which may be given as command line parameters.
const (
	EndDateFlag   string = "end"
	PasswordFlag  string = "passwd"
	StartDateFlag string = "start"
	UsernameFlag  string = "uname"
)

// Default values for all of the command line parameters.
const (
	DefaultEndDate   string = "2012-11-01"
	DefaultPassword  string = ""
	DefaultStartDate string = "2012-11-30"
	DefaultUsername  string = "tomharrison"
)

// Parameters and URL for the login form.
const (
	LoginFormEndpoint          string = "https://www.livestrong.com/login/"
	LoginFormPasswordParameter string = "login_password"
	LoginFormUsernameParameter string = "login_username"
)

// Parameters and URL for the CSV export.
const (
	CsvExportBaseUrl       string = "http://www.livestrong.com/thedailyplate/diary/csv/"
	CsvStartMonthParameter string = "start_Month"
	CsvStartDayParameter   string = "start_Day"
	CsvStartYearParameter  string = "start_Year"
	CsvEndMonthParameter   string = "end_Month"
	CsvEndDayParameter     string = "end_Day"
	CsvEndYearParameter    string = "end_Year"
	CsvFileTypeParameter   string = "ftype"
	CsvFormatParameter     string = "fltype"
)

// Parameters and URL for the weight graph/json page.
const (
	WeightFromMonthParameter    string = "from_Month"
	WeightFromDayParameter      string = "from_Day"
	WeightFromYearParameter     string = "from_Year"
	WeightToMonthParameter      string = "to_Month"
	WeightToDayParameter        string = "to_Day"
	WeightToYearParameter       string = "to_Year"
	WeightShowCaloriesParameter string = "show_net_cals_plz"
	WeightRefreshParameter      string = "refresh"
	WeightGraphEndpoint         string = "http://www.livestrong.com/thedailyplate/users/weight/"
)

type WeightRecord struct {
	Weight      float64 `json:"weight,omitifempty"`
	DisplayDate string  `json:"datestamp,omitifempty"`
}

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

	// Download and store data.
	login(opts.Username, opts.Password, client)

	dietResp := requestDietData(&start, &end, opts.Username, client)
	storeDietData(dietResp, entries)

	weightResp := requestWeightData(&start, &end, client)
	storeWeightData(weightResp, entries)
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
	params.Set(LoginFormUsernameParameter, username)
	params.Set(LoginFormPasswordParameter, password)

	resp, err := c.PostForm(LoginFormEndpoint, params)
	defer resp.Body.Close()

	if err != nil {
		panic(err)
	}
}

// Request diet data logged between the given start and end times, via the given
// http client. The response's payload will be CSV values each containing a date,
// calorie goal, gross calories consumed, total calories burned, and a net calorie
// count (gross minus burned).
func requestDietData(start *time.Time, end *time.Time, username string, c *http.Client) *http.Response {
	params := make(url.Values)
	params.Set(CsvStartMonthParameter, start.Format("01"))
	params.Set(CsvStartDayParameter, start.Format("02"))
	params.Set(CsvStartYearParameter, start.Format("2006"))
	params.Set(CsvEndMonthParameter, end.Format("01"))
	params.Set(CsvEndDayParameter, end.Format("02"))
	params.Set(CsvEndYearParameter, end.Format("2006"))
	params.Set(CsvFileTypeParameter, "overview")
	params.Set(CsvFormatParameter, "csv")

	endpoint := CsvExportBaseUrl + username
	resp, err := c.PostForm(endpoint, params)
	if err != nil {
		panic(err)
	}
	return resp
}

// Retrieve weight history by requesting the weight graph page, which has the desired
// data in the response body.
func requestWeightData(start *time.Time, end *time.Time, c *http.Client) *http.Response {
	params := make(url.Values)
	params.Set(WeightFromMonthParameter, start.Format("01"))
	params.Set(WeightFromDayParameter, start.Format("02"))
	params.Set(WeightFromYearParameter, start.Format("2006"))
	params.Set(WeightToMonthParameter, end.Format("01"))
	params.Set(WeightToDayParameter, end.Format("02"))
	params.Set(WeightToYearParameter, end.Format("2006"))
	params.Set(WeightShowCaloriesParameter, "")
	params.Set(WeightRefreshParameter, "Refresh")

	resp, err := c.PostForm(WeightGraphEndpoint, params)
	if err != nil {
		panic(err)
	}
	return resp
}

// storeDietData parses CSV data from the given http response's body into
// Entry records and stores them in MongoDB via the given repository.
func storeDietData(r *http.Response, entries *mgo.Collection) {
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
		update := bson.M{"$set": bson.M{
			"date":              entry.Date,
			"calorie_goal":      entry.CalorieGoal,
			"calories_consumed": entry.CaloriesConsumed,
			"calories_burned":   entry.CaloriesBurned,
			"net_calories":      entry.NetCalories,
		}}
		entries.Upsert(bson.M{"date": entry.Date}, update)
		lineCount += 1
	}
}

// Parses weight JSON data out of the given response body, and stores it in
// the given MongoDB collection.
func storeWeightData(r *http.Response, entries *mgo.Collection) {
	defer r.Body.Close()
	body, bodyErr := ioutil.ReadAll(r.Body)
	if bodyErr != nil {
		panic(bodyErr)
	}

	reg := regexp.MustCompile("var weight_data = (\\[[^\\[]+\\])")
	match := reg.FindSubmatch(body)
	if match == nil || len(match) != 2 {
		return
	}

	data := make([]WeightRecord, 0)
	weightErr := json.Unmarshal(match[1], &data)
	if weightErr != nil {
		panic(weightErr)
	}

	for _, record := range data {
		entry := newEntryFromWeightRecord(&record)
		update := bson.M{"$set": bson.M{
			"weight": entry.Weight,
			"date":   entry.Date,
		}}
		entries.Upsert(bson.M{"date": entry.Date}, update)
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

func newEntryFromWeightRecord(r *WeightRecord) *fit.Entry {
	e := &fit.Entry{Weight: r.Weight}
	e.Date, _ = time.Parse("2006-01-02", r.DisplayDate)
	return e
}
