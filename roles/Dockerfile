# roles/Dockerfile
FROM golang:1.22-alpine

WORKDIR /app

# Копируем зависимости и загружаем их
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Копируем остальные файлы проекта
COPY . .

# Переходим в директорию с main.go и собираем приложение
WORKDIR /app/cmd
RUN go build -o /roles_service

# Запускаем собранное приложение
CMD ["/roles_service"]