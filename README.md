# Примитивный worker-pool на go

Простой пул воркеров на go. Каждая задачa обрабатывается одним воркером с эмуляцией времени работы, а новые задачи ждут в общей очереди.

## Команды
- add добавлят воркера в пул.
- del <id> удаляет воркера с номером id.
- send <data> отправляет в очередь строку.
- q выход из приложения.

## Требования

Go 1.18+

## Установка и запуск

1. Клонируйте репозиторий:

   ```bash
   git clone https://github.com/artyomkorchagin/vkcontest2025go.git
   cd go-worker-pool
   ```
2. Соберите и запустите приложение:

   ```bash
   go build -o workerpool
   ./workerpool
   ```


Пример работы:

```text
add
Добавлен воркер №0
add
Добавлен воркер №1
send test1
Воркер 0 обрабатывает таску test1
send test2
Воркер 1 обрабатывает таску test2
send test3
Воркер 0 закончил с таском test1
Воркер 0 обрабатывает таску test3
Воркер 1 закончил с таском test2
Воркер 0 закончил с таском test3
del 0
Удален воркер №0
del 44
Нет воркера с номером №44
lalala
Неизвестная команда.
q
Завершение...
```

## Тестирование

В проекте есть базовые тесты:

```bash
go test -v
```

- TestAddAndRemoveWorker, TestSendDataEnqueues проверяют базовую логику.
- TestAsyncProcessing` проверяет асинхронную обработку задач воркерами.

---

![image](https://github.com/user-attachments/assets/63750a3e-de02-4b23-a5e8-d7ea177dcf38)
