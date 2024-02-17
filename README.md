
# EWallet - тестовое задание от Infotecs

По заданию необходимо разработать приложение EWallet реализующее систему обработки транзакций платёжной системы. Приложение должно быть реализовано в виде HTTP сервера, реализующее REST API. Сервер должен реализовать 4 метода и их логику:

- Создание кошелька;
- Перевод средств с одного кошелька на другой;
- Получение историй входящих и исходящих транзакций;
- Получение текущего состояния кошелька.

Также необходимо было обратить внимание при реализации на: 
- Безопасность: в приложении не должно быть уязвимостей, позволяющих произвольно менять данные в базе.
- Персистентность: данные и изменения не должны «теряться» теряться при перезапуске приложения. 

Текстовый файл с заданием находится в папке [/docs](https://github.com/egor-denisov/wallet-infotecs/tree/main/docs).

## Вопросы, возникшие при реализации

Во время реализации задания возникли вопросы, на которые ответила поддержка:

1) При реализации переводов мы даем пользователю переводить на свой же кошелек или возвращаем ошибку? *Рекомендуем возвращать ошибку с кодом HTTP 400.*
2) Нужно ли обязательно использовать числа с плавающей точкой или можно выбрать конкретное количество знаков после запятой для исключения ошибок связанных с типом данных float? *Следует использовать числа с плавающей точкой. Ошибки, связанные с неточностью представления вещественных чисел, предлагаем игнорировать как несущественные.*
3) При возврате всех транзакций, если они отсутствуют у пользователя возвращаем пустой массив или null? *Следует возвращать пустой массив. Это указано в спецификации API.*

## Ошибки при реализации

После сдачи задания, обнаружились проблемы:

1) При реализации переводов у пользователя есть возиожность переводить неположительную сумму. 
2) В техническом задании указано, что история выводится как для исходящих, так и для входящих транзакций. Я не заметил, условие про входящие, что привело к неправильной работе данного эндпоинта.

Все проблемы были исправлены.

## Как запустить?

Для запуска приложения в контейнере, необходимо выполнить команду:
```
make build
```

Для запуска вне контейнера, выполняем две команды:
```
make get

make run
```

## Доступные скрипты

- `make build` - запуск контейнеров

- `make run` - запуск приложения go

- `make get` - загрузка используемых пакетов

- `make test` - запуск тестов

- `make swag` - генерация новой спецификации


## Переменные окружения и конфигурация

Для корректной работы приложения необходимо указать переменные окружения в файле **.env** корневой директории. Ниже представлены сами переменные и краткое описание:

`DISABLE_SWAGGER_HTTP_HANDLER` - при существовании данной переменной Swagger не работает.

`GIN_MODE` - устанавливается в режим debug при необходимости отладки.

`HTTP_PORT` - порт по которому будет доступно приложение.

`POSTGRES_USER`, `POSTGRES_DB`, `POSTGRES_HOST`, `POSTGRES_PORT`, `POSTGRES_PASSWORD` - параметры для подключения к базе данной Postgresql.

Также присутствует файл [config.yaml](https://github.com/egor-denisov/wallet-infotecs/blob/main/config/config.yml) в котором указываются остальные данные (название и версия приложения, стандартный баланс и др.).

## Архитектура приложения

Приложение написано опираясь на принципы чистой архитектуры, приведу ниже, что в себе хранит каждая папка.

- `/cmd/app` - хранит файл для запуска приложения.
- `/config` - хранит в себе файлы конфигурации и средство для их чтения.
- `/docs` - документация Swagger и файл с заданием.
- `/internal/app` - сборка основных компонентов воедино.
- `/internal/controller` - хранит контроллеры.
- `/internal/entity` - хранит сущности бизнес логики.
- `/internal/usecase` - содержит бизнес логику проекта.
- `/internal/repository` - используется для работы с данными.
- `/migrations` - хранит sql скрипты для миграции.
- `/pkg` - содержит пакеты для внутреннего использования.

Данная архитектура выбрана для упрощения возможного расширения приложения и добавления нового функционала.

## Технологии, использованные в проекте

В этом проекте использовались:

- `postgresql` - в качестве основной бд

- `gin` - в качестве веб фреймворка

- `swagger` - для разработки спецификации

- `go-pg` - для работы с Postgresql

- `testify` - для тестирования

- `docker`, `docker compose` - для контейнеризации

- `make` - для сокращения скриптов запуска

