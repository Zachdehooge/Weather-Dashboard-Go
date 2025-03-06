package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/charmbracelet/lipgloss"
)

type Alert struct {
	Properties struct {
		Event       string `json:"event"`
		Headline    string `json:"headline"`
		Description string `json:"description"`
		Effective   string `json:"effective"`
		Expires     string `json:"expires"`
		Term        string `json:"term"`
		Definition  string `json:"definition"`
		SenderName  string `json:"senderName"`
	} `json:"properties"`
}

type AlertsResponse struct {
	Features []Alert `json:"features"`
}

type GeocodeResponse struct {
	Latt  string `json:"latt"`
	Longt string `json:"longt"`
}

var pointsData struct {
	Properties struct {
		ForecastOffice string `json:"cwa"`
		Radar          string `json:"radarStation"`
		GridX          int    `json:"gridX"`
		GridY          int    `json:"gridY"`
	} `json:"properties"`
}

var forecastData struct {
	Properties struct {
		Periods []struct {
			Name          string `json:"name"`
			Temperature   int    `json:"temperature"`
			ShortForecast string `json:"shortForecast"`
		} `json:"periods"`
	} `json:"properties"`
}

func AllAlerts() {
	resp, err := http.Get("https://api.weather.gov/alerts/active?status=actual&message_type=alert&&urgency=Immediate,Expected,Future&region_type=land&severity=Extreme,Severe,Moderate&limit=500")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	var data AlertsResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		fmt.Println("Error:", err)
		return
	}

	currentTime := time.Now()
	formattedTime := currentTime.Format("01-02-2006 | 15:04:05")

	currentTime2 := time.Now()
	formattedTime2 := currentTime2.Format("01-02-2006_15:04:05")

	currentUTC := currentTime.UTC()
	formattedTimeUTC := currentUTC.Format("15:04:05")

	// Open a file to write the output
	file, err := os.Create(fmt.Sprintf("All_Alerts_%s.txt", formattedTime2))
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	// Create a logger that writes to the file
	logger := log.New(file, "", 0)

	logger.Printf("ALL ALERTS: [%s Local | %s UTC]", formattedTime, formattedTimeUTC)
	logger.Println("---------------------------------------------------------------------------")

	for _, alert := range data.Features {

		logger.Println(alert.Properties.SenderName)
		logger.Print()
		logger.Println(alert.Properties.Event)
		logger.Println(alert.Properties.Description)

		//fmt.Println(alert.Properties.Event)

		//fmt.Println(alert.Properties.Description)

		effectiveTime, err := time.Parse(time.RFC3339, alert.Properties.Effective)
		if err != nil {
			fmt.Println("Error parsing effective time:", err)
			continue
		}

		utcExpiresTime1 := effectiveTime.UTC()
		readableEffectiveTime := utcExpiresTime1.Format("\nEffective: January 02, 2006, 03:04 PM MST")

		logger.Println(readableEffectiveTime)

		//fmt.Println(readableEffectiveTime)

		expiresTime, err := time.Parse(time.RFC3339, alert.Properties.Expires)
		if err != nil {
			fmt.Println("Error parsing expires time:", err)
			continue
		}

		utcExpiresTime := expiresTime.UTC()
		readableExpiresTime := utcExpiresTime.Format("Expires: January 02, 2006, 03:04 PM MST")

		logger.Println(readableExpiresTime + "\n ---------------------------------------------------------------------------")
		//fmt.Println(readableExpiresTime)

	}
}

