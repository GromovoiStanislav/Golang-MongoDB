package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

func main() {

	// Вызов первого запроса на заполнение базы
	_, err := http.Post("http://localhost:3000/api/products/populate", "", nil)
	if err != nil {
		fmt.Println("Ошибка при выполнении первого запроса:", err)
		return
	}

	// Запросы на чтение
	url1 := "http://localhost:3000/api/products/frontend"

	// Параметры запроса для второго запроса
	params2 := url.Values{}
	params2.Add("s", "aut")
	params2.Add("sort", "desc")
	url2 := "http://localhost:3000/api/products/backend?" + params2.Encode()

	// Параметры запроса для третьего запроса
	params3 := url.Values{}
	params3.Add("sort", "asc")
	params3.Add("page", "3")
	url3 := "http://localhost:3000/api/products/backend?" + params3.Encode()

	// Выполнение запросов
	result1, err1 := performRequest(url1)
	result2, err2 := performRequest(url2)
	result3, err3 := performRequest(url3)

	// Обработка результатов и ошибок
	handleResult("Первый запрос", result1, err1)
	handleResult("Второй запрос", result2, err2)
	handleResult("Третий запрос", result3, err3)
}

func performRequest(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	result, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func handleResult(desc string, result []byte, err error) {
	if err != nil {
		fmt.Printf("Ошибка при выполнении %s: %v\n", desc, err)
	} else {
		fmt.Printf("Результат %s: %s\n", desc, string(result))
	}
}
