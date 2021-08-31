package main 

import (
	"fmt"
	"os"
	"jokes"
	"cache"
	"errors"
)

func main() {
	for {
		joke_type := getJokeType()

		fmt.Printf("You've select: %s\n", joke_type)
		fmt.Printf("Now, how many %s jokes would you like?\n", joke_type)

		num_jokes := getJokeNum()

		// get jokes
		err := GetRandomJokes(num_jokes, joke_type)
		if err != nil {
			fmt.Print(err.Error())
			os.Exit(1)
		}
	}
}

func getJokeType() string {
	fmt.Println("What kind of joke would you like?")
	fmt.Print("Select joke category: ")
	for i := 1; i<6; i++ {
		fmt.Printf(" %v %s", i, jokes.Joke_types[i])
	}
	fmt.Print("\n")
	var n int
	_, err := fmt.Scanf("%d", &n)
	joke_type, ok := jokes.Joke_types[n]
	if err != nil || !ok {
		fmt.Println(errors.New("Invalid joke type!"))
		os.Exit(1)
	}
	return joke_type
}


func getJokeNum() int {
	var num_jokes int
	if _, err := fmt.Scanf("%d", &num_jokes); err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
	if num_jokes < 0 {
		fmt.Print(errors.New("Invalid num_jokes value!"))
		os.Exit(1)
	}
	return num_jokes
}

func GetRandomJokes(num_jokes int, joke_type string) (error){
	ch := make(chan jokes.JokeResult, num_jokes)
	c := cache.NewCache()
	jkres_list := make([]jokes.JokeResult, num_jokes)
	
	// get the underlying struct based on joke_type
	api, err := jokes.JokeRouter(joke_type)
	if err != nil{
		fmt.Print(err.Error())
		os.Exit(1)
	}

	cnt := 0
	for {
		if cnt == num_jokes {
			break
		}

		go api.GetARandomJoke(ch, joke_type)
		jkres := <-ch
		if jkres.GetError() != nil {
			return jkres.GetError()
		}

		if c.CheckVisited((jkres.GetId())) {
			continue
		} else {
			jkres_list[cnt] = jkres
			fmt.Println((cnt+1))
			fmt.Println("	Setup: ", jkres_list[cnt].GetSetup())
			fmt.Println("	Punchline: ", jkres_list[cnt].GetPunchline())
			cnt += 1
		}
	}
	return nil
}
