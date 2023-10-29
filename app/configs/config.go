package configs

import (
	"os"
	"strconv"

	"github.com/davecgh/go-spew/spew"
	"github.com/rs/zerolog/log"
)

type Config struct {
	AppAddress               string
	MongoURI                 string
	Debug                    bool
	MediaPath                string
	GithubID                 string
	GithubSecret             string
	BackendSchema            string
	BackendPort              string
	BackendDomain            string
	ClientAddress            string
	MobileAppName            string
	AccessCheckerURL         *string
	AccessCheckToken         *string
	AccessTokenCacheLifeTime int
	MaximumFileSize          int
}

func (c *Config) BackendHost() string {
	host := c.BackendSchema + "://" + c.BackendDomain
	if c.BackendPort != "" {
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

	debug := os.Getenv("DEBUG") == "true"

	clientAddress := appAddress
	if envClientAddress := os.Getenv("CLIENT_ADDRESS"); envClientAddress != "" {
		clientAddress = envClientAddress
	}

	backendDomain := os.Getenv("BACKEND_DOMAIN")
	backendSchema := os.Getenv("BACKEND_SCHEMA")
	if backendDomain == "" || backendSchema == "" {
		log.Fatal().Msg("BACKEND_DOMAIN or BACKEND_SCHEMA is not set")
	}

	var accessCheckerURL *string
	if envAccessCheckURL := os.Getenv("ACCESS_CHECK_URL"); envAccessCheckURL != "" {
		accessCheckerURL = &envAccessCheckURL
	}

	var accessCheckToken *string
	if envAccessCheckToken := os.Getenv("ACCESS_CHECK_TOKEN"); envAccessCheckToken != "" {
		accessCheckToken = &envAccessCheckToken
	}

	accessTokenCacheLifeTime := 60

	if envAccessTokenCacheLifeTime := os.Getenv("ACCESS_TOKEN_CACHE_LIFE_TIME"); envAccessTokenCacheLifeTime != "" {
		val, err := strconv.Atoi(envAccessTokenCacheLifeTime)
		if err != nil {
			accessTokenCacheLifeTime = val
		} else {
			log.Warn().Msgf("ACCESS_TOKEN_CACHE_LIFE_TIME is not a number, init with default value: %d", accessTokenCacheLifeTime)
		}
	}

	maximumFileSize := 1024 * 1024 * 10
	if envMaximumFileSize := os.Getenv("MAXIMUM_FILE_SIZE"); envMaximumFileSize != "" {
		val, err := strconv.Atoi(envMaximumFileSize)
		if err != nil {
			maximumFileSize = val
		} else {
			log.Warn().Msgf("MAXIMUM_FILE_SIZE is not a number, init with default value: %d", maximumFileSize)
		}
	}

	backendPort := os.Getenv("BACKEND_PORT")

	config := Config{
		AppAddress:               appAddress,
		MongoURI:                 mongoURI,
		Debug:                    debug,
		MediaPath:                "./media",
		GithubID:                 envGithubID,
		GithubSecret:             envGithubSecret,
		ClientAddress:            clientAddress,
		BackendSchema:            backendSchema,
		BackendDomain:            backendDomain,
		BackendPort:              backendPort,
		AccessCheckerURL:         accessCheckerURL,
		AccessCheckToken:         accessCheckToken,
		AccessTokenCacheLifeTime: accessTokenCacheLifeTime,
		MaximumFileSize:          maximumFileSize,
		MobileAppName:            "orgnote",
	}
	log.Info().Msgf("Config: %+v", spew.Sdump(config))

	return config
}
