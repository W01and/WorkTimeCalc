/*
*	author: 	Sergeev Maksim
*	description: Расчет рабочего времени
*	date: 		10.02.2023
*	version: 	2.0
 */

package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

var obed int = 0
var stand int = 0

// Функция для создания файла настроек
func createSettingsFile() {
	file, err := os.Create("settings.txt")
	if err != nil {
		log.Fatal(err)
	}
	data := "## Время обеда в минутах, пример:\r\n## Obed = 45\r\nObed = 45\r\n" +
		"\r\n## Рабочее время в минутах в день, пример:\r\n## WorkTime = 390\r\nWorkTime = 390"
	file.WriteString(data)
	defer file.Close()
	fmt.Println("Был создан файл настроек по умолчанию settings.txt")
	fmt.Println("Обед составляет 45 минут в день")
	fmt.Println("Рабочее время составляет 390 минут в день")
	fmt.Println("Если что-то не так, отредактируйте файл настроек")
	duration := 10 * time.Second
	time.Sleep(duration)
	os.Exit(0)
}

// Функция для считывания файла настроек
func readSettingsFile(file *os.File) (int, int) {
	obed := 0
	stand := 0
	// Создаём читателя, связанного с файлом
	reader := bufio.NewReader(file)
	for {
		// Читаем строку из файла
		line, err := reader.ReadString('\n')
		line = strings.TrimSuffix(line, "\r\n")

		if strings.Contains(line, "##") {
			continue
		}
		if strings.Contains(line, "Obed = ") {
			line = strings.TrimPrefix(line, "Obed = ")
			obed, err = strconv.Atoi(line)
			if err != nil {
				fmt.Println("ошибка преобразования из string в int\n", err)
			}
		}
		if strings.Contains(line, "WorkTime = ") {
			line = strings.TrimPrefix(line, "WorkTime = ")
			stand, err = strconv.Atoi(line)
			if err != nil {
				fmt.Println("ошибка преобразования из string в int\n", err)
			}
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("Ошибка чтения из файла", err)
			duration := 5 * time.Second
			time.Sleep(duration)
			os.Exit(0)
		}
	}
	return obed, stand
}

