package main 

import (
	"fmt"
	"os"
	"jokes"
	"cache"
	"errors"
)

func main() {
	jokemap := map[int]string{
		1: "general",
		2: "programming",
		3: "Pun",
		4: "Spooky",
		5: "Christmas",
	}

	for {
		fmt.Println("What kind of joke would you like?")
		fmt.Println("Select joke category:  1) General 2) Programming 3) Pun 4) Spooky 5) Christmas ")
		var n int
		_, err := fmt.Scanf("%d", &n)
		joke_type, ok := jokemap[n]
		if err != nil || !ok {
			fmt.Print(errors.New("Invalid joke type!"))
			os.Exit(1)
		}

		fmt.Printf("You've select: %s\n", joke_type)
		fmt.Printf("Now, how many %s jokes would you like?\n", joke_type)
		var num_jokes int
		_, err = fmt.Scanf("%d", &num_jokes)
		if err != nil {
			fmt.Print(err.Error())
			os.Exit(1)
		}
		if num_jokes < 0 {
			fmt.Print(errors.New("Invalid num_jokes value!"))
			os.Exit(1)
		}

		// get jokes
		err = GetRandomJokes(num_jokes, joke_type)
		if err != nil {
			fmt.Print(err.Error())
			os.Exit(1)
		}
	}
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
		if jkres.Geterror() != nil {
			return jkres.Geterror()
		}

		if c.CheckVisited((jkres.Getid())) {
			continue
		} else {
			jkres_list[cnt] = jkres
			fmt.Println((cnt+1))
			fmt.Println("	Setup: ", jkres_list[cnt].Getsetup())
			fmt.Println("	Punchline: ", jkres_list[cnt].Getpunchline())
			cnt += 1
		}
	}
	return nil
}
