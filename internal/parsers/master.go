package parsers

import "multimessenger_bot/internal/entities"

func ParseMasterData(data string) (*entities.Master, error) {
	return &entities.Master{
		ID:          "123",
		Name:        "Test",
		Description: "Test",
		Images: []string{
			"https://bot-dev-domain.com/masters/images/maria_ernandes/1.png",
		},
	}, nil
}
