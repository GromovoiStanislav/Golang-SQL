package main

import (
	"encoding/json"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
)

type User struct {
	ID       uint `gorm:"primaryKey"`
	Name     string
	Username string
	Email    string
	Phone    string
	Website  string
	Address  UserAddress `gorm:"foreignKey:UserID"` // one-to-one
	Company  UserCompany `gorm:"foreignKey:UserID"` // one-to-one
	Posts    []Post      `gorm:"foreignKey:UserID"` // one-to-many
}

type UserAddress struct {
	ID      uint `gorm:"primaryKey"`
	UserID  uint // Внешний ключ
	Street  string
	Suite   string
	City    string
	Zipcode string
	Lat     string
	Lng     string
}

type UserCompany struct {
	ID          uint `gorm:"primaryKey"`
	UserID      uint // Внешний ключ
	Name        string
	CatchPhrase string
	Bs          string
}

type Post struct {
	ID       uint `gorm:"primaryKey"`
	UserID   uint // Внешний ключ
	Title    string
	Body     string
	Comments []Comment // one-to-many
}

type Comment struct {
	ID     uint `gorm:"primaryKey"`
	PostID uint `gorm:"foreignKey:PostID"` // Внешний ключ
	Name   string
	Email  string
	Body   string
}

func autoMigrate(db *gorm.DB) {
	// Автомиграция - создание таблиц
	db.AutoMigrate(&User{})
	db.AutoMigrate(&UserAddress{})
	db.AutoMigrate(&UserCompany{})
	db.AutoMigrate(&Post{})
	db.AutoMigrate(&Comment{})
	fmt.Println("Таблицы успешно созданы.")
}

func seedUsers(db *gorm.DB) {
	// Выполним запрос к API и получим данные пользователей
	resp, err := http.Get("https://jsonplaceholder.typicode.com/users")
	if err != nil {
		panic("Не удалось выполнить запрос к API")
	}
	defer resp.Body.Close()

	var usersData []struct {
		ID       uint
		Name     string
		Username string
		Email    string
		Phone    string
		Website  string
		Address  struct {
			Street  string
			Suite   string
			City    string
			Zipcode string
			Geo     struct {
				Lat string
				Lng string
			}
		}
		Company struct {
			Name        string
			CatchPhrase string
			Bs          string
		}
	}

	err = json.NewDecoder(resp.Body).Decode(&usersData)
	if err != nil {
		panic("Не удалось декодировать JSON")
	}
	// Сохраняем данные пользователей в таблицы
	for _, userData := range usersData {
		user := User{
			Name:     userData.Name,
			Username: userData.Username,
			Email:    userData.Email,
			Phone:    userData.Phone,
			Website:  userData.Website,
		}
		db.Create(&user)

		address := UserAddress{
			UserID:  user.ID,
			Street:  userData.Address.Street,
			Suite:   userData.Address.Suite,
			City:    userData.Address.City,
			Zipcode: userData.Address.Zipcode,
			Lat:     userData.Address.Geo.Lat,
			Lng:     userData.Address.Geo.Lng,
		}
		db.Create(&address)

		company := UserCompany{
			UserID:      user.ID,
			Name:        userData.Company.Name,
			CatchPhrase: userData.Company.CatchPhrase,
			Bs:          userData.Company.Bs,
		}
		db.Create(&company)
	}

	fmt.Println("Данные успешно сохранены в базе данных.")
}

func seedPosts(db *gorm.DB) {
	// Выполним запрос к API и получим данные постов
	resp, err := http.Get("https://jsonplaceholder.typicode.com/posts")
	if err != nil {
		panic("Не удалось выполнить запрос к API")
	}
	defer resp.Body.Close()

	var postsData []struct {
		UserID uint
		ID     uint
		Title  string
		Body   string
	}

	err = json.NewDecoder(resp.Body).Decode(&postsData)
	if err != nil {
		panic("Не удалось декодировать JSON")
	}

	// Сохраняем данные постов в таблицу
	for _, postData := range postsData {
		post := Post{
			UserID: postData.UserID,
			Title:  postData.Title,
			Body:   postData.Body,
		}
		db.Create(&post)
	}

	fmt.Println("Данные постов успешно сохранены в базе данных.")
}

func seedComments(db *gorm.DB) {
	// Выполним запрос к API и получим данные комментариев
	resp, err := http.Get("https://jsonplaceholder.typicode.com/comments")
	if err != nil {
		panic("Не удалось выполнить запрос к API")
	}
	defer resp.Body.Close()

	var commentsData []struct {
		PostID uint
		ID     uint
		Name   string
		Email  string
		Body   string
	}

	err = json.NewDecoder(resp.Body).Decode(&commentsData)
	if err != nil {
		panic("Не удалось декодировать JSON")
	}

	// Сохраняем данные комментариев в таблицу
	for _, commentData := range commentsData {
		comment := Comment{
			PostID: commentData.PostID,
			Name:   commentData.Name,
			Email:  commentData.Email,
			Body:   commentData.Body,
		}
		db.Create(&comment)
	}

	fmt.Println("Данные комментариев успешно сохранены в базе данных.")
}

func usersALL(db *gorm.DB) {
	var users []User
	db.Preload("Address").Preload("Company").Find(&users)

	// fmt.Println(users)

	for _, user := range users {
		fmt.Printf("ID: %d\n", user.ID)
		fmt.Printf("Name: %s\n", user.Name)
		fmt.Printf("Username: %s\n", user.Username)
		fmt.Printf("Email: %s\n", user.Email)

		// Вывод информации о адресе
		fmt.Println("Address:")
		fmt.Printf("  Street: %s\n", user.Address.Street)
		fmt.Printf("  Suite: %s\n", user.Address.Suite)
		fmt.Printf("  City: %s\n", user.Address.City)
		fmt.Printf("  Zipcode: %s\n", user.Address.Zipcode)

		// Вывод информации о компании
		fmt.Println("Company:")
		fmt.Printf("  Name: %s\n", user.Company.Name)
		fmt.Printf("  CatchPhrase: %s\n", user.Company.CatchPhrase)
		fmt.Printf("  Bs: %s\n", user.Company.Bs)

		fmt.Println("-------------")
	}
}

func usersPart(db *gorm.DB, userID uint) {

	var users []struct {
		ID       uint
		Name     string
		Username string
		City     string
		Zipcode  string
		Company  string
	}

	db.Table("users").
		Select("users.id, users.name, users.username, user_addresses.city, user_addresses.zipcode, user_companies.name").
		Joins("LEFT JOIN user_addresses ON users.id = user_addresses.user_id").
		Joins("LEFT JOIN user_companies ON users.id = user_companies.user_id").
		Where("users.id = ?", userID).
		Scan(&users)

	for _, user := range users {
		fmt.Printf("ID: %d\n", user.ID)
		fmt.Printf("Name: %s\n", user.Name)
		fmt.Printf("Username: %s\n", user.Username)
		fmt.Printf("City: %s\n", user.City)
		fmt.Printf("Zipcode: %s\n", user.Zipcode)
		fmt.Printf("Company Name: %s\n", user.Company)
		fmt.Println("-------------")
	}

}

func main() {
	// Настроим соединение с базой данных PostgreSQL
	dsn := "host=localhost user=postgres password=root dbname=jsonplaceholder port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Не удалось подключиться к базе данных")
	}

	//// Автомиграция - создание таблиц
	// autoMigrate(db)

	//// загрузка данных
	//seedUsers(db)
	//seedPosts(db)
	//seedComments(db)

	//// Чтение данных
	usersALL(db)
	usersPart(db, 1)

	fmt.Println("END")
}
