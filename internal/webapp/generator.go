package webapp

import (
	"bytes"
	"html/template"
	"multimessenger_bot/internal/entities"
)

func GenerateMassterCard(master *entities.Master) (string, error) {

	var allPaths []string
	allPaths = append(allPaths, "./webapp/templates/master_card.tmpl")

	var processed bytes.Buffer
	template := template.Must(template.New("").ParseFiles(allPaths...))
	if err := template.ExecuteTemplate(&processed, "master_card", master); err != nil {
		return "", err
	}

	return processed.String(), nil
}
