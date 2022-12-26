package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/PeepoFrog/TwiTuiGo/internal/model"
)

// todo
// Створити функціх для викликів твітч апі для
// 1. Вибір топ 40 ігор присвоєння їх до слайсу з можливістю приєжнання наступних 30 ігор
func GetGames(authToTwitch *model.AuthToTwitch, cursor string) (model.Games, string) {
	client := http.Client{}
	req, err := http.NewRequest("GET", "https://api.twitch.tv/helix/games/top?after="+cursor, nil)
	if err != nil {
		log.Fatalln(err)
	}
	req.Header = http.Header{
		"Authorization": {authToTwitch.AccessToken},
		"Client-Id":     {authToTwitch.ClientID},
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	var games model.Games
	json.NewDecoder(resp.Body).Decode(&games)
	cursor = games.Pagination.Cursor
	defer resp.Body.Close()
	return games, cursor
}

// 2. Вибір в ігровій категорії списка стрімерів та присвоєння до слайсу + приєднання наступних
func GetStreamsFromSelectedGame(authToTwitch *model.AuthToTwitch, cursor, gameID string) (model.Streamers, string) {
	streams, cursor := getStreams(authToTwitch, "", "", gameID, "", cursor)
	return streams, cursor
}
func getStreams(authToTwitch *model.AuthToTwitch, userID, userLogin, gameID, allOrAlive, cursor string) (model.Streamers, string) {
	client := http.Client{}
	var apiURL = "https://api.twitch.tv/helix/streams?"
	if userID != "" {
		apiURL += "&user_id=" + userID
	}
	if userLogin != "" {
		apiURL += "&user_login=" + userLogin
	}
	if gameID != "" {
		apiURL += "&game_id=" + gameID
	}
	if allOrAlive != "" {
		apiURL += "&type=" + allOrAlive
	}
	if cursor != "" {
		apiURL += "&after=" + cursor
	}
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(apiURL)
	req.Header = http.Header{
		"Authorization": {authToTwitch.AccessToken},
		"Client-Id":     {authToTwitch.ClientID},
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	var streams model.Streamers
	json.NewDecoder(resp.Body).Decode(&streams)
	return streams, streams.Pagination.Cursor
}

// 3. Чи можливо зробити виклик для получення інформації про стрім який онлайн конкретно знайти допступну якість

// 4. Пошук стрімера
