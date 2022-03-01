package db

// Перечисление всех видов нарушений, и их коды
const (
	OBGON       = iota // Обгон в неположеном месте, 		 код 1
	PREVISHENIE        // Превышение скорости, 	             код 2
	KRASNIY            // Проезд на красный цвет светофора,  код 3
	PARKOVKA           // Парковка в неположенном месте,     код 4
	NEPROPUSTIL        // Непропустил пешехода,      		 код 5
)

// Наша "БД", не будем использовать какую-либо СУБД, т.к. тут мы создаем учебную версию системы
type DB struct {
	Size       uint         // Размер нашей БД (сколько нарушений там на данный момент)
	Violations []*Violation // Срез указателей на нарушения
}

// Структура Нарушение, в ней есть вся информация о конкретном нарушении, нарушения будут храниться в DB
type Violation struct {
	ID           int    // ID нарушения (потом будем сохранять все нарушения за все время)
	Person              // вложенная структура Person, там есть поля Имя Фамилия Отчество
	StateNumAuto string // гос. номер автомобиля нарушителя
	DataTime     string // потом возможно переформатируем в тип time
	//TypeViolation Violation	// вид нарушения ()
	TypeViolation int // вид нарушения (коды нарушений, сверху они перечислены)
	Fine          int // Размер штрафа
}

type Person struct {
	ID                       int    // ID каждого человека, чтоб удобно было искать у него все его нарушения
	Name, Surname, Otchestvo string // поля структуры, Имя Фамилия Отчество
}

func main() {

}
