package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Configuration
const (
	defaultServerURL = "http://localhost:8000/v1"
	configFile       = "cli-config.json"
)

// Config represents CLI configuration
type Config struct {
	ServerURL string `json:"server_url"`
	APIKey    string `json:"api_key,omitempty"`
	Username  string `json:"username,omitempty"`
}

// API Models matching your server
type User struct {
	ID               string    `json:"id"`
	Username         string    `json:"username"`
	Age              *int      `json:"age,omitempty"`
	Gender           *string   `json:"gender,omitempty"`
	Role             string    `json:"role"`
	OrganizationID   *string   `json:"organization_id,omitempty"`
	OrganizationName *string   `json:"organization_name,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type Task struct {
	ID         string    `json:"id"`
	TaskTitle  string    `json:"title"`
	TaskDesc   *string   `json:"description,omitempty"`
	IsFinished bool      `json:"is_finished"`
	UserID     string    `json:"user_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type Organization struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Request/Response models
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	APIKey string `json:"api_key"`
	User   User   `json:"user"`
}

type CreateUserRequest struct {
	Username string  `json:"username"`
	Password string  `json:"password"`
	Age      *int    `json:"age,omitempty"`
	Gender   *string `json:"gender,omitempty"`
}

type CreateTaskRequest struct {
	Title       string  `json:"title"`
	Description *string `json:"description,omitempty"`
}

type UpdateTaskRequest struct {
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	IsFinished  *bool   `json:"is_finished,omitempty"`
}

// Implement list.Item interface for Task
func (t Task) FilterValue() string { return t.TaskTitle }
func (t Task) Title() string       { return t.TaskTitle }
func (t Task) Description() string {
	status := "❌ Unfinished"
	if t.IsFinished {
		status = "✅ Finished"
	}
	desc := ""
	if t.TaskDesc != nil && *t.TaskDesc != "" {
		desc = *t.TaskDesc
	}
	return fmt.Sprintf("%s | %s | Created: %s", status, desc, t.CreatedAt.Format("2006-01-02 15:04"))
}

// Application states
type state int

const (
	loginView state = iota
	registerView
	menuView
	taskListView
	taskCreateView
	taskEditView
	taskDetailView
	userProfileView
	organizationView
)

// Main model
type Model struct {
	state  state
	config Config
	client *APIClient
	user   *User

	// UI components
	list       list.Model
	textInputs []textinput.Model
	focusIndex int

	// Data
	tasks        []Task
	selectedTask *Task

	// Messages
	message  string
	errorMsg string
}

// API Client
type APIClient struct {
	baseURL string
	apiKey  string
	client  *http.Client
}

func NewAPIClient(baseURL, apiKey string) *APIClient {
	return &APIClient{
		baseURL: baseURL,
		apiKey:  apiKey,
		client:  &http.Client{Timeout: 30 * time.Second},
	}
}

func (c *APIClient) makeRequest(method, endpoint string, body interface{}) (*http.Response, error) {
	var buf bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, c.baseURL+endpoint, &buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	if c.apiKey != "" {
		req.Header.Set("Authorization", "ApiKey "+c.apiKey)
	}

	return c.client.Do(req)
}

func (c *APIClient) Login(username, password string) (*LoginResponse, error) {
	loginReq := LoginRequest{Username: username, Password: password}

	resp, err := c.makeRequest("POST", "/login", loginReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("login failed: %s", string(body))
	}

	var loginResp LoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
		return nil, err
	}

	return &loginResp, nil
}

func (c *APIClient) Register(username, password string, age *int, gender *string) (*User, error) {
	createReq := CreateUserRequest{
		Username: username,
		Password: password,
		Age:      age,
		Gender:   gender,
	}

	resp, err := c.makeRequest("POST", "/user", createReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("registration failed: %s", string(body))
	}

	var user User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (c *APIClient) GetTasks() ([]Task, error) {
	resp, err := c.makeRequest("GET", "/tasks", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get tasks")
	}

	var tasks []Task
	if err := json.NewDecoder(resp.Body).Decode(&tasks); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (c *APIClient) CreateTask(title string, description *string) (*Task, error) {
	createReq := CreateTaskRequest{Title: title, Description: description}

	resp, err := c.makeRequest("POST", "/tasks", createReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to create task: %s", string(body))
	}

	var task Task
	if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
		return nil, err
	}

	return &task, nil
}

func (c *APIClient) UpdateTask(taskID string, req UpdateTaskRequest) (*Task, error) {
	resp, err := c.makeRequest("PUT", "/tasks/"+taskID, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to update task: %s", string(body))
	}

	var task Task
	if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
		return nil, err
	}

	return &task, nil
}

