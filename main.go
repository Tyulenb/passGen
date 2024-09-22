package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/AlecAivazis/survey/v2"
)

type Seed struct {
	Key     int `survey:"key-digit"`
	Length  int
	Service string
}

type User struct {
	Name string `survey:"name"`
	Key  string `survey:"key"`
}

type Storage struct {
	Service  string
	Password string
}

func main() {
	var currentSession User

	if checkNewUser() {
		registration()
	}

	answer := ""
	survey.AskOne(StartMenu, &answer)
	if answer[5:] == "up" {
		currentSession = registration()
	} else {
		currentSession = login()
	}

	if (currentSession == User{}) {
		log.Fatal("currentSession is empty")
	}

	for { //infinite loop
		survey.AskOne(MainSelect, &answer)

		if answer[2] == 'n' {
			pasSeed := Seed{}
			survey.Ask(GeneratePassQuestions, &pasSeed)
			pass := generatePass(&pasSeed)

			if pasSeed.Service == "" { //Without saving password into storage case
				fmt.Println("Your password is ", pass)
			}

			storageByte, err := readFile("user_data//" + currentSession.Name + "_storage.json")
			if err != nil {
				log.Fatal(err)
			}

			var userPasswords []Storage
			if string(storageByte) != "[]" { // Error occures if storageByte is empty
				userPasswords, err = decrypt(string(storageByte), []byte(currentSession.Key))
				if err != nil {
					log.Fatal(err)
				}
			}

			userPasswords = append(userPasswords, Storage{pasSeed.Service, pass})

			storageString, err := encrypt(userPasswords, []byte(currentSession.Key))
			if err != nil {
				log.Fatal(err)
			}

			err = writeFile("user_data//"+currentSession.Name+"_storage.json", []byte(storageString))
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("Your password is ", pass, " , successfully stored")

		} else if answer[2] == 't' {
			storageByte, err := readFile("user_data//" + currentSession.Name + "_storage.json")
			if err != nil {
				log.Fatal(err)
			}

			userPasswords, err := decrypt(string(storageByte), []byte(currentSession.Key))
			if err != nil {
				log.Fatal(err)
			}
			for _, v := range userPasswords {
				fmt.Println(v.Service, " - ", v.Password)
			}
			break
		} else {
			break
		}
	}
}

func registration() User {
	name := ""
	survey.AskOne(Registration, &name, survey.WithValidator(survey.Required))
	key := generateKey()
	hashedKey, err := hashKey([]byte(key))
	if err != nil {
		log.Fatal(err)
	}
	err = createNewUserData(name, string(hashedKey))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("YOUR KEY IS ", key, " KEEP IT SAFE IT WILL BE USED FOR YOUR FUTURE AUTHORIZATIONS!")
	return User{name, key}
}

func login() User {
	user := User{}
	survey.Ask(Authorization, &user)
	config, err := readFile("user_data//config.json")
	if err != nil {
		log.Fatal(err)
	}
	allUsers := []User{}
	err = json.Unmarshal(config, &allUsers)
	if err != nil {
		log.Fatal(err)
	}
	for _, v := range allUsers {
		if v.Name == user.Name {
			if compareHashKey([]byte(user.Key), []byte(v.Key)) {
				return user
			}
		}
	}
	return User{}
}

/*
	file -> read -> decrypt -> add -> encrypt -> write -> file
	user -> reg -> config.json {name: key(bcrypt)} -> name.json name of pass storage of user -> crypts with the key
	user -> login -> name + key -> check key in json with bcrypt
	storage crypts with the user key
	key is users password and key of his crypted storage same time

	//TO DO fix bug with displaying passwords and with displaying printing info
*/
