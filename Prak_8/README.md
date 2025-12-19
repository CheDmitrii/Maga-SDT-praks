# Практическая работа №8 — CRUD API с MongoDB, Чебыкин Д.К, ПИМО-0125

---

## Описание проекта

Цели:

1. Освоить базовые операции работы с MongoDB из Go-приложения.
2. Научиться создавать коллекции, индексы, использовать TTL для автоматического удаления данных.
3. Реализовать REST API с CRUD-операциями (создание, чтение, обновление, удаление заметок).
4. Научиться использовать текстовый поиск и фильтрацию данных.

**Требования:**

- Go 1.21+
- MongoDB
- Git
- Docker / Docker Compose
- curl или Postman для тестирования API

---


## Структура проекта
```
Prak_8/
├── cmd/
│ └── server/
│   └── main.go
├── internal/
│ ├── db/
│ │ └── mongo.go
│ └── notes/
│   ├── handlers.go
│   ├── repo.go
│   ├── repo_test.go
│   └── model.go
├── go.mod
└── docker-compose.yml
```


---

## Запуск проекта

### 1. Запуск MongoDB через Docker


docker-compose.yaml
```
version: '3.9'

services:
  mongo:
    image: mongo:7
    container_name: mongo-dev
    ports:
      - "27017:27017"
    environment:
      MONGO_INITDB_ROOT_USERNAME: root
      MONGO_INITDB_ROOT_PASSWORD: secret
    volumes:
      - mongo_data:/data/db
volumes:
  mongo_data:
```

команда запуска docker-compose.yaml
```bash
docker compose up -d
docker compose ps
```
MongoDB будет доступна на localhost:27017.

### 2. Настройка переменных окружения
Создайте файл .env на основе примера:

```ini
MONGO_URI=mongodb://root:secret@localhost:27017/?authSource=admin
MONGO_DB=pz8
HTTP_ADDR=:8080
```

### 3. Запуск API

```bash
go run ./cmd/api
```

Сервер слушает :8080.
Health check: GET /health → {"status":"ok"}


## Проверка работы API

### 1. Создать заметтку

```bash
curl -s -X POST http://localhost:8080/api/v1/notes \
  -H "Content-Type: application/json" \
  -d "{\"title\":\"First note\",\"content\":\"Hello Mongo!\"}" | jq
```
Результат:
![Создать заметтку](foto/create_note.png)

### 2. Получить список заметок (все, с пагинацией, по тексту*)

все
```bash
curl -s "http://localhost:8080/api/v1/notes"

```
Результат:
![Получить список заметок](foto/get_all_notes.png)
---

с погинацией
```bash
curl -s "http://localhost:8080/api/v1/notes?limit=5&skip=1"

```
Результат:
![Получить список заметок с пагинацией](foto/get_notes_pagination.png)
---
по тексту*
#### задание со звездочкой
```bash
curl -s "http://localhost:8080/api/v1/notes?limit=5&skip=0&q=first"

```
Результат:
![Получить список заметок по тексту](foto/get_notes_by_text.png)

### 3. Получить заметку по ID

```bash
curl -s http://localhost:8080/api/v1/notes/<id>
```
Результат:
![Получить заметку по ID](foto/get_by_id.png)

### 4. Частичное обновление заметки

```bash
curl -s -X PATCH http://localhost:8080/api/v1/notes/<id> \
  -H "Content-Type: application/json" \
  -d '{"content":"Updated content"}'

```
Результат:
![Частичное обновление заметки](foto/update_note.png)

### 5. Удаление заметки

```bash
curl -i -X DELETE http://localhost:8080/api/v1/notes/<id>
```
Результат:
![Удаление заметки](foto/delete_note.png)

### 6. TTL-индекс

#### создаем заметку с ttl
```bash
curl -s -X POST http://localhost:8080/api/v1/notes \
  -H "Content-Type: application/json" \
  -d "{\"title\":\"First note\",\"content\":\"Hello Mongo!\", \"expiresAt\": \"2025-12-18T23:33:00+03:00\"}" | jq
```
Результат:
![создаем заметку с ttl](foto/create_ttl.png)
---

