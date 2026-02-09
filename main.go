package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path"
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
	Filename string    `json:"-"`
}

type OutputSuggestions struct {
	Matcher  string    `json:"matcher"`
	Packages []Package `json:"packages"`
}

func main() {
	entries, err := os.ReadDir("./matchers/")
	if err != nil {
		log.Default().Fatal(err)
	}
	testableSuggestions := []TestableSuggestions{}
	for _, e := range entries {
		jsonFile, err := os.Open(path.Join("./matchers/", e.Name()))
		if err != nil {
			log.Default().Fatal(err)
		}
		byteValue, err := io.ReadAll(jsonFile)
		if err != nil {
			log.Default().Fatal(err)
		}
		suggestion := TestableSuggestions{}
		err = json.Unmarshal(byteValue, &suggestion)
		if err != nil {
			log.Default().Fatal(err)
		}
		jsonFile.Close()
		suggestion.Filename = e.Name()
		testableSuggestions = append(testableSuggestions, suggestion)
	}
	var outputSuggestions []OutputSuggestions
	for _, suggestion := range testableSuggestions {
		regex, err := regexp.Compile(suggestion.Matcher)
		if err != nil {
			log.Default().Fatal(err)
		}
		if len(suggestion.TestOK) < 1 || len(suggestion.TestFail) < 1 {
			log.Default().Fatal("No tests...")
		}
		for _, testOk := range suggestion.TestOK {
			if !regex.MatchString(testOk) {
				log.Default().Fatalf("%s (%s) FAIL with %s.",
					suggestion.Filename,
					suggestion.Matcher,
					testOk)
			}
		}
		for _, testFail := range suggestion.TestFail {
			if regex.MatchString(testFail) {
				log.Default().Fatalf("%s (%s) FAIL with %s.",
					suggestion.Filename,
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
