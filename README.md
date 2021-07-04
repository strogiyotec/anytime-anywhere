# anytime-anywhere
Script to automate the booking process for **AnytimeFitness Gym**


## Status
From the First July Gym Booking were abandoned. Everything goes to normal so no need in automation booking anymore 

## Motivation
I got tired from making my gym appoints (COVID is still a thing). This cli will do it for me


## How to use ? 

I am too lazy to create authorization logic(maybe will implement later on) ,for now go to your gym page and copy session token and sig from cookies
![where to get cookies from](https://raw.githubusercontent.com/strogiyotec/anytime-anywhere/master/images/cookies.png)

When you got these two cookies follow the instruction below

1. Clone this repository `git clone https://github.com/strogiyotec/anytime-anywhere.git`
2. `cd anytime-anywhere`
3. Set 2 environmental variables (example for bash and zsh)
    - `export ANYTIME_SESSION_TOKEN="YOUR_TOKEN"`
    - `export ANYTIME_SESSION_SIG="YOU_SIG"`
4. Run `go run main.go`
5. By default it will try to make an appoint for MONDAY,WEDNESDAY,FRIDAY(my gym days) at 9:00 A.M
6. As **AnytimeFitness** only allows booking two days in advance , it only possible to book one appointment in one day(appointment for Wednesday will be available on Monday and so on). So you can make a cron job that will run on days you are interested in

## TODO

1. [ ] Login (maybe ?)
2. [X] Move session to env variables
3. [ ] Make preferred days configurable
4. [ ] Make preferred time configurable

