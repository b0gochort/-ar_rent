# Car Rental Service

## Описание

Этот проект представляет собой сервис аренды автомобилей с функциональностью проверки статуса автомобиля, получения цены за аренду, создания сессии аренды и получения статистики использования автомобилей.

## Требования

- Go 1.18+
- PostgreSQL

## Настройка

1. Склонируйте репозиторий:
    ```sh
    git clone git@github.com:b0gochort/car_rent.git
    cd car-rental-service
    ```

2. Настройте подключение к базе данных PostgreSQL в файле конфигурации `config.yaml`:
    ```yaml
    postgres:
      user: youruser
      password: yourpassword
      dbname: yourdb
      host: localhost
      port: 5432
    server:
      port: 8080
    ```

Для запуска сервиса выполните:
```sh
go run main.go


API
Проверка статуса автомобиля
URL: /get-status
Метод: GET
Параметры:
car_id (integer) - ID автомобиля

curl -X GET "http://localhost:8080/get-status?car_id=2"


URL: /get-price
Метод: GET
Параметры:
rent_days (integer) - количество дней аренды
Пример запроса:

curl -X GET "http://localhost:8080/get-price?rent_days=1"


Создание сессии аренды
URL: /create-session
Метод: POST
Тело запроса:
{
    "price": 1000,
    "user_id": 1,
    "car_id": 1,
    "rent_days": 1
}

curl -X POST "http://localhost:8080/create-session" \
     -H "Content-Type: application/json" \
     -d '{
           "price": 1000,
           "user_id": 1,
           "car_id": 1,
           "rent_days": 1
         }'

URL: /get-report
Метод: GET
Параметры:
month (integer) - месяц
year (integer) - год

curl -X GET "http://localhost:8080/get-report?month=6&year=2024"







