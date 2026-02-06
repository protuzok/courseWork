package main

import (
	"bufio"
	"courseWork/shared"
	"fmt"
	"os"
	"strconv"
)

func printOptions() {
	fmt.Println("")
	fmt.Println("============================")
	fmt.Println("0) Вийти з програми")
	fmt.Println("1) Додати поле в таблицю")
	fmt.Println("2) Вивести всю таблицю")
	fmt.Println("3) Видалити поле(-я) з таблиці за id")
	fmt.Println("4) Змінити поле таблиці за id")
	fmt.Println("5) Відсортувати таблицю за run_100m")
	fmt.Println("6) Згрупувати таблицю за max press_cnt та min jump_distance, сортування за іменем")
	fmt.Println("7) Вибрати усіх атлетів з run_3km з відхиленням ±7,359% від середнього")
	fmt.Println("8) Вибрати усіх атлетів з min press_cnt та визначити для них відхилення результату run_100m")
	fmt.Println("9) Вивести людей, у яких загальний результат буде найкращим за всіма показниками")
}

func takeOption() (int, error) {
	input := takeInput("Оберіть опцію:")

	option, err := strconv.ParseInt(input, 10, 64)
	if err != nil {
		return -1, err
	}

	return int(option), nil
}

func takeInput(instruction string) (input string) {
	for {
		fmt.Println("")
		fmt.Println(instruction)
		fmt.Printf("-> ")
		scanner := bufio.NewScanner(os.Stdin)

		if scanner.Scan() {
			input = scanner.Text()

			if input == "" {
				fmt.Println("Пустий рядок. Спробуйте ще раз")
				fmt.Println("")

				continue
			}
		}

		break
	}

	return input
}

func printTable(athletes []shared.Athlete) {
	fmt.Println("Таблиця атлетів:")
	fmt.Printf("%-5s %-20s %-20s %-10s %-10s %-10s %-10s \n", "id", "name", "surname", "run100m", "run3km", "pressCnt", "jumpDist")

	for _, a := range athletes {
		fmt.Printf("%-5d %-20s %-20s %-10.2f %-10.2f %-10d %-10.2f \n", a.Id, a.Name, a.Surname, a.Run100m, a.Run3km, a.PressCnt, a.JumpDistance)
	}
}

func printTableTask4(rows []shared.Task4Row) {
	fmt.Println("Таблиця атлетів:")
	fmt.Printf("%-20s %-10s %-10s %-10s \n", "name", "run100m", "pressCnt", "deviation")

	for _, a := range rows {
		fmt.Printf("%-20s %-10d %-10.2f %-10.2f \n", a.Name, a.PressCnt, a.Run100m, a.Deviation)
	}
}
