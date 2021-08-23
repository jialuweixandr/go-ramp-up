package jokes

import (
	"testing"
)

func TestGetARandomJokeType(t *testing.T) {
	joke_type := "programming"
	ch := make(chan JokeResult)
	go GetARandomJoke(ch, joke_type) // Question: why were you blocking when I omitted "go" keyword??
	jkres := <- ch

	if jkres.Type != joke_type ||  jkres.Error != nil {
		t.Fatalf("Either wrong joke type or error somewhere down the line...Expected: %v, Actual: %v", joke_type, jkres.Type) // TODO
	}
}


// TODO: maybe check the exact error text? What's this error..?
func TestGetARandomJokeEmpty(t *testing.T) {
	joke_type := ""
	ch := make(chan JokeResult)
	go GetARandomJoke(ch, joke_type) // Question: why were you blocking when I omitted "go" keyword??
	jkres := <- ch

	if jkres.Error == nil {
		t.Fatalf("Expect error with empty joke_type string. Actual: No error.")
	}
}


func TestGetRandomJokes(t *testing.T) {
	joke_type := "general"
	num_jokes := 3

	var jkres_list, error = GetRandomJokes(num_jokes, joke_type)

	if len(jkres_list) != num_jokes || error != nil {
		t.Fatalf("Unexpected result length. Expected: %v, Actual: %v", len(jkres_list), num_jokes)
	}

	for _, jkres := range jkres_list {
		if jkres.Type != joke_type ||  jkres.Error != nil {
			t.Fatalf("Either wrong joke type or error somewhere down the line...Expected: %v, Actual: %v", joke_type, jkres.Type) // TODO
		}
	}
}

// TODO: what would happen if joke_type is bad?
func TestGetRandomJokesBad(t *testing.T) {
	joke_type := "general"
	num_jokes := -1

	var _, error = GetRandomJokes(num_jokes, joke_type)

	if error == nil{
		t.Fatalf("Expect error. Actual: No error.")
	}
}