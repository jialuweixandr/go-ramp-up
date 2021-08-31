package jokes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"errors"
	"strings"
)


var Joke_types = map[int]string{
	1: "General",
	2: "Programming",
	3: "Pun",
	4: "Spooky",
	5: "Christmas",
}

///////////////////////////////////////////////////////////////
// generic interface and struct for both joke apis - an abstraction
type JokeAPI interface {
	GetARandomJoke(ch chan JokeResult, joke_type string)
}

type JokeResult struct {
	id int
	site int
	setup string
	punchline string
	error error
}

func (jkres JokeResult) GetId () int {
	return jkres.id
}

func (jkres JokeResult) GetError () error {
	return jkres.error
}

func (jkres JokeResult) GetSetup () string {
	return jkres.setup
}

func (jkres JokeResult) GetPunchline () string {
	return jkres.punchline
}

///////////////////////////////////////////////////////////////
// router determines which api to call
func JokeRouter(joke_type string) (JokeAPI, error) {
	if joke_type == Joke_types[1] || joke_type == Joke_types[2] {
		return JokeResultSite1{}, nil
	} else if joke_type == Joke_types[3] || joke_type == Joke_types[4] || joke_type == Joke_types[5] {
		return JokeResultSite2{}, nil
	} else {
		return nil, errors.New("Invalid joke type!")
	}
}

////////////////////////////////////////////////////////////////
// JSON structure and function for joke api 1
type JokeResultSite1 struct {
	Id        int    `json:"id"`
	Type      string `json:"type"`
	Setup     string `json:"setup"`
	Punchline string `json:"punchline"`
}

func (jr1 JokeResultSite1) GetARandomJoke(ch chan JokeResult, joke_type string) {
	if joke_type != Joke_types[1] && joke_type != Joke_types[2] {
		ch <- JokeResult{error: errors.New("Bad joke type! Shold be programming or general")}
	}

	path := fmt.Sprintf("https://official-joke-api.appspot.com/jokes/%s/random", strings.ToLower(joke_type))
	response, err := http.Get(path)
	if err != nil {
		ch <- JokeResult{error: err}
		os.Exit(1)
	}

	// parse response
	var jkres1_arr []JokeResultSite1
	dec := json.NewDecoder(response.Body)
	if err = dec.Decode(&jkres1_arr); err != nil {
		ch <- JokeResult{error: err}
	 	os.Exit(1)
	}


	var jkres JokeResult = JokeResult{}
	jkres.id = jkres1_arr[0].Id
	jkres.setup = jkres1_arr[0].Setup
	jkres.punchline = jkres1_arr[0].Punchline
	jkres.site = 1
	jkres.error = nil

	ch <- jkres
}


/////////////////////////////////////////////////////////////////
//JSON structure and function for joke api 2
type JokeResultSite2 struct {
	Category string `json:"category"`
	Id int `json:"id"`
	Setup string `json:"setup"`
	Delivery string `json:"delivery"`
}

func (jr2 JokeResultSite2) GetARandomJoke(ch chan JokeResult, joke_type string) {
	if joke_type != Joke_types[3] && joke_type != Joke_types[4] && joke_type != Joke_types[5] {
		ch <- JokeResult{error: errors.New("Bad joke type! Shold be programming or general")}
	}

	path := fmt.Sprintf("https://v2.jokeapi.dev/joke/%s?type=twopart", joke_type)
	response, err := http.Get(path)
	if err != nil {
		ch <- JokeResult{error: err}
		os.Exit(1)
	}

	// parse response
	var jkres2 JokeResultSite2
	dec := json.NewDecoder(response.Body)
	if err = dec.Decode(&jkres2); err != nil {
		ch <- JokeResult{error: err}
	 	os.Exit(1)
	}

	var jkres JokeResult = JokeResult{}
	jkres.id = jkres2.Id
	jkres.setup = jkres2.Setup
	jkres.punchline = jkres2.Delivery
	jkres.site = 2
	jkres.error = nil

	ch <- jkres
}