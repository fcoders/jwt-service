// Copyright 2019 Foo Coders (www.foocoders.io).
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
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/fcoders/jwt-service/api"
	"github.com/fcoders/jwt-service/core/authentication"
	"github.com/fcoders/jwt-service/routes"
	"github.com/fcoders/jwt-service/services"
	"github.com/fcoders/jwt-service/settings"
	"github.com/fcoders/logger"
	"github.com/gin-gonic/gin"
	"github.com/yvasiyarov/gorelic"
)

// HTTPService represents the HTTP service that is initiated when the server starts.
// By having it contained in a separate type, we can easily handle all it's events.
type HTTPService struct {
	engine    *gin.Engine
	waitGroup *sync.WaitGroup
	agent     *gorelic.Agent
}

// Init creates a new instance of the HTTP engine
func (service *HTTPService) Init() {
	service.waitGroup = &sync.WaitGroup{}

	// api messages
	api.InitErrorMessages()

	gin.SetMode(gin.ReleaseMode)

	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(gin.Logger())

	service.engine = engine
}

// Start starts the HTTP service
func (service *HTTPService) Start() {

	routes.InitRoutes(service.engine)
	port := fmt.Sprintf(":%v", settings.Get().App.HTTPPort)

	server := &http.Server{
		Addr:    port,
		Handler: service.engine,
	}

	go server.ListenAndServe()
	service.waitGroup.Add(1)

	log := logger.GetLogger()
	log.Infof("%s service started!", settings.AppName)
	log.Infof("Version %s commit %s", settings.Version, settings.CommitHash)
}

// Stop ends the HTTP service execution and release all the resources
func (service *HTTPService) Stop(cause string) {

	log := services.Get().Logger
	log.Infof("Shutdown requested with signal '%s'", strings.ToUpper(cause))
	authentication.CloseCacheConnections()

	log.Infof("%s service is now ready to exit, bye!", settings.AppName)
	service.waitGroup.Done()
}
