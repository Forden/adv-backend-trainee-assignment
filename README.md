# Тестовое задание для backend-стажёра в команду Advertising

#### Запуск:
```$ make && make migrate```

Поднимает контейнер PostgreSQL 13 с юзернеймом `username`, паролем `password` и базой `avito-backend` и запускает сервер API на 8888 порту. Оба контейнера работают в сети хоста (`--network host`).

#### Перезапуск сервера API
```$ make restart_api```

#### Очистка после проверки
```$ make clean```

#### Миграции
``$ make migrate```
Для миграций используется [golang-migrate](https://github.com/golang-migrate/migrate). Миграции хранятся в папке `migrations`.

#### Описание методов
Сервер создан на основе OpenAPI спецификации, хранящейся в `swagger.yml` файле.
