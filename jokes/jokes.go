package jokes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"errors"
)

////////////////////////////////////////////////////////////////
// generic interface and struct for both sides - an abstraction
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


func (jkres JokeResult) Getid () int {
	return jkres.id
}

func (jkres JokeResult) Geterror () error {
	return jkres.error
}


/////////////////////////////////////////////////////////////////
// JSON result from site 1
type JokeResultSite1 struct {
	Id        int    `json:"id"`
	Type      string `json:"type"`
	Setup     string `json:"setup"`
	Punchline string `json:"punchline"`
}


//////////////////////////////////////////////////////////////////
// JSON result from site 2
type JokeResultSite2 struct {
	Category string `json:"category"`
	Id 		int `json:"id"`
	Setup     string `json:"setup"`
	Delivery string `json:"delivery"`
}





// // A thread-safe map
// type Cache struct {
// 	dict map[int]bool
// 	mux  sync.Mutex
// }


// // check if id is in map
// func (c *Cache) checkVisited(id int) bool {
// 	c.mux.Lock()
// 	defer c.mux.Unlock()
// 	_, ok := c.dict[id]
// 	if ok == false {
// 		c.dict[id] = true
// 		return false
// 	}
// 	return true
// }


// TODO: make errors consistent
// call the joke api and return
func (jr1 JokeResultSite1) GetARandomJoke(ch chan JokeResult, joke_type string) {
	if joke_type != "programming" && joke_type != "general" {
		//fmt.Println("Bad joke type! Shold be programming or general")
		ch <- JokeResult{error: errors.New("Bad joke type! Shold be programming or general")}
	}

	// build url with query
	//query := "general"
	path := fmt.Sprintf("https://official-joke-api.appspot.com/jokes/%s/random", joke_type)
	//fmt.Println("path: ", path)

	// call API
	//response, err := http.Get("https://official-joke-api.appspot.com/random_joke")
	//response, err := http.Get("https://official-joke-api.appspot.com/jokes/programming/random")
	response, err := http.Get(path)
	if err != nil {
		// fmt.Print(err.Error())
		ch <- JokeResult{error: err}
		os.Exit(1)
	}

	// parse the response
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		ch <- JokeResult{error: err}
		os.Exit(1)
	}


	// json parsing
	var jkres1_arr []JokeResultSite1
	json.Unmarshal(responseData, &jkres1_arr)
	if len(jkres1_arr) == 0 {
		ch <- JokeResult{error: errors.New("bad array length 0")}
		os.Exit(1)
	}


	// prepare return value
	var jkres JokeResult = JokeResult{}
	jkres.id = jkres1_arr[0].Id
	jkres.setup = jkres1_arr[0].Setup
	jkres.punchline = jkres1_arr[0].Punchline
	jkres.site = 1
	jkres.error = nil

	// put result struct into the channel - the asynchronous way of returning(?)
	ch <- jkres
}


// // TODO: make errors consistent
// // call the joke api and return
// func (jr2 JokeResultSite2) GetARandomJoke(ch chan JokeResult, joke_type string) {
// 	if joke_type != "Christmas" && joke_type != "Pun" && joke_type != "Spooky" {
// 		//fmt.Println("Bad joke type! Shold be programming or general")
// 		ch <- JokeResult{Error: errors.New("Bad joke type! Shold be programming or general")}
// 	}

// 	// return value
// 	var jkres JokeResult = JokeResult{}

// 	// build url with query
// 	//query := "general"
// 	path := fmt.Sprintf("https://v2.jokeapi.dev/joke/%s?type=twopart", joke_type)
// 	//fmt.Println("path: ", path)

// 	// call API
// 	//response, err := http.Get("https://official-joke-api.appspot.com/random_joke")
// 	//response, err := http.Get("https://official-joke-api.appspot.com/jokes/programming/random")
// 	response, err := http.Get(path)
// 	if err != nil {
// 		jkres.Error = err
// 		ch <- jkres
// 		fmt.Print(err.Error())
// 		os.Exit(1)
// 	}

// 	// parse the response
// 	responseData, err := ioutil.ReadAll(response.Body)
// 	if err != nil {
// 		jkres.Error = err
// 		ch <- jkres
// 		fmt.Print(err.Error())
// 		os.Exit(1)
// 	}

// 	// json parsing
// 	var jkres JokeResult
// 	json.Unmarshal(responseData, &jkres)
// 	// if len(jkres_arr) == 0 {
// 	// 	ch <- JokeResult{Error: errors.New("bad array length 0")}
// 	// 	os.Exit(1)
// 	// }

// 	// put result struct into the channel - the asynchronous way of returning(?)
// 	ch <- jkres
// }



// // TODO: figure out what to do with joke_type.. maybe make it into an array?
// // get n non-repeating jokes, concurrently
// func GetRandomJokes(num_jokes int, joke_type string) ([]JokeResult, error){
// 	// parameters error checking
// 	if joke_type != "programming" && joke_type != "general" {
// 		return nil, errors.New("Bad joke type! Shold be programming or general")
// 	}
// 	if num_jokes < 0 {
// 		return nil, errors.New("Invalid num_jokes value!")
// 	}

// 	// create a channel
// 	ch := make(chan JokeResult, num_jokes)
// 	// return value
// 	jkres_list := make([]JokeResult, num_jokes)
// 	// thread-safe dict
// 	cache := Cache{dict: make(map[int]bool)}


// 	// TODO:change this! 
// 	api := JokeResultSite1{}


// 	// asynchrounously call GetARandomJoke. n non-repeaing jokes. A better way to write this?
// 	cnt := 0
// 	for {
// 		if cnt == num_jokes {
// 			break
// 		}

// 		// call go routine and error check
// 		go api.GetARandomJoke(ch, joke_type)
// 		jkres := <-ch
// 		if jkres.error != nil {
// 			fmt.Print(jkres.error.Error())
// 			os.Exit(1)
// 		}

// 		// check duplicate. increment only when it's a new joke
// 		if cache.checkVisited((jkres.id)) {
// 			fmt.Println("old joke... fetching you a new one")
// 			continue
// 		} else {
// 			jkres_list[cnt] = jkres
// 			// fmt.Println((cnt))
// 			// fmt.Println("	Setup: ", jkres_list[cnt].Setup)
// 			// fmt.Println("	Punchline: ", jkres_list[cnt].Punchline)
// 			cnt += 1
// 		}
// 	}
// 	return jkres_list, nil
// }
