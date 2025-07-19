package globals

import (
	"github.com/spf13/viper"
	"log"
)

type GeneralQuestionsConfig struct {
	GeneralQuestions []struct {
		Answer    string   `mapstructure:"answer"`
		Questions []string `mapstructure:"questions"`
	} `mapstructure:"general_questions"`
}

var GeneralQuestions GeneralQuestionsConfig

func LoadGeneralQuestionsConfig() {
	viper.SetConfigName("general_questions")
	viper.AddConfigPath("./config")
	viper.SetConfigType("toml") // if you switch to YAML: set to "yaml"

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading general_questions TOML: %v", err)
	}

	if err := viper.Unmarshal(&GeneralQuestions); err != nil {
		log.Fatalf(" Error unmarshalling general_questions config: %v", err)
	}

	log.Println("Loaded general_questions fallback successfully.")
}
