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
	Age   int    `gorm:"column:age"`
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
	user := User{Name: "Name", Email: "name@example.com", Age: 25}
	db.Create(&user)

	// Чтение пользователя по ID
	var readUser User
	db.First(&readUser, user.ID)
	fmt.Printf("ID: %d, Name: %s, Email: %s\n", readUser.ID, readUser.Name, readUser.Email)

	// Обновление данных пользователя
	db.Model(&readUser).Update("Name", "Новое имя")

	// Обновление по условию
	db.Model(&User{}).Where("id = ?", 1).Update("Name", "Новое имя")
	//или:
	userToUpdate := User{ID: 1}
	db.Model(&userToUpdate).Update("Name", "Новое имя")

	// Обновление нескольких записей по условию
	db.Model(&User{}).Where("1 = 1").Update("Age", 18)
	db.Model(&User{}).Where("age IS NULL OR age = 0").Update("Age", 18)

	// Удаление пользователя
	//db.Delete(&readUser)

	// Удаление по условию
	db.Where("id = ?", 5).Delete(&User{})

	var users []User
	// Поиск всех записей:
	db.Find(&users)
	fmt.Println(users)

	// Поиск нескольких записей с условиями:
	db.Where("age > ?", 18).Find(&users)
	fmt.Println(users)

	// Лимитированный поиск записей:
	db.Limit(2).Find(&users)
	fmt.Println(users)

	// Сортировка записей:
	db.Order("age desc, name asc").Find(&users)
	fmt.Println(users)

	// Выбор конкретных столбцов:
	type UserProjection struct {
		Name  string
		Email string
	}
	var result []UserProjection
	db.Model(User{}).Find(&result)
	fmt.Println(result)
	// или:
	var resultMap []map[string]interface{}
	db.Model(User{}).Select("name, email").Find(&resultMap)
	fmt.Println(resultMap)

	fmt.Println("END")
}
