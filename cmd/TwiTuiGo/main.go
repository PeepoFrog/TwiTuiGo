package main

import (
	"fmt"
	"github.com/PeepoFrog/TwiTuiGo/internal/controller"
	"github.com/PeepoFrog/TwiTuiGo/internal/model"
	// "github.com/PeepoFrog/TwiTuiGo/internal/tui/tview"
	"github.com/PeepoFrog/TwiTuiGo/internal/tui/bubbletea"
	// "github.com/PeepoFrog/TwiTuiGo/internal/tui/tuiTesting"
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
	authToTwitch.ClientID = os.Getenv("ClientID")
	authToTwitch.ClientSecret = os.Getenv("TwitchToken")
	authToTwitch.AccessToken = os.Getenv("AccessToken")
	// usingtest()
	// tviewUserInterface.AuthToTwitch = authToTwitch
	// tviewUserInterface.Run()
	bubbleteaUserInterface.AuthToTwitch = authToTwitch

	bubbleteaUserInterface.Run()
	// tuiTesting.AuthToTwitch = authToTwitch
	// tuiTesting.RunV2()

}

func usingtest() {
	gamesCursor := ""
	streamsCursor := ""
	games, _ := controller.GetGames(&authToTwitch, gamesCursor)
	gamesCursor = games.Pagination.Cursor
	printGamesResponse(&games)
	games, _ = controller.GetGames(&authToTwitch, gamesCursor)
	gamesCursor = games.Pagination.Cursor
	printGamesResponse(&games)
	streams, _ := controller.GetStreamsFromSelectedGame(&authToTwitch, streamsCursor, "29595")
	streamsCursor = streams.Pagination.Cursor
	printStreamsResponse(&streams)
	streams, _ = controller.GetStreamsFromSelectedGame(&authToTwitch, streamsCursor, "29595")
	streamsCursor = streams.Pagination.Cursor
	printStreamsResponse(&streams)
}
func printGamesResponse(games *model.Games) {
	for a, b := range games.Data {
		fmt.Println(
			a,
			"ID:", b.ID,
			"Name:", b.Name,
			"IGBDid:", b.IGBDid)
	}
	fmt.Println(games.Pagination.Cursor)
}
func printStreamsResponse(games *model.Streamers) {
	for a, b := range games.Data {
		fmt.Println(
			a,
			"gameID:", b.GameID,
			"Game:", b.GameName,
			"UserName:", b.UserName,
			"Vievers", b.ViewerCount)
	}
	fmt.Println(games.Pagination.Cursor)
}
