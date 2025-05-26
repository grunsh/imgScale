# Что это?
Это финальный проект курса Go в Otus. Задание [тут](https://github.com/OtusGolang/final_project/blob/master/03-image-previewer.md).

# Что с приложением?

main.go в каталоге server.

Конфигурацию можно передать через переменные окружения:
PORT=8081
STORAGE_TYPE=file (memory)

# Docker для тестирования
В каталоге docker
docker build -t nginx-with-images .
docker run -d --name nginx-container -p 8080:80 my-nginx-image

# Makefile

Я не осилил. У меня Windows и из-за болезней и нагрузки на работе я не успел. Понять и простить...
