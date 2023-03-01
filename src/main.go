/*
*	author: 	Sergeev Maksim
*	description: Расчет рабочего времени
*	date: 		10.02.2023
*	version: 	1.0
 */

package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
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
	data := "## Время обеда в минутах, пример:\r\n## Obed = 45\r\nObed = \r\n" +
		"\r\n## Рабочее время в минутах в день, пример:\r\n## WorkTime = 390\r\nWorkTime = "
	file.WriteString(data)
	defer file.Close()
	fmt.Println("Заполните пожалуйста файл настроек settings.txt")
	duration := 5 * time.Second
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
		stroka++
		// Уберем последний символ перевода строки
		line = strings.TrimSuffix(line, "\n")
		line = strings.TrimSuffix(line, "\r")
		//fmt.Println("Считанная строка: ", line)
		numberOfColons := strings.Count(line, ":") // количество двоеточий в строке
		// Разделим строку на части
		parts := strings.Split(line, " ")

		// Запишем данные в переменные
		if numberOfColons == 2 {

			var time2Hour, time2Min int
			time1Hour, time1Min = strToTime(parts[1], stroka)
			fmt.Println("time1Hour,time1Min: ", time1Hour, time1Min)

			time2Hour, time2Min = strToTime(parts[2], stroka)
			fmt.Println("time2Hour,time2Min: ", time2Hour, time2Min)

			ost += calcDifference(time1Hour, time1Min, time2Hour, time2Min)
			fmt.Println("ost1 = ", ost)
		}
		if numberOfColons == 1 {
			date = parts[0] // дата незаконченного дня
			fmt.Println("Дата незаконченного дня: ", date)

			time1Hour, time1Min = strToTime(parts[1], stroka)
			fmt.Println("time1Hour,time1Min: ", time1Hour, time1Min)

			notFinishDay = true
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
	return date, time1Hour, time1Min, notFinishDay, ost, stroka
}

// Функция создания файла для учёта рабочего времени
func createDataFile() {
	fmt.Println("\nСоздаем новый файл с рабочим временем.")
	file, err := os.Create("data.txt")
	if err != nil {
		log.Fatal(err)
	}
	/* var data, time1 string
	fmt.Println("Введите дату и время прихода на работу в формате дд.мм.гггг чч:мм ")
	fmt.Scan(&data, &time1)
	file.WriteString(data + " " + time1) */
	date := inputDate(file)
	inputPrihod(file, date)
	fmt.Println("\nХотите ввести время ухода с работы? да-1/нет-0")
	var flag bool
	fmt.Scan(&flag)
	if flag {
		inputUhod(file, date)
		/*var time2 string
		fmt.Println("Введите время ухода с работы в формате чч:мм ")
		fmt.Scan(&time2)
		file.WriteString(" " + time2) */
	}
	defer file.Close()
	os.Exit(0)
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

func inputDate(file *os.File) string {
	var day, month, year string
	fmt.Println("Введите дату: день месяц год.")
	fmt.Println("Например, 5 12 2020")
	fmt.Print("-> ")
	fmt.Scan(&day)
	n, err := fmt.Scanf("%s %s", &month, &year)
	if err != nil || n != 2 {
		// handle invalid input
		fmt.Println(n, err)
		//return
	}
	file.WriteString("\r\n")
	file.WriteString(day + "." + month + "." + year)
	return fmt.Sprint(day + "." + month + "." + year)
}

func inputPrihod(file *os.File, date string) string {
	var time1 string
	fmt.Println("Введите время прихода на работу ", date, " в формате чч:мм ")
	fmt.Print("-> ")
	fmt.Scan(&time1)
	time1 = strings.Replace(time1, "^", ":", -1)
	file.WriteString(" " + time1)
	return time1
}

func inputUhod(file *os.File, date string) string {
	var time2 string
	fmt.Println("Введите время ухода с работы ", date, " в формате чч:мм ")
	fmt.Print("-> ")
	fmt.Scan(&time2)
	time2 = strings.Replace(time2, "^", ":", -1)
	file.WriteString(" " + time2)
	return time2
}

func calcDifference(time1Hour int, time1Min int, time2Hour int, time2Min int) int {
	fmt.Println("++++++++++++calcDifference============")
	fmt.Println("time2Hour = ", time2Hour)
	fmt.Println("time2Min = ", time2Min)
	fmt.Println("time1Hour = ", time1Hour)
	fmt.Println("time1Min = ", time1Min)
	fmt.Println("obed = ", obed)
	fmt.Println("stand = ", stand)
	fmt.Println("==========calcDifference--------------")
	return (time2Hour*60 + time2Min) - (time1Hour*60 + time1Min) - obed - stand
}

func main() {
	path, _ := os.Executable()
	fmt.Println("Путь к программе: ", path)

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
	var date, time1, time2 string

	stroka := 0

	// Cчитываем данные из файла data.txt
	date, time1Hour, time1Min, notFinishDay, totalDifference, stroka = readDataFile()

	f, err := os.OpenFile("data.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	if notFinishDay == true {
		/*	fmt.Print("Введите время ухода с работы ", date, " : ")
			var time2 string
			fmt.Scan(&time2)*/
		time2 = inputUhod(f, date)
		time2Hour, time2Min = strToTime(time2, stroka)
		fmt.Println("time2Hour,time2Min: ", time2Hour, time2Min)
		difference = calcDifference(time1Hour, time1Min, time2Hour, time2Min)
		totalDifference += calcDifference(time1Hour, time1Min, time2Hour, time2Min) //(time2Hour*60 + time2Min) - (time1Hour*60 + time1Min) - obed - stand

		f.WriteString(" " + fmt.Sprint(difference))
		fmt.Println("totalDifference33 = ", totalDifference)
		fmt.Println("\nХотите ввести время прихода на работу? да-1/нет-0")
		var flag1 bool
		fmt.Scan(&flag1)
		if flag1 {
			//var date, time1 string
			stroka++
			date = inputDate(f)
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
				//f.WriteString("\r\n" + date + " " + time1 + " " + time2 + " " + fmt.Sprint(difference))
				f.WriteString(" " + fmt.Sprint(difference))
				fmt.Println("difference4 = ", difference)
				fmt.Println("totalDifference4 = ", totalDifference)
			}
		}

		if err := f.Close(); err != nil {
			log.Fatal(err)
		}
		fmt.Println("difference2 = ", difference)
		fmt.Println("totalDifference2 = ", totalDifference)
	} else {
		fmt.Println("totalDifference3 = ", totalDifference)
		//var date, time1 string
		stroka++
		date = inputDate(f)
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
			//f.WriteString("\r\n" + date + " " + time1 + " " + time2 + " " + fmt.Sprint(difference))
			f.WriteString(" " + fmt.Sprint(difference))
			fmt.Println("difference4 = ", difference)
			fmt.Println("totalDifference4 = ", totalDifference)
		}

		//	fmt.Scan(&data, &time)

		defer f.Close()
	}
	fmt.Println("totalDifference5 = ", totalDifference)
	duration := 5 * time.Second
	time.Sleep(duration)
}
