/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"
)

// randomCmd represents the random command
var randomCmd = &cobra.Command{
	Use:   "random",
	Short: "Get a random dad joke",
	Long:  `This command fetches a random dad joke from the icanhazdadjoke api`,
	Run: func(cmd *cobra.Command, args []string) {
		jokeTerm, _ := cmd.Flags().GetString("term")

		if jokeTerm != "" {
			getRandomJokeWithTerm(jokeTerm)
		} else {
			getRandomJoke()
		}
	},
}

func init() {
	rootCmd.AddCommand(randomCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	randomCmd.PersistentFlags().String("term", "", "A search term")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// randomCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

type Joke struct {
	ID     string `json:"id"`
	Joke   string `json:"joke"`
	Status int    `json:"status"`
}

type SearchResult struct {
	Results    json.RawMessage `json:"results"`
	SearchTerm string		   `json:"search_term"`
	Status     int			   `json:"status"`
	TotalJokes int             `json:"total_jokes"`
}

func getRandomJoke() {
	fmt.Println("Get random dad joke :P")
	url := "https://icanhazdadjoke.com/"
	responseBytes := getJokeData(url)
	joke := Joke{}

	if err := json.Unmarshal(responseBytes, &joke); err != nil {
		fmt.Printf("Failed to unmarshal responseBytes. %v", err)
	}

	fmt.Println(string(joke.Joke))
}

func getRandomJokeWithTerm(jokeTerm string) {
	total, jokes := getJokeDataWithTerm(jokeTerm)
	randomiseJokeList(total, jokes)
}

func getJokeDataWithTerm(jokeTerm string) (totalJokes int, jokeList []Joke) {
	url := fmt.Sprintf("https://icanhazdadjoke.com/search?term=%s", jokeTerm)
	responseBytes := getJokeData(url)

	jokeListRaw := SearchResult{}

	if err := json.Unmarshal(responseBytes, &jokeListRaw); err != nil {
		log.Printf("Could not unmarshal responseBytes. %v", err)
	}

	var jokes []Joke
	if err := json.Unmarshal(jokeListRaw.Results, &jokes); err != nil {
		log.Printf("Failed to unmarshal response. %v", err)
	}

	return jokeListRaw.TotalJokes, jokes
}

func getJokeData(baseApi string) []byte {
	request, err := http.NewRequest(
		http.MethodGet,
		baseApi,
		nil,
	)

	if err != nil {
		log.Printf("Failed to request dadjoke. %v", err)
	}

	request.Header.Add("Accept", "application/json")
	request.Header.Add("User-Agent", "Dadjoke CLI (https://github.com/nicklpeterson/dadjoke)")

	response, err := http.DefaultClient.Do(request);

	if err != nil {
		log.Printf("Failed to make request. %v", err)
	}

	responseByte, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("Failed to read response body. %v", err)
	}

	return responseByte
}

func randomiseJokeList(length int, jokeList []Joke) {
	rand.Seed(time.Now().Unix())
	min := 0
	max := length - 1

	if length <= 0 {
		err := fmt.Errorf("No Jokes found")
		fmt.Println(err.Error())
	} else {
		randomNumber := min + rand.Intn(max - min)
		fmt.Println(jokeList[randomNumber].Joke)
	}
}
