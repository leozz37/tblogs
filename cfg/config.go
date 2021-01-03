package cfg

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/ezeoleaf/tblogs/models"
	"gopkg.in/yaml.v2"
)

// Config contains both API and APP configuration from config file
type Config struct {
	API APIConfig `yaml:"api"`
	APP APPConfig `yaml:"app"`
}

// APIConfig contains only the API configuration from config file
type APIConfig struct {
	Host string `yaml:"url"`
	Key  string `yaml:"key"`
}

// APPConfig contains only the APP configuration from config file
type APPConfig struct {
	SavedPosts     []models.Post `yaml:"saved_posts"`
	FollowingBlogs []int         `yaml:"following_blogs"`
	FirstUse       bool          `yaml:"first_use"`
	LastLogin      time.Time     `yaml:"last_login"`
	CurrentLogin   time.Time     `yaml:"current_login"`
	FilteredWords  []string      `yaml:"filtered_words"`
}

const configPath = "./cfg/config.yml"

func parseFlags() (string, error) {
	configPath := configPath
	// Validate the path first
	if err := validateConfigPath(configPath); err != nil {
		return "", err
	}

	// Return the configuration path
	return configPath, nil
}

func validateConfigPath(path string) error {
	s, err := os.Stat(path)
	if err != nil {
		return err
	}
	if s.IsDir() {
		return fmt.Errorf("'%s' is a directory, not a normal file", path)
	}
	return nil
}

// // GetAPPConfig returns only the APP config
// func GetAPPConfig() models.APPConfig {
// 	return config.APP
// }

// // GetAPIConfig returns only the API config for the application
// func GetAPIConfig() models.APIConfig {
// 	return config.API
// }

// // GetConfig returns the entire config for the application
// func GetConfig() models.Config {
// 	return config
// }

// UpdateConfig updates the yml file with the latest config
func (c *Config) UpdateConfig() {
	d, err := yaml.Marshal(c)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	cfgPath, err := parseFlags()
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(cfgPath, d, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

// ResetAPPConfig reset the entire AppConfig and return to the initial state
func (c *Config) ResetAPPConfig() {
	c.APP = APPConfig{}
	c.UpdateConfig()
}

func setNewFile() (string, error) {
	from, err := os.Open("./cfg/config.example.yml")
	if err != nil {
		log.Fatal(err)
	}
	defer from.Close()

	to, err := os.OpenFile(configPath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer to.Close()

	_, err = io.Copy(to, from)
	if err != nil {
		log.Fatal(err)
	}

	return parseFlags()
}

// Setup prepare and create the config file for the application
func Setup() *Config {
	cfgPath, err := parseFlags()
	if err != nil {
		log.Println(err)
	}

	if cfgPath == "" {
		cfgPath, err = setNewFile()
		if err != nil {
			log.Fatal(err)
		}
	}

	file, err := os.Open(cfgPath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Init new YAML decode
	d := yaml.NewDecoder(file)

	var config Config

	// Start YAML decoding from file
	if err := d.Decode(&config); err != nil {
		panic(err)
	}

	// Update last login date
	config.UpdateLoginDate()

	return &config
}

// UpdateLoginDate update the current login date with time.Now() in UTC
// It also saves the last login with the value of current login before updating it
func (c *Config) UpdateLoginDate() {
	var loc, _ = time.LoadLocation("UTC")
	now := time.Now().In(loc)

	c.APP.LastLogin = c.APP.CurrentLogin
	c.APP.CurrentLogin = now

	c.UpdateConfig()
}
