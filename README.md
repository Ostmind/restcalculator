# api-calculator

Проект 1: Калькулятор с REST API

## Обязательные Cookie
Token: <Value> - разделение запросов от пользователей по идентификатору

## Описание ручек

POST (application/json)

localhost:8080/{operation}


Принимают числа в формате a, b, c...
Возвращают результат сложения/умножения для данных чисел

## Примеры запросов и ответы


curl --request POST \
--url 'http://localhost:8080/+?=' \
--header 'Content-Type: application/json' \
--header 'User-Agent: insomnia/11.1.0' \
--cookie Token=asd \
--data '{
"Value": [
3,
4,
5,
4
]
}'