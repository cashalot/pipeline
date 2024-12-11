package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"sync"
)

// Функция фильтрации отрицательных чисел с логированием
func filterNegative(ints <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for v := range ints {
			// Логирование действия
			fmt.Println("Processing value in filterNegative:", v)
			if v >= 0 {
				out <- v
				// Логирование успешной фильтрации
				fmt.Println(v, "is not negative, passing to next stage.")
			} else {
				// Логирование игнорирования отрицательных значений
				fmt.Println(v, "is negative and ignored.")
			}
		}
	}()
	return out
}

// Функция фильтрации чисел, не кратных 3 с логированием
func filterNonMultipleOfThree(ints <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for v := range ints {
			// Логирование действия
			fmt.Println("Processing value in filterNonMultipleOfThree:", v)
			if v != 0 && v%3 == 0 {
				out <- v
				// Логирование успешной фильтрации
				fmt.Println(v, "is a multiple of 3, passing to next stage.")
			} else if v != 0 {
				// Логирование игнорирования значений, не кратных 3
				fmt.Println(v, "is not a multiple of 3 and ignored.")
			}
		}
	}()
	return out
}

// Этап получения данных от пользователя с логированием
// Для завершения введи "exit"
func inputData() <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		scanner := bufio.NewScanner(os.Stdin)
		for {
			fmt.Print("Enter a number (or type 'exit' to stop): ")
			scanner.Scan()
			input := scanner.Text()

			if input == "exit" {
				// Логирование завершения ввода
				fmt.Println("Exiting input stage.")
				break
			}

			num, err := strconv.Atoi(input)
			if err != nil {
				// Логирование ошибки ввода
				fmt.Println("Please enter a valid number.")
				continue
			}

			// Логирование полученного числа
			fmt.Println("Received number:", num)
			out <- num // Отправляем число в канал
		}
	}()
	return out
}

// Этап вывода данных с логированием
func outputData(ints <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for v := range ints {
		// Логирование вывода данных
		fmt.Println("Final received:", v)
	}
}

func main() {
	var wg sync.WaitGroup
	wg.Add(1)

	// Начало работы пайплайна
	fmt.Println("Pipeline started.")
	finalOutput := filterNonMultipleOfThree(filterNegative(inputData()))

	go outputData(finalOutput, &wg)

	wg.Wait()

	// Завершение пайплайна
	fmt.Println("Pipeline finished.")
}
