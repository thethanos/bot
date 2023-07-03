package server

import "strconv"

func getParamInt(param string, defaultValue int) (int, error) {
	if len(param) == 0 {
		return defaultValue, nil
	}
	return strconv.Atoi(param)
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
