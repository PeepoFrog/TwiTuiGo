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
	gamesCursor := ""
	var authToTwitch model.AuthToTwitch
	authToTwitch.ClientID = os.Getenv("ClientID")
	authToTwitch.ClientSecret = os.Getenv("TwitchToken")
	authToTwitch.AccessToken = os.Getenv("AccessToken")

	games, gamesCursor := controller.GetGames(&authToTwitch, gamesCursor)
	printBroadcastsResponse(&games)
	games, gamesCursor = controller.GetGames(&authToTwitch, gamesCursor)
	printBroadcastsResponse(&games)

}
func printBroadcastsResponse(games *model.Games) {
	for a, b := range games.Data {
		fmt.Println(
			a,
			"ID:", b.ID,
			"Name:", b.Name,
			"IGBDid:", b.IGBDid)
	}
	fmt.Println(games.Pagination.Cursor)
}
