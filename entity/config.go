package entity

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/herryg91/localhost-proxy/pkg/editor"
)

type Config struct {
	Port   int              `json:"port"`   // 7000
	Editor string           `json:"editor"` // vi | nano
	Routes []LocalhostRoute `json:"routes"`
}

type LocalhostRoute struct {
	Name        string `json:"name"`
	Pathname    string `json:"path"`
	Destination string `json:"dest"`
}

func (Config) FromDefaultConfig() *Config {
	return &Config{
		Port:   7000,
		Editor: "vi",
		Routes: []LocalhostRoute{
			{
				Name:        "user-api",
				Pathname:    "/user",
				Destination: "localhost:18080",
			},
			{
				Name:        "product-api",
				Pathname:    "/product",
				Destination: "localhost:18081",
			},
		},
	}
}

func (Config) FromFile() *Config {
	homeDir, _ := os.UserHomeDir()
	configDir := homeDir + "/.lprx"
	configFileLoc := configDir + "/config.json"

	// if ~/.dply not exist, then mkdir ~/.dply
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		os.Mkdir(configDir, 0755)
	}

	// if ~/.dply/config.json not exist, then create default config
	if _, err := os.Stat(configFileLoc); os.IsNotExist(err) {
		Config{}.FromDefaultConfig().SaveConfig()
		fmt.Println("config.json file isn't exist, generate default config file: `" + configFileLoc + "`")
	}

	if _, err := os.Stat(configFileLoc); os.IsNotExist(err) {
		return nil
	}

	configFromFile, _ := ioutil.ReadFile(configFileLoc)
	s := Config{}
	_ = json.Unmarshal(configFromFile, &s)

	needToRewrite := false
	if s.Port == 0 {
		needToRewrite = true
		s.Port = 7000
	}
	if s.Editor == "" {
		needToRewrite = true
		s.Editor = "vi"
	}
	if needToRewrite {
		s.SaveConfig()
		fmt.Println("config.json file is broken, repairing to default config file: `" + configFileLoc + "`")
	}

	return &s
}

func (s *Config) SaveConfig() error {
	homeDir, _ := os.UserHomeDir()
	configDir := homeDir + "/.lprx"
	configFileLoc := configDir + "/config.json"

	configJsonMarshalled, _ := json.Marshal(&s)

	err := ioutil.WriteFile(configFileLoc, configJsonMarshalled, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

var ErrConfigNothingChange = fmt.Errorf("Nothing to changes")

func (s *Config) UpdateByEditor() error {
	current_data, _ := json.MarshalIndent(s, "", "    ")

	editor_app := editor.EditorApp(s.Editor)
	updated_data, err := editor.Open(editor_app, "tmp_config_edit", current_data)
	if err != nil {
		return fmt.Errorf("Error editor: " + err.Error())
	}

	// if nothing to change
	if string(current_data) == string(updated_data) {
		return ErrConfigNothingChange
	}

	newConfig := &Config{}
	err = json.Unmarshal(updated_data, &newConfig)
	if err != nil {
		return fmt.Errorf("Error unmarshal: " + string(updated_data))
	}

	s = newConfig
	err = s.SaveConfig()
	if err != nil {
		return fmt.Errorf("Failed to save config: " + err.Error())
	}

	return nil
}
