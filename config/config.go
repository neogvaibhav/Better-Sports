package config

import (
	"log"
	"os"

	//Use this to fetch the package for environment variables
	"github.com/joho/godotenv"
)

// This will act as a file to store all the constants
type Config struct {
	serverHost   string
	mongoURL     string
	databaseName string
}

// To set the value of config object
func (c *Config) initialize() {
	//fetch value from local env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Can't fetch env variables")
		os.Exit(0)
	}
	c.mongoURL = os.Getenv("MONGO_URL")
	c.serverHost = os.Getenv("PORT")
	c.databaseName = os.Getenv("DB_NAME")
}

func (c *Config) GetMongoURL() string {
	return c.mongoURL
}

func (c *Config) GetDatabaseName() string {
	return c.databaseName
}

func (c *Config) GetServerHost() string {
	return c.serverHost
}

func NewConfig() *Config {
	config := new(Config)
	config.initialize()
	return config
}
