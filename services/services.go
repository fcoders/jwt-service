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

package services

import (
	"github.com/fcoders/jwt-service/core/cache"
	"github.com/fcoders/jwt-service/core/cache/redis"
	"github.com/fcoders/jwt-service/logger"

	"github.com/facebookgo/inject"
)

// AppLoader holds the services instances configured for the App
type AppLoader struct {
	Cache  cache.Connector `inject:""`
	Logger logger.Logger   `inject:""`
}

var (
	loader *AppLoader
)

// DefaultLogger returns the default logger configured
func DefaultLogger() logger.Logger {
	if loader != nil {
		return loader.Logger
	}
	return nil
}

// Get returns the current application services
func Get() *AppLoader {
	return loader
}

// Init handles the dependency injection
func Init() (err error) {

	loader = new(AppLoader)

	// set logger
	log := logger.NewSimpleLogger()
	logger.SetLogger(log)

	// cache
	cache := new(redis.Pool)

	// instances for service container
	var graph inject.Graph
	if err = graph.Provide(
		&inject.Object{Value: loader},
		&inject.Object{Value: log},
		&inject.Object{Value: cache},
	); err != nil {
		return
	}

	err = graph.Populate()
	return
}
