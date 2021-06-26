# anytime-anywhere
Script to automate the booking process for **AnytimeFitness Gym**

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

## TODO

1. [ ] Login (maybe ?)
2. [X] Move session to env variables

