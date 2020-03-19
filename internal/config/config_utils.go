package config

import (
	"github.com/subosito/gotenv"
	"path/filepath"
	"runtime"
)

// RootDir returns root dir of project
func RootDir() string {
	_, b, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(b), "../..")
}

// Load .env file in root path of module and return err
func LoadEnv() error {
	envFile := filepath.Join(RootDir(), envFileName)
	return gotenv.Load(envFile)
}

// MustLoadEnv is the same as LoadEnv but panics if an error occurs
func MustLoadEnv() {
	err := LoadEnv()
	if err != nil {
		panic(err)
	}
}
