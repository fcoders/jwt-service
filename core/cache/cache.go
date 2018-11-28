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

package cache

// Connector is the interface used to hande the cache system configured.
// Using this interface, you can create connections with Redis, Memcache, and so on.
type Connector interface {
	Init(params ...string)
	GetValue(key string) (interface{}, error)
	SetValue(key string, value string, params ...interface{}) error
	Close()
}
