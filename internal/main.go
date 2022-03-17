package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

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

var (
	tmplReg  *template.Template
	tmplLog  *template.Template
	tmplMain *template.Template
)

// Структура Нарушение, необходима для отображения полей таблицы, будем выводить на экран в сайте
type Violation struct {
	Num  string    `json:"StateNumAuto"` // гос. номер автомобиля нарушителя
	Date time.Time `json:"DataTime"`     // потом возможно переформатируем в тип time
	Type int       `json:"Type"`         // вид нарушения (коды нарушений, сверху они перечислены)
	Fine int       `json:"Fine"`         // Размер штрафа
}

// Структура, чтоб в странице mainpg вывести ФИО и номер водителя
type UserOfSite struct {
	Name       string
	Surname    string
	Otchestvo  string
	NumAuto    string
	Violations []Violation // Список его нарушений
}

var newUser UserOfSite // Это будет пользователь, зашедший на сайт

func isCorrect(num string) bool { // Функция проверки правильности ввода гос.номера (с 2й по 4й д.б. цифры)
	if len(num) != 6 {
		return false
	}
	if _, err := strconv.Atoi(num[1:4]); err == nil {
		return true
	}
	return false
}

// Обработчик регистрации
func registr(w http.ResponseWriter, r *http.Request) {
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
		tmplReg.Execute(w, nil)
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

		type Correct struct { // Структура для шаблона, оповещение пользвателя о том, что он неправильно что-то ввел
			IsnotOk1 bool
			IsnotOk2 bool
		}

		// Проверка, правильно ли введена форма ввода номера
		if !isCorrect(stateNum) {

			tmplReg.Execute(w, Correct{IsnotOk1: true, IsnotOk2: false})
			return
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
			http.Redirect(w, r, "/login", 301) // Если не встречаем пользователя с данным номером авто, вносим его в БД
		} else { // и редиректим в форму входа
			tmplReg.Execute(w, Correct{IsnotOk1: false, IsnotOk2: true})
		}
	}
}

// Обработчик входа в учетную запись
func login(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:mysql@/tpo") // Подключаемся к БД
	if err != nil {
		panic(err)
	}
	defer db.Close()

	inp, err := db.Query("SELECT * FROM `users`")
	if err != nil {
		panic(err)
	}

	defer inp.Close()

	fmt.Println("Метод:", r.Method) // получаем информацию о методе запроса
	if r.Method == "GET" {
		tmplLog.Execute(w, nil)
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

		type input struct { // Определим структуру, в ней будут данные из таблицы users (будет проверка на нового пользователя)
			StateNumAuto string
			Name         string
			Surname      string
			Otchestvo    string
			Password     string
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
		// Проверка, есть ли в БД пользователь с этими данными
		for i, _ := range allUsers {
			if (allUsers[i].StateNumAuto == stateNum) && (allUsers[i].Name == name) && (allUsers[i].Surname == surname) && (allUsers[i].Otchestvo == otchestvo) && (allUsers[i].Password == password) {
				flag = true
			}
		}

		if flag {
			// Если есть такой пользователь
			http.Redirect(w, r, "/mainpg", 301)
			newUser.Name = name
			newUser.Surname = surname
			newUser.Otchestvo = otchestvo
			newUser.NumAuto = stateNum
		} else {
			type IsOk struct {
				IsnotCor bool
			}
			tmplLog.Execute(w, IsOk{IsnotCor: true})
		}
	}
}

func mainpg(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:mysql@/tpo") // Подключаемся к БД
	if err != nil {
		panic(err)
	}
	defer db.Close()
	//inp, err := db.Query("SELECT * FROM `violations` WHERE StateNumAuto=?", newUser.NumAuto)
	inp, err := db.Query("SELECT * FROM `violations` WHERE StateNumAuto = ?", newUser.NumAuto)
	if err != nil {
		panic(err)
	}
	defer inp.Close()

	//if r.Method == "GET" {
	//tmplMain.Execute(w, nil)
	//} else {
	r.ParseForm()
	//allViolation := []Violation{}

	type Violation struct {
		Num  string `json:"StateNumAuto"` // гос. номер автомобиля нарушителя
		Date string `json:"DataTime"`     // потом возможно переформатируем в тип time
		Type string `json:"Type"`         // вид нарушения (коды нарушений, сверху они перечислены)
		Fine int    `json:"Fine"`         // Размер штрафа
	}
	/*
		type UserOfSite struct {
			Name       string
			Surname    string
			Otchestvo  string
			NumAuto    string
			Violations []Violation // Список его нарушений
		}
		}
		//allViolation := []input{}
		i := 0
		for inp.Next() {
			p := UserOfSite{}
			//p := input{}
			//err := inp.Scan(&p.Num, &p.Type, &p.Fine)
			//err := inp.Scan(&p.Num, &p.Date, &p.Type, &p.Fine)
			err := inp.Scan(&p.Violations[i].Num, &p.Date, &p.Type, &p.Fine)
			if err != nil {
				fmt.Println(err)
				continue
			}
			//allViolation = append(allViolation, p)
			newUser.Violations = append(newUser.Violations, p)
		}*/
	/*allViolation := []Violation{}
	for inp.Next() {
		p := Violation{}
		err := inp.Scan(&p.Num, &p.Date, &p.Type, &p.Fine)
		//err := inp.Scan(&p.Num, &p.Date, &p.Type, &p.Fine)
		if err != nil {
			fmt.Println(err)
			continue
		}
		allViolation = append(allViolation, p)
		tmplMain.Execute(w, allViolation)
	}*/
	//fmt.Println("HERE!")
	AllVio := []Violation{}
	for inp.Next() {
		p := Violation{}
		err := inp.Scan(&p.Num, &p.Date, &p.Type, &p.Fine)
		if err != nil {
			fmt.Println(err)
			continue
		}
		//fmt.Println(p)
		AllVio = append(AllVio, p)
	}
	tmplMain.Execute(w, AllVio)
	//}
}

// Обработчик нажатия на кнопку Оплатить на главной форме
func deleteVio(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:mysql@/tpo") // Подключаемся к БД
	if err != nil {
		panic(err)
	}
	defer db.Close()

	inp, err := db.Query("DELETE FROM `violations` WHERE StateNumAuto = ?", newUser.NumAuto)
	if err != nil {
		panic(err)
	}
	defer inp.Close()

	http.Redirect(w, r, "/mainpg", 301)
}

func handleFunc() {
	var err error

	tmplLog, _ = template.ParseFiles("front/login.html")
	tmplReg, _ = template.ParseFiles("front/registrate.html")
	tmplMain, _ = template.ParseFiles("front/mainpg.html")

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("front/static/"))))
	http.HandleFunc("/registrate", registr) // Обработчик страницы регистрации
	http.HandleFunc("/login", login)        // Обработчик страницы входа в учетную запись
	http.HandleFunc("/mainpg", mainpg)
	http.HandleFunc("/deleteVio", deleteVio) // Обработчик оплаты штрафа, на главной странице

	err = http.ListenAndServe(":9090", nil) // устанавливаем порт для прослушивания
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

	handleFunc()
}
