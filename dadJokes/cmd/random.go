/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>

In summary, this code defines a CLI command 'random' that fetches a random dad joke from the icanhazdadjoke API
by making an HTTP GET request and decoding the JSON response into a Joke struct.
The fetched joke is then printed to the console.
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/spf13/cobra"
)

// randomCmd represents the random command
// this command calls the getRandomJoke function
var randomCmd = &cobra.Command{
	Use:   "random",
	Short: "Fetches a random dad joke!",
	Long:  `this command calls the icanhazdadjoke api and retrieves a random dad joke.`,
	Run: func(cmd *cobra.Command, args []string) {

		jokeTerm, _ := cmd.Flags().GetString("term")

		if jokeTerm != "" {
			getRandomJokeWithTerm(jokeTerm)
		} else {
			getRandomJoke()
		}
	},
}

// Here you will define your flags and configuration settings.
// let us define a flag for the random command which allows us to pass in a string for a term
func init() {
	rootCmd.AddCommand(randomCmd)
	randomCmd.PersistentFlags().String("term", "", "A search term to find a dad joke")

	// the persistent flag is going to be of string type and we must pass in :
	// the name of the flag which in this case is "term"
	// the default value --> in our case empty string
	// and a description of the flag functionality

}

// a call to the api returns a joke(string), ID(string), and status(int)
// let's create a joke structure to hold these values
type Joke struct {
	ID     string `json:"id"`
	Joke   string `json:"joke"`
	Status int    `json:"status"`
}

type searchResult struct {
	Results    json.RawMessage `json:"results"`
	searchTerm string          `json:"search_term"`
	Status     int             `json:"status"`
	TotalJokes int             `json:"total_jokes"`
}

// this func gets called within the random cmd
func getRandomJoke() {
	url := "https://icanhazdadjoke.com/" // stores the url into a var
	responseBytes := getJokeData(url)    // calls the getJokeData func and stores the responseBytes into a var
	joke := Joke{}                       // creates a joke structure

	// json.Umarshal decodes json data and transforms it into a go data struct
	// in this case we pass in the responseBytes and it transforms it into a joke struct
	if err := json.Unmarshal(responseBytes, &joke); err != nil {
		fmt.Printf("Could not unmarshal reponseBytes.", err)
	}

	fmt.Println(string(joke.Joke))
}

func getRandomJokeWithTerm(jokeTerm string) {
	total, results := getJokeDataWithTerm(jokeTerm)
	randomJokeList(total, results)
}

// func responsible for making the http request to the icanhazdadjoke API
func getJokeData(baseAPI string) []byte {
	request, err := http.NewRequest(http.MethodGet, baseAPI, nil)

	if err != nil {
		log.Println("could not request a dad joke", err)
	}

	request.Header.Add("Accept", "application/json")
	request.Header.Add("User-Agent", "DadJokeCLI/1.0 (local dir)")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Println("could not request a dad joke", err)
	}

	responseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println("could not read response body", err)
	}

	return responseBytes
}

// func responsible for making the http request to the icanhazdadjoke API (WITH TERM)
func getJokeDataWithTerm(jokeTerm string) (totalJokes int, jokeList []Joke) {
	url := fmt.Sprintf("https://icanhazdadjoke.com/search?term=%s", jokeTerm)
	responseBytes := getJokeData(url)

	jokeListRaw := searchResult{}

	if err := json.Unmarshal(responseBytes, &jokeListRaw); err != nil {
		fmt.Printf("Could not unmarshal reponseBytes. %v", err)
	}

	jokes := []Joke{}

	if err := json.Unmarshal(jokeListRaw.Results, &jokes); err != nil {
		fmt.Printf("Could not unmarshal reponseBytes. %v", err)
	}

	return jokeListRaw.TotalJokes, jokes
}

func randomJokeList(length int, jokeList []Joke) {
	// takes in two parameters
	// int length to specify the ceiling for our random range
	// []Joke jokeList which is the list of jokes we want to choose from

	rand.Seed(time.Now().UnixNano())

	min := 0
	max := length - 1

	if length <= 0 {
		err := fmt.Errorf("length must be greater than 0")
		fmt.Println(err.Error())
	} else {
		randomNum := min + rand.Intn(max-min)
		fmt.Println(jokeList[randomNum].Joke)
	}

	// this function generates a random number
	// it then checks if the number is less than or equal to 0
	// if so we throw an error --> else we print out the random joke
}
