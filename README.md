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
- [ ] перенести API SPEC в Swagger;
- [ ] добавить мониторинг;


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

