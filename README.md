# 📜 Tender API

Tender API - это API, которая позволяет создать, редактировать тендеры. 

Текущий функионал:
- создание тендера;
- откат версии тендера;
- получение всех тендеров;
- получение тендеров пользователя;
- редактирование тендера.

Планируется добавить:
- взаимодействие с предложениями тендера. Создание отклика на тендер, выдача прав на реализацию задачи в тендера конкретному юзеру;
- аутентификация/авториция;

### Техническая часть
API реализовано на Go веб-фреймворке [Gin](https://github.com/gin-gonic/gin). 

СУБД - PostgreSQL.

Инструмент для миграции - [goose](https://github.com/pressly/goose).


## REST API

### Получение всех тендеров

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

### Получение тендеров, созданных указанным сотрудником

```
GET /api/tenders/my/?username=<employee_username>
```

Параметр `username` указывать обязательно. Также сотрудник с таким `username` должен существовать.

### Примеры ответа

#### Сотрудник существует, у него есть связанные тендеры

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

#### Не указан `username`

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

#### Сотрудника с таким `username` нет

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

#### Если у сотрудника нет тендеров

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

### Создание тендера

```
REQUEST:
POST /api/tenders/new
REQUEST BODY

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
REQUEST BODY
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

#### Нет какого-то/каких-то полей

Если нет каких-то полей, то вернется код `400` и пустой тендер. Также будет сообщение об
ошибке каких полей не хватает.

```
REQUEST:
POST /api/tenders/new
REQUEST BODY
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

#### Указан несуществующий `username`

Если указан несуществующий `username`, то вернется код `400` и пустой тендер

```
REQUEST:
POST /api/tensers/new
REQUEST BODY
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

#### Указана несуществующая организация

Если указана несуществующая организация, то вернется пустой тендер и код `400`

```
REQUEST:
POST /api/tenders/new
REQUEST BODY
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

#### Сотрудник не ответсвенный за орагнизацию

Если сотрудник существует и существует организация, но сотрудник не связано с нею, то вернется код 400 и пустой тендер

```
REQUEST:
POST /api/tenders/new
REQUEST BODY 
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

#### Неверный статус тендера

Тендер нельзя создать сразу со статусом `PUBLISHED` или как-либо другим, отличным от `CREATED`. Если переданный статус неудовлетворяет этому условию, то возвращается код `400` и пустой тендер

```
REQUEST:
POST /api/tenders/new
REQUEST BODY
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

#### Что-то пошло не так

В случае ошибки на сервере вернется пустой тендер и код `500`

```
REQUEST:
POST /api/tenders/new
REQUEST BODY
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

### Обновление тендера
```
REQUEST:
PATCH /api/tenders/:tenderId/edit
REQUEST BODY
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
#### `tenderId` не число

Если `tenderId` в URL не является целым положительным числом, то вернется код `404` и пустой тендер

```
REQUEST:
PATCH /api/tenders/qwe/edit
REQUEST BODY
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

#### Неизвестный статус тендера

Если статус тендера не из списка `CREATED`, `PUBLISHED`, `CLOSED`, то вернется код `400` и пустой тендер

```
REQUEST:
PATCH /api/tenders/1/edit
REQUEST BODY
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