package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	a "TestTaskPERX_v2.0/arithmetic"
	"github.com/gorilla/mux"
)

func main() {
	var tempN int
	// Получаем из параметров кол-во одновременно выполняемых задач
	flag.IntVar(&tempN, "N", 1, "Количество одновременно выполняемых задач")
	flag.Parse()
	a.SetWorkers(tempN)

	r := mux.NewRouter()                                   // Создаем роутер для маршрутизации
	r.HandleFunc("/AddTask", a.AddTask).Methods("POST")    // Объявление обработчика на добавление задачи
	r.HandleFunc("/ListTasks", a.ListTasks).Methods("GET") // Объявление обработчика на вывод списка всех задач
	r.HandleFunc("/", a.MainPage).Methods("GET")           // Объявление обработчика главной страницы
	fmt.Println("Сервер запущен на порту :80")             // Информация о запущенном сервере
	log.Fatal(http.ListenAndServe(":80", r))               // Запуск сервера и проверка ошибок

	a.Wg.Wait()
}
