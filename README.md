# WB_orders_microservice

Сервис получает заказы из топика kafka, сохраняет в PostgreSQL, с доступом к ним через REST API. Так же реализовано кэширование, с заполнением при запуске.

Стек:
 - Go
 - Kafka
 - PostgreSQL
 - Kafka UI (порт 8088)

Запуск:
  - Создать файл .env в корне проекта, скопировать данные из файла example.env и поодставить свои данные в места с пометкой "YOUR_"
  - Запуск через доскер командой docker-compose up -d

Использование:
  - Kafka UI доступна по адресу localhost:8088
  - Перейти во вкладку topics -> выбрать "orders" -> Produce Message -> вставить в поле "value" json из файла example_order.json
  - Доступ к заказам можно осуществлять через localhost:3000 по order_uid заказа