#### прроверяем создание
```bash
curl -s "http://localhost:8080/api/v1/notes"
```
Результат:
![прроверяем создание](foto/check_exist_ttl.png)
---

#### прроверяем спустя время что удалилась заметка
```bash
curl -s "http://localhost:8080/api/v1/notes"
```
Результат:
![прроверяем спустя время что удалилась заметка](foto/delete_note_ttl.png)


### 6. Aggregation pipeline(статистика)
```bash
curl -s "http://localhost:8080/api/v1/notes/stats"
```
Результат:
![Aggregation pipeline(статистика)](foto/stat.png)

### 7. Тестирование
```bash
go test ./...
```
Результат:
![Тестирование](foto/test.png)


---
# Контрольные вопросы по MongoDB

## 1. Чем документная модель MongoDB принципиально отличается от реляционной? Когда она удобнее?

**Отличия:**
- **Документы вместо таблиц:** В MongoDB данные хранятся в коллекциях документов (JSON/BSON), а не в строках таблиц.
- **Гибкая схема:** Документы в одной коллекции могут иметь разные поля и вложенные структуры.
- **Вложенные структуры:** Можно хранить массивы, объекты и сложные данные прямо внутри документа.
- **Отсутствие JOIN:** Связи между документами не поддерживаются напрямую, JOIN делается через агрегатные запросы.

**Когда удобнее:**
- Динамическая или меняющаяся структура данных (заметки, профили пользователей).
- Хранение вложенных объектов (адреса, массивы тегов, файлы).
- Высокая скорость вставки и выборки.
- Горизонтальное масштабирование и работа с большими объемами данных.

---

## 2. Что такое ObjectID и зачем нужен _id? Как корректно парсить/валидировать его в Go?

**ObjectID:**
- Специальный уникальный идентификатор MongoDB для каждого документа.
- Включает время создания, случайное значение и счетчик.
- Используется как `_id` — основной ключ документа, уникальный в коллекции.

**В Go:**
```go
import "go.mongodb.org/mongo-driver/bson/primitive"

oid, err := primitive.ObjectIDFromHex(idHex)
if err != nil {
    // некорректный ObjectID
}
```
**Особенности:**

- Позволяет гарантировать уникальность документов.
- Корректно сериализуется и десериализуется через BSON.



## 3. Какие операции CRUD предоставляет драйвер MongoDB и какие операторы обновления вы знаете?

**CRUD:**
- **Create:** `InsertOne`, `InsertMany`
- **Read:** `FindOne`, `Find` (с фильтрами, сортировкой, пагинацией)
- **Update:** `UpdateOne`, `UpdateMany`, `FindOneAndUpdate`
- **Delete:** `DeleteOne`, `DeleteMany`

**Операторы обновления:**
- `$set` — обновление или добавление поля
- `$inc` — увеличение числового значения
- `$push` — добавление элемента в массив
- `$pull` — удаление элемента из массива
- `$unset` — удаление поля из документа

---

## 4. Как устроены индексы в MongoDB? Как создать уникальный индекс и чем он грозит при вставке?

**Индексы в MongoDB:**
- Используются для ускорения поиска и сортировки документов.
- Типы индексов:
    - **Обычный** — по одному полю.
    - **Составной** — по нескольким полям.
    - **Уникальный** — запрещает дублирование значений.
    - **TTL (Time-To-Live)** — документы автоматически удаляются после заданного времени.

**Создание уникального индекса:**
```js
db.notes.createIndex({ title: 1 }, { unique: true })
```

**Особенности и риски:**

- Уникальный индекс гарантирует, что в коллекции не будет дубликатов указанного поля.
- При попытке вставить документ с уже существующим значением в этом поле произойдет ошибка duplicate key.


## 5. Почему важно использовать context.WithTimeout при вызовах к базе? Что произойдет при его срабатывании?

**Зачем использовать context.WithTimeout:**
- Предотвращает зависание операции при сетевых сбоях или медленных ответах базы.
- Ограничивает максимальное время выполнения запроса к MongoDB.

**Что происходит при срабатывании таймаута:**
- Операция прерывается автоматически.
- Возвращается ошибка `context.DeadlineExceeded`.
- Это позволяет корректно обрабатывать ошибки, освобождать ресурсы и не блокировать приложение.
