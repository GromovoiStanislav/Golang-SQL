package main

import (
  "time"
  
  "gorm.io/gorm"
  "gorm.io/driver/postgres"
)

type Author struct {
  Name  string
  Email string
}


type Blog1 struct {
  ID      int
  Author  Author `gorm:"embedded"`
  Upvotes int32
}
/* эквивалентно
type Blog struct {
  ID    int64
    Name  string
    Email string
  Upvotes  int32
}
*/




type Blog2 struct {
ID      int
Author  Author `gorm:"embedded;embeddedPrefix:author_"`
Upvotes int32
}
/* эквивалентно
type Blog struct {
  ID          int64
  AuthorName  string
  AuthorEmail string
  Upvotes     int32
}
*/



type User struct {
  gorm.Model
  Name string
  Age  uint8
  Birthday time.Time
}
/* эквивалентно
type User struct {
  ID        uint           `gorm:"primaryKey"`
  CreatedAt time.Time
  UpdatedAt time.Time
  DeletedAt gorm.DeletedAt `gorm:"index"`
  Name string
  Age  uint8
  Birthday time.Time
}
*/




func main() {

  dsn := "host=localhost user=postgres password=root dbname=golang port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
  if err != nil {
    panic("failed to connect database")
  }

  // Миграция схем
  db.AutoMigrate(&Blog1{})
  db.AutoMigrate(&Blog2{})
  db.AutoMigrate(&User{})


}