func (c *APIClient) DeleteTask(taskID string) error {
	resp, err := c.makeRequest("DELETE", "/tasks/"+taskID, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to delete task: %s", string(body))
	}

	return nil
}

func (c *APIClient) GetUser() (*User, error) {
	resp, err := c.makeRequest("GET", "/user", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get user")
	}

	var user User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

// Styles
var (
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1)

	appStyle = lipgloss.NewStyle().Padding(1, 2)

	statusMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"}).
				Render

	errorMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#FF0000", Dark: "#FF6B6B"}).
				Render

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262")).
			Render

	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF06B7"))
	blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#777777"))
)

func initialModel() Model {
	// Load config
	config := loadConfig()

	// Create text inputs
	inputs := make([]textinput.Model, 4)

	// Username input
	inputs[0] = textinput.New()
	inputs[0].Placeholder = "Username"
	inputs[0].Focus()
	inputs[0].PromptStyle = focusedStyle
	inputs[0].TextStyle = focusedStyle

	// Password input
	inputs[1] = textinput.New()
	inputs[1].Placeholder = "Password"
	inputs[1].EchoMode = textinput.EchoPassword
	inputs[1].EchoCharacter = '•'

	// Title input
	inputs[2] = textinput.New()
	inputs[2].Placeholder = "Task title..."
	inputs[2].CharLimit = 100
	inputs[2].Width = 50

	// Description input
	inputs[3] = textinput.New()
	inputs[3].Placeholder = "Task description (optional)..."
	inputs[3].CharLimit = 500
	inputs[3].Width = 50

	// Create list
	l := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	l.Title = "📋 Server Tasks"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)

	model := Model{
		config:     config,
		textInputs: inputs,
		list:       l,
	}

	// Determine initial state
	if config.APIKey != "" {
		model.state = menuView
		model.client = NewAPIClient(config.ServerURL, config.APIKey)
	} else {
		model.state = loginView
	}

	return model
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.state {
		case loginView:
			return m.updateLoginView(msg)
		case registerView:
			return m.updateRegisterView(msg)
		case menuView:
			return m.updateMenuView(msg)
		case taskListView:
			return m.updateTaskListView(msg)
		case taskCreateView:
			return m.updateTaskCreateView(msg)
		case taskEditView:
			return m.updateTaskEditView(msg)
		case taskDetailView:
			return m.updateTaskDetailView(msg)
		case userProfileView:
			return m.updateUserProfileView(msg)
		}

	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width - 4)
		m.list.SetHeight(msg.Height - 8)
		return m, nil
	}

	return m, nil
}

func (m Model) updateLoginView(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch keypress := msg.String(); keypress {
	case "ctrl+c", "q":
		return m, tea.Quit
	case "tab", "shift+tab", "enter", "up", "down":
		s := msg.String()
		if s == "enter" && m.focusIndex == len(m.textInputs[:2])-1 {
			// Attempt login
			username := m.textInputs[0].Value()
			password := m.textInputs[1].Value()

			if username == "" || password == "" {
				m.errorMsg = "Please enter both username and password"
				return m, nil
			}

			// Create client and attempt login
			client := NewAPIClient(m.config.ServerURL, "")
			loginResp, err := client.Login(username, password)
			if err != nil {
				m.errorMsg = fmt.Sprintf("Login failed: %v", err)
				return m, nil
			}

			// Save config and update state
			m.config.APIKey = loginResp.APIKey
			m.config.Username = username
			saveConfig(m.config)

			m.client = NewAPIClient(m.config.ServerURL, loginResp.APIKey)
			m.user = &loginResp.User
			m.state = menuView
			m.message = "Login successful!"
			m.errorMsg = ""

			return m, nil
		}

		if s == "up" || s == "shift+tab" {
			m.focusIndex--
		} else {
			m.focusIndex++
		}

		if m.focusIndex > len(m.textInputs[:2])-1 {
			m.focusIndex = 0
		} else if m.focusIndex < 0 {
			m.focusIndex = len(m.textInputs[:2]) - 1
		}

		for i := 0; i <= len(m.textInputs[:2])-1; i++ {
			if i == m.focusIndex {
				m.textInputs[i].Focus()
				m.textInputs[i].PromptStyle = focusedStyle
				m.textInputs[i].TextStyle = focusedStyle
			} else {
				m.textInputs[i].Blur()
				m.textInputs[i].PromptStyle = blurredStyle
				m.textInputs[i].TextStyle = blurredStyle
			}
		}
		return m, nil

	case "r":
		m.state = registerView
		m.focusIndex = 0
		m.errorMsg = ""
		return m, nil
	}

	var cmd tea.Cmd
	m.textInputs[m.focusIndex], cmd = m.textInputs[m.focusIndex].Update(msg)
	return m, cmd
}

