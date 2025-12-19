### Проектирование REST API (CRUD для заметок). Разработка структуры
### Чебыкин Д.К., ПИМО-01-25.

#### Суть практической 
   Освоить принципы проектирования REST API.
   Спроектировать и реализовать CRUD-интерфейс (Create, Read, Update, Delete) для сущности «Заметка».
   Подготовить основу для интеграции с базой данных и JWT-аутентификацией

#### Структура проекта

```bash
Prak_11/
├── cmd/
│   └── api/
│       └── main.go                 # Точка входа в приложение. Инициализирует репозитории, обработчики и запускает HTTP-сервер.
├── internal/                       # Приватный код, который не должен импортироваться внешними проектами (бизнес-логика, реализация).
│   ├── api/
│   │   └── openapi.yaml            # Спецификация REST API в формате OpenAPI (Swagger). Документация для клиентских разработчиков.
│   ├── core/                       # Слой бизнес-сущностей и интерфейсов (абстракций).
│   │   ├── note.go                 # Определение структуры (модели) Note (Заметка) — ключевой сущности приложения.
│   │   └── service/                # Слой бизнес-логики (правила, валидация).
│   │       └── note_service.go     # Реализация сервиса для работы с заметками. Содержит логику, вызываемую обработчиками (Handlers).
│   ├── http/                       # Слой HTTP/API. Отвечает за взаимодействие по сети.
│   │   ├── handlers/               # Обработчики HTTP-запросов. Транслируют HTTP-запросы в вызовы бизнес-логики (Service).
│   │   │   └── notes.go            # Реализация CRUD-обработчиков для ресурса /notes (CreateNote, GetNote и т.д.).
│   │   └── router.go               # Настройка маршрутизации (роутер). Определяет, какой Handler вызывается для каждого пути (URL).
│   └── repo/                       # Слой данных (Repository). Абстрагирует работу с хранилищем.
│       └── note_mem.go             # Реализация репозитория заметок с хранением данных в оперативной памяти (In-Memory).
└── go.mod                          # Файл управления зависимостями Go-проекта.
```
---

#### Инструкция запуска

**Подготовка проекта**  
```bash
mkdir notes-api
cd notes-api
go mod init example.com/notes-api
go get github.com/go-chi/chi/v5
```

## Коды

### Модель данных

```go
package core


import "time"


type Note struct {
  ID        int64
  Title     string
  Content   string
  CreatedAt time.Time
  UpdatedAt *time.Time
}

```
### In-memory репозиторий
``` go
package repo
import (
  "sync"
  "example.com/notes-api/internal/core"
)
type NoteRepoMem struct {
  mu    sync.Mutex
  notes map[int64]*core.Note
  next  int64
}


func NewNoteRepoMem() *NoteRepoMem {
  return &NoteRepoMem{notes: make(map[int64]*core.Note)}
}


func (r *NoteRepoMem) Create(n core.Note) (int64, error) {
  r.mu.Lock(); defer r.mu.Unlock()
  r.next++
  n.ID = r.next
  r.notes[n.ID] = &n
  return n.ID, nil
}
```

### HTTP-обработчик
``` go
package handlers


import (
  "encoding/json"
  "net/http"
  "example.com/notes-api/internal/core"
  "example.com/notes-api/internal/repo"
)
type Handler struct {
  Repo *repo.NoteRepoMem
}
func (h *Handler) CreateNote(w http.ResponseWriter, r *http.Request) {
  var n core.Note
  if err := json.NewDecoder(r.Body).Decode(&n); err != nil {
    http.Error(w, "Invalid input", http.StatusBadRequest)
    return
  }
  id, _ := h.Repo.Create(n)
  n.ID = id
  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(http.StatusCreated)
  json.NewEncoder(w).Encode(n)
}
```
### Маршрутизация
``` go
package httpx


import (
  "github.com/go-chi/chi/v5"
  "example.com/notes-api/internal/http/handlers"
)


func NewRouter(h *handlers.Handler) *chi.Mux {
  r := chi.NewRouter()
  r.Post("/api/v1/notes", h.CreateNote)
  return r
}

``` 
### Точка входа
``` go
package main


import (
  "log"
  "net/http"
  "example.com/notes-api/internal/http"
  "example.com/notes-api/internal/http/handlers"
  "example.com/notes-api/internal/repo"
)


func main() {
  repo := repo.NewNoteRepoMem()
  h := &handlers.Handler{Repo: repo}
  r := httpx.NewRouter(h)


  log.Println("Server started at :8080")
  log.Fatal(http.ListenAndServe(":8080", r))
}
``` 

### Запуск проекта 
``` bash
go run ./cmd/api
``` 
![screen](foto/start_app.png)


### Создание заметки 
``` bash
curl -X POST http://localhost:8080/api/v1/notes \
-H "Content-Type: application/json" \
-d '{"title":"Первая заметка", "content":"Это тест"}'

```
![screen2](foto/create_note.png)

#### Контрольные вопросы и ответы

1. **Что означает аббревиатура REST и в чём её суть?**  
   - Аббревиатура: REST расшифровывается как REpresentational State Transfer (Передача репрезентативного состояния).
   - Суть: Это архитектурный стиль для создания распределённых систем, таких как веб-сервисы. Суть REST заключается в том, что взаимодействие между клиентом и сервером происходит вокруг ресурсов (например, notes, users). Клиент взаимодействует с ресурсом, используя стандартные методы HTTP, и получает его представление (обычно в формате JSON или XML), после чего переходит в новое состояние.

2. **Как связаны CRUD-операции и методы HTTP?**
    CRUD-операции напрямую сопоставляются с HTTP-методами: **Create** → \`POST\`, **Read** → \`GET\`, **Update** → \`PUT\` / \`PATCH\`, **Delete** → \`DELETE\`.

3. **Для чего нужна слоистая архитектура (handler → service → repository)?**
    Она нужна для **разделения ответственности**, что упрощает **тестирование** и повышает **гибкость** и **поддерживаемость** кода.

4. **Что означает принцип «stateless» в REST API?**
    **Stateless** (без состояния) означает, что сервер **не хранит информацию о сессии** клиента между запросами. Каждый запрос должен содержать всю необходимую информацию для своей полной обработки.

5. **Почему важно использовать стандартные коды ответов HTTP?**
    Стандартные коды (2xx, 4xx, 5xx) обеспечивают **единообразие** и **предсказуемость**. Клиент может однозначно определить результат: успех (2xx), ошибка клиента (4xx) или ошибка сервера (5xx).

6. **Как можно добавить аутентификацию в REST API?**
    Наиболее популярный способ — использование **Bearer Токенов** (часто **JSON Web Tokens, JWT**). Токен передаётся в заголовке \`Authorization: Bearer <token>\`.

7. **В чём преимущество версионирования API (например, \`/api/v1/\`)?**
    Версионирование позволяет **развивать API** и вносить несовместимые изменения (в новой версии, \`/v2/\`) без нарушения работы **старых клиентов**, которые продолжают использовать предыдущую версию (\`/v1/\`).