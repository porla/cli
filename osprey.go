package main

import (
	"os"

	"osprey/config"
	"osprey/http"
	"osprey/ui"
	"osprey/utils"

	tea "github.com/charmbracelet/bubbletea"
	"gopkg.in/yaml.v3"
)

func main() {
	f, err := os.ReadFile(config.ConfigFilePath)
	utils.CheckError(err)
	err = yaml.Unmarshal(f, &config.Config)
	utils.CheckError(err)
	http.InitHTTPClient()
	p := tea.NewProgram(ui.InitialModel())
	_, err = p.Run()
	utils.CheckError(err)
}
