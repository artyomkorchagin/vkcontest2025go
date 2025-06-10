package main

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
	"time"
)

func TestAddAndRemoveWorker(t *testing.T) {
	wp := NewWorkerPool()

	wp.AddWorker()
	wp.AddWorker()
	wp.AddWorker()

	if len(wp.workers) != 3 {
		t.Errorf("ожидаем 3 воркера, получили %d", len(wp.workers))
	}

	wp.RemoveWorker(1)
	time.Sleep(10 * time.Millisecond)

	if _, exists := wp.workers[1]; exists {
		t.Errorf("воркер с номером 1 не был удалён")
	}
	if len(wp.workers) != 2 {
		t.Errorf("ожидаем 2 воркера после удаления, получили %d", len(wp.workers))
	}

	wp.RemoveWorker(99)
}

func TestSendDataEnqueues(t *testing.T) {
	wp := NewWorkerPool()

	wp.SendData("task1")
	wp.SendData("task2")

	if len(wp.inputCh) != 2 {
		t.Errorf("ожидаем 2 таски в очереди, получили %d", len(wp.inputCh))
	}

	first := <-wp.inputCh
	second := <-wp.inputCh
	if first != "task1" || second != "task2" {
		t.Errorf("задачи вышли не в том порядке: got %q, %q", first, second)
	}
}

func TestAsyncProcessing(t *testing.T) {
	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe error: %v", err)
	}
	os.Stdout = w

	wp := NewWorkerPool()
	wp.AddWorker()

	taskData := "42"
	wp.SendData(taskData)

	time.Sleep(6 * time.Second)

	w.Close()
	os.Stdout = oldStdout
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatalf("ошибка чтения потока: %v", err)
	}
	out := buf.String()

	startMsg := "Воркер 0 обрабатывает таску " + taskData
	finishMsg := "Воркер 0 закончил с таском " + taskData

	if !strings.Contains(out, startMsg) {
		t.Errorf("ожидали %q, но его нет в выводе:\n%s", startMsg, out)
	}
	if !strings.Contains(out, finishMsg) {
		t.Errorf("ожидали %q, но его нет в выводе:\n%s", finishMsg, out)
	}
}
