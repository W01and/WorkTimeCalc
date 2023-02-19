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

// Функция создания файла для учёта рабочего времени
func createDataFile() {
	fmt.Println("\nСоздаем новый файл с рабочим временем.")
	file, err := os.Create("data.txt")
	if err != nil {
		log.Fatal(err)
	}
	var data, time1 string
	fmt.Println("Введите дату и время прихода на работу в формате дд.мм.гггг чч:мм ")
	fmt.Scan(&data, &time1)
	file.WriteString(data + " " + time1)
	fmt.Println("\nХотите ввести время ухода с работы? да-1/нет-0")
	var flag bool
	fmt.Scan(&flag)
	if flag {
		var time2 string
		fmt.Println("Введите время ухода с работы в формате чч:мм ")
		fmt.Scan(&time2)
		file.WriteString(" " + time2)
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

func main() {
	path, _ := os.Executable()
	fmt.Println("Путь к программе: ", path)
	obed := 0
	stand := 0

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
	var ost int = 0 // остаток
	var notFinishDay bool = false
	var date string

	// Откройте файл для чтения рабочего времени
	file2, err := os.Open("data.txt")
	// Если файла нет - создадим файл
	if err != nil {
		createDataFile()
	}

	// Создаём читателя, связанного с файлом
	reader := bufio.NewReader(file2)
	stroka := 0
	for {
		// Читаем строку из файла
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}
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
			time1Hour, time1Min = strToTime(parts[1], stroka)
			fmt.Println("time1Hour,time1Min: ", time1Hour, time1Min)

			time2Hour, time2Min = strToTime(parts[2], stroka)
			fmt.Println("time2Hour,time2Min: ", time2Hour, time2Min)

			ost += (time2Hour*60 + time2Min) - (time1Hour*60 + time1Min) - obed - stand
			fmt.Println("ost = ", ost)
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
			return
		}
	}
	defer file2.Close()

	if notFinishDay == true {
		fmt.Print("Введите время ухода с работы ", date, " : ")
		var time2 string
		fmt.Scan(&time2)

		time2Hour, time2Min = strToTime(time2, stroka)
		fmt.Println("time2Hour,time2Min: ", time2Hour, time2Min)

		ost += (time2Hour*60 + time2Min) - (time1Hour*60 + time1Min) - obed - stand
		f, err := os.OpenFile("data.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}
		if _, err := f.Write([]byte(" " + time2 + "\n")); err != nil {
			f.Close() // ignore error; Write error takes precedence
			log.Fatal(err)
		}
		if err := f.Close(); err != nil {
			log.Fatal(err)
		}
		fmt.Println("ost = ", ost)
	} else {
		fmt.Println("ost = ", ost)
		var data, time string
		fmt.Println("Введите дату и время прихода на работу в формате дд.мм.гггг чч:мм ")
		fmt.Scan(&data, &time)
		f, err := os.OpenFile("data.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}
		f.WriteString("\r\n" + data + " " + time)
		defer f.Close()
	}

	duration := 5 * time.Second
	time.Sleep(duration)
}
