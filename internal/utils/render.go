package utils

import (
	"bytes"
	"fmt"
	"text/template"
)

func Render(path string, data interface{}) ([]byte, error) {
	tmp, err := template.ParseFiles(path)
	if err != nil {
		return []byte("Ошибка!"), fmt.Errorf("[utils.render.parse] %w", err)
	}

	var configBuffer bytes.Buffer
	if err := tmp.Execute(&configBuffer, data); err != nil {
		return []byte("Ошибка!"), fmt.Errorf("[utils.render.execute] %w", err)
	}

	return []byte(configBuffer.String()), nil
}
