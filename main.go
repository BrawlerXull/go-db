package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func main() {

	data := []map[string]interface{}{}
	var count int

	http.HandleFunc("/all", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			jsonData, err := json.Marshal(data)
			if err != nil {
				http.Error(w, "Unable to marshal JSON", http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, err = w.Write(jsonData)
			if err != nil {
				http.Error(w, "Unable to write response", http.StatusInternalServerError)
			}
		} else {
			http.Error(w, "Method not supported", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/send", func(w http.ResponseWriter, r *http.Request) {
		jsonData := make(map[string]interface{})
		if r.Method == http.MethodPost {
			decoder := json.NewDecoder(r.Body)

			if err := decoder.Decode(&jsonData); err != nil {
				http.Error(w, "Unable to decode JSON", http.StatusBadRequest)
				return
			}

			jsonData["id"] = count
			count++

			data = append(data, jsonData)
			fmt.Println(data)

			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "Successfully decoded")
		} else {
			http.Error(w, "Unsupported HTTP method", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/delete", func(w http.ResponseWriter, r *http.Request) {
		requestData := make(map[string]interface{})
		if r.Method == http.MethodPost {
			decoder := json.NewDecoder(r.Body)

			if err := decoder.Decode(&requestData); err != nil {
				http.Error(w, "Unable to decode JSON", http.StatusBadRequest)
				return
			}
			newData := []map[string]interface{}{}

			for _, value := range data {
				match := true
				for key, reqDataValue := range requestData {
					if itemValue, exists := value[key]; exists {
						if itemValue == reqDataValue {
							match = false
							break
						}
					}
				}
				if match {
					newData = append(newData, value)
				}
			}
			data = newData
			fmt.Println(data)

			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "Successfully deleted matching data")
		} else {
			http.Error(w, "Unsupported HTTP method", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/update", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPatch {
			decoder := json.NewDecoder(r.Body)
			requestData := make(map[string]interface{})
			if err := decoder.Decode(&requestData); err != nil {
				http.Error(w, "Unable to decode JSON", http.StatusBadRequest)
				return
			}
			if len(requestData) == 2 && requestData["toChange"] != nil && requestData["finalChange"] != nil {

				toChange := requestData["toChange"].(map[string]interface{})
				finalChange := requestData["finalChange"].(map[string]interface{})

				for key := range toChange {
					for _, dataValue := range data {
						if valueToChange, exists := toChange[key]; exists {
							if finalValue, finalExists := finalChange[key]; finalExists {
								if dataValue[key] == valueToChange {
									dataValue[key] = finalValue
								}
							}
						}
					}
				}

				fmt.Println(data)
				w.WriteHeader(http.StatusOK)
				fmt.Fprint(w, "Successfully updated data")
			} else {
				http.Error(w, "Bad Request", http.StatusBadRequest)
			}
		} else {
			http.Error(w, "Unsupported HTTP method", http.StatusMethodNotAllowed)
		}
	})

	fmt.Println("Listening on port 8080")
	http.ListenAndServe(":8080", nil)
}
