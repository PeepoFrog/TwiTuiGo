package main

import (
	"fmt"
	"log"
	"os"

	"github.com/PeepoFrog/TwiTuiGo/internal/controller"
	"github.com/PeepoFrog/TwiTuiGo/internal/model"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println(err)
		log.Fatalf("Error loading .env file")
	}
	var authToTwitch model.AuthToTwitch
	authToTwitch.ClientID = os.Getenv("ClientID")
	authToTwitch.ClientSecret = os.Getenv("TwitchToken")
	authToTwitch.AccessToken = os.Getenv("AccessToken")

	// gamesCursor := ""
	// games, gamesCursor := controller.GetGames(&authToTwitch, gamesCursor)
	// printGamesResponse(&games)
	// games, gamesCursor = controller.GetGames(&authToTwitch, gamesCursor)
	// printGamesResponse(&games)

	streams, cursor := controller.GetStreamsFromSelectedGame(&authToTwitch, "", "29595")
	printStreamsResponse(&streams)
	streams, _ = controller.GetStreamsFromSelectedGame(&authToTwitch, cursor, "29595")
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
