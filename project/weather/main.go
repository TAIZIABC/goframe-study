package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// wttr.in JSON API 响应结构
type WttrResponse struct {
	CurrentCondition []CurrentCondition `json:"current_condition"`
	NearestArea      []NearestArea      `json:"nearest_area"`
}

type CurrentCondition struct {
	TempC         string        `json:"temp_C"`
	TempF         string        `json:"temp_F"`
	FeelsLikeC    string        `json:"FeelsLikeC"`
	Humidity      string        `json:"humidity"`
	WindspeedKmph string        `json:"windspeedKmph"`
	WindDir       string        `json:"winddir16Point"`
	Pressure      string        `json:"pressure"`
	Visibility    string        `json:"visibility"`
	UVIndex       string        `json:"uvIndex"`
	WeatherDesc   []WeatherDesc `json:"weatherDesc"`
	WeatherIconURL []WeatherDesc `json:"weatherIconUrl"`
}

type WeatherDesc struct {
	Value string `json:"value"`
}

type NearestArea struct {
	AreaName []AreaValue `json:"areaName"`
	Country  []AreaValue `json:"country"`
	Region   []AreaValue `json:"region"`
}

type AreaValue struct {
	Value string `json:"value"`
}

// weatherEmoji 根据天气描述返回对应的 emoji
func weatherEmoji(desc string) string {
	d := strings.ToLower(desc)
	switch {
	case strings.Contains(d, "sunny"), strings.Contains(d, "clear"):
		return "☀️"
	case strings.Contains(d, "partly cloudy"):
		return "⛅"
	case strings.Contains(d, "cloudy"), strings.Contains(d, "overcast"):
		return "☁️"
	case strings.Contains(d, "rain"), strings.Contains(d, "drizzle"):
		return "🌧️"
	case strings.Contains(d, "thunder"):
		return "⛈️"
	case strings.Contains(d, "snow"):
		return "❄️"
	case strings.Contains(d, "fog"), strings.Contains(d, "mist"), strings.Contains(d, "haze"):
		return "🌫️"
	case strings.Contains(d, "wind"):
		return "💨"
	default:
		return "🌤️"
	}
}

func fetchWeather(city string) (*WttrResponse, error) {
	apiURL := fmt.Sprintf("https://wttr.in/%s?format=j1", url.PathEscape(city))

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "go-weather-cli/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("网络请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API 返回状态码: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	var data WttrResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("解析数据失败: %w", err)
	}

	if len(data.CurrentCondition) == 0 {
		return nil, fmt.Errorf("未找到 \"%s\" 的天气数据", city)
	}

	return &data, nil
}

func displayWeather(data *WttrResponse) {
	c := data.CurrentCondition[0]

	// 地点信息
	location := ""
	if len(data.NearestArea) > 0 {
		area := data.NearestArea[0]
		parts := []string{}
		if len(area.AreaName) > 0 && area.AreaName[0].Value != "" {
			parts = append(parts, area.AreaName[0].Value)
		}
		if len(area.Region) > 0 && area.Region[0].Value != "" {
			parts = append(parts, area.Region[0].Value)
		}
		if len(area.Country) > 0 && area.Country[0].Value != "" {
			parts = append(parts, area.Country[0].Value)
		}
		location = strings.Join(parts, ", ")
	}

	// 天气描述
	desc := "N/A"
	if len(c.WeatherDesc) > 0 {
		desc = c.WeatherDesc[0].Value
	}
	emoji := weatherEmoji(desc)

	fmt.Println()
	fmt.Println("  ┌────────────────────────────────────────┐")
	fmt.Printf("  │  📍 %s\n", location)
	fmt.Println("  ├────────────────────────────────────────┤")
	fmt.Printf("  │  %s  %s\n", emoji, desc)
	fmt.Println("  ├────────────────────────────────────────┤")
	fmt.Printf("  │  🌡️  温度:     %s°C (体感 %s°C)\n", c.TempC, c.FeelsLikeC)
	fmt.Printf("  │  💧 湿度:     %s%%\n", c.Humidity)
	fmt.Printf("  │  💨 风速:     %s km/h (%s)\n", c.WindspeedKmph, c.WindDir)
	fmt.Printf("  │  👁️  能见度:   %s km\n", c.Visibility)
	fmt.Printf("  │  🔽 气压:     %s hPa\n", c.Pressure)
	fmt.Printf("  │  🔆 紫外线:   %s\n", c.UVIndex)
	fmt.Println("  └────────────────────────────────────────┘")
}

func main() {
	fmt.Println("╔════════════════════════════════════════╗")
	fmt.Println("║       天气查询工具 (weather-cli)       ║")
	fmt.Println("╠════════════════════════════════════════╣")
	fmt.Println("║  输入城市名查询天气 (支持中英文)       ║")
	fmt.Println("║  示例: Beijing / 上海 / Tokyo / London ║")
	fmt.Println("║  输入 q 退出                           ║")
	fmt.Println("╚════════════════════════════════════════╝")

	// 如果有命令行参数，直接查询
	if len(os.Args) > 1 {
		city := strings.Join(os.Args[1:], " ")
		fmt.Printf("\n  查询中: %s ...\n", city)
		data, err := fetchWeather(city)
		if err != nil {
			fmt.Fprintf(os.Stderr, "  ✗ %v\n", err)
			os.Exit(1)
		}
		displayWeather(data)
		return
	}

	// 交互模式
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("\n🏙️  城市: ")
		if !scanner.Scan() {
			break
		}

		city := strings.TrimSpace(scanner.Text())
		if city == "" {
			continue
		}
		if city == "q" || city == "quit" || city == "exit" {
			fmt.Println("再见！")
			break
		}

		fmt.Printf("  查询中: %s ...\n", city)
		data, err := fetchWeather(city)
		if err != nil {
			fmt.Printf("  ✗ %v\n", err)
			continue
		}
		displayWeather(data)
	}
}
