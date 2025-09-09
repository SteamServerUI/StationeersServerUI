package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type Translation struct {
	Key   string
	Value string
}

type LanguageData struct {
	Language    string
	MissingKeys []Translation
}

func flattenJSON(prefix string, data map[string]interface{}, result map[string]string, orderedKeys *[]string) {
	for key, value := range data {
		newKey := key
		if prefix != "" {
			newKey = prefix + "." + key
		}
		*orderedKeys = append(*orderedKeys, newKey)
		switch v := value.(type) {
		case map[string]interface{}:
			flattenJSON(newKey, v, result, orderedKeys)
		case string:
			result[newKey] = v
		}
	}
}

func loadJSONFile(filename string) (map[string]interface{}, map[string]string, []string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return make(map[string]interface{}), make(map[string]string), []string{}, nil
		}
		return nil, nil, nil, err
	}
	var jsonData map[string]interface{}
	if err := json.Unmarshal(data, &jsonData); err != nil {
		return nil, nil, nil, err
	}
	flatData := make(map[string]string)
	orderedKeys := []string{}
	flattenJSON("", jsonData, flatData, &orderedKeys)
	return jsonData, flatData, orderedKeys, nil
}

func main() {
	localizationDir := "./UIMod/onboard_bundled/localization/"

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Load en-US.json
		_, enFlat, enOrderedKeys, err := loadJSONFile(filepath.Join(localizationDir, "en-US.json"))
		if err != nil {
			http.Error(w, fmt.Sprintf("Error loading en-US.json: %v", err), http.StatusInternalServerError)
			return
		}

		// Load de-DE.json and sv-SE.json
		languages := []string{"de-DE", "sv-SE"}
		data := make([]LanguageData, 0, len(languages))
		for _, lang := range languages {
			_, langFlat, _, err := loadJSONFile(filepath.Join(localizationDir, lang+".json"))
			if err != nil {
				http.Error(w, fmt.Sprintf("Error loading %s.json: %v", lang, err), http.StatusInternalServerError)
				return
			}
			missing := []Translation{}
			for _, key := range enOrderedKeys {
				if _, exists := langFlat[key]; !exists {
					if enVal, exists := enFlat[key]; exists {
						missing = append(missing, Translation{Key: key, Value: enVal})
					}
				}
			}
			data = append(data, LanguageData{
				Language:    lang,
				MissingKeys: missing,
			})
		}

		// Parse and execute template
		tmpl, err := template.New("index").Parse(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Translation Manager</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        h1 { text-align: center; }
        .lang-section { margin-bottom: 40px; }
        .lang-section h2 { color: #333; }
        table { width: 100%; border-collapse: collapse; margin-top: 10px; }
        th, td { border: 1px solid #ddd; padding: 8px; text-align: left; }
        th { background-color: #f2f2f2; }
    </style>
</head>
<body>
    <h1>Translation Manager - Missing Keys</h1>
    {{range .}}
    <div class="lang-section">
        <h2>{{.Language}}</h2>
        {{if .MissingKeys}}
        <table>
            <tr>
                <th>Key</th>
                <th>en-US Value</th>
            </tr>
            {{range .MissingKeys}}
            <tr>
                <td>{{.Key}}</td>
                <td>{{.Value}}</td>
            </tr>
            {{end}}
        </table>
        {{else}}
        <p>No missing translations for {{.Language}}.</p>
        {{end}}
    </div>
    {{end}}
</body>
</html>
`)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error parsing template: %v", err), http.StatusInternalServerError)
			return
		}
		if err := tmpl.Execute(w, data); err != nil {
			http.Error(w, fmt.Sprintf("Error executing template: %v", err), http.StatusInternalServerError)
		}
	})

	fmt.Println("Server starting on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
