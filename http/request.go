package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/tidwall/gjson"
)

const (
	GET_TIMES_API_URL    = "https://api.muuvlabs.com/anytime/clubs/2148"
	RESERVE_TIME_API_URL = "https://api.muuvlabs.com/anytime/reservations"
)

type reservationData struct {
	clubIdent int
	startInt  int64
}

func ReserverTime() {
	currentDay := time.Now().Weekday()
	if currentDay != time.Monday &&
		currentDay != time.Wednesday &&
		currentDay != time.Saturday {
		fmt.Println("The reservation days for me are Monday,Wednesday,Saturday")
	} else {
		reservationData := reservationTime(time.Now().Weekday())
		sendReservationRequest(reservationData)
	}
}

func sendReservationRequest(data *reservationData) {
	postData, _ := json.Marshal(map[string]interface{}{
		"club_ident": data.clubIdent,
		"start_int":  data.startInt,
	})
	requestBody := bytes.NewBuffer(postData)
	req, err := http.NewRequest("POST", RESERVE_TIME_API_URL, requestBody)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Cookie", "session=eyJ1c2VyIjoiNjA3MGU5OWRkYjhhZmYwMDA2YmJlZTg0In0=; session.sig=Mh7UlBMpHEGcGn4orX9fp9y3Ta8")
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	fmt.Println("response Status:", resp.Status)
}

func reservationTime(day time.Weekday) *reservationData {
	response, e := http.Get(GET_TIMES_API_URL)
	if e != nil {
		panic(e)
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(response.Body)
	body := buf.String()
	slots := gjson.Get(body, "time_slots").Array()
	nextDate := nextDate(day)
	fmt.Printf("Next date is %s\n", nextDate)
	for _, timeSlot := range slots {
		date := timeSlot.Get("date").String()
		//if it's Saturday then book for Monday on 9 A.M
		if date == nextDate && day == time.Saturday && timeSlot.Get("start_time").String() == "9:00 AM" {
			clubIdent, _ := strconv.ParseInt(timeSlot.Get("club_ident").String(), 10, 32)
			startInt, _ := strconv.ParseInt(timeSlot.Get("start_int").String(), 10, 64)
			return &reservationData{clubIdent: int(clubIdent), startInt: startInt}
		}
		//TODO MONDAY,WEDNESDAY
	}
	return nil
}

//convert date to the format for anytime fitness api
func nextDate(currentDay time.Weekday) string {
	currentTime := time.Now().AddDate(0, 0, 2)
	month := strconv.Itoa(int(currentTime.Month()))
	if len(month) == 1 {
		month = "0" + month
	}
	day := strconv.Itoa(currentTime.Day())
	if len(day) == 1 {
		day = "0" + day
	}
	year := strconv.Itoa(currentTime.Year())
	return year + "-" + month + "-" + day
}
