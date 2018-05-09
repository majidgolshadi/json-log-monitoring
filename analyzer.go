package json_log_monitoring

import (
	"errors"
	"encoding/json"
)

type analyzer struct {
	counter map[string]int
	alwaysValid bool
}

func CreateAnalyzer() *analyzer {
	a := &analyzer{
		counter: make(map[string]int),
		alwaysValid: true,
	}

	return a
}

func (a *analyzer) Analyze(data []byte) error {
	var jsonStruct map[string]string
	if json.Unmarshal(data, &jsonStruct) != nil {
		a.alwaysValid = false
		return errors.New("invalid data")
	}

	key := string(jsonStruct["eirName"])
	a.counter[key] = a.counter[key] + 1

	return nil
}

func (a *analyzer) ResetCounting() {
	a.alwaysValid = true
	a.counter = make(map[string]int)
}

func (a *analyzer) getResult() (bool, map[string]int) {
	return a.alwaysValid, a.counter
}
