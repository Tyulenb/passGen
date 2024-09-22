package main

import (
	"encoding/json"
	"log"
	"os"
)

// Reads file
func readFile(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileStats, err := file.Stat()
	if err != nil {
		return nil, err
	}
	data := make([]byte, fileStats.Size())
	file.Read(data)
	return data, nil
}

// Writes file
func writeFile(path string, data []byte) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	file.Write(data)
	return nil
}

// Check if the program was executed for the first time
// Creates all necessary files
func checkNewUser() bool {
	_, err := os.Stat("user_data")
	if err != nil {
		if os.IsNotExist(err) {
			os.Mkdir("user_data", 0644)
			file, err := os.Create("user_data//config.json")
			if err != nil {
				log.Fatal(err)
			}
			defer file.Close()
			file.WriteString("[]")
			return true
		}
	}
	return false
}

// Creates user files for new user
// Updates config
func createNewUserData(name, hashKey string) error {
	data, err := readFile("user_data//config.json")
	if err != nil {
		return err
	}
	config := []User{}
	err = json.Unmarshal(data, &config)
	if err != nil {
		return err
	}
	config = append(config, User{name, hashKey})
	jsonConfig, err := json.Marshal(config)
	if err != nil {
		return nil
	}
	err = writeFile("user_data//config.json", jsonConfig)
	if err != nil {
		return err
	}

	file, err := os.Create("user_data//" + name + "_storage.json")
	if err != nil {
		return err
	}
	defer file.Close()
	file.WriteString("[]")

	return nil
}