func (m Model) updateRegisterView(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch keypress := msg.String(); keypress {
	case "ctrl+c":
		return m, tea.Quit
	case "esc":
		m.state = loginView
		m.focusIndex = 0
		m.errorMsg = ""
		return m, nil
	case "enter":
		// TODO: Implement registration logic
		m.message = "Registration feature coming soon!"
		return m, nil
	}

	var cmd tea.Cmd
	m.textInputs[m.focusIndex], cmd = m.textInputs[m.focusIndex].Update(msg)
	return m, cmd
}

func (m Model) updateMenuView(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch keypress := msg.String(); keypress {
	case "ctrl+c", "q":
		return m, tea.Quit
	case "1", "t":
		m.state = taskListView
		m.loadTasks()
		return m, nil
	case "2", "p":
		m.state = userProfileView
		return m, nil
	case "3", "o":
		m.state = organizationView
		return m, nil
	case "l":
		// Logout
		m.config.APIKey = ""
		m.config.Username = ""
		saveConfig(m.config)
		m.state = loginView
		m.user = nil
		m.client = nil
		m.message = "Logged out successfully"
		return m, nil
	}
	return m, nil
}

func (m Model) updateTaskListView(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch keypress := msg.String(); keypress {
	case "ctrl+c":
		return m, tea.Quit
	case "esc", "b":
		m.state = menuView
		return m, nil
	case "n":
		m.state = taskCreateView
		m.textInputs[2].SetValue("")
		m.textInputs[3].SetValue("")
		m.textInputs[2].Focus()
		m.focusIndex = 2
		m.errorMsg = ""
		return m, nil
	case "enter":
		if len(m.tasks) > 0 {
			m.selectedTask = &m.tasks[m.list.Index()]
			m.state = taskDetailView
		}
		return m, nil
	case "r":
		m.loadTasks()
		m.message = "Tasks refreshed!"
		return m, nil
	case "t":
		if len(m.tasks) > 0 {
			selected := m.list.Index()
			task := m.tasks[selected]
			newStatus := !task.IsFinished

			updateReq := UpdateTaskRequest{IsFinished: &newStatus}
			_, err := m.client.UpdateTask(task.ID, updateReq)
			if err != nil {
				m.errorMsg = fmt.Sprintf("Failed to update task: %v", err)
			} else {
				m.loadTasks()
				status := "unfinished"
				if newStatus {
					status = "finished"
				}
				m.message = fmt.Sprintf("Task marked as %s!", status)
			}
		}
		return m, nil
	case "d":
		if len(m.tasks) > 0 {
			selected := m.list.Index()
			task := m.tasks[selected]

			err := m.client.DeleteTask(task.ID)
			if err != nil {
				m.errorMsg = fmt.Sprintf("Failed to delete task: %v", err)
			} else {
				m.loadTasks()
				m.message = "Task deleted successfully!"
			}
		}
		return m, nil
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m Model) updateTaskCreateView(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch keypress := msg.String(); keypress {
	case "ctrl+c":
		return m, tea.Quit
	case "esc":
		m.state = taskListView
		return m, nil
	case "tab", "shift+tab", "up", "down":
		if keypress == "up" || keypress == "shift+tab" {
			if m.focusIndex > 2 {
				m.focusIndex--
			}
		} else {
			if m.focusIndex < 3 {
				m.focusIndex++
			}
		}

		for i := 2; i <= 3; i++ {
			if i == m.focusIndex {
				m.textInputs[i].Focus()
			} else {
				m.textInputs[i].Blur()
			}
		}
		return m, nil
	case "enter":
		title := strings.TrimSpace(m.textInputs[2].Value())
		if title == "" {
			m.errorMsg = "Task title is required"
			return m, nil
		}

		var description *string
		if desc := strings.TrimSpace(m.textInputs[3].Value()); desc != "" {
			description = &desc
		}

		_, err := m.client.CreateTask(title, description)
		if err != nil {
			m.errorMsg = fmt.Sprintf("Failed to create task: %v", err)
		} else {
			m.state = taskListView
			m.loadTasks()
			m.message = "Task created successfully!"
		}
		return m, nil
	}

	var cmd tea.Cmd
	m.textInputs[m.focusIndex], cmd = m.textInputs[m.focusIndex].Update(msg)
	return m, cmd
}

func (m Model) updateTaskEditView(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Similar to create view but for editing
	switch keypress := msg.String(); keypress {
	case "ctrl+c":
		return m, tea.Quit
	case "esc":
		m.state = taskDetailView
		return m, nil
	case "enter":
		if m.selectedTask == nil {
			return m, nil
		}

		title := strings.TrimSpace(m.textInputs[2].Value())
		if title == "" {
			m.errorMsg = "Task title is required"
			return m, nil
		}

		var description *string
		if desc := strings.TrimSpace(m.textInputs[3].Value()); desc != "" {
			description = &desc
		}

		updateReq := UpdateTaskRequest{Title: &title, Description: description}
		_, err := m.client.UpdateTask(m.selectedTask.ID, updateReq)
		if err != nil {
			m.errorMsg = fmt.Sprintf("Failed to update task: %v", err)
		} else {
			m.loadTasks()
			m.state = taskDetailView
			m.message = "Task updated successfully!"
		}
		return m, nil
	}

	var cmd tea.Cmd
	m.textInputs[m.focusIndex], cmd = m.textInputs[m.focusIndex].Update(msg)
	return m, cmd
}

func (m Model) updateTaskDetailView(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch keypress := msg.String(); keypress {
	case "ctrl+c":
		return m, tea.Quit
	case "esc", "b":
		m.state = taskListView
		return m, nil
	case "e":
		if m.selectedTask != nil {
			m.textInputs[2].SetValue(m.selectedTask.TaskTitle)
			if m.selectedTask.TaskDesc != nil {
				m.textInputs[3].SetValue(*m.selectedTask.TaskDesc)
			} else {
				m.textInputs[3].SetValue("")
			}
			m.textInputs[2].Focus()
			m.focusIndex = 2
			m.state = taskEditView
			m.errorMsg = ""
		}
		return m, nil
	case "t":
		if m.selectedTask != nil {
			newStatus := !m.selectedTask.IsFinished
			updateReq := UpdateTaskRequest{IsFinished: &newStatus}
			_, err := m.client.UpdateTask(m.selectedTask.ID, updateReq)
			if err != nil {
				m.errorMsg = fmt.Sprintf("Failed to update task: %v", err)
			} else {
				m.selectedTask.IsFinished = newStatus
				m.loadTasks()
				status := "unfinished"
				if newStatus {
					status = "finished"
				}
				m.message = fmt.Sprintf("Task marked as %s!", status)
			}
		}
		return m, nil
	}
	return m, nil
}

func (m Model) updateUserProfileView(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch keypress := msg.String(); keypress {
	case "ctrl+c":
		return m, tea.Quit
	case "esc", "b":
		m.state = menuView
		return m, nil
	}
	return m, nil
}

func (m Model) View() string {
	switch m.state {
	case loginView:
		return m.loginView()
	case registerView:
		return m.registerView()
	case menuView:
		return m.menuView()
	case taskListView:
		return m.taskListView()
	case taskCreateView:
		return m.taskCreateView()
	case taskEditView:
		return m.taskEditView()
	case taskDetailView:
		return m.taskDetailView()
	case userProfileView:
		return m.userProfileView()
	}
	return ""
}

func (m Model) loginView() string {
	var content strings.Builder

	content.WriteString(titleStyle.Render("🔐 GO Task Manager - Server Login"))
	content.WriteString("\n\n")
	content.WriteString(fmt.Sprintf("Server: %s\n\n", m.config.ServerURL))

	content.WriteString("Username:\n")
	content.WriteString(m.textInputs[0].View())
	content.WriteString("\n\n")

	content.WriteString("Password:\n")
	content.WriteString(m.textInputs[1].View())
	content.WriteString("\n\n")

	if m.errorMsg != "" {
		content.WriteString(errorMessageStyle(m.errorMsg))
		content.WriteString("\n\n")
	}

	if m.message != "" {
		content.WriteString(statusMessageStyle(m.message))
		content.WriteString("\n\n")
	}

	content.WriteString(helpStyle("(enter) login • (r) register • (ctrl+c) quit"))

	return appStyle.Render(content.String())
}

func (m Model) registerView() string {
	var content strings.Builder

	content.WriteString(titleStyle.Render("📝 Register New Account"))
	content.WriteString("\n\n")
	content.WriteString("Registration feature coming soon...\n\n")
	content.WriteString(helpStyle("(esc) back to login"))

	return appStyle.Render(content.String())
}

func (m Model) menuView() string {
	var content strings.Builder

	username := "User"
	if m.user != nil {
		username = m.user.Username
	}

	content.WriteString(titleStyle.Render(fmt.Sprintf("🏠 Welcome, %s!", username)))
	content.WriteString("\n\n")

	if m.user != nil {
		content.WriteString(fmt.Sprintf("Role: %s\n", m.user.Role))
		if m.user.OrganizationName != nil {
			content.WriteString(fmt.Sprintf("Organization: %s\n", *m.user.OrganizationName))
		}
		content.WriteString("\n")
	}

	content.WriteString("📋 Main Menu\n\n")
	content.WriteString("1. (t) Task Management\n")
	content.WriteString("2. (p) User Profile\n")
	content.WriteString("3. (o) Organization\n\n")

	if m.message != "" {
		content.WriteString(statusMessageStyle(m.message))
		content.WriteString("\n\n")
	}

	content.WriteString(helpStyle("Select option • (l) logout • (q) quit"))

	return appStyle.Render(content.String())
}

func (m Model) taskListView() string {
	var content strings.Builder

	content.WriteString(titleStyle.Render("📋 Server Tasks"))
	content.WriteString("\n\n")

	if len(m.tasks) == 0 {
		content.WriteString("📝 No tasks found. Press 'n' to create your first task!\n\n")
	} else {
		content.WriteString(m.list.View())
		content.WriteString("\n")
	}

	if m.message != "" {
		content.WriteString(statusMessageStyle(m.message))
		content.WriteString("\n\n")
	}

	if m.errorMsg != "" {
		content.WriteString(errorMessageStyle(m.errorMsg))
		content.WriteString("\n\n")
	}

	help := "(n)ew • (r)efresh • (b)ack • (q)uit"
	if len(m.tasks) > 0 {
		help = "(↑↓) navigate • (enter) details • (n)ew • (t)oggle • (d)elete • (r)efresh • (b)ack"
	}
	content.WriteString(helpStyle(help))

	return appStyle.Render(content.String())
}

func (m Model) taskCreateView() string {
	var content strings.Builder

	content.WriteString(titleStyle.Render("✏️ Create New Task"))
	content.WriteString("\n\n")

	content.WriteString("Title:\n")
	content.WriteString(m.textInputs[2].View())
	content.WriteString("\n\n")

	content.WriteString("Description (optional):\n")
	content.WriteString(m.textInputs[3].View())
	content.WriteString("\n\n")

	if m.errorMsg != "" {
		content.WriteString(errorMessageStyle(m.errorMsg))
		content.WriteString("\n\n")
	}

	content.WriteString(helpStyle("(tab) next field • (enter) create • (esc) cancel"))

	return appStyle.Render(content.String())
}

func (m Model) taskEditView() string {
	var content strings.Builder

	content.WriteString(titleStyle.Render("✏️ Edit Task"))
	content.WriteString("\n\n")

	content.WriteString("Title:\n")
	content.WriteString(m.textInputs[2].View())
	content.WriteString("\n\n")

	content.WriteString("Description:\n")
	content.WriteString(m.textInputs[3].View())
	content.WriteString("\n\n")

	if m.errorMsg != "" {
		content.WriteString(errorMessageStyle(m.errorMsg))
		content.WriteString("\n\n")
	}

	content.WriteString(helpStyle("(enter) save • (esc) cancel"))

	return appStyle.Render(content.String())
}

func (m Model) taskDetailView() string {
	if m.selectedTask == nil {
		return "No task selected"
	}

	var content strings.Builder

	content.WriteString(titleStyle.Render("📋 Task Details"))
	content.WriteString("\n\n")

	content.WriteString(lipgloss.NewStyle().Bold(true).Render("Title: "))
	content.WriteString(m.selectedTask.TaskTitle)
	content.WriteString("\n\n")

	if m.selectedTask.TaskDesc != nil && *m.selectedTask.TaskDesc != "" {
		content.WriteString(lipgloss.NewStyle().Bold(true).Render("Description: "))
		content.WriteString(*m.selectedTask.TaskDesc)
		content.WriteString("\n\n")
	}

	content.WriteString(lipgloss.NewStyle().Bold(true).Render("Status: "))
	if m.selectedTask.IsFinished {
		content.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#04B575")).Render("✅ Finished"))
	} else {
		content.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#FF6B6B")).Render("❌ Unfinished"))
	}
	content.WriteString("\n\n")

	content.WriteString(lipgloss.NewStyle().Bold(true).Render("Created: "))
	content.WriteString(m.selectedTask.CreatedAt.Format("January 2, 2006 at 3:04 PM"))
	content.WriteString("\n")

	content.WriteString(lipgloss.NewStyle().Bold(true).Render("Updated: "))
	content.WriteString(m.selectedTask.UpdatedAt.Format("January 2, 2006 at 3:04 PM"))
	content.WriteString("\n\n")

	if m.message != "" {
		content.WriteString(statusMessageStyle(m.message))
		content.WriteString("\n\n")
	}

	if m.errorMsg != "" {
		content.WriteString(errorMessageStyle(m.errorMsg))
		content.WriteString("\n\n")
	}

	content.WriteString(helpStyle("(e)dit • (t)oggle status • (b)ack"))

	return appStyle.Render(content.String())
}

func (m Model) userProfileView() string {
	var content strings.Builder

	content.WriteString(titleStyle.Render("👤 User Profile"))
	content.WriteString("\n\n")

	if m.user != nil {
		content.WriteString(lipgloss.NewStyle().Bold(true).Render("Username: "))
		content.WriteString(m.user.Username)
		content.WriteString("\n\n")

		content.WriteString(lipgloss.NewStyle().Bold(true).Render("Role: "))
		content.WriteString(m.user.Role)
		content.WriteString("\n\n")

		if m.user.Age != nil {
			content.WriteString(lipgloss.NewStyle().Bold(true).Render("Age: "))
			content.WriteString(strconv.Itoa(*m.user.Age))
			content.WriteString("\n\n")
		}

		if m.user.Gender != nil {
			content.WriteString(lipgloss.NewStyle().Bold(true).Render("Gender: "))
			content.WriteString(*m.user.Gender)
			content.WriteString("\n\n")
		}

		if m.user.OrganizationName != nil {
			content.WriteString(lipgloss.NewStyle().Bold(true).Render("Organization: "))
			content.WriteString(*m.user.OrganizationName)
			content.WriteString("\n\n")
		}

		content.WriteString(lipgloss.NewStyle().Bold(true).Render("Joined: "))
		content.WriteString(m.user.CreatedAt.Format("January 2, 2006"))
		content.WriteString("\n\n")
	}

	content.WriteString(helpStyle("(b)ack to menu"))

	return appStyle.Render(content.String())
}

// Helper functions
func (m *Model) loadTasks() {
	if m.client == nil {
		return
	}

	tasks, err := m.client.GetTasks()
	if err != nil {
		m.errorMsg = fmt.Sprintf("Failed to load tasks: %v", err)
		return
	}

	m.tasks = tasks
	m.updateTaskList()
	m.errorMsg = ""
}

func (m *Model) updateTaskList() {
	items := make([]list.Item, len(m.tasks))
	for i, task := range m.tasks {
		items[i] = task
	}
	m.list.SetItems(items)
}

// Config functions
func loadConfig() Config {
	config := Config{
		ServerURL: defaultServerURL,
	}

	data, err := os.ReadFile(configFile)
	if err != nil {
		return config
	}

	json.Unmarshal(data, &config)
	return config
}

func saveConfig(config Config) {
	data, _ := json.MarshalIndent(config, "", "  ")
	os.WriteFile(configFile, data, 0644)
}

func main() {
	model := initialModel()

	// Load user data if logged in
	if model.client != nil {
		user, err := model.client.GetUser()
		if err == nil {
			model.user = user
		}
	}

	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
