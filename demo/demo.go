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

	var n int
	fmt.Print("Enter a number: ")
	_, err := fmt.Scanf("%d", &n)
	fmt.Println("Generating ", n, " jokes....................")
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	
	//api := jokes.JokeResultSite1{}
	// todo: change this into an array

	// ch := make(chan jokes.JokeResult)
	// go api.GetARandomJoke(ch, "programming")
	// res := <- ch
	// fmt.Println(res)


	jokes, err := GetRandomJokes(n, "programming")
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
	fmt.Println(jokes)


}


// TODO!!!!: add a router function in the jokes module to determine which api to call....
// get n non-repeating jokes, concurrently
func GetRandomJokes(num_jokes int, joke_type string) ([]jokes.JokeResult, error){
	// parameters error checking 
	if joke_type != "programming" && joke_type != "general" {
		return nil, errors.New("Bad joke type! Shold be programming or general")
	}
	if num_jokes < 0 {
		return nil, errors.New("Invalid num_jokes value!")
	}

	// create a channel
	ch := make(chan jokes.JokeResult, num_jokes)
	// return value
	jkres_list := make([]jokes.JokeResult, num_jokes)



	// thread-safe dict
	// TODO: change this!!! Don't expose internal structure
	// cache := Cache{dict: make(map[int]bool)}
	c := cache.NewCache()


	// TODO:change this! 
	api := jokes.JokeResultSite1{}


	// asynchrounously call GetARandomJoke. n non-repeaing jokes. A better way to write this?
	cnt := 0
	for {
		if cnt == num_jokes {
			break
		}

		// call go routine and error check
		go api.GetARandomJoke(ch, joke_type)
		jkres := <-ch
		if jkres.Geterror() != nil {
			fmt.Print(jkres.Geterror().Error())
			os.Exit(1)
		}

		// check duplicate. increment only when it's a new joke
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