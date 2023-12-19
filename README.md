# qon
QuickON backend


## Локальный запуск приложения

Для локального запуска приложения достаточно запустить `make docker-app-up` на Linux системе.
При этом в папку frontend будет скачан проект frontend'a из ветки develop, а из него будет собран новый образ frontend'a.

После запуска для доступа к приложению переходим на http://localhost:8080/, тут будет доступен как frontend, так и backend.
Чтобы получить swagger-спецификацию проекта переходим на страницу http://localhost:8081/swagger/index.html
