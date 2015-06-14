package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type weatherData struct {
	Name string `json:"name"`
	Main struct {
		Kelvin float64 `json:"temp"`
	} `json:main`
}

func main() {
	fmt.Println("hello")
	http.HandleFunc("/", hello)
	http.HandleFunc("/weather/", weather)
	http.ListenAndServe(":54321", nil)
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello"))
}

func weather(w http.ResponseWriter, r *http.Request) {
	city := strings.SplitN(r.URL.Path, "/", 3)[2]
	fmt.Fprintln(w, "city:", city)
	data, err := query(city)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(*data)
}

func query(city string) (*weatherData, error) {
	resp, err := http.Get("http://api.openweathermap.org/data/2.5/weather?q=" + city)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	d := new(weatherData)

	if err := json.NewDecoder(resp.Body).Decode(d); err != nil {
		return nil, err
	}
	return d, nil
}
