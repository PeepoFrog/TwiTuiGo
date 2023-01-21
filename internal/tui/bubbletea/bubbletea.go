package bubbleteaTUI

//6 _ % \ ` ~  ^
import (
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"github.com/PeepoFrog/TwiTuiGo/internal/controller"
	myModels "github.com/PeepoFrog/TwiTuiGo/internal/model"
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

/* CUSTOM ITEM THATS INSERTS INTO LIST*/
type Task struct {
	status         status
	title          string
	description    string
	viewers        string
	gameStruct     myModels.Game
	streamerStruct myModels.Streamer
}

func NewTask(status status, title, description string, gameid string, gamestruct myModels.Game) Task {
	return Task{status: status, title: title, description: description, gameStruct: gamestruct, viewers: ""}
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
	loaded                    bool
	focused                   status
	lists                     []list.Model
	err                       error
	quitting                  bool
	gameStruct                myModels.Games
	gameList                  []myModels.Game
	broadcastStruct           myModels.Streamers
	broadcastList             []myModels.Streamer
	SelectedGameInColumn      myModels.Game
	SelectedBroadcastinColumn myModels.Streamer
	gamesCursor               string
	gamesCursorState          bool
	broadcastsCursor          string
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
		m.broadcastsCursor = gameStruct.Pagination.Cursor
		m.broadcastList = gameStruct.Data
		m.broadcastStruct = gameStruct
		var listtoadd []list.Item
		// for _, b := range gameStruct.Data {
		// 	listtoadd = append(listtoadd, Task{status: broadcastsColumn, title: b.UserName, description: b.GameName + " viewers: " + strconv.Itoa(b.ViewerCount), viewers: "300"})
		// }
		for _, b := range m.broadcastStruct.Data {
			listtoadd = append(listtoadd, Task{status: broadcastsColumn, title: b.UserName, description: b.GameName + " viewers: " + strconv.Itoa(b.ViewerCount), viewers: "300", streamerStruct: b})
		}
		listtoadd = append(listtoadd, Task{status: broadcastsColumn, title: "LOAD MORE"})
		m.lists[broadcastsColumn].SetItems(listtoadd)
		// m.lists[broadcastsColumn].SetSize(50, 30)

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
func (m *Model) SelectBroadcast() tea.Msg {
	selectedItem := m.lists[m.focused].SelectedItem().(Task)

	// print(selectedItem.Title())
	if selectedItem.title == "LOAD MORE" {
		m.LoadMoreBroadcasts(selectedItem.gameStruct.ID)
	} else {
		m.SelectedBroadcastinColumn = selectedItem.streamerStruct
		print(m.SelectedBroadcastinColumn.UserName + " HERE ]]]]]]]]]]]]]]]]]")
		// m.RunStreamlink(m.SelectedBroadcastinColumn.UserName)
		m.RunStreamlink(m.SelectedBroadcastinColumn.UserName)
	}
	return nil
}
func (m *Model) LoadMoreBroadcasts(id string) tea.Msg {
	broadcastStruct, err := controller.GetStreamsFromSelectedGame(&AuthToTwitch, m.broadcastsCursor, m.SelectedGameInColumn.ID)
	m.broadcastsCursor = broadcastStruct.Pagination.Cursor
	if err != nil {
		print(err)
	}
	var listtoadd []list.Item

	m.broadcastList = append(m.broadcastList, broadcastStruct.Data...)
	for _, b := range m.broadcastList {
		listtoadd = append(listtoadd, Task{status: broadcastsColumn, title: b.UserName, description: b.GameName + " viewers: " + strconv.Itoa(b.ViewerCount), viewers: "300", streamerStruct: b})
	}
	listtoadd = append(listtoadd, Task{status: broadcastsColumn, title: "LOAD MORE"})
	m.lists[broadcastsColumn].SetItems(listtoadd)
	return nil
}
func (m *Model) RunStreamlink(bname string) {
	stream := "twitch.tv/" + bname
	// print(stream + "here here here")
	cmd := exec.Command("./streamlink", stream, "720p60")
	_, err := cmd.Output()
	if err != nil {
		print(err)
	}
	// print(string(output))
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
		Task{status: broadcastsColumn, title: "title", description: "description", viewers: "viewers", streamerStruct: myModels.Streamer{}},
	})
	// Init favorites column
	m.lists[favoritesColumn].Title = "FAVORITES"

	m.lists[favoritesColumn].SetItems([]list.Item{
		Task{status: favoritesColumn, title: "title", description: "description"},
	})

	// m.lists[favoritesColumn].SetHeight(20)
	// var a list.ItemDelegate
	// a.Height()
	// m.lists[broadcastsColumn].SetHeight(10)
	// m.lists[broadcastsColumn].SetSize(40, 10)

}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.initLists(msg.Width, msg.Height)
		if !m.loaded {
			columnStyle.Width(msg.Width/divisor - 5)
			focusedStyle.Width(msg.Width/divisor - 5)
			columnStyle.Height(msg.Height - divisor)
			focusedStyle.Height(msg.Height - divisor)
			m.loaded = true
		} else {
			columnStyle.Width(msg.Width/divisor - 5)
			focusedStyle.Width(msg.Width/divisor - 5)
			columnStyle.Height(msg.Height - divisor)
			focusedStyle.Height(msg.Height - divisor)
		}
		m.initLists(msg.Width, msg.Height)
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
				print(id)
				m.LoadBroadcastsFromSelectedGame(id)
				return m, nil
			}
			if m.focused == broadcastsColumn {
				print("pepeg broadcasts")
				m.SelectBroadcast()
				// streamName := m.SelectedBroadcastinColumn.UserName
				// m.RunStreamlink(streamName)
				return m, nil
			}
			if m.focused == favoritesColumn {
				return m, nil
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
