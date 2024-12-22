# 💸 Tender API

Tender API - это API, которое позволяет создавать, редактировать тендеры. 

💻 Текущий функионал:
- создание тендера;
- откат версии тендера;
- получение всех тендеров;
- получение тендеров пользователя;
- редактирование тендера.

📜 TODO:
- [ ] взаимодействие с предложениями тендера. Создание отклика на тендер, выдача прав на реализацию задачи в тендера конкретному юзеру;
- [ ] подключить [auth сервис](https://github.com/sariya23/sso) ;
- [x] перенести API doc в Swagger;
- [ ] добавить мониторинг;
- [x] разделить `SERVER_ADDRESS`, а потом собирать конечный адрес в `http.Server{}`;
- [x] добавить в env timeout для сервера;
- [ ] поправить API доку в PUT Rollback Swagger. Нет описания;


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
POSTGRESS_CONN=postgres://POSTGRES_USERNAME:POSTGRES_PASSWORD@POSTGRES_HOST:POSTGRES_PORT/DB_NAME - строка подключения к postgresql
POSTGRES_JDBC_CONN=jdbc:postgresql://POSTGRES_HOST:POSTGRES_PORT/DB_NAME - строка подключения к postgresql в формате jdbc
POSTGRES_USERNAME=POSTGRES_USERNAME - имя пользователя postgres
POSTRGRES_PASSWORD=POSTRGRES_PASSWORD - пароль от postgres
POSTGRES_HOST=POSTGRES_HOST - адрес postgres
POSTGRES_PORT=POSTGRES_PORT - порт postgres
POSTGRES_DB=POSTGRES_DB - название БД в postgres
```

### 🐳 Запуск через докер 

Создайте в корне проекта файл `docker.env` и создайте переменные окружения, описанные выше, только со своими значениями. 

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

