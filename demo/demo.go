package main 

import (
	"fmt"
	"os"
	"jokes"
)

// input: take number n from command line
// output: a list of n jokes: Setup and punchline (format tbd)
func main() {

	var n int
	fmt.Print("Enter a number: ")
	_, err := fmt.Scanf("%d", &n)
	fmt.Println("Generating ", n, " jokes...")
	fmt.Println("------")

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	

	// jokes.GetARandomJoke(ch, "programming")
	// res := <- ch
	// fmt.Println(res)


	jokes, err := jokes.GetRandomJokes(n, "programming")
	fmt.Println(jokes)


}