# (ENG) REST API application for working with user balances

## Description
This microservice provides the ability to register new users. When registering, a wallet with an initial balance of 0 is created. The wallet owner can perform the following actions:

- Check their balance (by default, the balance is displayed in rubles, but the user also has the option to convert their balance to any fiat currency).
- Transfer funds to another user.
- View transaction history.

Transaction history can be filtered by the following parameters:

- By the exact transaction amount (e.g., 200 rubles).
- Within a time frame (e.g., from 1.02 to 4.04).
- In a monetary interval (e.g., transactions amounting from 100 to 500 rubles).

Transaction history can also be sorted:
- In ascending order.
- In descending order.
  
For the convenience of users, pagination of pages is implemented. When making a request, you need to specify the page number and the maximum number of transactions per page.

The deposit and withdraw methods are implemented through the administrator API since such operations can only be performed by service administrators.

# (RU) REST API приложение для работы с балансом пользователей

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
  
### Used Tools and Technologies // Используемые инструменты и технологии 
* Golang
* Gin Web Framework
* Docker
* PostgreSQL
* Swagger
## Running the Microservice // Запуск микросервиса
1. Clone the repository // Клонируйте репозиторий
 ```bash
 git clone https://github.com/Vatset/BankApp.git
```
2. Create a .env file // Создайте .env файл</br>
   Example: // Пример:
```bash
DB_PASSWORD=yourdbpass

SALT_PASSWORD=salt@1219kq

SIGNING_KEY=signingkeyjju777@
```
3.  Preparing the database for operation//Подготовка бд к работе<br>
  *Download and run the Docker application beforehand*// *Предварительно скачайте и запустите приложение docker*<br>
Получение последней версии postgres
```bash   
docker pull postgres
```
Run the Docker container with the name "balance-app," using the previously downloaded PostgreSQL // Запуск Docker контейнера с именем "balance-app", используя ранее скачанный образ PostgreSQL. 
```bash
docker run --name=balance-app -e POSTGRES_PASSWORD="yourdbpass" -p 5436:5432 -d --rm postgres
```
Execute database migrations//Выполнение миграций базы данных
```bash 
migrate -path ./schema  -database 'postgres://postgres:yourdbpass@localhost:5436/postgres?sslmode=disable' up
```
4.Launch the project//Запускаем проект
```bash   
go run cmd/main.go
```

##  Examples of Requests and Responses//Примеры запросов и ответов
After running the project, the API documentation will be available at // После запуска проекта по адресу http://localhost:8080/swagger/index.html будет доступна документация API
### Registration//Регистрация [POST]
```bash   
/auth/sign-up
```
Request//Запрос
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
### Authorization//Авторизация [POST]
```bash   
/auth/sign-in
```
Request//Запрос
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
### Add Funds to User's Wallet // Пополнитель кошелек пользователя [POST]
```bash   
/api_admin/balance/deposit
```
Request//Запрос
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
### Withdraw Money from User's Wallet//Списать деньги с кошелька пользователя [POST]
```bash   
/api_admin/balance/withdraw
```
Request//Запрос
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
### View Wallet Balance // Посмотреть баланс кошелька [GET]
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
### View Wallet Balance in Another Currency//Посмотреть баланс кошелька в другой валюте
```bash   
/api/balance?currency=USD
```
```bash   
{
  "balance": 0.50038,
  "currency": "USD"
}
```
### Transfer to Another User//Перевод другому пользователю [POST]
```bash   
/api/balance/transfer
```
Request//Запрос
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
### Transaction History//История транзакций [GET]
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
