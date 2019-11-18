package arithmetic

import (
	"bytes"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

var (
	w                = httptest.NewRecorder()
	buffer           = new(bytes.Buffer)
	r                = httptest.NewRequest(http.MethodPost, "/add", buffer)
	num              int
	processingStatus int
)

func setN(tempN int) {
	n = tempN
	for i := 0; i < n; i++ {
		Wg.Add(1)
		go arithmetic(i)
	}
}

func add(n int, d float32, n1 float32, i float32, TTL float32) {
	params := "{\"n\": " + strconv.Itoa(int(n)) + ", \"d\": " + strconv.FormatFloat(float64(d), 'f', -1, 32) +
		", \"n1\": " + strconv.FormatFloat(float64(n1), 'f', -1, 32) + ", \"i\": " + strconv.FormatFloat(float64(i), 'f', -1, 32) +
		", \"TTL\": " + strconv.FormatFloat(float64(TTL), 'f', -1, 32) + "}"
	buffer.WriteString(params)
	AddTask(w, r)
}

func TestRun10tasks(t *testing.T) {
	setN(3) // Устанавливаем количество одновременно выполняемых задач в 3
	processingStatus = 3

	for i := 0; i < 10; i++ {
		rand.Seed(time.Now().UnixNano())
		num = rand.Intn(10) + 2
		add(num, 5.13, 13.23, 5.05, 13.13)
	}

	// Проверка корректности переданных параметров для задач
	if taskQueue[0].Status == "Processing" && taskQueue[1].Status == "Processing" && taskQueue[2].Status == "Processing" && taskQueue[3].Status == "Wait" {
		// Все Ок
	} else {
		t.Errorf("Ошибка очередности выполнения задач")
	}
}

func TestRun100tasks(t *testing.T) {
	setN(5) // Устанавливаем количество одновременно выполняемых задач в 5
	processingStatus = 5

	for i := 0; i < 100; i++ {
		rand.Seed(time.Now().UnixNano())
		num = rand.Intn(10) + 2
		add(num, 5.13, 13.23, 5.05, 13.13)
	}

	// Проверка корректности переданных параметров для задач
	if taskQueue[0].Status == "Processing" && taskQueue[1].Status == "Processing" && taskQueue[2].Status == "Processing" && taskQueue[3].Status == "Processing" && taskQueue[4].Status == "Processing" && taskQueue[5].Status == "Wait" {
		// Все Ок
	} else {
		t.Errorf("Ошибка очередности выполнения задач")
	}
}
func TestRun1000tasks(t *testing.T) {
	setN(10) // Устанавливаем количество одновременно выполняемых задач в 10
	processingStatus = 10

	for i := 0; i < 1000; i++ {
		rand.Seed(time.Now().UnixNano())
		num = rand.Intn(10) + 2
		add(num, 5.13, 13.23, 5.05, 13.13)
	}

	// Проверка корректности переданных параметров для задач
	if taskQueue[0].Status == "Processing" && taskQueue[1].Status == "Processing" && taskQueue[2].Status == "Processing" && taskQueue[3].Status == "Processing" && taskQueue[4].Status == "Processing" && taskQueue[5].Status == "Processing" && taskQueue[6].Status == "Processing" && taskQueue[7].Status == "Processing" && taskQueue[8].Status == "Processing" && taskQueue[9].Status == "Processing" && taskQueue[10].Status == "Wait" {
		// Все Ок
	} else {
		t.Errorf("Ошибка очередности выполнения задач")
	}
}

func TestCorrectList(t *testing.T) {
	setN(3) // Устанавливаем количество одновременно выполняемых задач в 3
	processingStatus = 3

	// Добавление тестовых (эталонных) параметров задач
	var sample = make([]*Task, 10)
	for i := 0; i < len(sample); i++ {
		sample[i] = &Task{}
		sample[i].N, sample[i].D, sample[i].N1, sample[i].I, sample[i].TTL = 11+(1*i), 11+float32(i), 111+float32(i), 11.11+float32(i), 11+float32(i)
		add(sample[i].N, sample[i].D, sample[i].N1, sample[i].I, sample[i].TTL)
	}

	// Проверка задач, отправленных на выполнение
	for i := 0; i < processingStatus; i++ {
		if taskQueue[i].Status == "Processing" {
			// OK
		} else {
			t.Errorf("Ошибка вывода списка задач")
		}
	}
	if taskQueue[processingStatus].Status == "Processing" {
		t.Errorf("Ошибка вывода списка задач")
	}

	// Проверка корректности вывода списка путем сравнения с эталоном sample
	for i := 0; i < len(sample); i++ {
		if taskQueue[i].N != sample[i].N && taskQueue[i].D != sample[i].D && taskQueue[i].N1 != sample[i].N1 && taskQueue[i].I != sample[i].I && taskQueue[i].TTL != sample[i].TTL {
			t.Errorf("Ошибка вывода списка задач")
		}
	}
}

func TestTaskQueue(t *testing.T) {
	setN(5) // Устанавливаем количество одновременно выполняемых задач в 5
	processingStatus = 5
	amountMissions := 100

	for i := 0; i < amountMissions; i++ {
		rand.Seed(time.Now().UnixNano())
		num = rand.Intn(10) + 2
		add(num, 5.13, float32(i), 1.05, 1000.13)
	}

	time.Sleep(time.Millisecond * time.Duration(1500*1000))
}
