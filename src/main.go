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

	// Откройте файл для чтения
	file, err := os.Open("data.txt")
	if err != nil {
		fmt.Println("Не могу прочитать файл2 \n", err)
		return
	}

	// Структура для хранения данных
	var time1Hour, time1Min int
	var time2Hour, time2Min int
	//	var data, time1, time2 []string
	var ost int = 0 // остаток
	// Создаём читателя, связанного с файлом
	obed := 45
	stand := 390
	reader := bufio.NewReader(file)
	for {
		// Читаем строку из файла
		line, err := reader.ReadString('\n')
		// Уберем последний символ перевода строки
		line = strings.TrimSuffix(line, "\n")
		line = strings.TrimSuffix(line, "\r")
		fmt.Println("Считанная строка: ", line)
		// Разделим строку на части
		parts := strings.Split(line, " ")
		// Запишем данные в переменные
		if len(parts) == 3 {
			time1x := strings.Split(parts[1], ":")
			time1Hour, err = strconv.Atoi(time1x[0])
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
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("ошибка", err)
			return
		}
	}
	fmt.Println("ost = ", ost)
	defer file.Close()

}
