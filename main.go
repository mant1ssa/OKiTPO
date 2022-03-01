package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
)

func sayhelloName(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() // Анализирует переданные параметры url, затем анализирует пакет ответа для тела POST (тела запроса)
	// внимание: без вызова метода ParseForm последующие данные не будут получены
	fmt.Println(r.Form) // печатает информацию на сервере
	fmt.Println("Путь: ", r.URL.Path)
	fmt.Println("Схема: ", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])
	for k, v := range r.Form {
		fmt.Println("Ключ: ", k)
		fmt.Println("Значение: ", strings.Join(v, ""))
	}
	fmt.Fprintf(w, "Привет astaxie!") // пишет данные в ответ
}

func login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Метод:", r.Method) // получаем информацию о методе запроса
	if r.Method == "GET" {
		t, _ := template.ParseFiles("login.html")
		t.Execute(w, nil)
	} else {
		r.ParseForm()
		// логическая часть процесса входа
		fmt.Println("Пользователь:", r.Form["username"])
		fmt.Println("Пароль:", r.Form["password"])
	}
}

func main() {
	http.HandleFunc("/", sayhelloName) // устанавливаем правила маршрутизатора
	http.HandleFunc("/login", login)
	err := http.ListenAndServe(":9090", nil) // устанавливаем порт для прослушивания
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
