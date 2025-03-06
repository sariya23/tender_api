# 💸 Tender API

Tender API - это API, которое позволяет создавать, редактировать тендеры. 

💻 Текущий функионал:
- создание тендера;
- откат версии тендера;
- получение всех тендеров;
- получение тендеров пользователя;
- редактирование тендера.


### Техническая часть
API реализовано на Go веб-фреймворке [Gin](https://github.com/gin-gonic/gin). 

СУБД - PostgreSQL.

Инструмент для миграции - [goose](https://github.com/pressly/goose).


## ⚙️ REST API

Сейчас доступны следующие эндпоинты:
- `GET /api/ping`
- `GET /api/tenders/`
- `GET /api/tenders/my`
- `POST /api/tenders/new`
- `PATCH /api/tenders/{tenderId}/edit`
- `PUT /api/tenders/{tenderId}/rollback/{version}`

Подробная документация размещена в SwaggerHub: https://app.swaggerhub.com/apis/sariya/tender_api/1.0.0


## 🚀 Локальный запуск

Конфигурация приложения происходит через переменные окружения:

```
SERVER_ADDRESS=0.0.0.0 - адрес сервера
SERVER_PORT=SERVER_PORT - порт сервера
TIMEOUT=TIMEOUT - таймаут на чтение у сервера 
POSTGRESS_CONN=postgres://POSTGRES_USERNAME:POSTGRES_PASSWORD@POSTGRES_HOST:POSTGRES_PORT/DB_NAME - строка подключения к БД
POSTGRESS_CONN_OUTSIDE=postgres://POSTGRES_USERNAME:POSTGRES_PASSWORD@POSTGRES_HOST:POSTGRES_PORT/DB_NAME - строка подключения с БД, которая крутится к контейнере. Добавить нужно и в docker.env и в local.env для совместимости
POSTGRES_JDBC_CONN=jdbc:postgresql://POSTGRES_HOST:POSTGRES_PORT/DB_NAME - строка подключения к БД в формате jdbc
POSTGRES_JDBC_CONN_OUTSIDE=jdbc:postgresql://POSTGRES_HOST:POSTGRES_PORT/DB_NAME - строка подключения к БД, которая крутится в контейнере. Добавить нужно и в docker.env и в local.env для совместимости  
POSTGRES_USERNAME=POSTGRES_USERNAME - имя пользователя postgres
POSTRGRES_PASSWORD=POSTRGRES_PASSWORD - пароль от postgres
POSTGRES_HOST=POSTGRES_HOST - адрес postgres
POSTGRES_PORT=POSTGRES_PORT - порт postgres
POSTGRES_DB=POSTGRES_DB - название БД в postgres
```

Пример находится в `doc/local-example.env`.

### 🐳 Запуск через докер 

Создайте в корне проекта файл `docker.env` и создайте переменные окружения, описанные выше, только со своими значениями. 

**Хост БД должен быть не `localhost`, а название сервиса docker-compose, т.е `db`**

Далее выполните команды:

```
docker-compose --env-file ./docker.env build
docker-compose --env-file ./docker.env up -d
```

Контейнерам нужно некоторое время, чтобы раздуплиться, поэтому нужно подождать секунду 20, чтобы можно было отправлять запросы.

### 💀 Запуск без докера

Создайте в корне проекта файл `local.env` и создайте переменные окружения, описанные выше, только со своими значениями. 

Далее создайте базу данных PostgreSQL с теми же параметрами, которые указаны в env-файле. 

Потом нужно накатить миграции. В корне проекта выполнить:

```
make ENV=local migrate
```

После можно запускать приложение. В корне проекта выполнить:

```
make ENV=local run
```

