// Copyright 2018 Foo Coders (www.foocoders.io).
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package settings

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/fcoders/jwt-service/common"
	"gopkg.in/yaml.v2"
)

// Settings is the structure used to hold configuration from settings.yml
type Settings struct {
	App struct {
		HTTPPort int `yaml:"http_port"`
		LogLevel int `yaml:"log_level"`
	} `yaml:"app"`
	JWT struct {
		TokenExpiration int `yaml:"token_expiration"`
	} `yaml:"jwt"`
	Redis struct {
		Address  string `yaml:"address"`
		Password string `yaml:"password"`
	} `yaml:"redis"`
	NewRelic struct {
		Enabled    bool   `yaml:"enabled"`
		LicenseKey string `yaml:"license_key"`
	} `yaml:"new_relic"`
	Slack struct {
		Enabled    bool   `yaml:"enabled"`
		WebhookURL string `yaml:"webhook_url"`
	} `yaml:"slack"`
	Proxy struct {
		Enabled bool   `yaml:"enabled"`
		Address string `yaml:"address"`
	} `yaml:"proxy"`
}

// App modes
const (
	ModeDebug   = "debug"
	ModeTest    = "test"
	ModeRelease = "release"
)

// Log destinations
const (
	LogDestinationConsole = "console"
	LogDestinationFile    = "file"
)

var cfg *Settings

// Init loads the content of a file named .env in the application's working
// path and returns a Settings object
func Init(settingsPath string) (err error) {
	if common.Exists(settingsPath) {
		err = loadSettingsFromFile(settingsPath)
	} else {
		return errors.New("File 'settings.yml' not found. Cannot load settings")
	}

	// get current host name
	HostName, _ = os.Hostname()

	return
}

// Get returns the configuration loaded
func Get() *Settings {
	return cfg
}

// GetHTTPClient returns a HTTP client with or without proxy configured
func GetHTTPClient() (client *http.Client) {
	if Get().Proxy.Enabled {
		proxyURL, _ := url.Parse(Get().Proxy.Address)
		client = &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyURL)}}
	} else {
		client = http.DefaultClient
	}
	return
}

// load settings from yaml file
func loadSettingsFromFile(file string) error {
	cfg = new(Settings)

	fileContent, errRead := ioutil.ReadFile(file)
	if errRead != nil {
		return errRead
	}

	return yaml.Unmarshal(fileContent, cfg)
}
