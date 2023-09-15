package server

import (
	"strconv"

	"golang.org/x/exp/constraints"
)

func getParam[T constraints.Integer](param string, defaultValue T) (T, error) {
	if len(param) == 0 {
		return defaultValue, nil
	}

	res, err := strconv.Atoi(param)
	if err != nil {
		return 0, err
	}

	return T(res), nil
}

type ID struct {
	ID string `json:"id"`
}

type Name struct {
	Name string `json:"name"`
}

type URL struct {
	URL string `json:"url"`
}
