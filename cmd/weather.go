package weather

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// TODO: Implement weather data fetch (PullWeather Func)
// TODO: Parse weather data
// TODO: Implement SLOGs

func PullWeather() {

	fmt.Println("Pulling weather data...")

	response, err := http.Get("https://api.weather.gov/openapi.json")
	if err != nil {
		fmt.Println("Error fetching weather data:", err)
		return
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	var data WeatherData
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return
	}

}
