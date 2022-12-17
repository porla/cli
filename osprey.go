package main

import (
	"os"

	"osprey/config"
	"osprey/i18n"
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
	config.Currenti18n = i18n.LoadLanguage(config.Config.I18nLanguage)
	p := tea.NewProgram(ui.InitialModel())
	_, err = p.Run()
	utils.CheckError(err)
}
