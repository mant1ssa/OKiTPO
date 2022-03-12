package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

// Перечисление всех видов нарушений, и их коды (В БД поле Type таблицы Violations)
const (
	OBGON       = iota // Обгон в неположеном месте, 		 код 1
	PREVISHENIE        // Превышение скорости, 	             код 2
	KRASNIY            // Проезд на красный цвет светофора,  код 3
	PARKOVKA           // Парковка в неположенном месте,     код 4
	NEPROPUSTIL        // Непропустил пешехода,      		 код 5
)

// Структура Нарушение, необходима для отображения полей таблицы, будем выводить на экран в сайте
type Violation struct {
	Num  string `json:"StateNumAuto"` // гос. номер автомобиля нарушителя
	Date string `json:"DataTime"`     // потом возможно переформатируем в тип time
	Type int    `json:"Type"`         // вид нарушения (коды нарушений, сверху они перечислены)
	Fine int    `json:"Fine"`         // Размер штрафа
}

func isCorrect(num string) bool { // Функция проверки правильности ввода гос.номера (с 2й по 4й д.б. цифры)
	if _, err := strconv.Atoi(num[1:4]); err == nil {
		return true
	}
	return false
}

func login(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:mysql@/tpo") // Подключаемся к БД
	if err != nil {
		panic(err)
	}
	defer db.Close()

	type input struct { // Определим структуру, в ней будут данные из таблицы users (будет проверка на нового пользователя)
		StateNumAuto string
		Name         string
		Surname      string
		Otchestvo    string
		Password     string
	}

	inp, err := db.Query("SELECT * FROM `users`")
	if err != nil {
		panic(err)
	}

	defer inp.Close()

	fmt.Println("Метод:", r.Method) // получаем информацию о методе запроса
	if r.Method == "GET" {
		t, err := template.ParseFiles("internal/front/registrate.html")
		if err != nil {
			fmt.Fprintf(w, err.Error())
			return
		}
		t.Execute(w, nil)
	} else {
		r.ParseForm()
		// логическая часть процесса входа
		fmt.Println("Пользователь:", r.Form["user_name"])
		fmt.Println("Пароль:", r.Form["user_password"])
		fmt.Println("Номер авто:", r.Form["user_stateNum"])

		// В эти переменные мы получаем введенные пользователем данные, r.FormValue() - метод, внутрь которой можно запихнуть название формочки для ввода, в registrate.html есть аргументы name="...."
		name := r.FormValue("user_name")
		surname := r.FormValue("user_surname")
		otchestvo := r.FormValue("user_otchestvo")
		stateNum := r.FormValue("user_stateNum")
		password := r.FormValue("user_password")

		if !isCorrect(stateNum) {
			fmt.Fprintf(w, "Введите правильный номер авто!")
			http.Redirect(w, r, "/login", 301)
		}

		// Проверка того, что в БД уже нет пользователя с данным гос.номером. Он не может повторяться
		allUsers := []input{}
		for inp.Next() {
			p := input{}
			err := inp.Scan(&p.StateNumAuto, &p.Name, &p.Surname, &p.Otchestvo, &p.Password)
			if err != nil {
				fmt.Println(err)
				continue
			}
			allUsers = append(allUsers, p)
		}
		flag := false
		for i, _ := range allUsers {
			if allUsers[i].StateNumAuto == stateNum {
				flag = true
			}
		}

		if !flag {
			// Здесь в таблицу users нашей БД вносим эти значения
			insert, err := db.Query(fmt.Sprintf("INSERT INTO `users` (`StateNumAuto`, `Name`, `Surname`, `Otchestvo`, `password`) VALUES ('%s', '%s', '%s', '%s', '%s')", stateNum, name, surname, otchestvo, password))
			if err != nil {
				panic(err)
			}
			defer insert.Close()
			http.Redirect(w, r, "/login", 301)
		} else {
			fmt.Fprintf(w, "Похоже вы уже зарегестрированы!")
			http.Redirect(w, r, "/login", 301)
		}
	}
}

func handleFunc() {
	http.HandleFunc("/login", login)
	err := http.ListenAndServe(":9090", nil) // устанавливаем порт для прослушивания
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func main() {
	db, err := sql.Open("mysql", "root:mysql@/tpo")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	fmt.Println("DB is opened")

	/*insert, err := db.Query("INSERT INTO `users` (`StateNumAuto`, `Name`, `Surname`, `Otchestvo`, `password`) VALUES ('a444xe', 'Ural', 'Sur', 'Otch', 'qwerty')")
	if err != nil {
		panic(err)
	}
	defer insert.Close()*/
	handleFunc()
}
