package main

import (
	"database/sql"
	"encoding/json"
	"fmt" //Пакет для форматированного вывода на консоль.
	"io"
	"log"
	"net/http" //Пакет для работы с HTTP-сервером.

	_ "github.com/lib/pq"
)

type DisplayInfo struct { //Структура
	ID_Display string `json:"id"`
	Diagonal   int    `json:"diagonal"`
	Resolution string `json:"resolution"`
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

var connStr = "user=postgres password=Waldronnezz1 dbname=shop sslmode=disable"

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
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}

}

func addMonitorHandler(w http.ResponseWriter, r *http.Request) { //обработчик HTTP-запросов получает данные из запроса, а затем вносит изменения в соответствующие мапы (displayInfoMap или monitorInfoMap)
	tempMonitor := MonitorInfo{}
	body, _ := io.ReadAll(r.Body)
	err := json.Unmarshal(body, &tempMonitor)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	result, err := db.Exec("insert into Monitors (Display_ID, Monitor_Gsync_Premium, Monitor_Curved) values ($1, $2, $3)", tempMonitor.DisplayID, tempMonitor.GSyncPremium, tempMonitor.IsCurved)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result.RowsAffected())

	w.Write([]byte("Новый Монитор добавлен."))
}

func addDisplayHandler(w http.ResponseWriter, r *http.Request) { //обработчик HTTP-запросов получает данные из запроса, а затем вносит изменения в соответствующие мапы (displayInfoMap или monitorInfoMap)
	tempDisplay := DisplayInfo{}
	body, _ := io.ReadAll(r.Body)
	err := json.Unmarshal(body, &tempDisplay)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	// displayInfoMap[tempDisplay.ID] = tempDisplay

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	result, err := db.Exec("insert into Displays (Display_Diagonal, Display_Resolution, Display_Type, Display_Gsync) values ($1, $2, $3, $4)", tempDisplay.Diagonal, tempDisplay.Resolution, tempDisplay.Type, tempDisplay.GSync)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result.RowsAffected())
	w.Write([]byte("Новый Дисплей добавлен."))
}

func removeDisplayHandler(w http.ResponseWriter, r *http.Request) { //Функция возвращают информацию о всех дисплеях в формате текста, который будет отправлен клиенту в ответ на HTTP-запрос.
	body, _ := io.ReadAll(r.Body)
	displayId := string(body)
	// delete(displayInfoMap, displayId)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	result, err := db.Exec("delete from Displays where ID_Display = $1", displayId)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result.RowsAffected())

	w.Write([]byte("Вы удалили дисплей."))
}

func removeMonitorHandler(w http.ResponseWriter, r *http.Request) { //Функция возвращают информацию о всех мониторах в формате текста, который будет отправлен клиенту в ответ на HTTP-запрос.
	body, _ := io.ReadAll(r.Body)
	monitorId := string(body)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	result, err := db.Exec("delete from Monitors where ID_Monitor = $1", monitorId)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result.RowsAffected())

	w.Write([]byte("Вы удалили монитор."))
}

func allDisplaysHandler(w http.ResponseWriter, r *http.Request) {

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("select * from Displays")
	if err != nil {
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

	for _, p := range products {
		fmt.Println(p.ID_Display, p.Diagonal, p.Resolution, p.Type, p.GSync)
	}
	out, _ := json.Marshal(products)
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

func allMonitorsHandler(w http.ResponseWriter, r *http.Request) {
	//for _, monitor := range monitorInfoMap {

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("select * from Monitors")
	if err != nil {
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
}

func getMonitorHandler(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	monitorID := string(body)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	var tempMonitor MonitorInfo
	monitor := db.QueryRow("select * from Monitors where ID_Monitor = $1", monitorID)
	monitor.Scan(&tempMonitor.ID_Monitor, &tempMonitor.DisplayID, &tempMonitor.GSyncPremium, &tempMonitor.IsCurved)
	defer db.Close()

	out, err := json.Marshal(tempMonitor)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}
