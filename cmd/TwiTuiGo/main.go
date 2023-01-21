package main

import (
	"fmt"
	"github.com/PeepoFrog/TwiTuiGo/internal/model"
	// "github.com/PeepoFrog/TwiTuiGo/internal/tui/tview"
	"github.com/PeepoFrog/TwiTuiGo/internal/tui/bubbletea"
	// "github.com/PeepoFrog/TwiTuiGo/internal/apiTesting"
	"github.com/PeepoFrog/TwiTuiGo/internal/tui/tuiTesting"
	"log"
	"os"

	"github.com/joho/godotenv"
)

var authToTwitch model.AuthToTwitch

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println(err)
		log.Fatalf("Error loading .env file")
	}
	//
	authToTwitch.ClientID = os.Getenv("ClientID")
	authToTwitch.ClientSecret = os.Getenv("TwitchToken")
	authToTwitch.AccessToken = os.Getenv("AccessToken")
	//
	bubbleteaTUI.AuthToTwitch = authToTwitch
	tuiTesting.AuthToTwitch = authToTwitch
	bubbleteaTUI.Run()
	// tuiTesting.Run()

}
