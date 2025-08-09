package main

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type IConfig interface {
	ReadConfig(path string)
	Env() string
	String(param string) string
	Int(param string) int
	List(param string) []string
}

func NewConfig(path string) *Config {
	conf := &Config{}
	conf.path = path
	conf.init()
	return conf
}

type Config struct {
	path string
	v    *viper.Viper
}

func (s *Config) init() {
	s.v = viper.New()
	s.ReadConfig(s.path)
}

func (s *Config) ReadConfig(path string) {
	s.v.SetConfigFile(path)

	// Handle errors reading the config file
	if err := s.v.ReadInConfig(); err != nil {
		panic(fmt.Errorf("can't open config file. %w", err))
	}
}

func (s *Config) Env() string {
	if os.Getenv("APPLICATION_ENV") != "" {
		return os.Getenv("APPLICATION_ENV")
	}
	return string(EnvDevelopment)
}

func (s *Config) String(param string) string {
	return s.v.GetString(param)
}

func (s *Config) Int(param string) int {
	return s.v.GetInt(param)
}

func (s *Config) List(param string) []string {
	return s.v.GetStringSlice(param)
}
