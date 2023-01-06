package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/PeepoFrog/TwiTuiGo/internal/model"
)

// todo
// Створити функціх для викликів твітч апі для
// 1. Вибір топ 40 ігор присвоєння їх до слайсу з можливістю приєжнання наступних 30 ігор
func GetGames(authToTwitch *model.AuthToTwitch, cursor string) (model.Games, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	var games model.Games

	req, err := http.NewRequest("GET", "https://api.twitch.tv/helix/games/top?after="+cursor, nil)
	if err != nil {
		log.Fatalln(err)
		return games, err
	}
	req.Header = http.Header{
		"Authorization": {authToTwitch.AccessToken},
		"Client-Id":     {authToTwitch.ClientID},
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return games, err
	}
	json.NewDecoder(resp.Body).Decode(&games)
	cursor = games.Pagination.Cursor
	defer resp.Body.Close()
	return games, nil
}

// 2. Вибір в ігровій категорії списка стрімерів та присвоєння до слайсу + приєднання наступних
func GetStreamsFromSelectedGame(authToTwitch *model.AuthToTwitch, cursor, gameID string) (model.Streamers, error) {
	streams, err := getStreams(authToTwitch, "", "", gameID, "", cursor)
	return streams, err
}
func getStreams(authToTwitch *model.AuthToTwitch, userID, userLogin, gameID, allOrAlive, cursor string) (model.Streamers, error) {
	client := &http.Client{Timeout: 10 * time.Second}
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
	var streams model.Streamers
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		log.Println(err)
		return streams, err
	}
	req.Header = http.Header{
		"Authorization": {authToTwitch.AccessToken},
		"Client-Id":     {authToTwitch.ClientID},
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return streams, err
	}
	// var streams model.Streamers
	json.NewDecoder(resp.Body).Decode(&streams)
	return streams, nil
}

// 3. Чи можливо зробити виклик для получення інформації про стрім який онлайн конкретно знайти допступну якість

// 4. Пошук стрімера