func StateAlerts() {

	var state string
	fmt.Print("Enter the state abbreviation (e.g. UT) or 4 to exit: ")
	fmt.Scanln(&state)
	if state != "4" {
		url := fmt.Sprintf("https://api.weather.gov/alerts/active/area/%s", state)
		resp, err := http.Get(url)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				fmt.Println("Error:", err)
			}
		}(resp.Body)

		var data AlertsResponse
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			fmt.Println("Error:", err)
			return
		}

		if len(data.Features) == 0 {
			fmt.Println("---------------------------------------------------------------------------")
			fmt.Println(lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#FF0000")).
				Render("No Events For This State"))
			fmt.Println("---------------------------------------------------------------------------")
		}

		for _, alert := range data.Features {

			fmt.Println("---------------------------------------------------------------------------")

			eventStyle := lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#5a02ff")).
				Render(alert.Properties.Event)

			fmt.Println(eventStyle)
			fmt.Println()
			//fmt.Println(alert.Properties.Headline)
			//fmt.Println(alert.Properties.Description)
			descriptionStyle := lipgloss.NewStyle().
				Width(50).
				Height(10).
				Border(lipgloss.RoundedBorder()).
				Padding(1, 2).
				Render(alert.Properties.Description)

			fmt.Println(descriptionStyle)
			fmt.Println()

			effectiveTime, err := time.Parse(time.RFC3339, alert.Properties.Effective)
			if err != nil {
				fmt.Println("Error parsing effective time:", err)
				continue
			}
			readableEffectiveTime := effectiveTime.Format("January 02, 2006, 03:04 PM MST")
			fmt.Println(lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#67ff02")).
				Render("Issued: " + readableEffectiveTime))

			expiresTime, err := time.Parse(time.RFC3339, alert.Properties.Expires)
			if err != nil {
				fmt.Println("Error parsing expires time:", err)
				continue
			}
			readableExpiresTime := expiresTime.Format("January 02, 2006, 03:04 PM MST")
			fmt.Println()
			fmt.Println(lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#FF0000")).
				Render("Expires: " + readableExpiresTime))

			fmt.Println("---------------------------------------------------------------------------")
		}
		StateAlerts()
	} else {
		return
	}
}

func Forecast() {

	// TODO: 1. Grab the coords passed from the user // 2. Grab the SPC outlook for day 1 and check the coords to determine if they are in a category, if so, return the value, otherwise declare no thunderstorms for that area today

	var cityState, cityState1 string

	fmt.Print("Enter City (e.g. Ringgold): ")
	fmt.Scan(&cityState)

	fmt.Print("Enter State (e.g. GA)): ")
	fmt.Scan(&cityState1)

	result := fmt.Sprintf("%s,%s", cityState, cityState1)

	baseURL := "https://geocode.xyz"
	params := url.Values{}
	params.Add("locate", result)
	params.Add("region", "US")
	params.Add("json", "1")

	reqURL := fmt.Sprintf("%s/?%s", baseURL, params.Encode())
	resp, err := http.Get(reqURL)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	var geocodeData GeocodeResponse
	if err := json.Unmarshal(body, &geocodeData); err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Uncomment to debug coords being passed
	//fmt.Printf("\nLatitude: %s, Longitude: %s\n\n", geocodeData.Latt, geocodeData.Longt)

	pointsURL := fmt.Sprintf("https://api.weather.gov/points/%s,%s", geocodeData.Latt, geocodeData.Longt)
	pointsResp, err := http.Get(pointsURL)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer pointsResp.Body.Close()

	if err := json.NewDecoder(pointsResp.Body).Decode(&pointsData); err != nil {
		fmt.Println("Error:", err)
		return
	}

	officeURL := fmt.Sprintf("https://api.weather.gov/offices/%s", pointsData.Properties.ForecastOffice)
	officeResp, err := http.Get(officeURL)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer officeResp.Body.Close()

	var officeData struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(officeResp.Body).Decode(&officeData); err != nil {
		fmt.Println("Error:", err)
		return
	}

	forecastStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#00BFFF")).
		Border(lipgloss.RoundedBorder()).
		Render(fmt.Sprintf("Forecast Office: %s | Nearest Radar Station: %s", officeData.Name, pointsData.Properties.Radar))
	fmt.Println(forecastStyle)

	forecastURL := fmt.Sprintf("https://api.weather.gov/gridpoints/%s/%d,%d/forecast", pointsData.Properties.ForecastOffice, pointsData.Properties.GridX, pointsData.Properties.GridY)
	forecastResp, err := http.Get(forecastURL)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer forecastResp.Body.Close()

	if err := json.NewDecoder(forecastResp.Body).Decode(&forecastData); err != nil {
		fmt.Println("Error:", err)
		return
	}

	for _, period := range forecastData.Properties.Periods {
		periodStyle := lipgloss.NewStyle().
			Bold(true).
			Border(lipgloss.RoundedBorder()).
			Foreground(lipgloss.Color("#00BFFF")).
			Render(fmt.Sprintf("%s: %d°F, %s\n", period.Name, period.Temperature, period.ShortForecast))

		fmt.Println(periodStyle)
	}

	var choice string

	fmt.Print("Do you want to look at another forecast? (y/n): ")
	fmt.Scanln(&choice)

	if choice == "y" {
		fmt.Println("---------------------------------------------------------------------------")
		Forecast()
	} else {
		fmt.Println("---------------------------------------------------------------------------")
	}

}
