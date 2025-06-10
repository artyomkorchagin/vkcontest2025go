package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Worker struct {
	id      int
	inputCh chan string
	stopCh  chan int
}

type WorkerPool struct {
	workers     map[int]*Worker
	inputCh     chan string
	mu          sync.Mutex
	workerCount int
}

func NewWorkerPool() *WorkerPool {
	return &WorkerPool{
		workers: make(map[int]*Worker),
		inputCh: make(chan string, 100),
	}
}

func (wp *WorkerPool) AddWorker() {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	id := wp.workerCount
	wp.workerCount++

	worker := &Worker{
		id:      id,
		inputCh: make(chan string),
		stopCh:  make(chan int),
	}
	wp.workers[id] = worker

	go func(w *Worker) {
		for {
			select {
			case data := <-w.inputCh:
				fmt.Printf("Воркер №%d обрабатывает таску №%s\n", w.id, data)
				time.Sleep(5 * time.Second)
				fmt.Printf("Воркер %d закончил с таском №%s\n", w.id, data)
				go wp.addTask(w)
			case <-w.stopCh:
				fmt.Printf("Удаляется воркер №%d\n", w.id)
				return
			}
		}
	}(worker)

	go wp.addTask(worker)

	fmt.Printf("Добавлен воркер №%d\n", id)
}

func (wp *WorkerPool) addTask(w *Worker) {
	for data := range wp.inputCh {
		w.inputCh <- data
	}
}

func (wp *WorkerPool) RemoveWorker(id int) {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	if worker, ok := wp.workers[id]; ok {
		worker.stopCh <- 0
		close(worker.stopCh)
		delete(wp.workers, id)
		wp.workerCount--
		fmt.Printf("Удален воркер №%d\n", id)
	} else {
		fmt.Printf("Нет воркера с номером №%d\n", id)
	}
}

func (wp *WorkerPool) SendData(data string) {
	wp.inputCh <- data
}

func main() {
	wp := NewWorkerPool()
	scanner := bufio.NewScanner(os.Stdin)
	for {
		if !scanner.Scan() {
			break
		}
		line := scanner.Text()
		if line == "" {
			continue
		}
		cmd := strings.Fields(line)

		switch cmd[0] {
		case "add":
			wp.AddWorker()
		case "del":
			id, err := strconv.Atoi(cmd[1])
			if err != nil {
				fmt.Println("Неверный номер воркера.")
				continue
			}
			wp.RemoveWorker(id)
		case "send":
			data := strings.Join(cmd[1:], " ")
			wp.SendData(data)
		case "q":
			fmt.Println("Завершение...")
			return
		default:
			fmt.Println("Неизвестная команда.")
		}
	}
}
