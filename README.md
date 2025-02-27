# AvitoTestTask

Avito Shop Service — это внутренний сервис для сотрудников

## Нагрузочное тестирование
Тестирование проводилось с помощью команды hey
```sh
hey -n 100000 -c 50 -H "Authorization: Bearer your_token_here" -m GET http://localhost:8080/api/info
```

<img alt="Нагрузочное тестирование" height="500" src="assets/load_test_results.png" width="400"/>

## Структура проекта

Основные директории:
- **/internal**
    - **/cmd** — точка входа
    - **/handlers** — HTTP-хендлеры
    - **/services** — бизнес-логика
    - **/repository** — взаимодействие с базой данных
    - **/config** — конфиг подключения к базе данных
- **/pkg** 
    - **/auth** — реализация JWT
  
## Переменные среды

Сервис использует следующие переменные среды, которые передаются через `docker-compose.yml`:

- **DATABASE_URL**  
  Строка подключения к основной базе данных (например, `postgres://postgres:password@db:5432/shop`).

- **TEST_DATABASE_URL**  
  Строка подключения к тестовой базе данных (например, `postgres://user:test@test_db:5432/shop`).

- **TEST_MODE**  
  Если установлено в `true`, сервис работает в тестовом режиме и использует `TEST_DATABASE_URL` для подключения к тестовой базе. В ней автоматически создаются 100к пользователей для тестирования нагрузок (надо подождать пока выполнится INSERT). Пользователи имеют вид: username: user`$i`, password: pass`$i`, где `$i` - номер пользователя от 1 до 100к.
  
  Если установлено в `false`, используется `DATABASE_URL`.

- **SERVER_PORT**  
  Порт, на котором запускается сервис.

- **CONFIG_PATH**  
  Путь к файлу конфигурации для подписи JWT.

> При изменении переменной `TEST_MODE` необходимо пересоздать сервис, чтобы приложение использовало новые переменные среды:
>```sh
>docker-compose up --force-recreate avito-shop-service -d
>```

## Запуск сервиса

Чтобы собрать и запустить сервис, выполните в корне проекта команду:

```sh
docker-compose up --build
```

Сервис будет доступен по адресу http://localhost:8080.

## Тестирование

### Unit-тесты

Unit-тесты запускаются автоматически при сборке контейнера.

### E2E-тесты

E2E-тесты запускаются вручную.

Чтобы запустить E2E-тесты, выполните следующие шаги:

- Включите тестовый режим(чтобы не захламлять основную бд)

- Зайдите в контейнер сервиса:
```sh
docker exec -it avito-shop-service sh
```

- Запустите E2E-тесты с тегом `e2e`:
```sh
go test ./... -tags=e2e
```

### Linter

Для запуска линтера:

- Установите линтер локально:

```sh
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

- Запустите команду в корне проекта:

```sh
golangci-lint run   
```