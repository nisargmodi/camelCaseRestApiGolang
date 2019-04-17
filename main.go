package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sort"
	"strings"

	"github.com/gorilla/mux"
)

var wordDict [5]string

func GetCamelCase(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	json.NewEncoder(w).Encode(wordBreak(params["input"], wordDict))
}

func wordBreak(s string, wordDict [5]string) []string {
	if len(wordDict) == 0 {
		return []string{}
	}

	dict := make(map[string]bool, len(wordDict))
	length := make(map[int]bool, len(wordDict))

	for _, w := range wordDict {
		dict[w] = true
		length[len(w)] = true
	}

	sizes := make([]int, 0, len(length))
	for k := range length {
		sizes = append(sizes, k)
	}
	sort.Ints(sizes)

	n := len(s)

	dp := make([]float64, len(s)+1)
	dp[0] = 1

	for i := 0; i <= n; i++ {
		if dp[i] == 0 {
			continue
		}

		for _, size := range sizes {
			if i+size <= n && dict[s[i:i+size]] {
				dp[i+size] += dp[i]
			}
		}
	}

	if dp[n] == 0 {
		return []string{}
	}

	res := make([]string, 0, int(dp[n]))

	// dfs
	var dfs func(int, string)
	dfs = func(i int, str string) {
		if i == len(s) {
			res = append(res, str[0:])
			return
		}

		for _, size := range sizes {
			if i+size <= len(s) && dict[s[i:i+size]] {
				if i == 0 {
					dfs(i+size, str+s[i:i+size])
				} else {
					dfs(i+size, str+strings.Title(s[i:i+size]))
				}
			}
		}
	}

	dfs(0, "")

	return res
}

func main() {
	router := mux.NewRouter()

	wordDict[0] = "cat"
	wordDict[1] = "cats"
	wordDict[2] = "and"
	wordDict[3] = "sand"
	wordDict[4] = "dog"

	router.HandleFunc("/camelcase/{input}", GetCamelCase).Methods("GET")
	log.Fatal(http.ListenAndServe(":8000", router))
}
