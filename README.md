# REST API приложение для работы с балансом пользователей

## Описание
Данный микросервис предоставляет возможность регистрации новых пользователей. При регистрации создается кошелек с начальным балансом 0. Владелец кошелька может осуществлять следующие действия:

- Посмотреть свой баланс (по умолчанию баланс отображается в рублях,но пользователь также имеет возможность конвертировать свой баланс в любую фиатную валюту)
- Осуществить перевод денежных средств другому пользователю
- Просмотреть историю операций

Историю транзакций можно фильтровать по следующим параметрам:

- По точной сумме транзакции (например, 200 рублей).
- В промежутке времени (например, с 1.02 по 4.04).
- В денежном интервале (например, транзакции на сумму от 100 до 500 рублей).

Историю транзакций также можно сортировать:

- По возрастанию.
- По убыванию.

Для удобства пользователей реализована пагинация страниц. При запросе необходимо указать номер страницы и максимальное количество транзакций на странице.

Методы `deposit` и `withdraw` реализованы через API администратора, поскольку подобные операции могут выполнять только администраторы сервиса.
  
### Используемые инструменты и технологии
* Golang
* Gin Web Framework
* Docker
* PostgreSQL
* Swagger
## Запуск микросервиса
1. Клонируйте репозиторий
 ```bash
 git clone https://github.com/Vatset/BankApp.git
```
2. Создайте .env файл
   Пример:
```bash
DB_PASSWORD=yourdbpass

SALT_PASSWORD=salt@1219kq

SIGNING_KEY=signingkeyjju777@
```
3. Подготовка бд к работе<br>
   *Предварительно скачайте и запустите приложение docker*<br>
Получение последней версии postgres
```bash   
docker pull postgres
```
Запуск Docker контейнера с именем "balance-app", используя ранее скачанный образ PostgreSQL. 
```bash
docker run --name=balance-app -e POSTGRES_PASSWORD="yourdbpass" -p 5436:5432 -d --rm postgres
```
Выполнение миграций базы данных
```bash 
migrate -path ./schema  -database 'postgres://postgres:yourdbpass@localhost:5436/postgres?sslmode=disable' up
```
4.Запускаем проект
```bash   
go run cmd/main.go
```

## Примеры запросов и ответов
После запуска проекта по адресу http://localhost:8080/swagger/index.html будет доступна документация API
### Регистрация [POST]
```bash   
/auth/sign-up
```
Запрос
```bash   
{
  "name": "Maria",
  "password": "123",
  "username": "mariaV"
}
```
```bash   
{
  "id": 1
}
```
### Авторизация [POST]
```bash   
/auth/sign-in
```
Запрос
```bash   
{
  "password": "123",
  "username": "mariaV"
}
```
```bash   
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTIwMzA3OTcsImlhdCI6MTY5MTk4NzU5NywidXNlcklkIjoxfQ.C853vXKfqi1gAbVnO4QbyCBmabpYdkiQbz4ooUbY800"
}
```
### Пополнитель кошелек пользователя [POST]
```bash   
/api_admin/balance/deposit
```
Запрос
```bash   
{
  "amount": 100,
  "description": "deposit",
  "id": 1
}
```
```bash   
{
  "new balance": 100
}
```
### Списать деньги с кошелька пользователя [POST]
```bash   
/api_admin/balance/withdraw
```
Запрос
```bash   
{
  "amount": 50,
  "description": "sms service",
  "id": 1
}
```
```bash   
{
  "new balance": 50
}
```
### Посмотреть баланс кошелька [GET]
```bash
Authorization: Bearer yJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTIwMzA3OTcsImlhdCI6MTY5MTk4NzU5NywidXNlcklkIjoxfQ.C853vXKfqi1gAbVnO4QbyCBmabpYdkiQbz4ooUbY800 
```
```bash   
/api/balance
```
```bash   
{
  "balance": 50
}
```
### Посмотреть баланс кошелька в другой валюте
```bash   
/api/balance?currency=USD
```
```bash   
{
  "balance": 0.50038,
  "currency": "USD"
}
```
### Перевод другому пользователю [POST]
```bash   
/api/balance/transfer
```
Запрос
```bash   
{
  "amount": 10,
  "description": "dolg",
  "id": 2
}
```
```bash   
{
  "status": "ok"
}
```
### История транзакций [GET]
```bash
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTIwMzA3OTcsImlhdCI6MTY5MTk4NzU5NywidXNlcklkIjoxfQ.C853vXKfqi1gAbVnO4QbyCBmabpYdkiQbz4ooUbY800

api/balance/history?sort_by=amount&sort_order=desc
```
```bash
{
  "data": [
    {
      "id": 1,
      "amount": 100,
      "description": "deposit",
      "date": "14-08-2023 07:45:38"
    },
    {
      "id": 3,
      "amount": -10,
      "description": "transfer to user 2 description: dolg",
      "date": "14-08-2023 07:56:10"
    },
    {
      "id": 2,
      "amount": -50,
      "description": "sms service",
      "date": "14-08-2023 07:48:11"
    }
  ],
  "page": 1,
  "total_pages": 1
}
```
```bash
api/balance/history?sort_field=amount_interval&start_value=-50&end_value=-10&limit=2&page=1
```
```bash
{
  "data": [
    {
      "id": 2,
      "amount": -50,
      "description": "sms service",
      "date": "14-08-2023 07:48:11"
    },
    {
      "id": 3,
      "amount": -10,
      "description": "transfer to user 2 description: dolg",
      "date": "14-08-2023 07:56:10"
    }
  ],
  "page": 1,
  "total_pages": 1
}
```
<p align="center">
  <img src="https://github.com/Vatset/BankApp/assets/88675235/7272b451-ea5e-4eb6-b8ae-b7d99d21c2d4" width="300" />
</p>
<p align="center">
  Technical Specification by https://github.com/avito-tech/autumn-2021-intern-assignment :slightly_smiling_face:
</p>  
