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
    curl -X POST http://localhost:8080/api/v1/calculate \
     -H "Content-Type: application/json" \
     -d '{"expression":"1+1"}'
````

* Некорректный запрос (вернет ошибку и статус код 422)
````bash
    curl -X POST http://localhost:8080/api/v1/calculate \
     -H "Content-Type: application/json" \
     -d '{"expression":"1++"}'
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