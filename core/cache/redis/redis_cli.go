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

package redis

import (
	"errors"
	"time"

	"github.com/garyburd/redigo/redis"
)

// Pool holds a connection pool to a Redis server.
type Pool struct {
	Connection *redis.Pool
}

// Init creates a connection pool to Redis.
func (p *Pool) Init(params ...string) {
	p.Connection = &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", params[0])
			if err != nil {
				return nil, err
			}

			// if we have a password, do AUTH
			if params[1] != "" {
				if _, err := c.Do("AUTH", params[1]); err != nil {
					c.Close()
					return nil, err
				}
			}

			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
}

// SetValue creates/replace a key/value pair on Redis.
func (p *Pool) SetValue(key string, value string, expiration ...interface{}) error {
	if p.Connection != nil {

		conn := p.Connection.Get()
		defer conn.Close()

		_, err := conn.Do("SET", key, value)

		if err == nil && expiration != nil {
			p.Connection.Get().Do("EXPIRE", key, expiration[0])
		}

		return err
	}

	return errors.New("Redis cache not initialized")
}

// GetValue retrieves an existing key/value pair from the server.
func (p *Pool) GetValue(key string) (interface{}, error) {
	if p.Connection != nil {
		conn := p.Connection.Get()
		defer conn.Close()
		return conn.Do("GET", key)
	}

	return nil, errors.New("Redis cache not initialized")
}

// Close ends all the connecctions with Redis server
func (p *Pool) Close() {
	if p.Connection != nil {
		p.Connection.Close()
	}
}
