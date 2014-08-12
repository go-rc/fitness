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

const (
	DefaultEndDate   string = "2012-11-01"
	DefaultPassword  string = ""
	DefaultStartDate string = "2012-11-30"
	DefaultUsername  string = "tomharrison"
	EndDateFlag      string = "end"
	PasswordFlag     string = "passwd"
	StartDateFlag    string = "start"
	UsernameFlag     string = "uname"
)

type Options struct {
	EndDate   string
	Password  string
	StartDate string
	Username  string
}

func main() {
	opts := getOptions()
	dbOpts := fit.ParseMongoOptions()

	cookieJar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: cookieJar}

	login(opts.Username, opts.Password, client)

	start, end := parseDateOptions(opts)
	dietResp, dietErr := requestDietData(&start, &end, client)
	if dietErr != nil {
		panic(dietErr)
	}

	_, collection := getMongo(dbOpts)
	repo := fit.NewEntryRepository(collection)
	storeDietData(dietResp, repo)
}

func getMongo(opts *fit.MongoOptions) (*mgo.Session, *mgo.Collection) {
	session, err := mgo.Dial(opts.Host)
	if err != nil {
		panic(err)
	}
	return session, session.DB(opts.Database).C("entries")
}

func getOptions() *Options {
	opts := &Options{}

	flag.StringVar(&opts.EndDate, EndDateFlag, DefaultEndDate, "End date")
	flag.StringVar(&opts.Password, PasswordFlag, DefaultPassword, "Password")
	flag.StringVar(&opts.StartDate, StartDateFlag, DefaultStartDate, "Start date")
	flag.StringVar(&opts.Username, UsernameFlag, DefaultUsername, "Username")

	flag.Parse()
	return opts
}

func parseDateOptions(opts *Options) (time.Time, time.Time) {
	start, _ := time.Parse("2006-01-02", opts.StartDate)
	end, _ := time.Parse("2006-01-02", opts.EndDate)
	return start, end
}

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

func requestDietData(start *time.Time, end *time.Time, c *http.Client) (*http.Response, error) {
	params := make(url.Values)
	params.Set("start_Month", start.Format("01"))
	params.Set("start_Day", start.Format("02"))
	params.Set("start_Year", start.Format("2006"))
	params.Set("end_Month", end.Format("01"))
	params.Set("end_Day", end.Format("02"))
	params.Set("end_Year", end.Format("2006"))
	params.Set("ftype", "overview")
	params.Set("fltype", "csv")

	return c.PostForm("http://www.livestrong.com/thedailyplate/diary/csv/tomharrison/", params)
}

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
