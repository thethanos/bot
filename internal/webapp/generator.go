package webapp

import (
	"bytes"
	"html/template"
	"multimessenger_bot/internal/entities"
)

type webapp struct {
	Header  string
	Masters []*entities.Master
}

func GenerateWebPage(header string, masters []*entities.Master) ([]byte, error) {

	webapp := webapp{
		Header:  header,
		Masters: masters,
	}

	allFiles := []string{"content.tmpl", "footer.tmpl", "header.tmpl", "page.tmpl"}

	var allPaths []string
	for _, tmpl := range allFiles {
		allPaths = append(allPaths, "./webapp/masters/templates/"+tmpl)
	}

	templates := template.Must(template.New("").ParseFiles(allPaths...))

	var processed bytes.Buffer
	if err := templates.ExecuteTemplate(&processed, "page", webapp); err != nil {
		return nil, err
	}

	return processed.Bytes(), nil
}
