package main

import (
	"bufio"
	"fmt"
	"os"
	"procontext/parser"
	"procontext/pkg/utils"
)

func main() {

	fmt.Println("Курсы валют по отношению к российскому рублю (RUB) за последние 90 дней:")
	//Парсим значения
	data := parser.Parse()
	//Сортируем по названию валюты
	utils.PrintSortedMapByKeys(data)

	fmt.Println("Нажмите клавишу Enter для завершения программы...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

}
