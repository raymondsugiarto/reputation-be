package config

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/gofiber/fiber/v2/log"
	"github.com/spf13/viper"
)

// Config :
type Config struct {
	Server   ServerList
	Database DatabaseList
	Groq     Groq
	Minimax  Minimax
}

var configuration Config

var (
	_, b, _, _   = runtime.Caller(0)
	basepath     = filepath.Dir(b)
	resourcepath = "/resources"
)

func init() {
	viper.AddConfigPath(basepath + resourcepath)
	viper.SetConfigType("yml")
	viper.SetConfigName("server.yml")
	errConf := viper.ReadInConfig()
	if errConf != nil {
		log.Fatalf("failed to load config: %v", errConf)
	}

	viper.SetConfigName("database.yml")
	errDatabase := viper.MergeInConfig()
	if errDatabase != nil {
		panic(fmt.Errorf("cannot load database config: %v", errDatabase))
	}

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	viper.Unmarshal(&configuration)

}

// GetConfig get config
func GetConfig() *Config {
	return &configuration
}
