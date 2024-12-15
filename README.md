# 💸 Tender API

Tender API - это API, которая позволяет создавать, редактировать тендеры. 

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

### Методы API
- Тендеры
    - [GET /api/tenders](#получение-всех-тендеров) 
    - [GET /api/tenders/my/?username=<employee_username>](#получение-тендеров-сотрудника)
    - [POST /api/tenders/new](#создание-тендера)
    - [PATCH /api/tenders/tenderId/edit ](#обновление-тендера)
    - [PUT /api/tenders/tenderId/rollback/version](#откат-версии-тендера)

### Содержание ответов
- [Получение всех тендеров](#получение-всех-тендеров)
    - [Ответ в случае, если тендеры есть](#тендеры-есть)
    - [Ответ в случае, если тендеров нет](#тендеров-нет)
    - [Ответ в случае ошибки на сервере](#что-то-пошло-не-так)
- [Получение тендеров сотрудника](#получение-тендеров-сотрудника)
    - [Ответ в случае, если у сотрудника есть тендеры](#у-сотрудника-есть-связанные-тендеры)
    - [Ответ в случае, если не указан username](#не-указан-username)
    - [Ответ в случае, если указан несуществующий сотрудник](#указан-несуществующий-сотрудник)
    - [Ответ в случае, если у сотрудника нет тендеров](#у-сотрудника-нет-тендеров)
    - [Ответ в случае ошибки на сервере](#что-то-пошло-не-так-1)
- [Создание тендера](#создание-тендера)
    - [Ответ в случае успешного создания тендера](#успешное-создание-тендера)
    - [Ответ в случае, если не переданы обязательные поля](#нет-обязательных-полей)
    - [Ответ в случае, если указан несуществующий сотрудник](#указан-несуществующий-сотрудник)
    - [Ответ в случае, если указана несуществующая организация](#указана-несуществующая-организация)
    - [Ответ в случае, если указанный сотрудник неответсвенный за указанную организацию](#сотрудник-не-ответсвенный-за-орагнизацию)
    - [Ответ в случае, если указан неверный статус тендера](#неверный-статус-тендера)
    - [Ответ в случае ошибки на сервере](#что-то-пошло-не-так-2)
- [Ответы при запросе на обновление тендера](#обновление-тендера)
    - [Ответ в случае успешного обновления тендера](#успешное-обновление-тендера)
    - [Ответ в случае, если параметр запроса некорректный](#tenderId-не-число)
    - [Ответ в случае, если статус тендера неверный](#неизвестный-статус-тендера)
    - [Ответ в случае, если невозможно изменить статус тендера при обновлении](#нельзя-изменить-статус)
    - [Ответ в случае, если пользователь неответсвенный за тендер](#указанный-в-теле-пользователь-неответсвенный-за-тендер)
    - [Ответ в случае, если указан несуществующий тендер к обновлению](#тендера-с-таким-tenderId-не-существует)
    - [Ответ в случае, если обновленного сотрудника не существуюет](#обновленного-сотрудника-не-существует)
    - [Ответ в случае, если обновленной организации не существует](#обновленной-организации-не-существует)
    - [Ответ в случае, если обновленный сотрудник неответсвенный за обновеленную организацию](#обновленный-сотрудник-неответсвенный-за-обновеленную-организацию)
    - [Ответ в случае ошибки на сервере](#что-то-пошло-не-так-3)
- [Ответы при запросе на откат версии тендера](#откат-версии-тендера)
    - [Ответ в случае успешного отката](#успешный-откат)
    - [Ответ в случае, если параметр запроса tenderId некорректный](#tenderId-не-число)
    - [Ответ в случае, если параметр запроса version некорректный](#version-не-число)
    - [Ответ в случае, если тендера с указанным tenderId несуществует](#тендер-с-таким-tenderId-не-существует)
    - [Ответ в случае, если у тендера нет указанной версии version](#у-тендера-нет-указанной-версии)
    - [Ответ в случае, если тендер пытается откатить сотрудник, неответсвенный за него](#пользователь-не-создатель-тендера)
    - [Ответ в случае ошибки на сервере](#что-то-пошло-не-так-4)

### Получение всех тендеров

[Назад к списку методов](#методы-API)

```
GET /api/tenders
``` 

Также можно указать параметр запроса `srv_type`, чтобы получить тендеры, относящиеся только к определенному типу услуг

```
GET /api/tender/?srv_type=development
```

По умолчанию в `srv_type` подставляется `all` - то есть вернутся все тендеры.

### Примеры ответа

#### Тендеры есть

Запрос может выполнить любой юзер. Возвращаются только те тендеры, которые опубликованы, то еть имеют статус `PUBLISHED`.  

```
REQUEST:
GET /api/tenders/

RESPONSE:
200 OK
{
    "tenders": {
        {
            "name": "tender 1",
            "description": "asd",
            "service_type": "sell",
            "status": "CREATED",
            "organization_id": 1,
            "creator_username": "sariya"
        },
        {
            ...
        },
        ...
    },
    "message": "ok"
}
```
[К содержанию](#содержание-ответов)

#### Тендеров нет

Когда тендеров нет, то возвращается пустой список тендеров и код `200`.

```
REQUEST:
GET /api/tenders/?srv_type=sell

RESPONSE:
200 OK
{
    "tenders": [],
    "message": "no tenders found with service type=<sell>"
}
```
[К содержанию](#содержание-ответов)

#### Что-то пошло не так

В случае ошибки на стороне сервера возвращается пустой список тендеров и код `500`

```
REQUEST:
GET /api/tenders/

RESPONSE:
500 INTERNAL SERVER ERROR
{
    "tenders": [],
    "message": "internal error"
}
```
[К содержанию](#содержание-ответов)

### Получение тендеров сотрудника

[Назад к списку методов](#методы-API)

```
GET /api/tenders/my/?username=<employee_username>
```

Параметр `username` указывать обязательно. Также сотрудник с таким `username` должен существовать.

### Примеры ответа

#### У сотрудника есть связанные тендеры

В случае успешного запроса вернется код `200` и список связанных тендеров сотрудника с
таким `username`

```
REQUEST:
GET /api/tenders/my/?username=nikita

RESPONSE:
200 OK
{
    "tenders": {
        {
            "name": "tender 1",
            "description": "asd",
            "service_type": "sell",
            "status": "CREATED",
            "organization_id": 1,
            "creator_username": "sariya"
        },
        {
            ...
        },
        ...
    },
    "message": "ok"
}
```
[К содержанию](#содержание-ответов)

#### Не указан username

Если `username` не указан, то возвращается пустой список тендеров и код `400`

```
REQUEST:
GET /api/tenders/my/

RESPONSE:
400 BAD REQUEST
{
    "tenders": [],
    "message": "username query parameter not specified"
}

```
[К содержанию](#содержание-ответов)

#### Указан несуществующий сотрудник

Если сотрудника с таким `username` нет, то возвращается код `404` и пустой список тендеров

```
REQUEST:
GET api/tenders/my/?username=qwe

RESPONSE:
404 NOT FOUND
{
    "tenders": [],
    "message": "employee with username=<qwe> not found"
}
```
[К содержанию](#содержание-ответов)

#### У сотрудника нет тендеров

Если сотрудник с таким `username` существует, но у него нет связанных тендеров, то возвращается код `200` и пустой список тендеров

```
REQUEST:
GET /api/tenders/my/?username=aboba

RESPONSE:
200 OK
{
    "tenders": [],
    "message": "not found tenders for employee with username=<aboba>"
}
```
[К содержанию](#содержание-ответов)

#### Что-то пошло не так

Если произошла ошибка на сервере, то вернется пустой список тендеров и код `500`

```
REQUEST:
GET /api/tenders/my/?username=sariya23

RESPONSE:
500 INTERNAL SERVER ERROR
{
    "tenders": [],
    "message": "internal error"
}
```
[К содержанию](#содержание-ответов)

### Создание тендера

[Назад к списку методов](#методы-API)

```
REQUEST:
POST /api/tenders/new
REQUEST BODY:

{
    "tender": {
        "name": "Tender 1",
        "description": "first created tender",
        "service_type": "sell",
        "status": "CREATED",
        "organization_id": 1,
        "creator_username": "sariya23"
    }
}

RESPONSE:
200 OK
{
    "tender": {
        "name": "Tender 1",
        "description": "first created tender",
        "service_type": "sell",
        "status": "CREATED",
        "organization_id": 1,
        "creator_username": "sariya23"
    },
    "message": "ok"
}
```

Важные уточнения:
- все поля обязательные;
- `organization_id` должен быть целым число от 0 и больше;
- сотрудник с таким `creator_username` должен существовать;
- организация с таким `organization_id` должна существовать;
- сотрудник с указанным `username` должен быть связан с организацией с id `organization_id`;
- статус `status` тендера должен быть `CREATED`.

### Примеры ответов

### Успешное создание тендера
```
REQUEST:
POST /api/tenders/new
REQUEST BODY:
{
    "tender": {
    "name": "Tender 2",
    "description": "first created tender",
    "service_type": "sell",
    "organization_id": 1,
    "status": "CREATED",
    "creator_username": "sariya"
    }
}

RESPONSE:
200 OK
{
    "tender": {
        "name": "Tender 2",
        "description": "first created tender",
        "service_type": "sell",
        "status": "CREATED",
        "organization_id": 1,
        "creator_username": "sariya"
    },
    "message": "ok"
}
```
[К содержанию](#содержание-ответов)

#### Нет обязательных полей

Если нет каких-то полей, то вернется код `400` и пустой тендер. Также будет сообщение об
ошибке каких полей не хватает.

```
REQUEST:
POST /api/tenders/new
REQUEST BODY:
{
    "tender": {
    "name": "Tender 1",
    "description": "first created tender",
    "service_type": "sell",
    "organization_id": 1,
    "creator_username": "aboba"
    }
}

RESPONSE:
400 BAD REQUEST
{
    "tender": {
        "name": "",
        "description": "",
        "service_type": "",
        "status": "",
        "organization_id": 0,
        "creator_username": ""
    },
    "message": "validation failed: Key: 'CreateTenderRequest.Tender.Status' Error:Field validation for 'Status' failed on the 'required' tag"
}
```
[К содержанию](#содержание-ответов)

#### Указан несуществующий сотрудник

Если указан несуществующий `username`, то вернется код `400` и пустой тендер

```
REQUEST:
POST /api/tensers/new
REQUEST BODY:
{
    "tender": {
    "name": "Tender 1",
    "description": "first created tender",
    "service_type": "sell",
    "organization_id": 1,
    "status": "CREATED",
    "creator_username": "zxc"
    }
}

RESPONSE:
400 BAD REQUEST
{
    "tender": {
        "name": "",
        "description": "",
        "service_type": "",
        "status": "",
        "organization_id": 0,
        "creator_username": ""
    },
    "message": "employee with username=<zxc> not found"
}
```
[К содержанию](#содержание-ответов)

#### Указана несуществующая организация

Если указана несуществующая организация, то вернется пустой тендер и код `400`

```
REQUEST:
POST /api/tenders/new
REQUEST BODY:
{
    "tender": {
        "name": "Tender 1",
        "description": "first created tender",
        "service_type": "sell",
        "organization_id": 10,
        "status": "CREATED",
        "creator_username": "sariya"
    }
}

RESPONSE:
400 BAD REQUEST
{
    "tender": {
        "name": "",
        "description": "",
        "service_type": "",
        "status": "",
        "organization_id": 0,
        "creator_username": ""
    },
    "message": "organization with id=<10> not found"
}
```
[К содержанию](#содержание-ответов)

#### Сотрудник не ответсвенный за орагнизацию

Если сотрудник существует и существует организация, но сотрудник не связано с нею, то вернется код 400 и пустой тендер

```
REQUEST:
POST /api/tenders/new
REQUEST BODY:
{
    "tender": {
        "name": "Tender 1",
        "description": "first created tender",
        "service_type": "sell",
        "organization_id": 2,
        "status": "CREATED",
        "creator_username": "aboba"
    }
}

RESPONSE:
400 BAD REQUEST
{
    "tender": {
        "name": "",
        "description": "",
        "service_type": "",
        "status": "",
        "organization_id": 0,
        "creator_username": ""
    },
    "message": "employee <aboba> not responsible for organization with id=<2>"
}
```
[К содержанию](#содержание-ответов)

#### Неверный статус тендера

Тендер нельзя создать сразу со статусом `PUBLISHED` или как-либо другим, отличным от `CREATED`. Если переданный статус неудовлетворяет этому условию, то возвращается код `400` и пустой тендер

```
REQUEST:
POST /api/tenders/new
REQUEST BODY:
{
    "tender": {
        "name": "Tender 1",
        "description": "first created tender",
        "service_type": "sell",
        "organization_id": 1,
        "status": "CLOSED",
        "creator_username": "sariya"
    }
}

RESPONSE:
400 BAD REQUEST
{
    "tender": {
        "name": "",
        "description": "",
        "service_type": "",
        "status": "",
        "organization_id": 0,
        "creator_username": ""
    },
    "message": "cannot create tender with status <CLOSED>"
}
```
[К содержанию](#содержание-ответов)

#### Что-то пошло не так

В случае ошибки на сервере вернется пустой тендер и код `500`

```
REQUEST:
POST /api/tenders/new
REQUEST BODY:
{
    "tender": {
    "name": "Tender 1",
    "description": "first created tender",
    "service_type": "sell",
    "organization_id": 1,
    "status": "CREATED",
    "creator_username": "sariya"
    }
}

RESPONSE:
500 INTERNAL SERVER ERROR
{
    "tender": {
        "name": "",
        "description": "",
        "service_type": "",
        "status": "",
        "organization_id": 0,
        "creator_username": ""
    },
    "message": "internal error"
}
```
[К содержанию](#содержание-ответов)

### Обновление тендера

[Назад к списку методов](#методы-API)

```
При обновлении тендера его версия увеличивается автоматически и возвращатся при получении будет именно она. Обновить тендер может только тот пользователь, который его создал.

REQUEST:
PATCH /api/tenders/:tenderId/edit
REQUEST BODY:
{
    "update_tender_data": {
        "name": "Updated tender name"
    },
    "username": "sariya"
}

RESPONSE:
200 OK 
{
    {
    "updated_tender": {
        "name": "Updated tender name",
        "description": "zxc",
        "service_type": "zxc",
        "status": "CREATED",
        "organization_id": 1,
        "creator_username": "qwe"
    },
    "message": "ok
}
}
```

Важные уточнения:
- `tenderId` должен быть целым положительным числом;
- если хотите обновить статус тендера, то учьтите, что тендер из статуса `PUBLISHED` нельзя перевести в `CREATED` и из статуса `CLOSED` в `CREATED`;
- статус тендера должен быть из этого списка: `CREATED`, `PUBLISHED`, `CLOSED`;
- при обновлении организации она должна существовать;
- при обновлении сотрудника он должен существовать;
- если обновляется и сотрудник и организация, то этот новый сотрудник должен быть ответсвенным за эту новую организацию;
- если обновляется только (в контексте сотрудник, организация) сотрудник, то обновленный сотрудник должен быть ответсвенным за орагнизацию, которая сейчас установлена в тендере;
- если обновляется только организация, то текущий сотрудник в тендере должен быть ответсвенным за эту обновленную организацию;
- если сотрудник с `username`, указанным в теле запроса не создавал тендер, то он не может его обновить;
- поле `username` в теле запроса обязательно.

### Примеры ответов 

#### Успешное обновление тендера
```
REQUEST:
PATCH /api/tenders/1/edit
REQUEST BODY:
{
    "update_tender_data": {
        "name": "Updated Tender"
    },
    "username": "sariya"
}

RESPONSE:
200 OK
{
    "updated_tender": {
        "name": "Updated Tender",
        "description": "first created tender",
        "service_type": "sell",
        "status": "PUBLISHED",
        "organization_id": 1,
        "creator_username": "sariya"
    },
    "message": "ok"
}
```
[К содержанию](#содержание-ответов)

#### tenderId не число

Если `tenderId` в URL не является целым положительным числом, то вернется код `404` и пустой тендер

```
REQUEST:
PATCH /api/tenders/qwe/edit
REQUEST BODY:
{
    "update_tender_data": {
        "status": "CLOSED"
    },
    "username": "aboba"
}

RESPONSE:
400 BAD REQUEST
{
    "updated_tender": {
        "name": "",
        "description": "",
        "service_type": "",
        "status": "",
        "organization_id": 0,
        "creator_username": ""
    },
    "message": "tenderId must be positive integer number"
}
```
[К содержанию](#содержание-ответов)

#### Неизвестный статус тендера

Если статус тендера не из списка `CREATED`, `PUBLISHED`, `CLOSED`, то вернется код `400` и пустой тендер

```
REQUEST:
PATCH /api/tenders/1/edit
REQUEST BODY:
{
    "update_tender_data": {
        "status": "qwe"
    },
    "username": "aboba"
}

RESPONSE:
400 BAD REQUEST
{
    "updated_tender": {
        "name": "",
        "description": "",
        "service_type": "",
        "status": "",
        "organization_id": 0,
        "creator_username": ""
    },
    "message": "tender status=<qwe> is unknown"
}
```
[К содержанию](#содержание-ответов)

#### Нельзя изменить статус

(
если хотите обновить статус тендера, то учьтите, что тендер из статуса `PUBLISHED` нельзя перевести в `CREATED` и из статуса `CLOSED` в `CREATED`.
)

Если в запросе прилетает обновление статуса на `CREATED`, а текущий статус тендера, например, `CLOSED`, то вернется код `400` и пустой тендер

```
REQUEST:
PATCH /api/tenders/1/edit
REQUEST BODY:
{
    "update_tender_data": {
        "status": "CREATED"
    },
    "username": "sariya"
}

RESPONSE:
400 BAD REQUEST
{
    "updated_tender": {
        "name": "",
        "description": "",
        "service_type": "",
        "status": "",
        "organization_id": 0,
        "creator_username": ""
    },
    "message": "cannot set this tender status. Cannot set tender status from PUBLISHED to CREATED and from CLOSED to CREATED"
}
```
[К содержанию](#содержание-ответов)

#### Указанный в теле пользователь неответсвенный за тендер

Если в теле запроса указан пользователь, который неответсвенный за обновляемый тендер, то вернется код `403` и пустой тендер

```
REQUEST:
PATCH /api/tenders/1/edit
REQUEST BODY:
{
    "update_tender_data": {
        "status": "CLOSED"
    },
    "username": "zxc"
}

RESPONSE:
403 FORBIDDEN
{
    "updated_tender": {
        "name": "",
        "description": "",
        "service_type": "",
        "status": "",
        "organization_id": 0,
        "creator_username": ""
    },
    "message": "employee with username=<zxc> not creator of tender with id=<1>"
}
```
[К содержанию](#содержание-ответов)

#### Тендера с таким tenderId не существует

Если в URL указать id тендера, которого не существует, то вернется код 404 и пустой тендер

```
REQUEST:
PATCH /api/tenders/10/edit
REQUEST BODY:
{
    "update_tender_data": {
        "status": "CLOSED"
    },
    "username": "zxc"
}

RESPONSE:
404 NOT FOUND
{
    "updated_tender": {
        "name": "",
        "description": "",
        "service_type": "",
        "status": "",
        "organization_id": 0,
        "creator_username": ""
    },
    "message": "tender with id=<10> not found"
}
```
[К содержанию](#содержание-ответов)

#### Обновленного сотрудника не существует
Если попытаться обновить сотрудника тендера, которого не сущесвует, то вернется код `400` и пустой тендер

```
REQUEST:
PATCH /api/tenders/1/edit
REQUEST BODY:
{
    "update_tender_data": {
        "creator_username": "shrek"
    },
    "username": "sariya"
}

RESPONSE:
400 BAD REQUEST
{
    "updated_tender": {
        "name": "",
        "description": "",
        "service_type": "",
        "status": "",
        "organization_id": 0,
        "creator_username": ""
    },
    "message": "updated employee with username=<shrek> not found"
}
```
[К содержанию](#содержание-ответов)

#### Обновленной организации не существует
Если попытаться обновить организацию на ту, которая не существует, то вернется код 400 и пустой тендер

```
REQUEST:
PATCH /api/tenders/1/edit
REQUEST BODY:
{
    "update_tender_data": {
        "organization_id": 20
    },
    "username": "sariya"
}

RESPONSE:
400 BAD REQUEST
{
    "updated_tender": {
        "name": "",
        "description": "",
        "service_type": "",
        "status": "",
        "organization_id": 0,
        "creator_username": ""
    },
    "message": "updated organization with id=<20> not found"
}
```
[К содержанию](#содержание-ответов)

#### Обновленный сотрудник неответсвенный за обновеленную организацию
Если обновеленный пользователь неответсвенный за новую организацию, то вернется код 400 и пустой тендер

```
REQUEST:
PATCH /api/tenders/1/edit
REQUEST BODY:
{
    "update_tender_data": {
        "creator_username": "aboba",
        "organization_id": 2
    },
    "username": "sariya"
}

RESPONSE:
400 BAD REQUEST
{
    "updated_tender": {
        "name": "",
        "description": "",
        "service_type": "",
        "status": "",
        "organization_id": 0,
        "creator_username": ""
    },
    "message": "employee with username=<aboba> not responsible for organization with id=<2>"
}
```
[К содержанию](#содержание-ответов)

#### Что-то пошло не так
В случае ошибки на сервере вернется код 500 и пустой тендер

```
REQUEST:
PATCH /api/tenders/1/edit
REQUEST BODY:
{
    "update_tender_data": {
        "name": "Updated Tender"
    },
    "username": "sariya"
}

RESPONSE:
500 INTERNAL SERVER ERROR
{
    "updated_tender": {
        "name": "",
        "description": "",
        "service_type": "",
        "status": "",
        "organization_id": 0,
        "creator_username": ""
    },
    "message": "internal error"
}
```
[К содержанию](#содержание-ответов)

### Откат версии тендера

[Назад к списку методов](#методы-API)

```
REQUEST:
PUT /api/tenders/tenderId/rollback/version
REQUEST BODY:
{
    "username": "kapi"
}

RESPONSE:
CODE
{
    {
    "rollback_tender": {
        ...
    },
    "message": "..."
}
}
```

При откате версии тендера автоматически деактивируется предыдущая версия и активной становится та версия, которая указана в запросе. Откатить тендер может тот пользователь, который его создал.

### Примеры ответов

#### Успешный откат
Откатить тендер получится если тендер с таким `tenderId` есть, у него есть версия, указанная в `version` и в теле передается пользователь, который создал этот тендер

```
REQUEST:
PUT /api/tenders/1/rollback/2
REQUEST BODY:
{
    "username": "kapi"
}

RESPONSE:
200 OK
{
    "rollback_tender": {
        "name": "Tender 1",
        "description": "first created tender",
        "service_type": "sell",
        "status": "PUBLISHED",
        "organization_id": 1,
        "creator_username": "kapi"
    },
    "message": "ok"
}
```
[К содержанию](#содержание-ответов)

#### tenderId не число

В случае если `tenderId` в URL не целое положительное число, то вернется код `404` и пустой тендер

```
REQUEST:
PUT api/tenders/qwe/rolback/2
REQUEST BODY:
{
    "username": "kapi"
}

RESPONSE:
404 NOT FOUND
{
    "rollback_tender": {
        "name": "",
        "description": "",
        "service_type": "",
        "status": "",
        "organization_id": 0,
        "creator_username": ""
    },
    "message": "tenderId must be positive integer number"
}
```
[К содержанию](#содержание-ответов)

#### version не число

В случае если `version` в URL не целое положительное число, то вернется код 404 и пустой тендер

```
REQUEST:
PUT /api/tenders/1/rollback/qwe
REQUEST BODY:
{
    "username": "kapi"
}

RESPONSE:
404 NOT FOUND
{
    "rollback_tender": {
        "name": "",
        "description": "",
        "service_type": "",
        "status": "",
        "organization_id": 0,
        "creator_username": ""
    },
    "message": "version must be positive integer number"
}
```
[К содержанию](#содержание-ответов)

#### Тендер с таким tenderId не существует
```
REQUEST:
PUT /api/tenders/10/rollback/2
REQUEST BODY:
{
    "username": "sariya"
}

RESPONSE:
404 NOT FOUND
{
    "rollback_tender": {
        "name": "",
        "description": "",
        "service_type": "",
        "status": "",
        "organization_id": 0,
        "creator_username": ""
    },
    "message": "tender with id=<10> not found"
}
```
[К содержанию](#содержание-ответов)

#### У тендера нет указанной версии
```
REQUEST:
PUT /api/tenders/1/rollback/20
REQUEST BODY:
{
    "username": "kapi"
}

RESPONSE:
404 NOT FOUND
{
    "rollback_tender": {
        "name": "",
        "description": "",
        "service_type": "",
        "status": "",
        "organization_id": 0,
        "creator_username": ""
    },
    "message": "tender with id=<1> doesnt have version=<20>"
}
```
[К содержанию](#содержание-ответов)

#### Пользователь не создатель тендера

Если пользователь, который не создавал тендер, попытается откатить его, то он получить код 403)

```
REQUEST:
PUT api/tenders/1/rollback/2
REQUEST BODY:
{
    "username": "ne kapi"
}

RESPONSE:
403 FORBIDDEN
{
    "rollback_tender": {
        "name": "",
        "description": "",
        "service_type": "",
        "status": "",
        "organization_id": 0,
        "creator_username": ""
    },
    "message": "employee with username=<ne kapi> not creator of tender with id=<1>"
}
```
[К содержанию](#содержание-ответов)

#### Что-то пошло не так
В случае ошибки на сервере вернется код 500 и пустой тендер

```
REQUEST:
PUT api/tenders/1/rollback/2
REQUEST BODY:
{
    "username": "kapi"
}

RESPONSE:
500 INTERNAL SERVER ERROR
{
    "rollback_tender": {
        "name": "",
        "description": "",
        "service_type": "",
        "status": "",
        "organization_id": 0,
        "creator_username": ""
    },
    "message": "internal error"
}
```
[К содержанию](#содержание-ответов)