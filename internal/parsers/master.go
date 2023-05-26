package parsers

import (
	"fmt"
	"multimessenger_bot/internal/entities"
	"time"
)

func ParseMasterData(data string) (*entities.Master, error) {
	return &entities.Master{
		ID:          fmt.Sprintf("%d", time.Now().Unix()),
		Name:        "Test",
		Description: "Test",
	}, nil
}
