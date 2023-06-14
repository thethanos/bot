package parsers

import (
	"bytes"
	"fmt"
	"io"
	"multimessenger_bot/internal/entities"
	"os"
	"time"

	"github.com/tealeg/xlsx"
)

func ParseMasterData(data []byte) (*entities.Master, error) {

	file, err := xlsx.OpenBinary(data)
	if err != nil {
		return nil, err
	}
	fmt.Println(file.Sheet)
	return &entities.Master{
		ID:          fmt.Sprintf("%d", time.Now().Unix()),
		Name:        "Test",
		Description: "Test",
	}, nil
}

func SaveFile(name, path, ext string, data []byte) string {
	if err := os.MkdirAll(fmt.Sprintf("%s", path), os.ModePerm); err != nil {
		return ""
	}

	out, err := os.Create(fmt.Sprintf("%s/%s.%s", path, name, ext))
	if err != nil {
		return ""
	}
	defer out.Close()

	if _, err = io.Copy(out, bytes.NewReader(data)); err != nil {
		return ""
	}

	return fmt.Sprintf("%s/%s.%s", path, name, ext)
}
