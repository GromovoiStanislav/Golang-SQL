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


	// Пример транзакции:
	err = ExampleTransaction(db)
	if err != nil {
		return
	}


	var users []User
	// Поиск всех записей:
	db.Find(&users)
	fmt.Println(users)

	// Поиск нескольких записей с условиями:
	db.Where("age > ?", 18).Find(&users)
	fmt.Println(users)

	// Поиск нескольких записей по списку:
	db.Where("id IN (?)", []uint{1, 4}).Find(&users)
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

	// общеее количество записей в таблице users
	var count int64
	db.Model(&User{}).Count(&count)
	fmt.Println("Всего человек:", count)

	// Средний возраст
	var averageAge float64
	db.Model(&User{}).Select("AVG(age) as average_age").Scan(&averageAge)
	fmt.Println("Средний возраст:", averageAge)

	// Минимальный возраст
	var minAge float64
	db.Model(&User{}).Select("MIN(age) as min_age").Scan(&minAge)
	fmt.Println("Самый молодой:", minAge)

	//Пагинация
	db.Order("id").Limit(2).Offset(2).Find(&users)
	fmt.Println(users)

	fmt.Println("END")
}

func ExampleTransaction(db *gorm.DB) error {
	// Создаем новую транзакцию
	tx := db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Примеры SQL-запросов, которые будут выполнены внутри транзакции
	user1 := User{Name: "User1", Email: "user1@example.com", Age: 19}
	user2 := User{Name: "User2", Email: "user2@example.com", Age: 21}

	// Вставляем записи в таблицу "users" внутри транзакции
	if err := tx.Create(&user1).Error; err != nil {
		// Если произошла ошибка, откатываем транзакцию
		tx.Rollback()
		return err
	}

	if err := tx.Create(&user2).Error; err != nil {
		// Если произошла ошибка, откатываем транзакцию
		tx.Rollback()
		return err
	}

	// Если все запросы прошли успешно, фиксируем транзакцию
	tx.Commit()

	return nil
}
