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

package main

import (
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/fcoders/jwt-service/services"
	"github.com/fcoders/jwt-service/settings"
)

func main() {
	appInit()
	httpServiceInit()
}

func appInit() {

	// load settings
	cfgFile := getAppPath() + "/settings.yml"
	if err := settings.Init(cfgFile); err != nil {
		panic(err)
	}

	// dependency manager
	err := services.Init()
	if err != nil {
		log.Panicf("Error initiation dependency manager: %s", err.Error())
	}

}

func httpServiceInit() {
	httpService := HTTPService{}
	httpService.Init()

	go httpService.Start()

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	sign := <-ch

	httpService.Stop(sign.String())
	os.Exit(0)
}

func getAppPath() string {
	if dir, err := filepath.Abs(filepath.Dir(os.Args[0])); err == nil {
		return dir
	}

	return ""
}
