package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/tidwall/gjson"
)

const (
	GET_TIMES_API_URL       = "https://api.muuvlabs.com/anytime/clubs/2148"
	RESERVE_TIME_API_URL    = "https://api.muuvlabs.com/anytime/reservations"
	GET_RESERVATION_API_URL = "https://api.muuvlabs.com/anytime/reservations?include_club=0&per_page=2&ran=1624740618984"
	PREFERRED_TIME          = "9:00 AM"
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
		day := time.Now().Weekday()
		nextDate := nextDate(day)
		fmt.Printf("Next booking date is %s\n", nextDate)
		if alreadyCreated(nextDate) {
			fmt.Printf("Reservation for %s already exists\n", nextDate)
		} else {
			reservationData := reservationTime(nextDate, day)
			if reservationData == nil {
				fmt.Printf(
					"There is no available time slot at %s on %s \n",
					PREFERRED_TIME,
					nextDate,
				)
			} else {
				fmt.Printf("Reservation made at %v", reservationData)
			}
			//	sendReservationRequest(reservationData)
		}
	}
}

//check if reservation already exists
func alreadyCreated(nextDate string) bool {
	req, err := http.NewRequest("GET", GET_RESERVATION_API_URL, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Cookie", "session=eyJ1c2VyIjoiNjA3MGU5OWRkYjhhZmYwMDA2YmJlZTg0In0=; session.sig=Mh7UlBMpHEGcGn4orX9fp9y3Ta8")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	reservations := gjson.Get(string(data), "results").Array()
	if reservations == nil {
		return false
	}
	return reservations[0].Get("date").String() == nextDate
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
	//TODO:move it to env variables
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

func reservationTime(nextDate string, day time.Weekday) *reservationData {
	response, e := http.Get(GET_TIMES_API_URL)
	if e != nil {
		panic(e)
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(response.Body)
	body := buf.String()
	slots := gjson.Get(body, "time_slots").Array()
	for _, timeSlot := range slots {
		date := timeSlot.Get("date").String()
		//My preferred time is 9:00 A.M because it's a hot summer and and it's not good to go out later
		if date == nextDate && timeSlot.Get("start_time").String() == PREFERRED_TIME {
			clubIdent, _ := strconv.ParseInt(timeSlot.Get("club_ident").String(), 10, 32)
			startInt, _ := strconv.ParseInt(timeSlot.Get("start_int").String(), 10, 64)
			return &reservationData{clubIdent: int(clubIdent), startInt: startInt}
		}
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
