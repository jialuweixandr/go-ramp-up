package jokes

import (
	"errors"
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