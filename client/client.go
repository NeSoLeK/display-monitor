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

// var monitors = make(map[string]Monitor)

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
	fmt.Println("Монитор успешно добавлен")
	return monitor
}

// func ViewAllMonitors() {
// 	fmt.Println("Список мониторов: ")
// 	for name, monitor := range monitors {
// 		fmt.Printf("Монитор: %s \n", name)
// 		fmt.Printf("Напряжение питания: %v \n", monitor.PowerColtage)
// 		fmt.Printf("Дисплей: \n")
// 		fmt.Println("---------------------------")
// 		// fmt.Printf("Название дисплея: ", monitor.Display.DName)
// 		fmt.Printf("Диагональ: %.2f дюймов\n", monitor.Display.Diagonal)
// 		fmt.Printf("Разрешение экрана: %s \n", monitor.Display.Resolution)
// 		fmt.Printf("Тип: %s \n", monitor.Display.Type)
// 		fmt.Printf("Поддержка G-Sync: %v \n", monitor.Display.GSync)
// 		fmt.Printf("Поддержка G-Sync Premium: %v \n", monitor.GSyncPremium)
// 		fmt.Printf("Изогнутый: %v \n", monitor.IsCurved)
// 	}
// }

func main() {

	client := &http.Client{}

	for {
		fmt.Println("Выберите действие: ")
		fmt.Println("1) Добавить дисплей")
		fmt.Println("2) Добавить монитор")
		fmt.Println("3) Посмотреть все мониторы")
		fmt.Println("4) Посмотреть все дисплеи")
		fmt.Println("5) Удалить монитор")
		fmt.Println("6) Удалить дисплей")
		fmt.Println("7) Посмотреть монитор по ID")
		fmt.Println("8) Выйти")

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
			// req.Header.Add("Key", newdispl.DName)
			response, err := client.Do(req)
			if err != nil {
				log.Fatal("[CLIENT] ERROR", err)
			}
			defer response.Body.Close()

		case 2:
			newmonitor := AddMonitor()
			newj, err := json.Marshal(newmonitor)
			if err != nil {
				log.Fatal(err)
			}
			req, _ := http.NewRequest("POST", "http://127.0.0.1:8080/addMonitor", bytes.NewReader(newj))
			// req.Header.Add("Key", newmonitor.MName)
			response, err := client.Do(req)
			if err != nil {
				log.Fatal("[CLIENT] ERROR", err)
			}
			defer response.Body.Close()

		case 3:
			// ViewAllMonitors()

			req, _ := http.NewRequest("GET", "http://127.0.0.1:8080/allMonitors", nil)
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
			response, err := client.Do(req)
			if err != nil {
				log.Fatal("[CLIENT] ERROR", err)
			}
			defer response.Body.Close()
		case 6:
			var delword string
			fmt.Println("Введите ID дисплея для удаления: ")
			fmt.Scanln(&delword)
			req, _ := http.NewRequest("POST", "http://127.0.0.1:8080/removeDisplay", bytes.NewReader([]byte(delword)))
			response, err := client.Do(req)
			if err != nil {
				log.Fatal("[CLIENT] ERROR", err)
			}
			defer response.Body.Close()
		case 7:
			var searchWord string
			fmt.Println("Введите ID дисплея для поиска: ")
			fmt.Scanln(&searchWord)
			req, _ := http.NewRequest("POST", "http://127.0.0.1:8080/getMonitor", bytes.NewReader([]byte(searchWord)))
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
			fmt.Println("Программа завершена")
			return
		default:
			fmt.Println("Выбирете другой параметр")
		}

	}

}
