package main

import (
	"bufio"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
)

type Dictionary struct {
	filePath   string
	addChan    chan entry
	removeChan chan string
}

type entry struct {
	word       string
	definition string
}

func NewDictionary(filePath string) *Dictionary {
	dict := &Dictionary{
		filePath:   filePath,
		addChan:    make(chan entry),
		removeChan: make(chan string),
	}

	go dict.handleAdditions()
	go dict.handleRemovals()

	return dict
}

func (d *Dictionary) Add(word, definition string) {
	d.addChan <- entry{word: word, definition: definition}
}

func (d *Dictionary) Remove(word string) {
	d.removeChan <- word
}

func main() {
	filePath := "C:/Users/MHD/OneDrive/Desktop/goo/Dictionnaires.txt"
	myDictionary := NewDictionary(filePath)

	http.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		word := r.FormValue("word")
		definition := r.FormValue("definition")

		if word == "" || definition == "" {
			http.Error(w, "Word or Definition missing", http.StatusBadRequest)
			return
		}

		myDictionary.Add(word, definition)
		w.WriteHeader(http.StatusCreated)
	})

	http.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		word := r.FormValue("word")
		if word == "" {
			http.Error(w, "Word missing", http.StatusBadRequest)
			return
		}

		file, err := os.Open(filePath)
		if err != nil {
			http.Error(w, "Dictionary file not found", http.StatusInternalServerError)
			return
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			parts := strings.Split(scanner.Text(), ":")
			if len(parts) == 2 && parts[0] == word {
				w.Write([]byte(parts[1]))
				return
			}
		}

		http.Error(w, "Word not found", http.StatusNotFound)
	})

	http.HandleFunc("/remove", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, "Method not allowed Mr Zied !", http.StatusMethodNotAllowed)
			return
		}

		word := r.FormValue("word")
		if word == "" {
			http.Error(w, "Word missing", http.StatusBadRequest)
			return
		}

		myDictionary.Remove(word)
		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("/list", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		file, err := os.Open(filePath)
		if err != nil {
			http.Error(w, "Dictionary file not found", http.StatusInternalServerError)
			return
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		result := make(map[string]string)
		for scanner.Scan() {
			parts := strings.Split(scanner.Text(), ":")
			if len(parts) == 2 {
				result[parts[0]] = parts[1]
			}
		}

		jsonBytes, err := json.Marshal(result)
		if err != nil {
			http.Error(w, "Error converting to JSON", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonBytes)
	})
	log.Fatal(http.ListenAndServe(":8080", nil))

}
func (d *Dictionary) handleAdditions() {
	// Logic to handle additions
}

func (d *Dictionary) handleRemovals() {
	// Logic to handle removals
}
