package main

import (
	"bufio"
	"courseWork/backend"
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

func printTable(athletes []backend.Athlete) {
	fmt.Println("Таблиця атлетів:")
	fmt.Printf("%-5s %-20s %-20s %-10s %-10s %-10s %-10s \n", "id", "name", "surname", "run100m", "run3km", "pressCnt", "jumpDist")

	for _, a := range athletes {
		fmt.Printf("%-5d %-20s %-20s %-10.2f %-10.2f %-10d %-10.2f \n", a.Id, a.Name, a.Surname, a.Run100m, a.Run3km, a.PressCnt, a.JumpDistance)
	}
}
