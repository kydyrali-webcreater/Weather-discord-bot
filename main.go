package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/go-resty/resty/v2"
)

const token = "MTIwMzY0NzQ2NjIwMTY3Nzg1NQ.GkHEuQ.LTyrqSHHrBRof3vpTebzccherqQxDEf6vdp3lc"
const weatherAPIKey = "e15fd09991adb196569025d0975d0cd4"
const weatherAPIEndpoint = "https://api.openweathermap.org/data/2.5/weather"

func main() {
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("Error creating Discord session,", err)
		return
	}

	dg.AddHandler(messageCreate)

	err = dg.Open()
	if err != nil {
		fmt.Println("Error opening Discord session,", err)
		return
	}

	fmt.Println("Bot is now running. Press Ctrl+C to exit.")
	select {}
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.HasPrefix(m.Content, "!weather") {
		go func() {
			// Extract location from the command
			location := strings.TrimSpace(strings.TrimPrefix(m.Content, "!weather"))

			// Fetch weather information
			weatherInfo, err := getWeatherInfo(location)
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "Error fetching weather information.")
				return
			}

			// Display weather information
			response := fmt.Sprintf("Weather in %s:\nTemperature: %.2f°C\nDescription: %s",
				location, weatherInfo.Main.Temp-273.15, weatherInfo.Weather[0].Description)
			s.ChannelMessageSend(m.ChannelID, response)
		}()
	} else if strings.HasPrefix(m.Content, "!help") {
		// Display help command
		help := "Available commands:\n" +
			"!weather [location] - Get current weather information for the specified location.\n" +
			"!help - Display this help message."
		s.ChannelMessageSend(m.ChannelID, help)
	}
}

func getWeatherInfo(location string) (*WeatherResponse, error) {
	client := resty.New()

	resp, err := client.R().
		SetQueryParam("q", location).
		SetQueryParam("appid", weatherAPIKey).
		Get(weatherAPIEndpoint)

	if err != nil {
		log.Printf("Error making request to weather API: %v", err)
		return nil, err
	}

	var weatherInfo WeatherResponse
	err = json.Unmarshal(resp.Body(), &weatherInfo)
	if err != nil {
		log.Printf("Error decoding weather API response: %v", err)
		return nil, err
	}

	return &weatherInfo, nil
}

// WeatherResponse represents the structure of the weather API response
type WeatherResponse struct {
	Main struct {
		Temp float64 `json:"temp"`
	} `json:"main"`
	Weather []struct {
		Description string `json:"description"`
	} `json:"weather"`
}
