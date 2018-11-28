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
// limitations under the License.s

package api

import (
	"net/http"

	"github.com/pquerna/ffjson/ffjson"
)

// Token represents a request/response
type Token struct {
	Token     string `json:"token,omitempty"`
	ExpiresIn int    `json:"expires_in,omitempty"`
}

// Claim represents a request/response
type Claim struct {
	Claims map[string]interface{} `json:"claims"`
}

// TokenValidResponse is the response used for a /token/validate operation
type TokenValidResponse struct {
	Description string `json:"description"`
}

// Response represents a HTTP Response
type Response struct {
	ContentType string
	Status      int
	ErrorCode   string
	Payload     []byte
}

// Send writes the HTTP response in a http.ResponseWriter
func (response *Response) Send(w http.ResponseWriter) {

	// default content type is 'application/json'
	if response.ContentType == "" {
		response.ContentType = "application/json"
	}

	w.Header().Set("Content-Type", response.ContentType)
	w.WriteHeader(response.Status)

	if response.Status >= http.StatusBadRequest {
		if response.Payload == nil && response.ErrorCode != "" {

			if msg, ok := ErrorMessages[response.ErrorCode]; ok {
				errData := ErrorData{
					Error:   response.ErrorCode,
					Message: msg,
				}
				response.Payload = errData.Marshall()
			}
		}
	}

	w.Write(response.Payload)
}

// ErrorData is used to represent API errors
type ErrorData struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

// Marshall encodes ErrorData content to a json byte array
func (errData *ErrorData) Marshall() []byte {
	if b, err := ffjson.Marshal(errData); err == nil {
		return b
	}
	return []byte("")
}
