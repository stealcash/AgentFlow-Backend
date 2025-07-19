package globals

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"log"
	"os"
	"strconv"
	"strings"
)

var RootDirPath string

var Config config

type config struct {
	Regional        configRegional
	App             configApp
	MainDatabase    configMainDatabase
	ElasticDatabase configElasticDatabase
	Server          configServer
	ChatGPTModels   []ChatGPTModel `mapstructure:"chatgpt_models"`
	TrustedProxies  []string       `mapstructure:"trusted_proxies"`
	AllowedOrigins  []string
	FrontEnd        configFrontEnd
}

type ChatGPTModel struct {
	Identifier    string `mapstructure:"identifier"`
	ModelName     string `mapstructure:"model_name"`
	ChatGPTAPIKey string `mapstructure:"chatgpt_api_key"`
}

type configRegional struct {
	Language string
	Timezone string
}
type configApp struct {
	Port         string
	Env          string
	JWTSecretKey string `mapstructure:"jwt_secret_key"`
	ApiUrl       string `mapstructure:"api_url"`
}

type configMainDatabase struct {
	User     string
	Password string
	Name     string
	Host     string
	Port     int
	SSLMode  string
}

type configElasticDatabase struct {
	User                      string
	Password                  string
	Name                      string
	Host                      string
	Port                      string
	RequiredElasticConnection bool
}

type configFrontEnd struct {
	Path string
}

type configServer struct {
	APIBaseURL string
}

func InitConfiguration() {
	// Load .env
	err := godotenv.Load()
	if err != nil {
		log.Println(" No .env file found, continuing...")
	}

	// Load config.toml
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(getEnv("CONFIG_FOLDER", "./config"))

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf(" Error reading config.toml: %s", err)
	}

	if err := viper.Unmarshal(&Config); err != nil {
		log.Fatalf(" Unable to decode config.toml: %s", err)
	}

	// Merge .env values
	Config.ElasticDatabase.Host = getEnv("ELASTIC_SEARCH_HOST", "")
	Config.ElasticDatabase.Port = getEnv("ELASTIC_SEARCH_PORT", "")
	Config.ElasticDatabase.User = getEnv("ELASTIC_SEARCH_USERNAME", "")
	Config.ElasticDatabase.Password = getEnv("ELASTIC_SEARCH_PASSWORD", "")
	Config.ElasticDatabase.RequiredElasticConnection = getEnvBool("ELASTIC_CONNECTION_REQUIRED", true)

	// AllowedOrigins — from .env only
	origins := getEnv("ALLOWED_ORIGINS", "")
	if origins != "" {
		Config.AllowedOrigins = parseCommaList(origins)
	}

	log.Println(" Configuration initialized successfully (TOML + ENV)")
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
func getEnvBool(key string, fallback bool) bool {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	result, err := strconv.ParseBool(val)
	if err != nil {
		log.Printf(" Invalid bool for %s: %s, using fallback %v", key, val, fallback)
		return fallback
	}
	return result
}

func parseCommaList(s string) []string {
	items := strings.Split(s, ",")
	var trimmed []string
	for _, item := range items {
		item = strings.TrimSpace(item)
		if item != "" {
			trimmed = append(trimmed, item)
		}
	}
	return trimmed
}
