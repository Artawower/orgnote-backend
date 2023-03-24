package configs

import (
	"os"

	"github.com/rs/zerolog/log"
)

type Config struct {
	AppAddress    string
	MongoURI      string
	Debug         bool
	MediaPath     string
	GithubID      string
	GithubSecret  string
	BackendSchema string
	BackendPort   string
	BackendDomain string
	ClientAddress string
}

func (c *Config) BackendHost() string {
	host := c.BackendSchema + "://" + c.BackendDomain
	if (c.BackendPort != "") {
		host += ":" + c.BackendPort
	}
	// TODO: master add version to environment and config
	return host + "/v1"
}

// TODO: master split into several functions
func NewConfig() Config {
	appAddress := "127.0.0.1:3000"
	if envAddr := os.Getenv("APP_ADDRESS"); envAddr != "" {
		appAddress = envAddr
	}

	mongoURI := "mongodb://127.0.0.1:27017"
	// TODO: master check correct enviroment variable
	if envMongoURL := os.Getenv("MONGO_URL"); envMongoURL != "" {
		mongoUser := os.Getenv("MONGO_USERNAME")
		mongoPassword := os.Getenv("MONGO_PASSWORD")
		mongoPort := os.Getenv("MONGO_PORT")
		mongoURI = "mongodb://" + mongoUser + ":" + mongoPassword + "@" + envMongoURL + ":" + mongoPort
	}

	envGithubID := os.Getenv("GITHUB_ID")
	envGithubSecret := os.Getenv("GITHUB_SECRET")

	if envGithubID == "" || envGithubSecret == "" {
		log.Warn().Msg("Github OAuth is not configured")
	}

	debug := os.Getenv("MODE") == "DEBUG"

	clientAddress := appAddress
	if envClientAddress := os.Getenv("CLIENT_ADDRESS"); envClientAddress != "" {
		clientAddress = envClientAddress
	}

	backendDomain := os.Getenv("BACKEND_DOMAIN")
	backendSchema := os.Getenv("BACKEND_SCHEMA")
	if (backendDomain == "" || backendSchema == "") {
		log.Fatal().Msg("BACKEND_DOMAIN or BACKEND_SCHEMA is not set")
	}
		
	backendPort := os.Getenv("BACKEND_PORT")

	config := Config{
		AppAddress:    appAddress,
		MongoURI:      mongoURI,
		Debug:         debug,
		MediaPath:     "./media",
		GithubID:      envGithubID,
		GithubSecret:  envGithubSecret,
		ClientAddress: clientAddress,
		BackendSchema: backendSchema,
		BackendDomain: backendDomain,
		BackendPort:   backendPort,

	}
	log.Info().Msgf("Config: %+v", config)

	return config
}
