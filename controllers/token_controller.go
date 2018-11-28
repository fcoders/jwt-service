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

package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/fcoders/jwt-service/api"
	"github.com/fcoders/jwt-service/services/token"
	"github.com/gin-gonic/gin"
)

// Generate handles the requests to generate a new token
func Generate() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiResponse := new(api.Response)
		request := new(api.Claim)

		clientID := c.Request.Header.Get("Auth-Client")
		if len(clientID) == 0 {
			apiResponse.Status = http.StatusBadRequest
			apiResponse.ErrorCode = api.ErrorInvalidClient
		} else {

			decoder := json.NewDecoder(c.Request.Body)
			if errDecode := decoder.Decode(&request); errDecode != nil {
				apiResponse.Status = http.StatusBadRequest
				apiResponse.ErrorCode = api.ErrorParsingRequest
			} else {
				apiResponse = token.Generate(request, clientID)
			}
		}

		apiResponse.Send(c.Writer)
	}
}

// Validate handles the requests to validate if a token is legal or not
func Validate() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiResponse := new(api.Response)
		request := new(api.Token)

		clientID := c.Request.Header.Get("Auth-Client")
		if len(clientID) == 0 {
			apiResponse.Status = http.StatusBadRequest
			apiResponse.ErrorCode = api.ErrorInvalidClient

		} else {

			decoder := json.NewDecoder(c.Request.Body)
			if errDecode := decoder.Decode(&request); errDecode != nil {
				apiResponse.Status = http.StatusBadRequest
				apiResponse.ErrorCode = api.ErrorParsingRequest
			} else {
				apiResponse = token.Validate(request, clientID)
			}
		}

		apiResponse.Send(c.Writer)
	}
}

// Destroy invalidates a not expired token, by saving it in a temporary cache until its expiration
func Destroy() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiResponse := new(api.Response)
		request := new(api.Token)

		clientID := c.Request.Header.Get("Auth-Client")
		if len(clientID) == 0 {
			apiResponse.Status = http.StatusBadRequest
			apiResponse.ErrorCode = api.ErrorInvalidClient
		} else {

			decoder := json.NewDecoder(c.Request.Body)

			if errDecode := decoder.Decode(&request); errDecode != nil {
				apiResponse.Status = http.StatusBadRequest
				apiResponse.ErrorCode = api.ErrorParsingRequest
			} else {
				apiResponse = token.Destroy(request, clientID)
			}
		}

		apiResponse.Send(c.Writer)
	}
}
