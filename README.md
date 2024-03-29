### Что это за проект
Проект представляет собой прототип консольного мессенджера, написанного на языке Go. Включает в себя клиента и сервера, 
поддерживающего n содеинений.

Клиенты и сервер должны быть на одной машине. 

### На данный момент возможно:
- зарегистрировать нового пользователя, используя придуманное имя и имейл
- авторизировать нового пользователя, используя только имя
- отправить сообщения от одного пользователя к другому

### Что планируется добавить:
- тесты + пайплайн
- корректное завершение горутин
- rabbitMQ, чтоб хранить сообщения, которые не удалось доставить по какой-то причине и пытаться их отправить повторно
- дату отправки сообщения
- поддержку групповых чатов
- проверку имейла и имени на допустимый формат
- поддержку криптографии
- нормальные креды для регистрации пользователя и их соответтсвующее хранение (хэш от пароля с солью и т.д.)
- поддержку юзер-интерфейса

### Как запустить проект:
#### Запуск инфрастуктуры и выполнение миграций:
```shell
 make all
```
#### Запуск сервера:
Из корня проекта выполните
```shell
 go run ./server/main.go
```

#### Запуск клиента:
Из корня проекта выполните
```shell
 go run ./client/main.go
```

#### Как работать с клиентом:
В ответ на сообщение 
```text
Enter your host for yor server:
```
введите хост и порт, на котором хотите запустить клиента. Например, localhost:8787. Внимание, localhost:8080 зарезервирован для сервера.
Затем пройдите процедуру авторизации. Введите свое имя. Если вы зарегистрированный пользователь и если имя было введено корректно, то далее вам будет предложено ввести имя вашего собеседника. 
Если он является зарегистрированным пользователем мессенджера, то далее вы можете вводить сообщение. 

Если вы не являетесь зарегистрированным пользователем, то вам будет предложено зарегистроваться. Введите ваше имя и предполагаемый имейл. 
Если процедура регистрации прошла успешно, далее вы действуете как после успешной авторизации. 
