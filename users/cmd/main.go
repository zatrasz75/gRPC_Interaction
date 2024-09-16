package main

import (
	"fmt"
	"os"
	"path/filepath"
	"zatrasz75/gRPC_Interaction/users/configs"
	"zatrasz75/gRPC_Interaction/users/internal/app"
	"zatrasz75/gRPC_Interaction/users/pkg/logger"
)

func main() {
	l := logger.NewLogger()

	// Получаем текущий рабочий каталог
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Ошибка при получении текущего рабочего каталога:", err)
		return
	}
	fmt.Println(cwd)
	// Построение абсолютного пути к файлу configs.yml
	configPath := filepath.Join(cwd, "configs", "configs.yml")
	fmt.Println(configPath)
	if _, err = os.Stat(configPath); os.IsNotExist(err) {
		configPath = filepath.Join(cwd, "..", "configs", "configs.yml")
	}
	fmt.Println(configPath)

	// Configuration
	cfg, err := configs.NewConfig(configPath)
	if err != nil {
		l.Fatal("ошибка при разборе конфигурационного файла", err)
	}

	fmt.Println("ConnStr --- ", cfg.PostgreSQL.ConnStr)

	app.Run(cfg, l)
}
