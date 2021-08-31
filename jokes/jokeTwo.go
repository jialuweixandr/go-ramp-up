package jokes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"errors"
)

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