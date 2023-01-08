package bubbleteaTUI

//6 _ % \ ` ~  ^
import (
	"fmt"
	"github.com/PeepoFrog/TwiTuiGo/internal/controller"
	myModels "github.com/PeepoFrog/TwiTuiGo/internal/model"
	"os"
	"strconv"

	"github.com/charmbracelet/bubbles/list"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var AuthToTwitch myModels.AuthToTwitch

type status int

const divisor = 3

const (
	gamesColumn status = iota
	broadcastsColumn
	favoritesColumn
)

/* MODEL MANAGEMENT */
var models []tea.Model

const (
	model status = iota
	form
)

/* STYLING */
var (
	columnStyle = lipgloss.NewStyle().
			Padding(1, 2).
			Border(lipgloss.HiddenBorder())
	focusedStyle = lipgloss.NewStyle().
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62"))
	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))
)

/* CUSTOM ITEM */

type Task struct {
	status      status
	title       string
	description string
	viewers     string

	gameStruct myModels.Game
}

func NewTask(status status, title, description string, gameid string, gamestruct myModels.Game) Task {
	return Task{status: status, title: title, description: description, gameStruct: gamestruct}
}

func (t *Task) Next() {
	if t.status == favoritesColumn {
		t.status = gamesColumn
	} else {
		t.status++
	}
}

// implement the list.Item interface
func (t Task) FilterValue() string {
	return t.title
}

func (t Task) Title() string {
	return t.title
}

func (t Task) Description() string {
	return t.description
}

/* MAIN MODEL */

type Model struct {
	loaded               bool
	focused              status
	lists                []list.Model
	err                  error
	quitting             bool
	gameStruct           myModels.Games
	gameList             []myModels.Game
	broadcastStruct      myModels.Streamers
	broadcastList        []myModels.Streamer
	SelectedGameInColumn myModels.Game
	gamesCursor          string
	gamesCursorState     bool
	broadcastsCursor     string
}

func New() *Model {
	return &Model{}
}
func (m *Model) LoadMoreGames() tea.Msg {
	print("worked")
	m.gamesCursorState = false
	gamesStruct, err := controller.GetGames(&AuthToTwitch, m.gamesCursor)
	if err != nil {
		panic(err)
	}
	m.gamesCursor = gamesStruct.Pagination.Cursor
	var listtoadd []list.Item
	m.gameList = append(m.gameList, gamesStruct.Data...)

	for _, b := range m.gameList {
		listtoadd = append(listtoadd, Task{status: gamesColumn, title: b.Name, description: "description", gameStruct: b})
	}
	listtoadd = append(listtoadd, Task{status: gamesColumn, title: "LOAD MORE"})
	m.lists[gamesColumn].SetItems(listtoadd)

	return nil
}
func (m *Model) LoadBroadcastsFromSelectedGame(id string) tea.Msg {
	if m.gamesCursorState != true {
		gameStruct, err := controller.GetStreamsFromSelectedGame(&AuthToTwitch, "", id)
		if err != nil {
			panic(err)
		}
		m.broadcastStruct = gameStruct
		var listtoadd []list.Item

		for _, b := range gameStruct.Data {
			v := strconv.Itoa(b.ViewerCount)
			listtoadd = append(listtoadd, Task{status: broadcastsColumn, title: b.UserName, description: b.GameName + " viewers: " + v})
		}

		m.lists[broadcastsColumn].SetItems(listtoadd)
	} else {
		m.LoadMoreGames()
	}
	return nil
}
func (m *Model) SelectGame() tea.Msg {
	selectedItem := m.lists[m.focused].SelectedItem().(Task)
	if selectedItem.title == "LOAD MORE" {
		m.gamesCursorState = true
		print("loadmorekekw")
	} else {
		print(selectedItem.gameStruct.Name + "game ID: " + selectedItem.gameStruct.ID + selectedItem.title)
		m.SelectedGameInColumn = selectedItem.gameStruct
	}
	return nil
}
func (m *Model) MoveToNext() tea.Msg {
	selectedItem := m.lists[m.focused].SelectedItem()
	selectedTask := selectedItem.(Task)
	m.lists[selectedTask.status].RemoveItem(m.lists[m.focused].Index())
	selectedTask.Next()
	m.lists[selectedTask.status].InsertItem(len(m.lists[selectedTask.status].Items())-1, list.Item(selectedTask))
	return nil
}

func (m *Model) DeleteCurrent() tea.Msg {
	if len(m.lists[m.focused].VisibleItems()) > 0 {
		selectedTask := m.lists[m.focused].SelectedItem().(Task)
		m.lists[selectedTask.status].RemoveItem(m.lists[m.focused].Index())
	}
	return nil
}

func (m *Model) Next() {
	if m.focused == favoritesColumn {
		m.focused = gamesColumn
	} else {
		m.focused++
	}
}

func (m *Model) Prev() {
	if m.focused == gamesColumn {
		m.focused = favoritesColumn
	} else {
		m.focused--
	}
}

