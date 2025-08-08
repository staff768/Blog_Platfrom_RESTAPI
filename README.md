

**Название проекта**: BlogPlatform REST API  
**Описание**: Простое REST API для управления постами блога (CRUD операции) с использованием Go и PostgreSQL.

### Функциональные возможности
- Создание постов (POST)
    
- Просмотр всех постов (GET)
    
- Просмотр конкретного поста по ID (GET)
    
- Обновление существующего поста (PUT)
    
- Удаление поста (DELETE)

### Технологический стек

- **Язык**: Go 1.21+
    
- **База данных**: PostgreSQL 15+
    
- **Зависимости**:
    
    - `github.com/lib/pq` - драйвер PostgreSQL
        
    - Стандартная библиотека Go (net/http, encoding/json, database/sql)

### Структура проекта
blogplatform/
├── conf/               # Конфигурация
│   ├── setting.go      # Загрузка конфигурации
│   └── setting.cfg     # Файл настроек
├── internal/
│   └── models/         # Модели и работа с БД
│       ├── error.go    # Кастомные ошибки
│       └── post.go     # Логика работы с постами
├── cmd/
│   ├── main.go         # Точка входа
│   ├── handlers.go     # Обработчики HTTP-запросов
│   └── routes.go       # Маршрутизация
└── go.mod              


### Быстрый старт

#### 1.Установите зависимости
```
go get github.com/lib/pq
```


#### 2.Настройте базу данных

Создайте БД и таблицу:

```
CREATE DATABASE blogplatform;

CREATE TABLE post (
    id SERIAL PRIMARY KEY,
    title VARCHAR(100) NOT NULL,
    content TEXT NOT NULL,
    category VARCHAR(50),
    tags VARCHAR(100)
);
```

#### 3. Настройте конфигурацию

Отредактируйте `conf/setting.cfg`:

```
{
    "ServerHost": "localhost",
    "ServerPort": ":8080",
    "PgHost": "localhost",
    "PgPort": "5432",
    "PgUser": "your_username",
    "PgPass": "your_password",
    "PgBase": "blogplatform"
}
```

#### 4. Запустите приложение

```
go run .\cmd\web\
```

### API Endpoints

| Метод  | Путь           | Описание            | Пример тела запроса                                                                   |
| ------ | -------------- | ------------------- | ------------------------------------------------------------------------------------- |
| POST   | /posts         | Создать пост        | `{"title":"Hello","content":"World","category":"general","tags":"new"}`               |
| GET    | /posts         | Получить все посты  | -                                                                                     |
| GET    | /posts/id?id=1 | Получить пост по ID | -                                                                                     |
| PUT    | /posts/id?id=1 | Обновить пост       | `{"title":"Updated","content":"New content","category":"online","tags":"table game"}` |
| DELETE | /posts/id?id=1 | Удалить пост        | -                                                                                     |