package jokes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"errors"
)

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

func (jkres JokeResult) Getid () int {
	return jkres.id
}

func (jkres JokeResult) Geterror () error {
	return jkres.error
}

func (jkres JokeResult) Getsetup () string {
	return jkres.setup
}

func (jkres JokeResult) Getpunchline () string {
	return jkres.punchline
}

///////////////////////////////////////////////////////////////
// router determines which api to call
func JokeRouter(joke_type string) (JokeAPI, error) {
	if joke_type == "programming" || joke_type == "general" {
		return JokeResultSite1{}, nil
	} else if joke_type == "Pun" || joke_type == "Spooky" || joke_type == "Christmas" {
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
	if joke_type != "programming" && joke_type != "general" {
		ch <- JokeResult{error: errors.New("Bad joke type! Shold be programming or general")}
	}

	path := fmt.Sprintf("https://official-joke-api.appspot.com/jokes/%s/random", joke_type)
	response, err := http.Get(path)
	if err != nil {
		ch <- JokeResult{error: err}
		os.Exit(1)
	}

	// parse response
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		ch <- JokeResult{error: err}
		os.Exit(1)
	}

	var jkres1_arr []JokeResultSite1
	json.Unmarshal(responseData, &jkres1_arr)
	if len(jkres1_arr) == 0 {
		ch <- JokeResult{error: errors.New("bad array length 0")}
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
	Id 		int `json:"id"`
	Setup     string `json:"setup"`
	Delivery string `json:"delivery"`
}

func (jr2 JokeResultSite2) GetARandomJoke(ch chan JokeResult, joke_type string) {
	if joke_type != "Christmas" && joke_type != "Pun" && joke_type != "Spooky" {
		ch <- JokeResult{error: errors.New("Bad joke type! Shold be programming or general")}
	}

	path := fmt.Sprintf("https://v2.jokeapi.dev/joke/%s?type=twopart", joke_type)
	response, err := http.Get(path)
	if err != nil {
		ch <- JokeResult{error: err}
		os.Exit(1)
	}

	// parse response
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		ch <- JokeResult{error: err}
		os.Exit(1)
	}

	var jkres2 JokeResultSite2
	json.Unmarshal(responseData, &jkres2)

	var jkres JokeResult = JokeResult{}
	jkres.id = jkres2.Id
	jkres.setup = jkres2.Setup
	jkres.punchline = jkres2.Delivery
	jkres.site = 2
	jkres.error = nil
	
	ch <- jkres
}