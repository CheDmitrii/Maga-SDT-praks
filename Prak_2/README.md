# Практическое задание № 2
# Чебыкин Д.К. ПИМО-01-21

---

# Структура Go-проекта

**Цели:**

-         Понять назначение ключевых директорий (cmd/, internal/, pkg/ и др.).

-         Научиться раскладывать код и артефакты проекта по «правильным» местам.

-         Собрать минимальный скелет проекта и запустить «пустой» main.go.

---

## Требования

- Go ≥ 1.21
- Git

## Функционал

Сервер реализует следующие публичные интерфейсы:

- `/` — возвращает текстовый ответ **"Hello, Go project structure!"**
- `/ping` — возвращает JSON со статусом `status` (`"ok"`) и текущем временем `time` в формате RFC3339
- `/fail` — имитирует ошибку

Проект имеет следующую структуру

```
Prak_2/
├── cmd/
│   └── myapp/
│       └── main.go
├── internal/
│   ├── app/
│   │   ├── app.go
│   │   └── handlers/
│   │       └── ping.go
│   └── utils/
│       ├── httpjson.go
│       └── logger.go
```

## Установка и запуск


Запуск сервера:

```
go run ./cmd/myapp
```

Сборка бинарника:

```
go build -o myapp ./cmd/myapp
./myapp
```

------

## Примеры запросов

```
curl http://localhost:8080/
```

![curl_hello_result.png](misc\curl_root_result.png)

```
curl http://localhost:8080/ping
```

![curl_user_result.png](misc\curl_ping_result.png)

```
curl http://localhost:8080/fail
```

![curl_health_result.png](misc\curl_fail_result.png)

![curluserresultpng](file://C:\Users\dimma\Desktop\ПИШ\технологии создания ПО\ПР2\myapp\misc\curl_ping_result.png?msec=1760021277831)

```
curl -i -H "X-Request-Id: demo-123" http://localhost:8080/ping
```

![curlhealthresultpng](misc/curl_ping_with_request_id_result.png)

Логи на стороне сервера:

![curlhealthresultpng](misc/log_result.png)