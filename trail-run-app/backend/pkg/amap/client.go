package amap

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

type Client struct {
	apiKey string
	client *http.Client
}

func NewClient() *Client {
	apiKey := os.Getenv("AMAP_KEY")
	if apiKey == "" {
		apiKey = "demo_key"
	}

	return &Client{
		apiKey: apiKey,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// WeatherResponse 高德天气API响应
type WeatherResponse struct {
	Status   string `json:"status"`
	Info     string `json:"info"`
	Lives    []Live `json:"lives"`
	Forecast []DayForecast `json:"forecasts"`
}

type Live struct {
	Province     string `json:"province"`
	City        string `json:"city"`
	Weather     string `json:"weather"`
	Temperature string `json:"temperature"`
	WindDirection string `json:"winddirection"`
	WindPower   string `json:"windpower"`
	Humidity    string `json:"humidity"`
}

type DayForecast struct {
	Date        string `json:"date"`
	Week       string `json:"week"`
	DayWeather string `json:"dayweather"`
	NightWeather string `json:"nightweather"`
	DayTemp    string `json:"daytemp"`
	NightTemp  string `json:"nighttemp"`
	DayWind    string `json:"daywind"`
	NightWind  string `json:"nightwind"`
	DayPower   string `json:"daypower"`
	NightPower string `json:"nightpower"`
}

// GetWeather 获取实时天气
func (c *Client) GetWeather(cityCode string) (*Live, error) {
	apiURL := "https://restapi.amap.com/v3/weather/weatherInfo"

	params := url.Values{}
	params.Set("key", c.apiKey)
	params.Set("city", cityCode)
	params.Set("extensions", "base")

	resp, err := c.client.Get(apiURL + "?" + params.Encode())
	if err != nil {
		return nil, fmt.Errorf("failed to get weather: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var result WeatherResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if result.Status != "1" || len(result.Lives) == 0 {
		return nil, fmt.Errorf("API error: %s", result.Info)
	}

	return &result.Lives[0], nil
}

// GetForecast 获取天气预报
func (c *Client) GetForecast(cityCode string) ([]DayForecast, error) {
	apiURL := "https://restapi.amap.com/v3/weather/weatherInfo"

	params := url.Values{}
	params.Set("key", c.apiKey)
	params.Set("city", cityCode)
	params.Set("extensions", "all")

	resp, err := c.client.Get(apiURL + "?" + params.Encode())
	if err != nil {
		return nil, fmt.Errorf("failed to get forecast: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var result WeatherResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if result.Status != "1" {
		return nil, fmt.Errorf("API error: %s", result.Info)
	}

	return result.Forecast, nil
}

// RegeoResponse 地理编码响应
type RegeoResponse struct {
	Status   string `json:"status"`
	Info     string `json:"info"`
	Regeocode RegeoCode `json:"regeocode"`
}

type RegeoCode struct {
	AddressComponent AddressComponent `json:"addressComponent"`
	POIs            []POI `json:"pois"`
}

type AddressComponent struct {
	Province    string `json:"province"`
	City        string `json:"city"`
	District    string `json:"district"`
	 Township   string `json:"township"`
}

type POI struct {
	Name     string `json:"name"`
	Location string `json:"location"`
	Distance string `json:"distance"`
	Type     string `json:"type"`
}

// GetNearbyPOIs 获取附近POI（停车场等）
func (c *Client) GetNearbyPOIs(lat, lng float64, keywords string) ([]POI, error) {
	apiURL := "https://restapi.amap.com/v3/place/around"

	params := url.Values{}
	params.Set("key", c.apiKey)
	params.Set("location", fmt.Sprintf("%f,%f", lng, lat)) // 高德是 lng,lat 顺序
	params.Set("keywords", keywords)
	params.Set("types", "150停车场")
	params.Set("radius", "3000")
	params.Set("offset", "10")

	resp, err := c.client.Get(apiURL + "?" + params.Encode())
	if err != nil {
		return nil, fmt.Errorf("failed to get nearby POIs: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var result struct {
		Status  string `json:"status"`
		Info    string `json:"info"`
		POIs    []POI `json:"pois"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if result.Status != "1" {
		return nil, fmt.Errorf("API error: %s", result.Info)
	}

	return result.POIs, nil
}

// ConvertAddress 地理编码（地址转坐标）
func (c *Client) ConvertAddress(address string) (lat, lng float64, err error) {
	apiURL := "https://restapi.amap.com/v3/geocode/geo"

	params := url.Values{}
	params.Set("key", c.apiKey)
	params.Set("address", address)

	resp, err := c.client.Get(apiURL + "?" + params.Encode())
	if err != nil {
		return 0, 0, fmt.Errorf("failed to geocode: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to read response: %w", err)
	}

	var result struct {
		Status  string `json:"status"`
		Geocodes []struct {
			Location string `json:"location"`
		} `json:"geocodes"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return 0, 0, fmt.Errorf("failed to parse response: %w", err)
	}

	if result.Status != "1" || len(result.Geocodes) == 0 {
		return 0, 0, fmt.Errorf("geocode failed")
	}

	// 解析 location (lng,lat)
	var l1, l2 float64
	fmt.Sscanf(result.Geocodes[0].Location, "%f,%f", &l1, &l2)
	return l2, l1, nil // 返回 lat, lng
}
