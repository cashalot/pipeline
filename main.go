package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"sync"
)

// Функция фильтрации отрицательных чисел
func filterNegative(ints <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for v := range ints {
			if v >= 0 {
				out <- v
			} else {
				fmt.Println(v, "is negative and ignored.")
			}
		}
	}()
	return out
}

// Функция фильтрации чисел, не кратных 3
func filterNonMultipleOfThree(ints <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for v := range ints {
			if v != 0 && v%3 == 0 {
				out <- v
			} else if v != 0 {
				fmt.Println(v, "is not a multiple of 3 and ignored.")
			}
		}
	}()
	return out
}

// Этап получения данных от пользователя
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
				break // Выход из цикла при вводе "exit"
			}

			num, err := strconv.Atoi(input)
			if err != nil {
				fmt.Println("Please enter a valid number.")
				continue
			}

			out <- num // Отправляем число в канал
		}
	}()
	return out
}

// Этап вывода данных
func outputData(ints <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for v := range ints {
		fmt.Println("Final received:", v)
	}
}

func main() {
	var wg sync.WaitGroup
	wg.Add(1)

	finalOutput := filterNonMultipleOfThree(filterNegative(inputData()))

	go outputData(finalOutput, &wg)

	wg.Wait()
}
