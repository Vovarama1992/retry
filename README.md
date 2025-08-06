# Retry Analytics

Микросервисная система для сбора пользовательских событий с сайта (визиты, клики и т.д.).  
Первый микросервис — `track-service`, принимает сырые события и сохраняет их в PostgreSQL.

---

## 📦 Структура

retry/
├── track-service/ # микросервис сбора событий
│ ├── cmd/ # точка входа
│ └── internal/ # бизнес-логика
├── migrations/ # SQL миграции
├── utils/ # JS-скрипты для вставки на сайт
├── .env.example # переменные окружения
├── docker-compose.yml # сборка и запуск всех сервисов
├── Makefile # команды для сборки/запуска
└── README.md

## 🚀 Запуск

```bash
make refresh

🌐 Endpoint

POST /track/visit
Content-Type: application/json

Пример тела запроса:


{
  "visit_id": "visit_abc123",
  "source": "utm:google",
  "timestamp": "2025-08-05T12:34:56Z"
}
🗃️ База данных
action_types — справочник типов событий (visit, click, …)

actions — таблица событий (visit_id, source, ip, timestamp, …)

🛠 Зависимости
Go 1.23+

Docker + Docker Compose

PostgreSQL 16

migrate/migrate — миграции

📌 TODO
 Swagger для track-service

 Валидация запросов

 Хранение click, scroll, и т.д.

 analytics-service для построения воронок