package aiModels

import "github.com/stealcash/AgentFlow/app/globals"

func GetChatGPTKeyByID(identifier string) (string, bool) {
	for _, model := range globals.Config.ChatGPTModels {
		if model.Identifier == identifier {
			return model.ChatGPTAPIKey, true
		}
	}
	return "", false
}

func GetChatGPTModelByID(identifier string) (*globals.ChatGPTModel, bool) {
	for _, model := range globals.Config.ChatGPTModels {
		if model.Identifier == identifier {
			return &model, true
		}
	}
	return nil, false
}
