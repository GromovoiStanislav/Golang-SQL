package main

import (
	"time"
  "fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)


type User struct {
  gorm.Model
  Name string
  Age  uint8
  Birthday time.Time
}


func main() {

  dsn := "host=localhost user=postgres password=root dbname=golang port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
  if err != nil {
    panic("failed to connect database")
  }


  {
    user := User{Name: "Jinzhu", Age: 18, Birthday: time.Now()}
    result := db.Create(&user) // передаем указатель на данные в Create
    /*
    user.ID             // возвращает первичный ключ добавленной записи
    result.Error        // возвращает ошибку
    result.RowsAffected // возвращает количество вставленных записей
    */


    fmt.Println("ID", user.ID)
    fmt.Println("Error", result.Error)
    fmt.Println("RowsAffected", result.RowsAffected)
  }

  {
    users := []User{
      User{Name: "Jinzhu", Age: 18, Birthday: time.Now()},
      User{Name: "Jackson", Age: 19, Birthday: time.Now()},
    }
  
    result := db.Create(users) // передайте фрагмент, чтобы вставить несколько строк

    
    fmt.Println("Error", result.Error)
    fmt.Println("RowsAffected", result.RowsAffected)

    for _, user := range users {
      fmt.Println("ID", user.ID)
    }
  }
}