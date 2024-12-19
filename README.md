# Запуск

# Примеры запросов (Postman или curl)

# Тело запроса если используете postman (на Url : localhost:8080/api/v1/calculate)
````json
    {
        "expression":"1+1"   
    }
````

* Корректный запрос (вернет 2 и статус код 200)
````bash
    curl.exe -X POST http://localhost:8000/api/v1/calculate --header "Content-Type: application/json" --data "{\"expression\":\"1+1\"}"
````
* Ответ

````json
{"result":2,"code":200}
````

* Некорректный запрос (вернет ошибку и статус код 422)
````bash
    curl.exe -X POST http://localhost:8000/api/v1/calculate --header "Content-Type: application/json" --data "{\"expression\":\"1++\"}"
````

* Ответ
````json
{"error":"Expression is not valid","code":422}
````

* Ошибка статус 500
````bash
    curl -X POST http://localhost:8000/api/v1/calculate --header "Content-Type: application/json" --data '{"expression":"10/0"}'
````

* Ответ
````json
{"error":"Invalid JSON format","code":500}
````

* Запуск приложения
````bash
    go run cmd/app/main.go
````

# Запуск тестов

````bash
    go test ./...
````

# Запуск тестов на покрытие
````bash
  cd internal
  go test -cover ./...
````