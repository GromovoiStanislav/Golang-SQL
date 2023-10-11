package main

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Создаем модель данных (структуру)
type User struct {
	ID    uint   `gorm:"primaryKey"`
	Name  string `gorm:"column:name"`
	Email string `gorm:"column:email"`
}

func main() {
	// Настроим соединение с базой данных PostgreSQL
	dsn := "host=localhost user=postgres password=root dbname=golang port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Не удалось подключиться к базе данных")
	}

	// Автомиграция - создание таблицы, если она не существует
	db.AutoMigrate(&User{})

	// Создание пользователя
	user := User{Name: "Имя пользователя", Email: "user@example.com"}
	db.Create(&user)

	// Чтение пользователя по ID
	var readUser User
	db.First(&readUser, user.ID)
	fmt.Printf("ID: %d, Name: %s, Email: %s\n", readUser.ID, readUser.Name, readUser.Email)

	// Обновление данных пользователя
	db.Model(&readUser).Update("Name", "Новое имя")

	// Удаление пользователя
	db.Delete(&readUser)
}
