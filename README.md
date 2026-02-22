# Subscription Service

REST сервис для агрегации данных об онлайн подписках пользователей.

## Функциональность

- CRUD операции над записями о подписках
- Подсчет суммарной стоимости подписок с фильтрацией

## Технологии

- Go 1.21
- PostgreSQL
- Gin Framework
- Docker & Docker Compose

## Запуск

1. Клонировать репозиторий
2. Скопировать `.env.example` в `.env` и настроить
3. Запустить Docker Compose:
```bash
docker-compose up -d