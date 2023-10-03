package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type Display struct {
	Diagonal   float64
	Resolution string
	Type       string
	GSync      bool
}

type Monitor struct {
	PowerColtage int
	Display      Display
	DisplayID    int
	GSyncPremium bool
	IsCurved     bool
}

type User struct {
	ID_User       int    `json:"id"`
	Username_User string `json:"username"`
	Password_User string `json:"password"`
	Email_User    string `json:"email"`
	Is_Admin_User bool   `json:"isadmin"`
}
type LoginStruct struct {
	User_Login    string `json:"login"`
	User_Password string `json:"password"`
}

func AddDisplay() Display {
	var display Display
	fmt.Println("Введите диагональ: ")
	fmt.Scanln(&display.Diagonal)
	fmt.Println("Введите разрешение экрана: ")
	fmt.Scanln(&display.Resolution)
	fmt.Println("Введите тип дисплея: ")
	fmt.Scanln(&display.Type)
	fmt.Println("Есть ли поддержка G-Sync: ")
	fmt.Scanln(&display.GSync)
	return display
}

func AddMonitor() Monitor {
	var monitor Monitor
	fmt.Println("Введите напряжение питания: ")
	fmt.Scanln(&monitor.PowerColtage)
	fmt.Println("Введите ID дисплея: ")
	fmt.Scanln(&monitor.DisplayID)
	fmt.Println("Поддерживает ли монитор G-Sync Premium: ")
	fmt.Scanln(&monitor.GSyncPremium)
	fmt.Println("Является ли монитор изогнутым: ")
	fmt.Scanln(&monitor.IsCurved)
	return monitor
}

func AddUser() User {
	var user User
	fmt.Println("Введите username: ")
	fmt.Scanln(&user.Username_User)
	fmt.Println("Введите password: ")
	fmt.Scanln(&user.Password_User)
	fmt.Println("Введите email: ")
	fmt.Scanln(&user.Email_User)
	fmt.Println("Является пользователь админом? (true/false): ")
	fmt.Scanln(&user.Is_Admin_User)
	return user
}

// var username_user string
// var password_user string
var TOKEN string = ""

