package bubbleteaUserInterface

import (
	"fmt"
	"github.com/PeepoFrog/TwiTuiGo/internal/controller"
	models "github.com/PeepoFrog/TwiTuiGo/internal/model"
	tea "github.com/charmbracelet/bubbletea"
	"os"
)

var AuthToTwitch models.AuthToTwitch
var GameListCursor = ""

type model struct {
	GamesData []models.Game
	// GameData    models.Game
	// GamesStruct models.Games
	choices  []string
	cursor   int
	selected map[int]struct{}
	err      error
}
type (
	// gameList models.Games
	gameList []models.Game
	errMsg   struct{ err error }
)

func (e errMsg) Error() string { return e.err.Error() }
func getGames() tea.Msg {
	games, err := controller.GetGames(&AuthToTwitch, GameListCursor)
	if err != nil {
		return errMsg{err}
	}
	GameListCursor = games.Pagination.Cursor
	return gameList(games.Data)
}
func initialModel() model {

	games, _ := controller.GetGames(&AuthToTwitch, GameListCursor)
	var c []string
	for _, b := range games.Data {
		c = append(c, b.Name)
	}
	return model{
		GamesData: games.Data,
		choices:   c,
		selected:  make(map[int]struct{}),
	}
}
func (m model) Init() tea.Cmd {
	return getGames
}
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case gameList:
		m.GamesData = msg
		// fmt.Println("worked")
		return m, nil

	case errMsg:
		m.err = msg
		return nil, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter", " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		}
	}
	return m, nil
}
func (m model) View() string {
	if m.err != nil {
		return fmt.Sprintf("\nWe had some trouble: %v\n\n", m.err)
	}
	s := "games \n \n"
	cursor := ">"
	checked := "x"
	if len(m.GamesData) > 1 {
		for i, choise := range m.GamesData {
			cursor = " "
			if m.cursor == i {
				cursor = ">"
			}
			checked = " "
			if _, ok := m.selected[i]; ok {
				checked = "x"
			}
			s += fmt.Sprintf("%s [%s] %s \n", cursor, checked, choise.Name)

		}
	}
	s += fmt.Sprintf("%s [%s] load more \n", cursor, checked)

	s += "\nPress q to exit"
	return s

}
func Run() {
	// p := tea.NewProgram(initialModel())
	// if err := p.Start(); err != nil {
	// 	fmt.Println("err", err)
	// 	os.Exit(1)
	// }
	if _, err := tea.NewProgram(model{}).Run(); err != nil {
		fmt.Printf("Uh oh, there was an error: %v\n", err)
		os.Exit(1)
	}
}
