package main

import (
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"fmt" //Пакет для форматированного вывода на консоль.
	"io"
	"log"
	"net/http" //Пакет для работы с HTTP-сервером.

	_ "github.com/lib/pq"
	"github.com/nats-io/nats.go"
)

type DisplayInfo struct { //Структура
	ID_Display string `json:"id"`
	Diagonal   int    `json:"diagonal"`
	Resolution string `json:"resolsution"`
	Type       string `json:"type"`
	GSync      bool   `json:"gsync"`
}

type MonitorInfo struct { //Структура
	ID_Monitor   string      `json:"id"`
	PowerVoltage int         `json:"powerVoltage"`
	Display      DisplayInfo `json:"display"`
	DisplayID    int         `json:"displayID"`
	GSyncPremium bool        `json:"gSyncPremium"`
	IsCurved     bool        `json:"isCurved"`
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

var natsURL string = "nats://95.165.107.100:4222"

var Nc *nats.Conn
var connStr = "user=postgres password=admin12345 dbname=shop sslmode=disable"
var tokens = make(map[string]User)

func main() { // запускает веб-сервер, который слушает на порту 8080

	fmt.Println("Запуск сервера...")
	startServer()

}

func startServer() {

	http.HandleFunc("/addDisplay", addDisplayHandler)       //обрабатывают запросы для добавления новых дисплеев в мапы.
	http.HandleFunc("/addMonitor", addMonitorHandler)       //обрабатывают запросы для добавления новых мониторов в мапы.
	http.HandleFunc("/removeDisplay", removeDisplayHandler) //обрабатывают запросы для удаления дисплеев из мапов.
	http.HandleFunc("/removeMonitor", removeMonitorHandler) //обрабатывают запросы для удаления мониторов из мапов.
	http.HandleFunc("/allDisplays", allDisplaysHandler)     //обрабатывают запросы для получения информации о всех дисплеях.
	http.HandleFunc("/allMonitors", allMonitorsHandler)     //обрабатывают запросы для получения информации о всех мониторах.
	http.HandleFunc("/getMonitor", getMonitorHandler)
	http.HandleFunc("/addUser", addUserHandler)
	http.HandleFunc("/login", loginUser)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}

}

func sendNats(theme string, text string) {
	Nc, err := nats.Connect(natsURL)
	fmt.Println(err)
	Nc.Publish(theme, []byte(text))
}

func checkLogin(token string) bool {
	_, ok := tokens[token]
	// If the key exists
	if ok {
		return true
	}
	return false
}

func loginUser(w http.ResponseWriter, r *http.Request) { //Функция возвращают информацию о всех дисплеях в формате текста, который будет отправлен клиенту в ответ на HTTP-запрос.
	tempLogin := LoginStruct{}
	body, _ := io.ReadAll(r.Body)
	err := json.Unmarshal(body, &tempLogin)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
		sendNats("error", err.Error())
	}
	defer db.Close()

	data := db.QueryRow("select * from Users where Username_User = $1", tempLogin.User_Login)
	var tempData User
	data.Scan(&tempData.ID_User, &tempData.Username_User, &tempData.Password_User, &tempData.Email_User, &tempData.Is_Admin_User)
	defer db.Close()
	loginPass := sha256.Sum256([]byte(tempLogin.User_Password))
	for fmt.Sprintf("%x", loginPass) == tempData.Password_User {
		userHashToken := sha256.Sum224([]byte(tempLogin.User_Login + fmt.Sprintf("%x", loginPass) + string(tempData.ID_User)))
		tokens[fmt.Sprintf("%x", userHashToken[:])] = tempData
		w.Write([]byte(userHashToken[:]))
		sendNats("log", fmt.Sprintf("Authorization success for ID_User: %v with token: %v", tempData.ID_User, fmt.Sprintf("%x", userHashToken[:])))
		return
	}
}

