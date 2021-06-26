package http

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/tidwall/gjson"
)

const (
	GET_TIMES_API_URL = "https://api.muuvlabs.com/anytime/clubs/2148"
)

func ReserverTime() {
	currentTime := time.Now()
	switch currentTime.Weekday() {
	case time.Monday:
		fmt.Println("Monday")
	case time.Wednesday:
		fmt.Println("Wednesday")
	case time.Friday:
		fmt.Println("Friday")
	default:
		fmt.Println("I only go to gym on Monday,Wednesday,Friday -_-")
	}
	reserveTime(time.Monday)
}

func reserveTime(day time.Weekday) {
	response, e := http.Get(GET_TIMES_API_URL)
	if e != nil {
		panic(e)
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(response.Body)
	body := buf.String()
	slots := gjson.Get(body, "time_slots").Array()
	nextDate := nextDate()
	for _, time := range slots {
		date := time.Get("date").String()
		if date == nextDate {
			fmt.Println(time.Get("start_time").String())
		}
	}
}

//convert date to the format for anytime fitness api
func nextDate() string {
	now := time.Now()
	month := strconv.Itoa(int(now.Month()))
	if len(month) == 1 {
		month = "0" + month
	}
	//just hardcoded , to remove
	day := "25" //TODO: I need to get a current day and create an appropriate month for that
	if len(day) == 1 {
		day = "0" + day
	}
	year := strconv.Itoa(now.Year())
	return year + "-" + month + "-" + day
}
