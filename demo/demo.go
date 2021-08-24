package main 

import (
	"fmt"
	"os"
	"jokes"
	"cache"
	"errors"
)

// input: take number n from command line
// output: a list of n jokes: Setup and punchline (format tbd)
func main() {

	var num_jokes int
	var joke_type_int int
	jkmap := map[int]string{
		1: "general",
		2: "programming",
		3: "Pun",
		4: "Spooky",
		5: "Christmas",
	}

	// get joke type
	fmt.Println("What kind of joke would you like?")
	fmt.Println("Select joke category:  1) General 2) Programming 3) Pun 4) Spooky 5) Christmas ")
	_, err := fmt.Scanf("%d", &joke_type_int)
	joke_type, ok := jkmap[joke_type_int]
	if err != nil || !ok {
		fmt.Print(err.Error()) // TODO: customize error
		os.Exit(1)
	}

	// get joke num
	fmt.Printf("You've select: %s\n", joke_type)
	fmt.Printf("Now, how many %s jokes would you like?\n", joke_type)
	_, err = fmt.Scanf("%d", &num_jokes)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	
	// api := jokes.JokeResultSite2{}
	// // todo: change this into an array
	// ch := make(chan jokes.JokeResult)
	// go api.GetARandomJoke(ch, "Christmas")
	// res := <- ch
	// fmt.Println(res)


	jokes, err := GetRandomJokes(num_jokes, joke_type)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
	fmt.Println(jokes)


}


// get n non-repeating jokes, concurrently
func GetRandomJokes(num_jokes int, joke_type string) ([]jokes.JokeResult, error){
	// error checking 
	if joke_type != "programming" && joke_type != "general" && joke_type != "Pun" && joke_type != "Christmas" && joke_type != "Spooky"{
		return nil, errors.New("Bad joke type! Joke type should be one of followwing: programming, general, Christmas, Pun, Spooky")
	}
	if num_jokes < 0 {
		return nil, errors.New("Invalid num_jokes value!")
	}

	ch := make(chan jokes.JokeResult, num_jokes)
	c := cache.NewCache() // tracks existing ids
	jkres_list := make([]jokes.JokeResult, num_jokes) // return value
	api := jokes.JokeRouter(joke_type) // find correct struct to call based on joke_type


	// asynchrounously call GetARandomJoke. n non-repeaing jokes. A better way to write this?
	cnt := 0
	for {
		if cnt == num_jokes {
			break
		}

		go api.GetARandomJoke(ch, joke_type)
		jkres := <-ch
		if jkres.Geterror() != nil {
			fmt.Print(jkres.Geterror().Error())
			os.Exit(1)
		}

		// check duplicate
		if c.CheckVisited((jkres.Getid())) {
			fmt.Println("old joke... fetching you a new one")
			continue
		} else {
			fmt.Println("Found a joke for you!")
			jkres_list[cnt] = jkres
			// fmt.Println((cnt))
			// fmt.Println("	Setup: ", jkres_list[cnt].Setup)
			// fmt.Println("	Punchline: ", jkres_list[cnt].Punchline)
			cnt += 1
		}
	}
	return jkres_list, nil
}