func addMonitorHandler(w http.ResponseWriter, r *http.Request) { //обработчик HTTP-запросов получает данные из запроса, а затем вносит изменения в соответствующие мапы (displayInfoMap или monitorInfoMap)
	tempToken := r.Header.Get("Authorization")
	if checkLogin(tempToken) && tokens[tempToken].Is_Admin_User == true {
		tempMonitor := MonitorInfo{}
		body, _ := io.ReadAll(r.Body)
		err := json.Unmarshal(body, &tempMonitor)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}

		db, err := sql.Open("postgres", connStr)
		if err != nil {
			sendNats("error", err.Error())
			log.Fatal(err)
		}
		defer db.Close()

		result, err := db.Exec("insert into Monitors (Display_ID, Monitor_Gsync_Premium, Monitor_Curved) values ($1, $2, $3)", tempMonitor.DisplayID, tempMonitor.GSyncPremium, tempMonitor.IsCurved)
		if err != nil {
			sendNats("error", err.Error())
			log.Fatal(err)
		}
		fmt.Println(result.RowsAffected())

		w.Write([]byte("Новый Монитор добавлен."))
		sendNats("log", fmt.Sprintf("Новый монитор добавлен. Администратор: %v", tokens[tempToken].Username_User))
		return
	}

	w.Write([]byte("Для добавления нужно обладать правами администратора!"))

}

func addDisplayHandler(w http.ResponseWriter, r *http.Request) { //обработчик HTTP-запросов получает данные из запроса, а затем вносит изменения в соответствующие мапы (displayInfoMap или monitorInfoMap)
	tempToken := r.Header.Get("Authorization")
	if checkLogin(tempToken) && tokens[tempToken].Is_Admin_User == true {
		tempDisplay := DisplayInfo{}
		body, _ := io.ReadAll(r.Body)
		err := json.Unmarshal(body, &tempDisplay)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}

		db, err := sql.Open("postgres", connStr)
		if err != nil {
			sendNats("error", err.Error())
			log.Fatal(err)
		}
		defer db.Close()

		result, err := db.Exec("insert into Displays (Display_Diagonal, Display_Resolution, Display_Type, Display_Gsync) values ($1, $2, $3, $4)", tempDisplay.Diagonal, tempDisplay.Resolution, tempDisplay.Type, tempDisplay.GSync)
		if err != nil {
			sendNats("error", err.Error())
			log.Fatal(err)
		}
		fmt.Println(result.RowsAffected())
		w.Write([]byte("Новый Дисплей добавлен."))
		sendNats("log", fmt.Sprintf("Новый дисплей добавлен. Администратор: %v", tokens[tempToken].Username_User))

		return
	}
	w.Write([]byte("Для добавления нужно обладать правами администратора!"))
}

func addUserHandler(w http.ResponseWriter, r *http.Request) { //обработчик HTTP-запросов получает данные из запроса, а затем вносит изменения в соответствующие мапы (displayInfoMap или monitorInfoMap)
	tempUser := User{}
	body, _ := io.ReadAll(r.Body)
	err := json.Unmarshal(body, &tempUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
		sendNats("error", err.Error())
	}
	defer db.Close()
	userpasswordhash := sha256.Sum256([]byte(tempUser.Password_User))
	result, err := db.Exec("insert into Users (Username_User, Password_User, Email_User, Is_Admin_User) values ($1, $2, $3, $4)", tempUser.Username_User, fmt.Sprintf("%x", userpasswordhash[:]), tempUser.Email_User, tempUser.Is_Admin_User)
	if err != nil {
		log.Fatal(err)
		sendNats("error", err.Error())
	}
	fmt.Println(result.RowsAffected())
	sendNats("log", fmt.Sprintf("Новый юзер зарегестрирован. User_Email: %v", tempUser.Email_User))

	w.Write([]byte("Пользователь зарегестрирован."))
}

func removeDisplayHandler(w http.ResponseWriter, r *http.Request) { //Функция возвращают информацию о всех дисплеях в формате текста, который будет отправлен клиенту в ответ на HTTP-запрос.
	tempToken := r.Header.Get("Authorization")
	if checkLogin(tempToken) && tokens[tempToken].Is_Admin_User == true {
		body, _ := io.ReadAll(r.Body)
		displayId := string(body)

		db, err := sql.Open("postgres", connStr)
		if err != nil {
			sendNats("error", err.Error())
			log.Fatal(err)
		}
		defer db.Close()

		result, err := db.Exec("delete from Displays where ID_Display = $1", displayId)
		if err != nil {
			sendNats("error", err.Error())
			log.Fatal(err)
		}
		fmt.Println(result.RowsAffected())

		w.Write([]byte("Вы удалили дисплей."))
		sendNats("log", fmt.Sprintf("Дисплей удалён. Администратор: %v", tokens[tempToken].Username_User))

		return
	}
	w.Write([]byte("Для удаления нужно обладать правами администратора!"))
}

