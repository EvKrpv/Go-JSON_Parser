JSON Парсер

Резюме реализации

Разработан полнофункциональный JSON парсер с поддержкой:
    Все типы JSON: объекты, массивы, строки, числа, boolean, null
    Рекурсивный парсинг вложенных структур
    PathIterator для обхода по пути с поддержкой объектов и массивов
    Обработка команд через STDIN/STDOUT
    Полная обработка ошибок

Архитектура решения
Основные компоненты:

1. JSONParser - ядро парсинга
2. PathIterator - система обхода по пути
3. Обработчик команд - интерфейс STDIN/STDOUT

Ключевые особенности:
1. Рекурсивный спуск для парсинга
2. Интерфейсы Go для полиморфизма JSON значений
3. Замыкания для реализации итератора
4. Строгая обработка ошибок на каждом этапе

Структура кода
```
// Основные типы данных
type JSONParser struct {
    input string
    pos   int
}

// Методы парсера
Parse() → parseValue() → 
  ├── parseString()
  ├── parseNumber() 
  ├── parseObject()
  ├── parseArray()
  ├── parseBoolean()
  └── parseNull()

// Система итерации
PathIterator() → createIterator()
```
Как работает парсинг

1. Лексический анализ
    Пропуск пробелов (skipSpace())
    Определение типа по первому символу
    Постепенное движение по строке (pos)

2. Синтаксический анализ
    parseValue() - диспетчер по первому символу
    Рекурсивные вызовы для вложенных структур
    Строгая проверка синтаксиса (кавычки, скобки, запятые)

3. PathIterator
    Поддержка путей: "foo.bar.0.baz"
    Работа с объектами (map[string]interface{})
    Работа с массивами ([]interface{})
    Детальные сообщения об ошибках

Запуск и работа с приложением

1. Компиляция и запуск
```
# Сохраните код в файл (например, json_parser.go)
go build json_parser.go
./json_parser
```
2. Формат ввода команд
Приложение работает в интерактивном режиме и понимает две команды:

Команда PARSE - разбор JSON
```
PARSE {"name": "Alice", "age": 25, "hobbies": ["reading", "swimming"]}
```
Команда ITERATE - обход данных
```
ITERATE
ITERATE hobbies
ITERATE address.city
```

3. Примеры работы
Пример 1: Простой объект
```
# Ввод:
PARSE {"name": "Alice", "age": 25, "active": true}

# Затем:
ITERATE
# Вывод:
name: Alice
age: 25
active: true
```
Пример 2: Вложенные объекты
```
# Ввод:
PARSE {"user": {"name": "Bob", "profile": {"age": 30}}}

# Затем:
ITERATE user.profile
# Вывод:
age: 30
```
Пример 3: Массивы
```
# Ввод:
PARSE {"users": ["Alice", "Bob", "Charlie"]}

# Затем:
ITERATE users
# Вывод:
0: Alice
1: Bob
2: Charlie
```

Пример 4: Комплексная структура
```
# Ввод:
PARSE {
  "users": [
    {"name": "Alice", "age": 25},
    {"name": "Bob", "age": 30}
  ],
  "count": 2
}

# Затем:
ITERATE users.0
# Вывод:
name: Alice
age: 25

# Или:
ITERATE users.1.name
# Вывод:
Bob
```
4. Полный сеанс работы
```
$ ./json_parser

# Парсим JSON
PARSE {"name": "John", "age": 30, "hobbies": ["golf", "chess"], "address": {"city": "Moscow", "zip": "12345"}}

# Итерируем корневой уровень
ITERATE
# Вывод:
name: John
age: 30
hobbies: [golf chess]
address: map[city:Moscow zip:12345]

# Итерируем массив hobbies
ITERATE hobbies
# Вывод:
0: golf
1: chess

# Итерируем вложенный объект
ITERATE address
# Вывод:
city: Moscow
zip: 12345

# Получаем конкретное поле
ITERATE address.city
# Вывод:
Moscow
```
5. Работа с файлами (потоковый ввод)

```
# Создайте файл с командами
echo 'PARSE {"users": ["Alice", "Bob"]}
ITERATE users' > commands.txt

# Передайте в приложение
cat commands.txt | ./json_parser
# Вывод:
0: Alice
1: Bob
```
6. Пример файла с несколькими операциями
commands.txt:
```
PARSE {"products": [{"name": "laptop", "price": 999}, {"name": "mouse", "price": 25}]}
ITERATE products.0
PARSE {"status": "success", "data": {"items": [1,2,3]}}
ITERATE data.items
```
```
go run json_parser.go < commands.txt
```
Особенности использования
    Поддерживаемые типы данных:

    Строки: "hello"
    Числа: 123, -45, 3.14
    Булевы значения: true, false
    null: null
    Объекты: {"key": "value"}
    Массивы: [1, 2, 3]

Синтаксис путей:

    ITERATE - корневой уровень
    ITERATE users - поле users
    ITERATE users.0 - первый элемент массива users
    ITERATE address.city - вложенное поле

Обработка ошибок:
```
# Невалидный JSON
PARSE {"name": "John"
# Вывод: Error: invalid JSON

# Несуществующий путь
ITERATE nonexistent
# Вывод: Error: key 'nonexistent' not found in object

# Выход за границы массива
ITERATE users.5
# Вывод: Error: index '5' is out of range
```

Советы по использованию:
    1. JSON должен быть в одной строке (без переносов)
    2. Кавычки обязательны для строковых ключей и значений
    3. Индексы массивов начинаются с 0
    4. Приложение сохраняет состояние между командами

Быстрый старт:
```
# 1. Компилируем
go build json_parser.go

# 2. Запускаем
./json_parser

# 3. Вводим команды:
PARSE {"name": "Test", "values": [10, 20, 30]}
ITERATE values
```

Особенности реализации

    1. Самодостаточность - не использует стандартный пакет encoding/json
    2. Производительность - однопроходный парсинг
    3. Расширяемость - легко добавить новые типы или функции
    4. Надёжность - полная валидация на каждом шаге








