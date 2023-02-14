/*
*	author: Sergeev Maksim
*
*	description: Расчет рабочего времени
*
*	date: 10.02.2023
 */

package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func main() {
	path, _ := os.Executable()
	fmt.Println("Путь к программе: ", path)
	obed := 0
	stand := 0

	file, err := os.Open("settings.txt") // For read access.
	if err != nil {
		file, err = os.Create("settings.txt")
		data := "## Время обеда в минутах, пример:\r\n## Obed = 45\r\nObed = \r\n" +
			"\r\n## Рабочее время в минутах в день, пример:\r\n## WorkTime = 390\r\nWorkTime = "
		file.WriteString(data)
		defer file.Close()
		os.Exit(0)
	} else {
		// Создаём читателя, связанного с файлом
		reader := bufio.NewReader(file)
		for {
			// Читаем строку из файла
			line, err := reader.ReadString('\n')
			line = strings.TrimSuffix(line, "\r\n")
			fmt.Println("cтрока: \"", line, "\"")

			//	comm := "##"
			if strings.Contains(line, "##") {
				continue
			}
			if strings.Contains(line, "Obed = ") {
				line = strings.TrimPrefix(line, "Obed = ")
				fmt.Println("Обрезанная cтрока: \"", line, "\"")
				obed, err = strconv.Atoi(line)
				if err != nil {
					fmt.Println("ошибка преобразования из string в int\n", err)
				}
			}
			if strings.Contains(line, "WorkTime = ") {
				line = strings.TrimPrefix(line, "WorkTime = ")
				fmt.Println("Обрезанная cтрока: \"", line, "\"")
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
				return
			}
		}
	}
	fmt.Println("Время обеда: ", obed)
	fmt.Println("Время работы: ", stand)
	// Откройте файл для чтения
	file2, err := os.Open("data.txt")
	if err != nil {
		fmt.Println("Не могу прочитать файл,", err)
		//		return
	}

	// Структура для хранения данных
	var time1Hour, time1Min int
	var time2Hour, time2Min int
	//	var data, time1, time2 []string
	var ost int = 0 // остаток

	// Создаём читателя, связанного с файлом
	reader := bufio.NewReader(file2)
	for {
		// Читаем строку из файла
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("Ошибка чтения из файла", err)
			return
		}
		// Уберем последний символ перевода строки
		line = strings.TrimSuffix(line, "\n")
		line = strings.TrimSuffix(line, "\r")
		fmt.Println("Считанная строка: ", line)
		// Разделим строку на части
		parts := strings.Split(line, " ")
		// Запишем данные в переменные
		if len(parts) == 3 {
			time1x := strings.Split(parts[1], ":")   //разделяем часы и минуты
			time1Hour, err = strconv.Atoi(time1x[0]) //преобразуем в int
			if err != nil {
				fmt.Println("ошибка", err)
			}
			time1Min, err = strconv.Atoi(time1x[1])
			if err != nil {
				fmt.Println("ошибка", err)
			}
			fmt.Println("time1Hour = ", time1Hour, "\n time1Min = ", time1Min)

			time2x := strings.Split(parts[2], ":")
			fmt.Println("time2x = ", time2x)
			time2Hour, err = strconv.Atoi(time2x[0])
			if err != nil {
				fmt.Println("ошибка", err)
			}
			time2Min, err = strconv.Atoi(time2x[1])
			if err != nil {
				fmt.Println("ошибка", err)
			}
			fmt.Println("time2Hour = ", time2Hour, "\n time2Min = ", time2Min)
			ost += (time2Hour*60 + time2Min) - (time1Hour*60 + time1Min) - obed - stand
			fmt.Println("ost = ", ost)
		}

	}
	fmt.Println("ost = ", ost)
	defer file.Close()

}
