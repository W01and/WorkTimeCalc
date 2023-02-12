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
	var data, time1, time2 []string
	var i int = 0
	// Создайте читателя, связанного с файлом
	reader := bufio.NewReader(file)
	for {
		// Читаем строку из файла
		line, err := reader.ReadString('\n')
		// Уберем последний символ перевода строки
		line = strings.TrimSuffix(line, "\n")
		// Разделите строку на части
		parts := strings.Split(line, " ")
		// Запишите данные в переменные
		data = append(data, parts[0])
		time1 = append(time1, parts[1])
		time2 = append(time2, parts[2])
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("ошибка", err)
			return
		}
		i++
	}
	// ost - остаток
	//	ost:=0
	for j := 0; j <= i; j++ {

		fmt.Printf("Data: %s, time1: %s, time2: %s\n", data[j], time1[j], time2[j])
	}

	defer file.Close()

}
