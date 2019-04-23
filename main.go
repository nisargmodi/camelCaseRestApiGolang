package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
)

var checked = make(map[string]bool)
var cache = make(map[string]bool)

type Configuration struct {
	app_id  string
	app_key string
}

func isWord(word string) bool {
	if checked[word] {
		return cache[word]
	}
	checked[word] = true
	lexicalCategory := "/lexicalCategory=suffix,noun,determiner,adverb,combining_form,idiomatic,predeterminer,particle,residual,adjective,preposition,prefix,other,verb,numeral,conjunction,pronoun,interjection,contraction"
	url := "https://od-api.oxforddictionaries.com:443/api/v1/inflections/en/" + word + lexicalCategory
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)

	file, _ := os.Open("conf.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("Error reading config! ", err)
	}

	req.Header.Set("app_id", configuration.app_id)
	req.Header.Set("app_key", configuration.app_key)
	res, error := client.Do(req)
	if error != nil {
		fmt.Printf("The HTTP request failed with error %s\n", error)
	} else if res.StatusCode != 404 {
		cache[word] = true
		return true
	}
	cache[word] = true
	return false
}

func GetCamelCase(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	json.NewEncoder(w).Encode(wordBreak(params["input"]))
}

func wordBreak(str string) []string {
	strLength := len(str)
	solution := make([][]int, strLength)

	for i := strLength - 1; i >= 0; i-- {
		for j := i + 1; j <= strLength; j++ {
			possibleWord := str[i:j]
			if j == strLength || len(solution[j]) > 0 {
				if ok := isWord(possibleWord); ok == true {
					solution[i] = append(solution[i], j)
				}
			}
		}
	}

	sentencePaths := [][]int{[]int{0}}
	sentences := make([]string, 0)

	for {
		nextSentencePaths := [][]int{}
		for _, sentencePath := range sentencePaths {
			sentencePathLength := len(sentencePath)
			if sentencePath[sentencePathLength-1] == strLength {
				lastPosition := sentencePathLength - 1
				temp := []string{}
				for i := 0; i < lastPosition; i++ {
					if i == 0 {
						temp = append(temp, str[sentencePath[i]:sentencePath[i+1]])
					} else {
						temp = append(temp, strings.Title(str[sentencePath[i]:sentencePath[i+1]]))
					}
				}
				sentences = append(sentences, strings.Join(temp, ""))
			} else {
				for _, j := range solution[sentencePath[sentencePathLength-1]] {
					newPath := append(sentencePath, j)
					nextSentencePaths = append(nextSentencePaths, newPath)
				}
			}
		}
		if len(nextSentencePaths) == 0 {
			break
		} else {
			sentencePaths = nextSentencePaths
		}
	}

	encountered := map[string]bool{}

	for v := range sentences {
		encountered[sentences[v]] = true
	}

	result := []string{}
	for key, _ := range encountered {
		result = append(result, key)
	}

	return result
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/camelcase/{input}", GetCamelCase).Methods("GET")
	log.Fatal(http.ListenAndServe(":8000", router))
}
