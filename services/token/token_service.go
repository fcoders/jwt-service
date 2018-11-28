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

package token

import (
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/fcoders/jwt-service/api"
	"github.com/fcoders/jwt-service/core/authentication"
	"github.com/fcoders/jwt-service/services"
	"github.com/pquerna/ffjson/ffjson"
)

// Generate generates a new token
func Generate(request *api.Claim, client string) *api.Response {

	httpResponse := new(api.Response)
	authBackend, errJWT := authentication.InitJWTAuthenticationBackend(services.Get().Cache)

	if errJWT != nil {
		httpResponse.Status = http.StatusBadRequest
		httpResponse.ErrorCode = api.ErrorInvalidClient
		return httpResponse
	}

	token, expiresIn, err := authBackend.GenerateToken(request.Claims, client)
	if err != nil {

		services.Get().Logger.Infof("Error generating token: %s", err)

		httpResponse.Status = http.StatusInternalServerError
		httpResponse.ErrorCode = api.ErrorCreatingToken
		return httpResponse
	}

	httpResponse.Status = http.StatusOK
	httpResponse.Payload, _ = ffjson.Marshal(api.Token{Token: token, ExpiresIn: expiresIn})

	return httpResponse
}

// Validate validates the token
func Validate(request *api.Token, client string) *api.Response {

	httpResponse := new(api.Response)
	authBackend, errJWT := authentication.InitJWTAuthenticationBackend(services.Get().Cache)

	if errJWT != nil {
		httpResponse.Status = http.StatusBadRequest
		httpResponse.ErrorCode = api.ErrorInvalidClient
		return httpResponse
	}

	// check if token is not in blacklist
	if authBackend.IsInBlacklist(request.Token) {
		httpResponse.Status = http.StatusBadRequest
		httpResponse.ErrorCode = api.ErrorInvalidToken

	} else {

		// parse token and check its validity
		token, err := authBackend.ParseToken(request.Token, client)
		tokenClaims := token.Claims.(jwt.MapClaims)

		if err == nil && token.Valid {
			// return claims data

			claims := make(map[string]interface{})
			for k, v := range tokenClaims {
				switch k {

				case "exp":
					expiresIn := authBackend.GetTokenRemainingValidity(tokenClaims["exp"])
					claims["expires_in"] = expiresIn

				default:
					claims[k] = v
				}
			}

			response, _ := ffjson.Marshal(api.Claim{Claims: claims})
			httpResponse.Payload = response
			httpResponse.Status = http.StatusOK

		} else {

			// token is not valid
			response, _ := ffjson.Marshal(api.ErrorData{Error: api.ErrorInvalidToken, Message: api.ErrorMessages[api.ErrorInvalidToken]})
			httpResponse.Status = http.StatusBadRequest
			httpResponse.Payload = response

		}
	}

	return httpResponse
}

// Destroy executes the logout
func Destroy(request *api.Token, client string) *api.Response {

	httpResponse := new(api.Response)
	authBackend, errJWT := authentication.InitJWTAuthenticationBackend(services.Get().Cache)

	if errJWT != nil {
		httpResponse.Status = http.StatusBadRequest
		httpResponse.ErrorCode = api.ErrorInvalidClient
		return httpResponse
	}

	token, errParse := authBackend.ParseToken(request.Token, client)
	if errParse != nil {
		httpResponse.Status = http.StatusBadRequest
		httpResponse.ErrorCode = api.ErrorParsingRequest
		return httpResponse
	}

	err := authBackend.Destroy(token)
	if err != nil {
		httpResponse.Status = http.StatusInternalServerError
		httpResponse.ErrorCode = api.ErrorRedis
	} else {
		httpResponse.Status = http.StatusOK
	}

	return httpResponse
}
