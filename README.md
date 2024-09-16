# gRPC_Interaction
Добрый день
Прошу прощения, на работе навалилось и из головы совсем вылетело про задачу...
Предлагаю вам сделать систему, состоящую из нескольких модулей. Модули имеют gRPC сервисы для внешних обращений и для внутреннего взаимодействия. Сервисы для внешних обращений аутентифицируют пользователя через JWT. Внутренние сервисы пусть используют другой механизм (аутентификация модуль к модулю), какой именно - на ваше усмотрение.

Модули:

Аутентификация:
- внешний сервис принимает логин и пароль и отдает JWT содержащий id пользователя.

Авторизация:
- внутренний сервис принимает id пользователя и отдает список ролей

Сведения о пользователе:
- внешний сервис отдает информацию о текущем (аутентифицированном пользователе), как то, имя, email, роли и дополнительные поля на ваше усмотрение
- внутренний сервис отдает ту же информацию для произвольного пользователя по его id

Уведомления:
- внутренний сервис списка уведомлений, возвращает список уведомлений ( id, текст и перечень id пользователей, кому предназначено) для пользователей, возвращаются еще не обработанные сообщения
- внутренний сервис отметки уведомления обработанным - принимает список id уведомлений и помечает их как обработанные
- внешний сервис отметки уведомления полученным - принимает id уведомления и отмечает, что текущий пользователь получил указанное уведомление

Отправка уведомлений:
Сервисов не имеет. Должен периодически обращаться к модулю уведомлений, получать новые, отправлять их в firebase и отмечать обработанными

Модули должны быть реализованы как отдельные приложения, в отдельных докер контейнерах, каждый со своей бд (если она требуется)

## Архитектура системы
1. Auth Service (Аутентификация)
   Задача: Принимает логин и пароль, возвращает JWT с id пользователя.
   Хранилище: PostgreSQL для хранения данных о пользователях.
   JWT: Используется для всех внешних сервисов.
2. Role Service (Авторизация)
   Задача: Принимает id пользователя и возвращает список его ролей.
   Хранилище: PostgreSQL для хранения информации о ролях и привязке ролей к пользователям.
   Внешняя аутентификация: По JWT, проверяем id пользователя.
   Внутренняя аутентификация: Используем, например, mTLS (взаимная TLS-аутентификация) для безопасности между сервисами.
3. User Service (Сведения о пользователе)
   Задача:
   Внешний сервис: Возвращает информацию о текущем (аутентифицированном) пользователе (имя, email, роли).
   Внутренний сервис: Предоставляет те же данные для произвольного пользователя по id.
   Хранилище: PostgreSQL для хранения профилей пользователей.
   Внешняя аутентификация: JWT.
   Внутренняя аутентификация: mTLS.
4. Notification Service (Уведомления)
   Задача:
   Внутренний сервис списка уведомлений: Возвращает список уведомлений для пользователей, которые еще не обработаны.
   Внутренний сервис отметки уведомления как обработанного.
   Внешний сервис отметки уведомления полученным: Для текущего пользователя.
   Хранилище: PostgreSQL для хранения уведомлений и их статусов.
   Внешняя аутентификация: JWT.
   Внутренняя аутентификация: mTLS.
5. Notification Sender (Отправка уведомлений)
   Задача: Не имеет публичного API. Периодически обращается к Notification Service для получения необработанных уведомлений, отправляет их через Firebase, отмечает уведомления как обработанные.
   Взаимодействие: Через gRPC с Notification Service.
   Периодичность: Можно реализовать через cron или системный таймер.

##   Общие детали:
   Аутентификация модуль-к-модулю
   Для внутреннего взаимодействия между модулями можно использовать mTLS (Mutual TLS). Это надежный способ аутентификации модуль к модулю, где каждый сервис имеет свой TLS-сертификат, который проверяется другим модулем. Таким образом, мы можем быть уверены, что запросы поступают от авторизованных модулей.

JWT для внешней аутентификации
Для внешнего доступа все запросы будут требовать JWT, который будет проверяться на стороне каждого микросервиса через промежуточный слой (middleware). JWT будет содержать информацию о пользователе, например, его id и роли.

Взаимодействие через gRPC
Все модули будут взаимодействовать между собой через gRPC. Это обеспечит высокую производительность и эффективное взаимодействие между сервисами.

Базы данных
Каждый модуль будет иметь свою отдельную базу данных. Это позволяет изолировать данные и масштабировать каждый модуль независимо от других. Мы будем использовать PostgreSQL для хранения данных.

Docker-композиция
Все модули будут разворачиваться в отдельных Docker-контейнерах с их собственными базами данных и gRPC-сервисами. Взаимодействие между контейнерами будет происходить через сеть Docker.

Пример схемы взаимодействия
Auth Service: Пользователь отправляет логин и пароль → сервис возвращает JWT.
User Service: Пользователь с JWT запрашивает информацию о себе (имя, email, роли и т.д.).
Role Service: Внутренние модули запрашивают список ролей пользователя по его id.
Служба уведомлений:
Внешний пользователь помечает уведомление как прочитанное.
Внутренний сервис отправки уведомлений помечает их как обработанные.
Notification Sender: Периодически запрашивает необработанные уведомления и отправляет их через Firebase, после чего отмечает как обработанные.
Пр
1. Аутентификация и авторизация пользователя
   Пользователь отправляет запрос на логин в Auth Service.
   Auth Service проверяет логин и пароль, генерирует JWT с id пользователя.
   Пользователь получает JWT и использует его для дальнейших запросов.
2. Запрос ролей пользователя
   User Service или другие модули отправляют gRPC-запрос в Role Service с id пользователя.
   Role Service возвращает список ролей для пользователя.
3. Уведомления
   Notification Service возвращает список уведомлений для пользователя.
   Пользователь помечает уведомление как прочитанное, и Notification Service обновляет статус.