func (m *Model) initLists(width, height int) {
	defaultList := list.New([]list.Item{}, list.NewDefaultDelegate(), width/divisor+5, height-7)
	defaultList.SetShowHelp(false)
	m.lists = []list.Model{defaultList, defaultList, defaultList}

	// Init To Do
	m.lists[gamesColumn].Title = "GAMES"
	gameslist, err := controller.GetGames(&AuthToTwitch, "")
	m.gamesCursor = gameslist.Pagination.Cursor
	m.gameStruct = gameslist
	m.gameList = gameslist.Data
	if err != nil {
		panic(err)
	}
	//init games column
	var listtoadd []list.Item
	for _, b := range gameslist.Data {
		listtoadd = append(listtoadd, Task{status: gamesColumn, title: b.Name, description: "description", gameStruct: b})
	}
	listtoadd = append(listtoadd, Task{status: gamesColumn, title: "LOAD MORE"})
	m.lists[gamesColumn].SetItems(listtoadd)
	// Init in broadcast column
	m.lists[broadcastsColumn].Title = "BROADCASTS"
	m.lists[broadcastsColumn].SetItems([]list.Item{
		Task{status: broadcastsColumn, title: "title", description: "description", viewers: "viewers"},
	})
	// Init favorites column
	m.lists[favoritesColumn].Title = "FAVORITES"

	m.lists[favoritesColumn].SetItems([]list.Item{
		Task{status: favoritesColumn, title: "title", description: "description"},
	})
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if !m.loaded {
			columnStyle.Width(msg.Width / divisor)
			focusedStyle.Width(msg.Width / divisor)
			columnStyle.Height(msg.Height - divisor)
			focusedStyle.Height(msg.Height - divisor)
			m.initLists(msg.Width, msg.Height)
			m.loaded = true
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		case "left", "h":
			m.Prev()
		case "right", "l":
			m.Next()
		case "enter":
			if m.focused == gamesColumn {
				m.SelectGame()
				id := m.SelectedGameInColumn.ID
				m.LoadBroadcastsFromSelectedGame(id)
				return m, nil
			}
			if m.focused == broadcastsColumn {
				return m, m.SelectGame
			}
			if m.focused == favoritesColumn {
				return m, m.MoveToNext
			}
			return m, m.MoveToNext
		case "n":
			models[model] = m // save the state of the current model
			// models[form] = NewForm(m.focused)
			return models[form].Update(nil)
		case "d":
			return m, m.DeleteCurrent
		}
	case Task:
		task := msg
		return m, m.lists[task.status].InsertItem(len(m.lists[task.status].Items()), task)
	}
	var cmd tea.Cmd
	m.lists[m.focused], cmd = m.lists[m.focused].Update(msg)
	return m, cmd
}

func (m Model) View() string {
	if m.quitting {
		return ""
	}
	if m.loaded {
		todoView := m.lists[gamesColumn].View()
		inProgView := m.lists[broadcastsColumn].View()
		doneView := m.lists[favoritesColumn].View()
		switch m.focused {
		case broadcastsColumn:
			return lipgloss.JoinHorizontal(
				lipgloss.Left,
				columnStyle.Render(todoView),
				focusedStyle.Render(inProgView),
				columnStyle.Render(doneView),
			)
		case favoritesColumn:
			return lipgloss.JoinHorizontal(
				lipgloss.Left,
				columnStyle.Render(todoView),
				columnStyle.Render(inProgView),
				focusedStyle.Render(doneView),
			)
		default:
			return lipgloss.JoinHorizontal(
				lipgloss.Left,
				focusedStyle.Render(todoView),
				columnStyle.Render(inProgView),
				columnStyle.Render(doneView),
			)
		}
	} else {
		return "loading..."
	}
}

/* FORM MODEL */
// type Form struct {
// 	focused     status
// 	title       textinput.Model
// 	description textarea.Model
// }

// func NewForm(focused status) *Form {
// 	form := &Form{focused: focused}
// 	form.title = textinput.New()
// 	form.title.Focus()
// 	form.description = textarea.New()
// 	return form
// }

// func (m Form) CreateTask() tea.Msg {
// 	task := NewTask(m.focused, m.title.Value(), m.description.Value())
// 	return task
// }

// func (m Form) Init() tea.Cmd {
// 	return nil
// }

// func (m Form) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
// 	var cmd tea.Cmd
// 	switch msg := msg.(type) {
// 	case tea.KeyMsg:
// 		switch msg.String() {
// 		case "ctrl+c", "q":
// 			return m, tea.Quit
// 		case "enter":
// 			if m.title.Focused() {
// 				m.title.Blur()
// 				m.description.Focus()
// 				return m, textarea.Blink
// 			} else {
// 				models[form] = m
// 				return models[model], m.CreateTask
// 			}
// 		}
// 	}
// 	if m.title.Focused() {
// 		m.title, cmd = m.title.Update(msg)
// 		return m, cmd
// 	} else {
// 		m.description, cmd = m.description.Update(msg)
// 		return m, cmd
// 	}
// }

// func (m Form) View() string {
// 	return lipgloss.JoinVertical(lipgloss.Left, m.title.View(), m.description.View())
// }

func Run() {
	// models = []tea.Model{New(), NewForm(games)}
	models = []tea.Model{New()}
	m := models[model]
	p := tea.NewProgram(m)
	if err := p.Start(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