func main() {

	client := &http.Client{}

	for TOKEN == "" {
		var loginForm LoginStruct

		fmt.Println("Введите логин")
		fmt.Scanln(&loginForm.User_Login)
		fmt.Println("Введите пароль")
		fmt.Scanln(&loginForm.User_Password)

		usr_data, err := json.Marshal(loginForm)
		if err != nil {
			log.Fatal(err)
		}

		req, _ := http.NewRequest("POST", "http://127.0.0.1:8080/login", bytes.NewReader(usr_data))
		response, err := client.Do(req)
		if err != nil {
			log.Fatal("[CLIENT] ERROR", err)
		}
		defer response.Body.Close()
		out, err := io.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}
		TOKEN = string(out)

		log.Println(string(out))

	}
	for {
		fmt.Println("Выберите действие: ")
		fmt.Println("1) Добавить дисплей")
		fmt.Println("2) Добавить монитор")
		fmt.Println("3) Посмотреть все мониторы")
		fmt.Println("4) Посмотреть все дисплеи")
		fmt.Println("5) Удалить монитор")
		fmt.Println("6) Удалить дисплей")
		fmt.Println("7) Посмотреть монитор по ID")
		fmt.Println("8) Регистрация")
		fmt.Println("9) Выйти")

		var choice int
		fmt.Scanln(&choice)

		switch choice {
		case 1:

			newdispl := AddDisplay()
			newj, err := json.Marshal(newdispl)
			if err != nil {
				log.Fatal(err)
			}
			req, _ := http.NewRequest("POST", "http://127.0.0.1:8080/addDisplay", bytes.NewReader(newj))
			req.Header.Add("Authorization", fmt.Sprintf("%x", TOKEN))
			response, err := client.Do(req)
			if err != nil {
				log.Fatal("[CLIENT] ERROR", err)
			}
			defer response.Body.Close()
			out, err := io.ReadAll(response.Body)
			if err != nil {
				log.Fatal(err)
			}
			log.Println(string(out))
		case 2:
			newmonitor := AddMonitor()
			newj, err := json.Marshal(newmonitor)
			if err != nil {
				log.Fatal(err)
			}
			req, _ := http.NewRequest("POST", "http://127.0.0.1:8080/addMonitor", bytes.NewReader(newj))
			req.Header.Add("Authorization", fmt.Sprintf("%x", TOKEN))
			response, err := client.Do(req)
			if err != nil {
				log.Fatal("[CLIENT] ERROR", err)
			}
			defer response.Body.Close()
			out, err := io.ReadAll(response.Body)
			if err != nil {
				log.Fatal(err)
			}
			log.Println(string(out))

		case 3:

			req, _ := http.NewRequest("GET", "http://127.0.0.1:8080/allMonitors", nil)
			req.Header.Add("Authorization", fmt.Sprintf("%x", TOKEN))
			response, err := client.Do(req)
			if err != nil {
				log.Fatal("[CLIENT] ERROR", err)
			}

			defer response.Body.Close()
			out, err := io.ReadAll(response.Body)
			if err != nil {
				log.Fatal(err)
			}
			log.Println(string(out))

		case 4:
			req, _ := http.NewRequest("GET", "http://127.0.0.1:8080/allDisplays", nil)
			req.Header.Add("Authorization", fmt.Sprintf("%x", TOKEN))
			response, err := client.Do(req)
			if err != nil {
				log.Fatal("[CLIENT] ERROR", err)
			}

			defer response.Body.Close()
			out, err := io.ReadAll(response.Body)
			if err != nil {
				log.Fatal(err)
			}
			log.Println(string(out))
		case 5:
			var delword string
			fmt.Println("Введите название монитора для удаления: ")
			fmt.Scanln(&delword)
			req, _ := http.NewRequest("POST", "http://127.0.0.1:8080/removeMonitor", bytes.NewReader([]byte(delword)))
			req.Header.Add("Authorization", fmt.Sprintf("%x", TOKEN))
			response, err := client.Do(req)
			if err != nil {
				log.Fatal("[CLIENT] ERROR", err)
			}
			defer response.Body.Close()
			out, err := io.ReadAll(response.Body)
			if err != nil {
				log.Fatal(err)
			}
			log.Println(string(out))
		case 6:
			var delword string
			fmt.Println("Введите ID дисплея для удаления: ")
			fmt.Scanln(&delword)
			req, _ := http.NewRequest("POST", "http://127.0.0.1:8080/removeDisplay", bytes.NewReader([]byte(delword)))
			req.Header.Add("Authorization", fmt.Sprintf("%x", TOKEN))
			response, err := client.Do(req)
			if err != nil {
				log.Fatal("[CLIENT] ERROR", err)
			}
			defer response.Body.Close()
			out, err := io.ReadAll(response.Body)
			if err != nil {
				log.Fatal(err)
			}
			log.Println(string(out))
		case 7:
			var searchWord string
			fmt.Println("Введите ID дисплея для поиска: ")
			fmt.Scanln(&searchWord)
			req, _ := http.NewRequest("POST", "http://127.0.0.1:8080/getMonitor", bytes.NewReader([]byte(searchWord)))
			req.Header.Add("Authorization", fmt.Sprintf("%x", TOKEN))
			response, err := client.Do(req)
			if err != nil {
				log.Fatal("[CLIENT] ERROR", err)
			}
			defer response.Body.Close()
			out, err := io.ReadAll(response.Body)
			if err != nil {
				log.Fatal(err)
			}
			log.Println(string(out))
		case 8:
			newUser := AddUser()
			newu, err := json.Marshal(newUser)
			if err != nil {
				log.Fatal(err)
			}
			req, _ := http.NewRequest("POST", "http://127.0.0.1:8080/addUser", bytes.NewReader(newu))
			req.Header.Add("Authorization", fmt.Sprintf("%x", TOKEN))
			response, err := client.Do(req)
			if err != nil {
				log.Fatal("[CLIENT] ERROR", err)
			}
			defer response.Body.Close()

			out, err := io.ReadAll(response.Body)
			if err != nil {
				log.Fatal(err)
			}
			log.Println(string(out))
		case 9:
			fmt.Println("Программа завершена")
			return
		default:
			fmt.Println("Выбирете другой параметр")
		}

	}

}
