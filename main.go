package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
)

type Package struct {
	Name    string `json:"name"`
	Manager string `json:"manager"`
}

type TestableSuggestions struct {
	Matcher  string    `json:"matcher"`
	TestOK   []string  `json:"test_ok"`
	TestFail []string  `json:"test_fail"`
	Packages []Package `json:"packages"`
}

type OutputSuggestions struct {
	Matcher  string    `json:"matcher"`
	Packages []Package `json:"packages"`
}

func main() {
	jsonFile, err := os.Open("./suggestion.json")
	if err != nil {
		log.Default().Fatal(err)
	}
	defer jsonFile.Close()
	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		log.Default().Fatal(err)
	}
	var testableSuggestions []TestableSuggestions
	var outputSuggestions []OutputSuggestions
	err = json.Unmarshal(byteValue, &testableSuggestions)
	if err != nil {
		log.Default().Fatal(err)
	}
	for _, suggestion := range testableSuggestions {
		regex, err := regexp.Compile(suggestion.Matcher)
		if err != nil {
			log.Default().Fatal(err)
		}
		if len(suggestion.TestOK) < 1 || len(suggestion.TestFail) < 1 {
			log.Default().Fatal("")
		}
		for _, testOk := range suggestion.TestOK {
			if !regex.MatchString(testOk) {
				log.Default().Fatalf("%s FAIL with %s.",
					suggestion.Matcher,
					testOk)
			}
		}
		for _, testFail := range suggestion.TestFail {
			if regex.MatchString(testFail) {
				log.Default().Fatalf("%s FAIL with %s.",
					suggestion.Matcher,
					testFail)
			}
		}
		outputSuggestions = append(outputSuggestions, OutputSuggestions{
			suggestion.Matcher, suggestion.Packages})
	}
	output, err := json.Marshal(outputSuggestions)
	if err != nil {
		log.Default().Fatal(err)
	}
	fmt.Print(string(output))
}
