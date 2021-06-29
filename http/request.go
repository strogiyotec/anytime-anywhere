package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/tidwall/gjson"
)

const (
	GetTimesApiUrl       = "https://api.muuvlabs.com/anytime/clubs/2148"
	ReserveTimeApiUrl    = "https://api.muuvlabs.com/anytime/reservations"
	GetReservationApiUrl = "https://api.muuvlabs.com/anytime/reservations?include_club=0&per_page=2&ran=1624740618984"
	PreferredTime        = "9:00 AM"
)

type reservationData struct {
	clubIdent int   // club id
	startInt  int64 //start time
}

func ReserveTime() {
	currentDay := time.Now().Weekday()
	if currentDay != time.Monday &&
		currentDay != time.Wednesday &&
		currentDay != time.Saturday {
		fmt.Println("The reservation days for me are Monday,Wednesday,Saturday")
	} else {
		nextDate := nextDate()
		fmt.Printf("Next booking date is %s\n", nextDate)
		if alreadyExists(nextDate) {
			fmt.Printf("Reservation for %s already exists\n", nextDate)
		} else {
			reservationData := reservationTime(nextDate)
			if reservationData == nil {
				fmt.Printf(
					"There is no available time slot at %s on %s \n",
					PreferredTime,
					nextDate,
				)
			} else {
				sendReservationRequest(reservationData)
				//check if it really was saved
				if alreadyExists(nextDate) {
					fmt.Printf("Reservation made at %v", reservationData)
				} else {
					fmt.Println("anytime tried to make a reservation, but it failed")
				}
			}
		}
	}
}

//check if reservation already exists
func alreadyExists(nextDate string) bool {
	req, err := http.NewRequest("GET", GetReservationApiUrl, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Cookie", userCookies())
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

//send an actual request to make a reservation
func sendReservationRequest(payload *reservationData) {
	postData, _ := json.Marshal(map[string]interface{}{
		"club_ident": payload.clubIdent,
		"start_int":  payload.startInt,
	})
	requestBody := bytes.NewBuffer(postData)
	req, err := http.NewRequest("POST", ReserveTimeApiUrl, requestBody)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Cookie", userCookies())
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	fmt.Println("response Status:", resp.Status)
}

//return a payload to create a reservation
func reservationTime(nextDate string) *reservationData {
	response, e := http.Get(GetTimesApiUrl)
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
		if date == nextDate && timeSlot.Get("start_time").String() == PreferredTime {
			clubIdent, _ := strconv.ParseInt(timeSlot.Get("club_ident").String(), 10, 32)
			startInt, _ := strconv.ParseInt(timeSlot.Get("start_int").String(), 10, 64)
			return &reservationData{clubIdent: int(clubIdent), startInt: startInt}
		}
	}
	return nil
}

func userCookies() string {
	token := os.Getenv("ANYTIME_SESSION_TOKEN")
	sig := os.Getenv("ANYTIME_SESSION_SIG")
	return fmt.Sprintf("session=%s; session.sig=%s", token, sig)
}

//convert date to the format for anytime fitness api
func nextDate() string {
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
