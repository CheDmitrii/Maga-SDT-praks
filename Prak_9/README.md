# Практическое задание №9, Чебыкин Д.К., ПИМО-01-25


## Тема:
# Реализация регистрации и входа пользователей. Хэширование паролей с bcrypt


# PZ9-AUTH

## Краткое описание проекта

**PZ9-AUTH** — это HTTP API‑сервис на Go, реализующий регистрацию и аутентификацию пользователей с хранением паролей в PostgreSQL через bcrypt. Проект демонстрирует работу с миграциями, безопасное хранение паролей и обработку ошибок при регистрации и входе.

## Требования

- Go ≥ 1.21
- PostgreSQL (локально или через Docker)
- curl/Postman для тестирования API

## Структура проекта

```
Prak_9/
├── assets/
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── core/
│   │   └── user.go
│   ├── http/
│   │   └── handlers/
│   │       └── auth.go
│   ├── platform/
│   │   └── config/
│   │       └── config.go
│   └── repo/
│       ├── postgres.go
│       └── user_repo.go
├── go.mod
└── README.md
```

---

## Как начать работу

### Инициализаци и установка зависимостей

```sh
cd Prak_9
go mod init example.com/pz9-auth
go get github.com/go-chi/chi/v5
go get gorm.io/gorm gorm.io/driver/postgres
go get golang.org/x/crypto/bcrypt
go get github.com/joho/godotenv
```

### Настройка переменных окружения

Создайте файл `.env` и укажите параметры подключения к базе данных и настройки сервера:

```env
DB_DSN=postgres://USERNAME:PASSWORD@HOST:PORT/DBNAME?sslmode=disable
ADDR=:8080
BCRYPT_COST=12
```

- `DB_DSN` — строка подключения к PostgreSQL.
- `ADDR` — адрес и порт для запуска HTTP-сервера.
- `BCRYPT_COST` — стоимость хэширования пароля (рекомендуется 10–14).

В переменной `DB_DSN` подставьте собственные значения на месте следующих переменных:

- ```USERNAME``` — имя пользователя базы данных.
- ```PASSWORD``` — пароль пользователя.
- ```HOST``` — адрес сервера базы данных (```localhost``` для локальной базы или IP/домен удалённого сервера).
- ```PORT``` — порт PostgreSQL (обычно ```5432```).
- ```DBNAME``` — имя базы данных.
- ```sslmode``` — режим SSL (```disable``` для локальной разработки, ```require``` для продакшена при необходимости).

### Запуск приложения

```sh
go run ./cmd/api
```

---

## AutoMigrate

При запуске приложения автоматически выполняется миграция таблицы пользователей:

```go
if err := users.AutoMigrate(); err != nil { log.Fatal("migrate:", err) }
```

Это создаёт таблицу `users` с уникальным индексом по email и полем для хранения хэша пароля.

---

## Скриншоты

### 1. Регистрация пользователя
```sh
curl -Method POST http://localhost:8080/auth/register `
  -ContentType "application/json" `
  -Body '{"email":"testtst@mail.ru","password":"qwerty1234"}'
```
![register_success](foto/register_success.png)

### 2. Регистрация с неправильными данными
```sh
curl -Method POST http://localhost:8080/auth/register `
  -ContentType "application/json" `
  -Body '{"email":"testtst@mail.ru","password":"qwer"}'
```
![register_duplicate](foto/register_bad_data.png)

### 3. Повторная регистрация (ошибка уникальности)
```sh
curl -Method POST http://localhost:8080/auth/register `
  -ContentType "application/json" `
  -Body '{"email":"testtst@mail.ru","password":"qwerty1234"}'
```
![register_duplicate](foto/register_dublicate.png)

### 4. Успешный вход
```sh
curl -Method POST http://localhost:8080/auth/login `
  -ContentType "application/json" `
  -Body '{"email":"testtst@mail.ru","password":"qwerty1234!"}'
```
![login_success](foto/login_suc.png)

### Ошибка входа (неверный пароль)
```sh
curl -Method POST http://localhost:8080/auth/login `
  -ContentType "application/json" `
  -Body '{"email":"user@example.com","password":"wrong"}'
```
![login_fail](foto/login_fail.png)

---

## Фрагменты кода

### Обработчик регистрации

```go
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
    var in registerReq
    if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
        writeErr(w, http.StatusBadRequest, "invalid_json"); return
    }
    in.Email = strings.TrimSpace(strings.ToLower(in.Email))
    if in.Email == "" || len(in.Password) < 8 {
        writeErr(w, http.StatusBadRequest, "email_required_and_password_min_8"); return
    }

    // bcrypt hash
    hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), h.BcryptCost)
    if err != nil {
        writeErr(w, http.StatusInternalServerError, "hash_failed"); return
    }

    u := core.User{Email: in.Email, PasswordHash: string(hash)}
    if err := h.Users.Create(r.Context(), &u); err != nil {
        if err == repo.ErrEmailTaken {
            writeErr(w, http.StatusConflict, "email_taken"); return
        }
        writeErr(w, http.StatusInternalServerError, "db_error"); return
    }

    writeJSON(w, http.StatusCreated, authResp{
        Status: "ok",
        User:   map[string]any{"id": u.ID, "email": u.Email},
    })
}
```

### Обработчик входа

```go
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
    var in loginReq
    if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
        writeErr(w, http.StatusBadRequest, "invalid_json"); return
    }
    in.Email = strings.TrimSpace(strings.ToLower(in.Email))
    if in.Email == "" || in.Password == "" {
        writeErr(w, http.StatusBadRequest, "email_and_password_required"); return
    }

    u, err := h.Users.ByEmail(context.Background(), in.Email)
    if err != nil {
        // не раскрываем, что именно не так
        writeErr(w, http.StatusUnauthorized, "invalid_credentials"); return
    }

    if bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(in.Password)) != nil {
        writeErr(w, http.StatusUnauthorized, "invalid_credentials"); return
    }

    writeJSON(w, http.StatusOK, authResp{
        Status: "ok",
        User:   map[string]any{"id": u.ID, "email": u.Email},
    })
}
```

### Вызовы bcrypt

- `bcrypt.GenerateFromPassword` — используется при регистрации для хэширования пароля.
- `bcrypt.CompareHashAndPassword` — используется при входе для проверки пароля.

---

## Краткие выводы

- **Почему нельзя хранить пароли в открытом виде?**  
  Открытые пароли — это критическая уязвимость. В случае утечки базы злоумышленники получат доступ ко всем аккаунтам.  
  Пароли всегда должны храниться только в виде хэшей.

- **Почему bcrypt?**  
  Bcrypt — современный алгоритм хэширования паролей, устойчивый к перебору и атакам на скорость.  
  Он позволяет настраивать "стоимость" вычисления, замедляя подбор пароля даже при наличии мощного оборудования.

---