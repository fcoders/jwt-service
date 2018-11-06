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

package logger

// Default log levels
const (
	LevelError = iota
	LevelInfo
	LevelDebug
)

// Logger presents a common interface for logger
type Logger interface {
	Debug(v ...interface{})
	Debugf(format string, v ...interface{})
	Info(v ...interface{})
	Infof(format string, v ...interface{})
	Error(v ...interface{})
	Errorf(format string, v ...interface{})
	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})
	Print(v ...interface{})
	SetLevel(level int)
	GetLevel() int
	LogToSlack(webHook, title, text string)
	LogErrorToSlack(webHook, title, text string)
}

var (
	logImpl Logger
)

func SetLogger(l Logger) {
	logImpl = l
}

// GetLogger returns the current logger defined for the service
func GetLogger() Logger {
	return logImpl
}

// SetLogLevel configures the application's log level
func SetLogLevel(level int) {
	if logImpl != nil {
		logImpl.SetLevel(level)
	}
}
