package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Year             string `json:"year"`
	Month            string `json:"month"`
	Price            int    `json:"price"`
	Duration         int    `json:"duration"`
	DurationAll      int    `json:"durationAll"`
	Excess           int    `json:"excess"`
	ExcessSum        int    `json:"excessSum"`
	LessonsCount     int    `json:"lessonsCount"`
	LessonsCompleted int    `json:"lessonsCompleted"`
	Sum              int    `json:"sum"`
	ExpectedSum      int    `json:"expectedSum"`
}

/* false - json to struct, true - struct to json */
func jsonConverter(byt []byte, data interface{}, flag bool) (Config, string) {
	var structure Config
	var jsonString string

	if flag == false {

		var config = Config{}

		err := json.Unmarshal(byt, &config)
		if err != nil {
			panic(err)
		}

		structure = config

	} else {
		config, err := json.Marshal(data)

		jsonString = string(config)

		if err != nil {
			panic(err)
		}
	}

	return structure, jsonString
}

func addChanges(config Config, counter int, duration int, lessonsPerDay int) Config {
	var newExcess int
	var minutePrice = config.Price / config.Duration

	if ((config.Duration * lessonsPerDay) - duration) < 0 {
		newExcess = int(float64(config.Excess) + math.Abs(float64((config.Duration*lessonsPerDay)-duration)))
	} else {
		newExcess = config.Excess
	}

	config.DurationAll = config.DurationAll + duration
	config.Excess = newExcess
	config.ExcessSum = int(newExcess * minutePrice)
	config.LessonsCompleted = counter
	config.Sum = config.Sum + int(duration*minutePrice)

	return config
}

func addCall(config Config) {
	var duration string
	var lessonsPerDay string
	var plugConverter []byte // Crutch
	var counter int = config.LessonsCompleted

	fmt.Print("Введите длительность разговора: ")
	fmt.Fscan(os.Stdin, &duration)
	duration = strings.TrimSpace(duration)
	intDuration, _ := strconv.Atoi(duration)

	fmt.Print("Введите количество уроков за разговор: ")
	fmt.Fscan(os.Stdin, &lessonsPerDay)
	lessonsPerDay = strings.TrimSpace(lessonsPerDay)
	lessonsPerDayInt, _ := strconv.Atoi(lessonsPerDay)

	counter += lessonsPerDayInt

	config = addChanges(config, counter, intDuration, lessonsPerDayInt)
	_, jsonString := jsonConverter(plugConverter, config, true)

	mydata := []byte(jsonString)

	err := ioutil.WriteFile("config.json", mydata, 0777)
	if err != nil {
		panic(err)
	}

	showMenu(config)

}

func getConfig() (year int, month time.Month, price string, duration string, durationAll int, excess int, excessSum int, lessonsCount string, sum int, expectedSum int) {
	currentTime := time.Now()
	month = currentTime.Month()
	year = currentTime.Year()

	fmt.Print("Введите цену одного занятия(без наименования валюты): ")
	fmt.Fscan(os.Stdin, &price)
	price = strings.TrimSpace(price)

	fmt.Print("Введите длительность занятия(в минутах): ")
	fmt.Fscan(os.Stdin, &duration)
	duration = strings.TrimSpace(duration)

	fmt.Print("Введите количество занятий: ")
	fmt.Fscan(os.Stdin, &lessonsCount)
	lessonsCount = strings.TrimSpace(lessonsCount)

	priceInt, _ := strconv.Atoi(price)
	lessonsCountInt, _ := strconv.Atoi(lessonsCount)
	expectedSum = priceInt * lessonsCountInt

	return year, month, price, duration, durationAll, excess, excessSum, lessonsCount, sum, expectedSum
}

func showMenu(config Config) {
	var menu string = "Выберите один из пунктов меню \n 1. Вывести статистику за текущий месяц. \n 2. Добавить данные о текущем звонке. \n 3. Выход"
	var numberInput string

	fmt.Println(menu)
	fmt.Print(": ")
	fmt.Fscan(os.Stdin, &numberInput)
	numberInput = strings.TrimSpace(numberInput)
	menuNumber, _ := strconv.Atoi(numberInput)

	switch menuNumber {
	case 1:
		fmt.Printf(" -----------------------------------------------\n %v %v\n -----------------------------------------------\n Завершено звонков: %v из %v\n Стоимость одного звонка: %v рублей\n Стоимость звонков за месяц: %v рублей\n Ожидаемая стоимость звонков: %v рублей\n Заявленная длительность звонка: %v минут\n Длительность всех разговоров: %v минут\n Избыток длительности: %v минут\n Переплата за избыток: %v рублей\n", config.Year, config.Month, config.LessonsCompleted, config.LessonsCount, config.Price, config.Sum, config.ExpectedSum, config.Duration, config.DurationAll, config.Excess, config.ExcessSum)
	case 2:
		addCall(config)
	case 3:
		return
	}
}

func logConfig() {
	logsByt, _ := ioutil.ReadFile("config.json")
	var logsString string = string(logsByt[:])
	logsString += "\n"

	logsFile, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		logsFile, err := os.Create("logs.txt")
		if err != nil {
			panic(err)
		}

		defer logsFile.Close()

		if _, err = logsFile.WriteString(logsString); err != nil {
			panic(err)
		}

		return
	}

	defer logsFile.Close()

	if _, err = logsFile.WriteString(logsString); err != nil {
		panic(err)
	}
}

func firstLoad() {

	year, month, price, duration, durationAll, excess, excessSum, lessonsCount, sum, expectedSum := getConfig()
	var lessonsCompleted int = 0

	mydata := []byte(fmt.Sprintf("{\"year\": \"%v\", \"month\": \"%v\", \"price\": %v, \"duration\": %v, \"durationAll\": %v, \"excess\": %v, \"excessSum\": %v, \"lessonsCount\": %v, \"lessonsCompleted\": %v, \"Sum\": %v, \"expectedSum\": %v}\n", year, month, price, duration, durationAll, excess, excessSum, lessonsCount, lessonsCompleted, sum, expectedSum))

	err := ioutil.WriteFile("config.json", mydata, 0777)
	if err != nil {
		panic(err)
	}
}

func main() {
	byt, err := ioutil.ReadFile("config.json")
	if err != nil || byt == nil {
		firstLoad()
	} else {
		config, _ := jsonConverter(byt, Config{}, false)
		// _, jsonString := jsonConverter(byt, structure, true)

		_, month, _ := time.Now().Date()

		if config.Month != month.String() {

			logConfig()

			fmt.Println("Начало нового месяца, повторите ввод данных")
			firstLoad()
		}

		showMenu(config)

	}

}
