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

package api

// Internal API error codes
const (
	ErrorRedis          = "err_redis"
	ErrorCreatingToken  = "err_creating_token"
	ErrorParsingRequest = "err_parsing_token"
	ErrorInvalidToken   = "invalid_token"
	ErrorInvalidClient  = "invalid_client"
)

// ErrorMessages has the descriptions associated to the API error codes
var ErrorMessages map[string]string

// InitErrorMessages initializes the default messages for the API error codes
func InitErrorMessages() {
	ErrorMessages = make(map[string]string)
	ErrorMessages[ErrorRedis] = "There was an error connecting to the Redis Server"
	ErrorMessages[ErrorCreatingToken] = "Error creating auth token"
	ErrorMessages[ErrorParsingRequest] = "Error parsing token"
	ErrorMessages[ErrorInvalidToken] = "Invalid token"
	ErrorMessages[ErrorInvalidClient] = "Invalid client"
}
