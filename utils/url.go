package utils

import "net/http"

func GetUrlHealthy(url string) bool {
	response, _ := http.Get(url)
	return response.StatusCode == 200
}
