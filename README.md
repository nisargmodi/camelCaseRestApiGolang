# camelCaseRestApiGolang
Sample golang REST API that formats a given string into camelCase format.

This API calls Oxford dictionary's API: https://od-api.oxforddictionaries.com:443/api/v1/inflections/en/ to check if each candiadte word is a valid dictionary word.

The algorithm runs in O(2^N) where N is the length of the string. That's because it generates a list of possible answers. So for just "onetwo" as input:

["onEtWo","oneTwO","oNEtWo","oNEtWO","oneTwo","oNetWo","onETwo","oNeTwo","oneTWo"]

are the list of words returned. Since the oxford API considers inputs like "Et" and "Wo" as valid words, the answer list is somehwat unusual for now.
The APIs can be used in a better way or some other dictionary API can be used to give more accurate results.

To run this:

1. Install Go
2. Clone the repo
3. From the terminal run: $ go build && camelCaseRestApiGolang
4. Open the browser and go to http://localhost:8000/camelcase/<input>
For example: http://localhost:8000/camelcase/onetwo
5. Note, since the solution is slow it can take upto 10 seconds for results to show up.
6. The Oxford API has a limit of 60 requests/minute after which it gives 403s so longer words may not work now.

Future work:

1. Add retry with exponential backoof + jitter mechanism to REST calls downstream to Oxford API to support longer words.
2. Tune the algorithm to make concurrent (& parallel) requests for faster output: check the earlier commits where go routines were tried but since tje current algo is very synchronous it didn't help
3. Read the secrets from the config files instead of hard-coding: check the earlier commits where I have tried this but it didn't work and needs some investigation