// Функция создания файла для учёта рабочего времени
func createDataFile() {
	fmt.Println("\nСоздаем новый файл с рабочим временем.")
	file, err := os.Create("data.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	date := inputDate(file, "0.0.2023")
	time1 := inputPrihod(file, date)
	fmt.Println("\nХотите ввести время ухода с работы? да-1/нет-0")
	var flag bool
	fmt.Scan(&flag)
	if flag {
		time1Hour, time1Min := strToTime(time1, 1)
		time2 := inputUhod(file, date)
		time2Hour, time2Min := strToTime(time2, 1)
		totalDifference := calcDifference(time1Hour, time1Min, time2Hour, time2Min)
		file.WriteString(" " + fmt.Sprint(totalDifference))
		printDifference(totalDifference)
	}
	duration := 5 * time.Second
	time.Sleep(duration)

	os.Exit(0)
}

// Функция для считывания времени прихода/ухода из файла
func readDataFile() (string, int, int, bool, int, int) {
	// Откройте файл для чтения рабочего времени
	file2, err := os.Open("data.txt")
	// Если файла нет - создадим файл
	if err != nil {
		createDataFile()
	}
	// Создаём читателя, связанного с файлом
	reader := bufio.NewReader(file2)
	stroka := 0
	var date string
	var time1Hour, time1Min int
	var ost int = 0 // остаток
	var notFinishDay bool = false

	for {
		// Читаем строку из файла
		line, err := reader.ReadString('\n')
		fmt.Println("len(line) = ", len(line))
		if len(line) > 10 {
			stroka++
			// Уберем последний символ перевода строки
			line = strings.TrimSuffix(line, "\n")
			line = strings.TrimSuffix(line, "\r")
			numberOfColons := strings.Count(line, ":") // количество двоеточий в строке
			// Разделим строку на части
			parts := strings.Split(line, " ")
			date = parts[0] // дата
			// Запишем данные в переменные
			if numberOfColons == 2 {

				var time2Hour, time2Min int
				time1Hour, time1Min = strToTime(parts[1], stroka)
				//fmt.Println("time1Hour,time1Min: ", time1Hour, time1Min)

				time2Hour, time2Min = strToTime(parts[2], stroka)
				//fmt.Println("time2Hour,time2Min: ", time2Hour, time2Min)

				ost += calcDifference(time1Hour, time1Min, time2Hour, time2Min)
				//fmt.Println("ost1 = ", ost)
			}

			if numberOfColons == 1 {
				fmt.Println("Дата незаконченного дня: ", date)

				time1Hour, time1Min = strToTime(parts[1], stroka)
				fmt.Println("time1Hour,time1Min: ", time1Hour, time1Min)

				notFinishDay = true
			}
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("Ошибка чтения из файла", err)
			os.Exit(0)
		}

	}
	defer file2.Close()
	fmt.Println("Последняя дата: ", date)
	return date, time1Hour, time1Min, notFinishDay, ost, stroka
}

// Функция для вывода переработки/недоработки на экран
func printDifference(totalDifference int) {
	fmt.Println("\r\n***************************")
	if totalDifference > 0 {
		fmt.Println("Вы ПЕРЕработали ", totalDifference, " мин.")
	} else if totalDifference < 0 {
		fmt.Println("Вы НЕДОработали ", math.Abs(float64(totalDifference)), " мин.")
	} else {
		fmt.Println("У вас нет долгов по времени")
	}
	fmt.Println("***************************")
}

// Функция перевода строки времени в часы и минуты
func strToTime(time string, stroka int) (int, int) {
	time1x := strings.Split(time, ":")        //разделяем часы и минуты
	time1Hour, err := strconv.Atoi(time1x[0]) //преобразуем в int
	if err != nil {
		fmt.Println("ошибка, stroka:", stroka, err)
	}
	time1Min, err := strconv.Atoi(time1x[1])
	if err != nil {
		fmt.Println("ошибка, stroka:", stroka, err)
	}
	return time1Hour, time1Min
}

func dateStringToInt(date string) (int, int, int) {
	parts := strings.Split(date, ".")
	intDay, err := strconv.Atoi(parts[0])
	if err != nil {
		panic(err)
	}
	intMonth, err := strconv.Atoi(parts[1])
	if err != nil {
		panic(err)
	}
	intYear, err := strconv.Atoi(parts[2])
	if err != nil {
		panic(err)
	}
	fmt.Println(intDay, intMonth, intYear)
	return intDay, intMonth, intYear
}

func inputDate(file *os.File, lastDate string) string {
	lastDay, _, _ := dateStringToInt(lastDate)
	var day, month, year string
	for {
		fmt.Println("Введите дату: день месяц год.")
		fmt.Println("Например, 25 12 2020")
		fmt.Print("-> ")
		fmt.Scan(&day)
		n, err := fmt.Scanf("%s %s", &month, &year)
		if err != nil || n != 2 {
			// handle invalid input
			fmt.Println(n, err)
			//return
		}
		intDay, err := strconv.Atoi(day)
		if err != nil {
			panic(err)
		}
		intMonth, err := strconv.Atoi(month)
		if err != nil {
			panic(err)
		}
		intYear, err := strconv.Atoi(year)
		if err != nil {
			panic(err)
		}
		if intYear < 100 {
			intYear += 2000
			year = "20" + year
		}

		if intDay <= 31 && intMonth <= 12 && intYear >= 2023 && intDay > lastDay {
			if intDay < 10 {
				day = "0" + fmt.Sprint(intDay)
			}
			if intMonth < 10 {
				month = "0" + fmt.Sprint(intMonth)
			}
			break
		}
		fmt.Println("\r\nВы сделали ошибку при вводе!")
	}
	file.WriteString("\r\n")
	file.WriteString(day + "." + month + "." + year)
	return fmt.Sprint(day + "." + month + "." + year)
}

func inputPrihod(file *os.File, date string) string {
	nameDayOfWeek := dateStringToNameOfDay(date)
	var time1 string
	for {
		fmt.Println("Введите время прихода на работу", nameDayOfWeek, date, "в формате чч:мм ")
		fmt.Print("-> ")
		fmt.Scan(&time1)
		time1 = strings.Replace(time1, "^", ":", -1)
		time1Hour, time1Min := strToTime(time1, 0)
		if time1Hour < 24 && time1Min < 60 {
			break
		}
	}
	file.WriteString(" " + time1)
	return time1
}

func inputUhod(file *os.File, date string) string {
	nameDayOfWeek := dateStringToNameOfDay(date)
	var time2 string
	for {
		fmt.Println("Введите время ухода с работы", nameDayOfWeek, date, "в формате чч:мм ")
		fmt.Print("-> ")
		fmt.Scan(&time2)
		time2 = strings.Replace(time2, "^", ":", -1)
		time2Hour, time2Min := strToTime(time2, 0)
		if time2Hour < 24 && time2Min < 60 {
			break
		}
	}
	file.WriteString(" " + time2)
	return time2
}

func calcDifference(time1Hour int, time1Min int, time2Hour int, time2Min int) int {
	return (time2Hour*60 + time2Min) - (time1Hour*60 + time1Min) - obed - stand
}

func dateToDay(d int, month int, year int) int {
	m := month - 2
	if m < 1 {
		m += 12
		year--
	}
	c := year / 100
	y := year - c*100

	dayOfTheWeek := (d + (13*m-1)/5 + y + y/4 + c/4 - 2*c + 777) % 7
	return dayOfTheWeek
}

func dateStringToNameOfDay(date string) string {
	intDay, intMonth, intYear := dateStringToInt(date)
	dayOfTheWeek := dateToDay(intDay, intMonth, intYear)
	var output string
	switch dayOfTheWeek {
	case 1:
		output = "в Понедельник"
	case 2:
		output = "во Вторник"
	case 3:
		output = "в Среду"
	case 4:
		output = "в Четверг"
	case 5:
		output = "в Пятницу"
	case 6:
		output = "в Субботу"
	case 0:
		output = "в Воскресенье"
	}
	return output
}

func main() {
	//	path, _ := os.Executable()
	//	fmt.Println("Путь к программе: ", path)

	// Открываем и считываем файл настроек
	file, err := os.Open("settings.txt") // For read access.
	if err != nil {
		createSettingsFile()
	} else {
		obed, stand = readSettingsFile(file)
	}
	defer file.Close()

	// Структура для хранения данных
	var time1Hour, time1Min int
	var time2Hour, time2Min int
	var difference int = 0 // остаток общий
	var totalDifference int = 0
	var notFinishDay bool = false
	var lastDate, date, time1, time2 string

	stroka := 0

	// Cчитываем данные из файла data.txt
	lastDate, time1Hour, time1Min, notFinishDay, totalDifference, stroka = readDataFile()
	fmt.Println("lastDate = ", lastDate)

	f, err := os.OpenFile("data.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	if notFinishDay == true {
		time2 = inputUhod(f, lastDate)
		time2Hour, time2Min = strToTime(time2, stroka)
		difference = calcDifference(time1Hour, time1Min, time2Hour, time2Min)
		totalDifference += calcDifference(time1Hour, time1Min, time2Hour, time2Min)

		f.WriteString(" " + fmt.Sprint(difference))
		printDifference(totalDifference)
		fmt.Println("\nХотите ввести время прихода на работу? да-1/нет-0")
		var flag1 bool
		fmt.Scan(&flag1)
		if flag1 {
			//var date, time1 string
			f.WriteString("\r\n")
			stroka++
			date = inputDate(f, lastDate)
			time1 = inputPrihod(f, date)
			fmt.Println("\nХотите ввести время ухода с работы? да-1/нет-0")
			var flag2 bool
			fmt.Scan(&flag2)
			if flag2 {
				time1Hour, time1Min = strToTime(time1, stroka)
				time2 = inputUhod(f, date)
				time2Hour, time2Min = strToTime(time2, stroka)
				difference = calcDifference(time1Hour, time1Min, time2Hour, time2Min)
				totalDifference += calcDifference(time1Hour, time1Min, time2Hour, time2Min)
				f.WriteString(" " + fmt.Sprint(difference))
			}
		}

		if err := f.Close(); err != nil {
			log.Fatal(err)
		}
	} else {

		printDifference(totalDifference)
		stroka++
		date = inputDate(f, lastDate)
		time1 = inputPrihod(f, date)
		fmt.Println("\nХотите ввести время ухода с работы? да-1/нет-0")
		var flag bool
		fmt.Scan(&flag)
		if flag {
			time1Hour, time1Min = strToTime(time1, stroka)
			time2 = inputUhod(f, date)
			time2Hour, time2Min = strToTime(time2, stroka)
			difference = calcDifference(time1Hour, time1Min, time2Hour, time2Min)
			totalDifference += calcDifference(time1Hour, time1Min, time2Hour, time2Min)
			f.WriteString(" " + fmt.Sprint(difference))
		}
		defer f.Close()
	}
	printDifference(totalDifference)
	duration := 5 * time.Second
	time.Sleep(duration)
}
