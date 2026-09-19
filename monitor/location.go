package monitor

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// locationMonitor periodically queries an IP geolocation API.
type locationMonitor struct {
	callback func(string)
	lastIP   string
}

// newLocationMonitor creates a new location monitoring instance.
func newLocationMonitor(callback func(string)) *locationMonitor {
	return &locationMonitor{
		callback: callback,
	}
}

func (lm *locationMonitor) start() {
	// Check immediately on startup
	lm.checkLocation()
	for {
		time.Sleep(15 * time.Minute)
		lm.checkLocation()
	}
}

type geoResponse struct {
	Status      string  `json:"status"`
	Query       string  `json:"query"`
	Country     string  `json:"country"`
	RegionName  string  `json:"regionName"`
	City        string  `json:"city"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
}

func (lm *locationMonitor) checkLocation() {
	client := http.Client{
		Timeout: 15 * time.Second,
	}

	resp, err := client.Get("http://ip-api.com/json/")
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	var geo geoResponse
	if err := json.Unmarshal(body, &geo); err != nil {
		return
	}

	if geo.Status != "success" || geo.Query == "" {
		return
	}

	// Only log if IP has changed (meaning network/location might have changed)
	// to avoid spamming the logs with the exact same location every 15 minutes.
	if geo.Query != lm.lastIP || lm.lastIP == "" {
		lm.lastIP = geo.Query
		ts := time.Now().Format("2006-01-02 15:04:05")
		url := fmt.Sprintf("https://maps.google.com/?q=%f,%f", geo.Lat, geo.Lon)
		logLine := fmt.Sprintf("[%s] [LOCATION] 📍 Lat: %f, Lng: %f | City: %s, %s, %s | %s",
			ts, geo.Lat, geo.Lon, geo.City, geo.RegionName, geo.Country, url)
		
		if lm.callback != nil {
			lm.callback(logLine)
		}
	}
}
