package jokes


import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"errors"
	"strings"
)


//////////////////////////////////////////////////////////////
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
