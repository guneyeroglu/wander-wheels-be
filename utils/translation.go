package utils

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

func GetTranslation(lang, key string) string {
	type Translations map[string]map[string]string
	var translations Translations

	jsonData, err := os.Open("translations.json")

	if err != nil {
		log.Fatal(err)
	}

	data, err := io.ReadAll(jsonData)

	if err != nil {
		log.Fatal(err)
	}

	jsonString := string(data)

	err = json.Unmarshal([]byte(jsonString), &translations)

	if err != nil {
		log.Fatal(err)
	}

	if translations[lang] != nil && translations[lang][key] != "" {
		return translations[lang][key]
	}

	return ""
}
