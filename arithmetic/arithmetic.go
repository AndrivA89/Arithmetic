package arithmetic

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"
)

/*
	TimeClock - Структура для времени
*/
type TimeClock struct {
	Hour   int `json:"Hour,omitempty"`
	Minute int `json:"Minute,omitempty"`
	Second int `json:"Second,omitempty"`
}

/*
	Task - структура для задачи
*/
type Task struct {
	Id        int       `json:"index"`               // id - Текущий номер в очереди
	Status    string    `json:"Status"`              // Status - Статус задачи
	TimeSet   TimeClock `json:"TimeSet,omitempty"`   // TimeSet - Время постановки задачи
	TimeStart TimeClock `json:"TimeStart,omitempty"` // TimeStart - Время старта задачи
	Iteration int       `json:"Iteration,omitempty"` // Iteration - Текущая итерация
	N         int       `json:"n"`                   // n - Количество элементов
	D         float32   `json:"d"`                   // D - Дельта между элементами последовательности
	N1        float32   `json:"n1"`                  // n1 - Стартовое значение
	I         float32   `json:"i"`                   // i - Интервал в секундах между итерациями
	TTL       float32   `json:"ttl"`                 // ttl - Время хранения результата в секундах
	TimeDone  TimeClock `json:"TimeDone,omitempty"`  // TimeDone - Время окончания задачи (в случае если задача завершена)
}

/*
	index - Переменная для соблюдения очереди задач
*/
var index int

/*
	n - Количество одновременных задач
*/
var n int

/*
	amount - Количество задач в очереди
*/
var amount int

/*
	taskQueue - Очередь задач
*/
var taskQueue []*Task

/*
	stopChan - Канал для последовательности выполнения задач и ожидания воркеров
*/
var stopChan = make(chan struct{}, 1)

/*
	Wg - WaitGroup для воркеров
*/
var Wg = &sync.WaitGroup{}

/*
	mutex - Блокировка одновременного использования данных
*/
var mutex = &sync.Mutex{}

/*
arithmetic - Функция для расчета арифметической прогрессии
*/
func arithmetic(i int) {
	for {
		_ = <-stopChan // Берем слот из канала ожидания
		defer Wg.Done()
		mutex.Lock()
		taskNow := taskQueue[index]
		log.Printf("Взята задача %v", taskNow.Id) // Проверка последовательности выполнения задач
		index++
		mutex.Unlock()
		result := taskNow.N1 // Переменная для хранения результата прогрессии
		t := time.Now()
		taskNow.TimeStart.Hour, taskNow.TimeStart.Minute, taskNow.TimeStart.Second = t.Clock()
		taskNow.Status = "Processing"
		// Цикл вычисления прогрессии
		for i := 0; i < taskNow.N; i++ {
			taskNow.Iteration = i
			time.Sleep(time.Duration(taskNow.I) * time.Second)
			result += taskNow.D
		}
		t = time.Now()
		taskNow.TimeDone.Hour, taskNow.TimeDone.Minute, taskNow.TimeDone.Second = t.Clock()
		taskNow.Status = "Done"
		Wg.Add(1)
		func(TaskToDel *Task) {
			_ = time.AfterFunc(time.Millisecond*time.Duration((taskNow.TTL)*1000), func() {
				delTaskTotime(TaskToDel)
				Wg.Done()
			})
		}(taskNow)
	}
}

/*
delTaskTotime - Удаление задачи из списка
*/
func delTaskTotime(TaskToDel *Task) {
	mutex.Lock()
	taskQueue = append(taskQueue[:TaskToDel.Id], taskQueue[TaskToDel.Id+1:]...)
	amount--
	for i := TaskToDel.Id; i < amount; i++ {
		taskQueue[i].Id = i
	}
	mutex.Unlock()
}

/*
errorMsg - Печать ошибки
*/
func errorMsg(err error, comment string) {
	if err != nil {
		log.Printf("Ошибка %v! Текст ошибки: %v", comment, err)
	}
}

/*
SetWorkers - Установка количества воркеров
*/
func SetWorkers(tempN int) {
	n = tempN
	for i := 0; i < n; i++ {
		Wg.Add(1)
		go arithmetic(i)
	}
}

/*
MainPage - Обработка главной страницы
*/
func MainPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(n) // Вывод количества одновременных задач
}

/*
AddTask - Обработчик добавления новой задачи
*/
func AddTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var m Task
	err := json.NewDecoder(r.Body).Decode(&m)
	errorMsg(err, "Получение JSOn - функция AddTask")
	t := time.Now()
	m.Status = "Wait"
	m.TimeSet.Hour, m.TimeSet.Minute, m.TimeSet.Second = t.Clock()
	m.Id = amount                     // iD задачи. Для проверки последовательности выполнения
	taskQueue = append(taskQueue, &m) // Добавляем задачу в очередь
	amount++                          // Увеличение количества задач
	go func() {
		stopChan <- struct{}{} // Отправка слота в канал ожидания для воркеров
	}()
}

/*
ListTasks - Обработка вывода списка задач
*/
func ListTasks(w http.ResponseWriter, r *http.Request) {
	for i := 0; i < amount; i++ {
		json.NewEncoder(w).Encode(taskQueue[i]) // Передача списка задач в JSOn
	}
}
