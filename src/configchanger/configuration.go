package configchanger

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strconv"

	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/config"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/loader"
)

func SaveConfigForm(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parse the form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, fmt.Sprintf("Error parsing form: %v", err), http.StatusInternalServerError)
		return
	}

	// Load existing configuration
	existingConfig, err := config.LoadConfig()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error loading existing configuration: %v", err), http.StatusInternalServerError)
		return
	}

	// Use reflection to update fields present in the form
	v := reflect.ValueOf(existingConfig).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		fieldName := t.Field(i).Tag.Get("json")
		formValues, exists := r.Form[fieldName] // Check if the field exists in the form data

		if !exists {
			continue // Skip fields not submitted in the form
		}

		// If the field exists, use the first value (even if it's empty)
		formValue := ""
		if len(formValues) > 0 {
			formValue = formValues[0]
		}

		field := v.Field(i)
		fieldType := field.Type()

		switch fieldType.Kind() {
		case reflect.String:
			field.SetString(formValue) // Set the value, even if it's empty to allow clearing the field
		case reflect.Ptr:
			if fieldType.Elem().Kind() == reflect.Bool {
				newBool := formValue == "true"       // Convert form value to bool
				field.Set(reflect.ValueOf(&newBool)) // Set the pointer to the new bool
			}
		}
	}

	// Save the updated config
	if err := loader.SaveConfig(existingConfig); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func SaveConfigRestful(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading request body: %v", err), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	// Parse the JSON data into a map
	var requestData map[string]interface{}
	if err := json.Unmarshal(body, &requestData); err != nil {
		http.Error(w, fmt.Sprintf("Error parsing JSON: %v", err), http.StatusBadRequest)
		return
	}

	// Load existing configuration
	existingConfig, err := config.LoadConfig()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error loading existing configuration: %v", err), http.StatusInternalServerError)
		return
	}

	// Use reflection to update fields present in the request data
	v := reflect.ValueOf(existingConfig).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		fieldName := t.Field(i).Tag.Get("json")
		value, exists := requestData[fieldName] // Check if the field exists in the request data

		if !exists {
			continue // Skip fields not submitted in the request
		}

		field := v.Field(i)
		fieldType := field.Type()

		switch fieldType.Kind() {
		case reflect.String:
			if strValue, ok := value.(string); ok {
				field.SetString(strValue)
			}
		case reflect.Ptr:
			if fieldType.Elem().Kind() == reflect.Bool {
				if boolValue, ok := value.(bool); ok {
					field.Set(reflect.ValueOf(&boolValue)) // Set the pointer to the new bool
				}
			}
		case reflect.Int:
			switch v := value.(type) {
			case float64: // JSON numbers are parsed as float64 by default
				field.SetInt(int64(v))
			case int:
				field.SetInt(int64(v))
			case int64:
				field.SetInt(v)
			case string:
				if intValue, err := strconv.ParseInt(v, 10, 64); err == nil {
					field.SetInt(intValue)
				}
			}
		case reflect.Map:
			// Handle map fields like Users (map[string]string)
			if mapValue, ok := value.(map[string]interface{}); ok {
				if field.Type() == reflect.TypeOf(map[string]string{}) {
					newMap := make(map[string]string)
					for k, v := range mapValue {
						if strVal, ok := v.(string); ok {
							newMap[k] = strVal
						}
					}
					field.Set(reflect.ValueOf(newMap))
				}
			}
		}
	}

	// Save the updated config
	if err := loader.SaveConfig(existingConfig); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return success response in JSON format
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success", "message": "Configuration updated successfully"})
}
