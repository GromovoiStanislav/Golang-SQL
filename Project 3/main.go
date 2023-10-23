package main

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Создаем модель данных (структуру)
type MyModel struct {
   //gorm.Model
	ID       uint   `gorm:"primaryKey"`
	Name     string `gorm:"column:name"`
	DeletedAt gorm.DeletedAt // Добавьте это поле для мягкого удаления
}

func main() {
	// Настроим соединение с базой данных PostgreSQL
	dsn := "host=localhost user=postgres password=root dbname=golang port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Не удалось подключиться к базе данных")
	}

	//// Автомиграция - создание таблицы, если она не существует
	//db.AutoMigrate(&MyModel{})

	// // Создание записей
	// db.Create(&MyModel{Name: "Name1"})
	// db.Create(&MyModel{Name: "Name2"})
	// db.Create(&MyModel{Name: "Name3"})

	// // Мягкое удаление
	// db.Delete(&MyModel{}, 1) // Удалить запись с ID 1
	// db.Where("name = ?", "Name3").Delete(&MyModel{})

	fmt.Println("\nFind:\n")
	{
		var models []MyModel
		db.Unscoped().Order("id").Find(&models) // Найти все записи, ВКЛЮЧАЯ удаленные
		fmt.Println(models)
		for _, model := range models {
			fmt.Println(model)
			//fmt.Printf("ID: %d, Name: %s, DeletedAt: %s\n", model.ID, model.Name, model.DeletedAt)
		}
	}

	fmt.Println("")
	{
		var models []MyModel
		db.Order("id").Find(&models) // Найти все записи, БЕЗ удаленных
		fmt.Println(models)
		for _, model := range models {
			fmt.Println(model)
			//fmt.Printf("ID: %d, Name: %s, DeletedAt: %s\n", model.ID, model.Name, model.DeletedAt)
		}
	}

	fmt.Println("")
	{
		var models []MyModel
		db.Where("name = ?", "Name3").Find(&models)
		fmt.Println(models) // []
	}

	fmt.Println("")
	{
		var models []MyModel
		db.Where("name = ?", "Name2").Find(&models)
		fmt.Println(models) // [{2 Name2 {0001-01-01 00:00:00 +0000 UTC false}}]
	}

	fmt.Println("")
	{
		var models []MyModel
		db.Where("id = ?", 1).Find(&models)
		fmt.Println(models) // []
	}

	fmt.Println("")
	{
		var models []MyModel
		db.Find(&models,1)  // []
		fmt.Println(models)
	}


	fmt.Println("")
	{
		var models []MyModel
		db.Unscoped().Find(&models,1)  // []
		fmt.Println(models) // [{1 Name1 {2023-10-23 17:01:00.06926 +0600 +06 true}}]
	}

	fmt.Println("")
	{
		var model MyModel
		db.Unscoped().Find(&model,1)  // []
		fmt.Println(model) // {1 Name1 {2023-10-23 17:01:00.06926 +0600 +06 true}}
	}



	fmt.Println("\nFirst:\n")



	fmt.Println("")
	{
		var model MyModel
		if err := db.Where("name = ?", "Name3").First(&model).Error; err != nil {
			if err != nil {
				fmt.Println("Запись не найдена") // Запись не найдена
			}
		}
		fmt.Println(model) // {0  {0001-01-01 00:00:00 +0000 UTC false}}
	}

	fmt.Println("")
	{
		var model MyModel
		if err := db.Unscoped().Where("name = ?", "Name3").First(&model).Error; err != nil {
			if err != nil {
				fmt.Println("Запись не найдена") // 
			}
		}
		fmt.Println(model) // {3 Name3 {2023-10-23 17:01:00.072546 +0600 +06 true}}
	}


	fmt.Println("")
	{
		var model MyModel
		if err := db.Where("name = ?", "Name2").First(&model).Error; err != nil {
			if err != nil {
				fmt.Println("Запись не найдена") //
			}
		}
		fmt.Println(model) // {2 Name2 {0001-01-01 00:00:00 +0000 UTC false}}
	}

	fmt.Println("")
	{
		var models []MyModel
		if err := db.Where("name = ?", "Name2").First(&models).Error; err != nil {
			if err != nil {
				fmt.Println("Запись не найдена") //
			}
		}
		fmt.Println(models) // [{2 Name2 {0001-01-01 00:00:00 +0000 UTC false}}]
	}

	fmt.Println("")


	{
		var model MyModel
		result := db.First(&model,1);
		if  result.Error != nil {
			if err != nil {
				fmt.Println("Запись не найдена") //
			}
		}
		fmt.Println(model) // {0  {0001-01-01 00:00:00 +0000 UTC false}}
	}

	fmt.Println("")
	{
		var model MyModel
		if err := db.Unscoped().First(&model,1).Error; err != nil {
			if err != nil {
				fmt.Println("Запись не найдена") //
			}
		}
		fmt.Println(model) // {1 Name1 {2023-10-23 17:01:00.06926 +0600 +06 true}}
	}


	fmt.Println("")
	{
		var model MyModel
		if err := db.First(&model,2).Error; err != nil {
			if err != nil {
				fmt.Println("Запись не найдена") //
			}
		}
		fmt.Println(model) // {2 Name2 {0001-01-01 00:00:00 +0000 UTC false}}
	}



	fmt.Println("\nTake:\n")

	// fmt.Println("")
	// {
	// 	var model MyModel
	// 	if err := db.Where("name = ?", "Name3").Take(&model).Error; err != nil {
	// 		if err != nil {
	// 			fmt.Println("Запись не найдена") // Запись не найдена
	// 		}
	// 	}
	// 	fmt.Println(model) // {0  {0001-01-01 00:00:00 +0000 UTC false}}
	// }

	// fmt.Println("")
	// {
	// 	var model MyModel
	// 	if err := db.Unscoped().Where("name = ?", "Name3").Take(&model).Error; err != nil {
	// 		if err != nil {
	// 			fmt.Println("Запись не найдена") // 
	// 		}
	// 	}
	// 	fmt.Println(model) // {3 Name3 {2023-10-23 17:01:00.072546 +0600 +06 true}}
	// }


	fmt.Println("END")
}

