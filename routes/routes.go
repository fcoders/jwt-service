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

package routes

import (
	"github.com/fcoders/jwt-service/controllers"

	"github.com/gin-gonic/gin"
)

// InitRoutes configures all the routes handlded by the application
func InitRoutes(engine *gin.Engine) {

	v1 := engine.Group("/v1")
	{
		token := v1.Group("/token")
		{
			token.POST("/generate", controllers.Generate())
			token.POST("/validate", controllers.Validate())
			token.POST("/destroy", controllers.Destroy())
		}
	}
}
