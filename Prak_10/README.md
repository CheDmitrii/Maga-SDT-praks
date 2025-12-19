### Практическое задание №10 Чебыкин Д.К., ПИМО-01-25


### Тема
## JWT-аутентификация: создание и проверка токенов. Middleware для авторизации

# PZ10-AUTH

## Краткое описание проекта

**PZ10-AUTH** — это учебный HTTP API-сервис на Go, реализующий полноценную аутентификацию и авторизацию пользователей с помощью JWT (access/refresh токены), разграничение прав по ролям (admin/user), а также защищённые эндпоинты для профиля и админ-статистики.

---

## Краткое описание: что сделано

В проекте реализованы:
- Регистрация и логин пользователей с выдачей access/refresh JWT-токенов.
- Хранение пользователей в памяти (user_mem.go), поддержка ролей (admin/user).
- Middleware для проверки подлинности токена (authn) и проверки прав (authz).
- Эндпоинты для получения профиля, обновления токенов, просмотра пользователей, админ-статистики.
- Вся логика JWT вынесена в отдельный пакет, поддерживается настройка TTL и секретов через переменные окружения.

---

## Структура проекта

```
Prak_10/
├── assets/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── core/
│   │   ├── service.go
│   │   └── user.go
│   ├── http/
│   │   ├── router.go
│   │   └── middleware/
│   │       ├── authn.go
│   │       └── authz.go
│   ├── platform/
│   │   ├── config/
│   │   │   └── config.go
│   │   └── jwt/
│   │       └── jwt.go
│   └── repo/
│       └── user_mem.go
├── go.mod
├── README.md
```

---

## Как начать работу



### Команда запуска

```bash
go run ./cmd/server
```

---

## Скриншоты

### HAPPY-PATH: ADMIN

#### Логин admin
```bash
curl -Method POST http://localhost:8080/api/v1/login `
  -Headers @{"Content-Type"="application/json"} `
  -Body '{"email":"admin@example.com","password":"secret123"}'
```
![admin_login](foto/login_admin.png)


#### Получить свой профиль /me
```bash
curl -Method GET http://localhost:8080/api/v1/me `
  -Headers @{"Authorization"="Bearer $ADMIN_TOKEN"}
```
![admin_me](foto/get_profile_admin.png)

#### Доступ к /admin/stats
```bash
curl -Method GET http://localhost:8080/api/v1/admin/stats `
  -Headers @{"Authorization"="Bearer $ADMIN_TOKEN"}
```
![admin_stats](foto/get_stat_admin.png)

#### Admin может смотреть любого юзера /users/{id}
```bash
curl -Method GET http://localhost:8080/api/v1/users/1 `
  -Headers @{"Authorization"="Bearer $ADMIN_TOKEN"}
```
![admin_get_user](foto/get_user_admin.png)

#### Admin может обновлять токены (refresh)
```bash
curl -Method POST http://localhost:8080/api/v1/refresh `
  -Headers @{"Content-Type"="application/json"} `
  -Body "{`"Refresh`": `"$ADMIN_REFRESH`"}"
```
![admin_refresh](foto/refresh_token_admin.png)

---

### HAPPY-PATH: USER

#### Логин user
```bash
curl -Method POST http://localhost:8080/api/v1/login `
  -Headers @{"Content-Type"="application/json"} `
  -Body '{"email":"user@example.com","password":"secret123"}'
```
![user_login](foto/login_user.png)



#### User получает свой профиль /me
```bash
curl -Method GET http://localhost:8080/api/v1/me `
  -Headers @{"Authorization"="Bearer $USER_TOKEN"}
```
![user_me](foto/get_user_user.png)

#### User может посмотреть ТОЛЬКО себя /users/{id}
```bash
curl -Method GET http://localhost:8080/api/v1/users/2 `
  -Headers @{"Authorization"="Bearer $USER_TOKEN"}
```
![user_get_self](foto/get_me_user.png)


#### User может обновить токены (refresh)
```bash
$userRefreshResponse = curl -Method POST http://localhost:8080/api/v1/refresh `
  -Headers @{"Content-Type"="application/json"} `
  -Body "{`"Refresh`": `"$USER_REFRESH`"}"
$userRefreshResponse.Content
```
![user_refresh](foto/refresh_token_user.png)

---

### НЕ HAPPY-PATH — проверка ошибок

#### ❌ User НЕ может зайти в /admin/stats
```bash
curl -Method GET http://localhost:8080/api/v1/admin/stats `
  -Headers @{"Authorization"="Bearer $USER_TOKEN"}
