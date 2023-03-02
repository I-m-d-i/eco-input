package configs

import (
	"encoding/json"
	"log"
	"os"
	"sync"
)

// Config config running applications
type Config struct {
	ConnStr string `json:"connect_string"`
}

var (
	_config     *Config
	_onceConfig sync.Once
)

// GetConfig получение объекта конфига
func GetConfig() *Config {
	_onceConfig.Do(func() {
		_config = new(Config)
		_config.load()
	})
	return _config
}

func (c *Config) load() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	confFile, err := os.Open("configs/config.json")
	if err != nil {
		log.Fatal(err)
	}

	dc := json.NewDecoder(confFile)
	if err := dc.Decode(&c); err != nil {
		log.Fatal("Read Config file: ", err)
	}

	if c.ConnStr == "" {
		log.Fatal("Can`t read connection string: ", c.ConnStr)
	}
}