func removeMonitorHandler(w http.ResponseWriter, r *http.Request) { //Функция возвращают информацию о всех мониторах в формате текста, который будет отправлен клиенту в ответ на HTTP-запрос.
	tempToken := r.Header.Get("Authorization")
	if checkLogin(tempToken) && tokens[tempToken].Is_Admin_User == true {

		body, _ := io.ReadAll(r.Body)
		monitorId := string(body)

		db, err := sql.Open("postgres", connStr)
		if err != nil {
			sendNats("error", err.Error())
			log.Fatal(err)
		}
		defer db.Close()

		result, err := db.Exec("delete from Monitors where ID_Monitor = $1", monitorId)
		if err != nil {
			sendNats("error", err.Error())
			log.Fatal(err)
		}
		fmt.Println(result.RowsAffected())

		w.Write([]byte("Вы удалили монитор."))
		sendNats("log", fmt.Sprintf("Монитор удалён. Администратор: %v", tokens[tempToken].Username_User))

		return
	}

	w.Write([]byte("Для удаления нужно обладать правами администратора!"))

}

func allDisplaysHandler(w http.ResponseWriter, r *http.Request) {
	tempToken := r.Header.Get("Authorization")
	if checkLogin(tempToken) {
		db, err := sql.Open("postgres", connStr)
		if err != nil {
			sendNats("error", err.Error())
			log.Fatal(err)
		}
		defer db.Close()

		rows, err := db.Query("select * from Displays")
		if err != nil {
			sendNats("error", err.Error())
			log.Fatal(err)
		}
		defer rows.Close()

		products := []DisplayInfo{}

		for rows.Next() {
			p := DisplayInfo{}
			err := rows.Scan(&p.ID_Display, &p.Diagonal, &p.Resolution, &p.Type, &p.GSync)
			if err != nil {
				fmt.Println(err)
				continue
			}
			products = append(products, p)
		}
		out, _ := json.Marshal(products)
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
		sendNats("log", fmt.Sprintf("all displays handler"))

		return
	}
	w.Write([]byte("Для просмотра нужно быть авторизованным!"))

}

func allMonitorsHandler(w http.ResponseWriter, r *http.Request) {
	tempToken := r.Header.Get("Authorization")
	if checkLogin(tempToken) {
		db, err := sql.Open("postgres", connStr)
		if err != nil {
			sendNats("error", err.Error())
			log.Fatal(err)
		}
		defer db.Close()

		rows, err := db.Query("select * from Monitors")
		if err != nil {
			sendNats("error", err.Error())
			log.Fatal(err)
		}
		defer rows.Close()

		products := []MonitorInfo{}

		for rows.Next() {
			p := MonitorInfo{}
			err := rows.Scan(&p.ID_Monitor, &p.Display.ID_Display, &p.GSyncPremium, &p.IsCurved)
			if err != nil {
				fmt.Println(err)
				continue
			}
			products = append(products, p)
		}

		out, _ := json.Marshal(products)

		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
		sendNats("log", fmt.Sprintf("all monitors handler"))

		return
	}
	w.Write([]byte("Для просмотра нужно быть авторизованным!"))

}

func getMonitorHandler(w http.ResponseWriter, r *http.Request) {
	tempToken := r.Header.Get("Authorization")
	if checkLogin(tempToken) {
		body, _ := io.ReadAll(r.Body)
		monitorID := string(body)

		db, err := sql.Open("postgres", connStr)
		if err != nil {
			sendNats("error", err.Error())
			log.Fatal(err)
		}
		defer db.Close()
		var tempMonitor MonitorInfo
		monitor := db.QueryRow("select * from Monitors where ID_Monitor = $1", monitorID)
		monitor.Scan(&tempMonitor.ID_Monitor, &tempMonitor.DisplayID, &tempMonitor.GSyncPremium, &tempMonitor.IsCurved)
		defer db.Close()

		out, err := json.Marshal(tempMonitor)
		if err != nil {
			sendNats("error", err.Error())
			log.Fatal(err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
		return
	}
	w.Write([]byte("Для просмотра нужно быть авторизованным!"))

}