```
_Ожидаем: 403 Forbidden_
![user_forbidden_admin_stats](foto/get_stat_user.png)

#### ❌ User НЕ может получить /users/1 (админа)
```bash
curl -Method GET http://localhost:8080/api/v1/users/1 `
  -Headers @{"Authorization"="Bearer $USER_TOKEN"}
```
_Ожидаем: 403 Forbidden_
![user_get_self](foto/get_not_me_user.png)

---

# Контрольные вопросы

## 1. Что такое клеймы JWT и чем отличаются registered, public, private? Почему важно `exp`?

**Клеймы (claims)** — это утверждения (набор данных) внутри JWT, которые описывают пользователя и контекст токена.

### Виды клеймов:
- **Registered claims** — стандартные клеймы, описанные в RFC 7519:
    - `iss` — издатель токена
    - `sub` — субъект (пользователь)
    - `aud` — аудитория
    - `exp` — время истечения
    - `iat` — время выпуска
    - `nbf` — не действителен до
- **Public claims** — пользовательские клеймы, зарегистрированные в публичном реестре или оформленные как URI, чтобы избежать конфликтов имён.
- **Private claims** — произвольные клеймы, используемые по договорённости между сервисами (например, `role`, `user_id`).

### Почему важен `exp`:
- Ограничивает время жизни токена
- Снижает риск компрометации (украденный токен нельзя использовать вечно)
- Является обязательным для безопасной stateless-аутентификации

Без `exp` токен может быть использован неограниченно долго, что является серьёзной уязвимостью.

---

## 2. Чем stateless-аутентификация на JWT отличается от сессионных cookie на сервере? Плюсы и минусы

### JWT (stateless):
**Как работает:**
- Сервер выдаёт токен
- Сервер не хранит состояние сессии
- Вся информация находится внутри JWT

**Плюсы:**
- Хорошо масштабируется (не нужен общий session store)
- Удобно для микросервисов
- Нет зависимости от хранилища сессий

**Минусы:**
- Сложно отзывать токены до `exp`
- JWT больше по размеру
- Ошибки в клеймах = уязвимости

---

### Сессионные cookie (stateful):
**Как работает:**
- Сервер хранит сессию
- Клиент передаёт только идентификатор сессии

**Плюсы:**
- Простая логика отзыва сессии
- Меньше данных на клиенте
- Проще реализовать сложные политики безопасности

**Минусы:**
- Требуется хранилище сессий (Redis, БД)
- Сложнее масштабировать
- Нужно решать sticky sessions или shared storage

---

## 3. Как устроена цепочка middleware и почему AuthZ должна идти после AuthN?

**Middleware** — это цепочка обработчиков, где каждый:
- получает запрос
- может его обработать или прервать
- передаёт управление дальше

### Типичная цепочка:
1. Logging
2. Recovery
3. Authentication (AuthN)
4. Authorization (AuthZ)
5. Business logic

### Почему AuthZ после AuthN:
- **AuthN** отвечает на вопрос: *кто ты?*
- **AuthZ** отвечает на вопрос: *что тебе можно?*

Без успешной аутентификации невозможно корректно проверить права доступа.  
AuthZ без AuthN не имеет смысла, так как нет субъекта для проверки.

---

## 4. RBAC vs ABAC: когда что выбирать? Примеры

### RBAC (Role-Based Access Control):
Доступ определяется **ролью пользователя**.

**Пример:**
- `admin` — полный доступ
- `manager` — редактирование
- `user` — просмотр

**Когда использовать:**
- Простая иерархия прав
- Небольшое количество ролей
- Админ-панели, корпоративные системы

**Минусы:**
- Плохо масштабируется при большом количестве условий
- Роли разрастаются

---

### ABAC (Attribute-Based Access Control):
Доступ определяется **атрибутами** пользователя, ресурса и контекста.

**Пример:**
- Пользователь может редактировать документ, если:
    - он владелец
    - документ в статусе `draft`
    - время запроса — рабочее

**Когда использовать:**
- Сложные бизнес-правила
- Много условий доступа
- Финансовые, облачные и enterprise-системы

**Минусы:**
- Сложнее реализовать и тестировать

---

## 5. Как безопасно хранить пароль и почему нужен bcrypt/argon2 вместо SHA-256? (соль/pepper)

### Неправильный подход:
```text
hash = SHA-256(password)
