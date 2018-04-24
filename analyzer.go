package json_log_monitoring

import (
	"regexp"
	"errors"
	"encoding/json"
)

type analyzer struct {
	countingRegex *regexp.Regexp
	counter map[string]int
	alwaysValid bool
}

func CreateAnalyzer(countingRegex string) *analyzer {
	a := &analyzer{
		countingRegex: regexp.MustCompile(countingRegex),
		counter: make(map[string]int),
		alwaysValid: true,
	}

	return a
}

func (a *analyzer) Analyze(data string) error {
	var jsonStruct map[string]interface{}
	if json.Unmarshal([]byte(data), &jsonStruct) != nil {
		a.alwaysValid = false
		return errors.New("invalid data")
	}

	submatch := a.countingRegex.FindStringSubmatch(data)
	for _, key := range submatch {
		a.counter[key] = a.counter[key] + 1
	}

	return nil
}

func (a *analyzer) ResetCounting() {
	a.alwaysValid = true
	a.counter = make(map[string]int)
}

func (a *analyzer) getResult() (bool, map[string]int) {
	return a.alwaysValid, a.counter
}
