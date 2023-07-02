package server

import "strconv"

func getParamInt(param string, defaultValue int) (int, error) {
	if len(param) == 0 {
		return defaultValue, nil
	}
	return strconv.Atoi(param)
}
