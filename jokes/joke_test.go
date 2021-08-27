package jokes

import (
	"testing"
)

func TestGetARandomJokeNormal(t *testing.T) {
	joke_type := "Programming"
	api, err := JokeRouter(joke_type)
	if err != nil {
		t.Fatalf("Invalid joke type!")
	}
	ch := make(chan JokeResult)
	go api.GetARandomJoke(ch, joke_type)
	jkres := <- ch
	if jkres.error != nil {
		t.Fatalf("Error somewhere down the line..")
	}
}

func TestJokeRouterEmpty(t *testing.T) {
	joke_type := ""
	_, err := JokeRouter(joke_type)
	if err == nil {
		t.Fatalf("Expect error with empty joke_type string")
	}
}