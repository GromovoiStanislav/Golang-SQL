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

func usersByIDList(db *gorm.DB, idList []uint) []User {
	var users []User
	db.Where("id IN (?)", idList).Find(&users)

	for _, user := range users {
		fmt.Printf("ID: %d\n", user.ID)
		fmt.Printf("Name: %s\n", user.Name)
		fmt.Printf("Username: %s\n", user.Username)
		fmt.Printf("Email: %s\n", user.Email)
		fmt.Println("-------------")
	}

	return users
}

func exampleTransaction(db *gorm.DB) error {
	// Создаем новую транзакцию
	tx := db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Примеры SQL-запросов, которые будут выполнены внутри транзакции
	user1 := User{Name: "User1", Email: "user1@example.com"}
	user2 := User{Name: "User2", Email: "user2@example.com"}

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

func GetUsersWithNoPosts(db *gorm.DB) []User {
	var users []User
	subquery := db.Model(&Post{}).Select("DISTINCT user_id")

	db.Not("id IN (?)", subquery).Find(&users)

	for _, user := range users {
		fmt.Printf("ID: %d\n", user.ID)
		fmt.Printf("Name: %s\n", user.Name)
		fmt.Printf("Username: %s\n", user.Username)
		fmt.Printf("Email: %s\n", user.Email)
		fmt.Println("-------------")
	}

	return users
}

type UserWithoutPosts struct {
	ID    uint
	Name  string
	Email string
}

func FindUsersWithoutPosts(db *gorm.DB) []UserWithoutPosts {
	var result []UserWithoutPosts

	subquery := db.Model(&Post{}).Select("DISTINCT user_id")
	db.Model(&User{}).
		Select("id, name, email").
		Where("id NOT IN (?)", subquery).
		Find(&result)

	for _, user := range result {
		fmt.Printf("ID: %d\n", user.ID)
		fmt.Printf("Name: %s\n", user.Name)
		fmt.Printf("Email: %s\n", user.Email)
		fmt.Println("-------------")
	}

	return result
}

type PostCountByUser struct {
	UserID    uint `gorm:"column:user_id"`
	PostCount int
}

func GetPostCountByUser(db *gorm.DB) []PostCountByUser {
	var result []PostCountByUser

	db.Model(&Post{}).
		Select("user_id, COUNT(*) as post_count").
		Group("user_id").
		Scan(&result)

	for _, user := range result {
		fmt.Printf("user_id: %d\n", user.UserID)
		fmt.Printf("post_count: %d\n", user.PostCount)
		fmt.Println("-------------")
	}

	return result
}

type UserDataWithPostCount struct {
	UserID    uint
	Name      string
	PostCount int
}

func GetUserDataWithPostCount(db *gorm.DB) []UserDataWithPostCount {
	var result []UserDataWithPostCount

	//db.Model(&Post{}).
	//	Select("user_id, users.name, COUNT(*) as post_count").
	//	Joins("JOIN users ON users.id = user_id").
	//	Group("user_id, users.name").
	//	Scan(&result)

	db.Model(&User{}).
		Select("users.id as user_id, users.name, COALESCE(COUNT(posts.id), 0) as post_count").
		Joins("LEFT JOIN posts ON users.id = posts.user_id").
		Group("users.id").
		Scan(&result)

	for _, user := range result {
		fmt.Printf("user_id: %d\n", user.UserID)
		fmt.Printf("user.name: %s\n", user.Name)
		fmt.Printf("post_count: %d\n", user.PostCount)
		fmt.Println("-------------")
	}

	return result
}

func GetUsersWithLimitAndOffset(db *gorm.DB, limit, offset int) []User {
	var users []User
	db.Order("name").Limit(limit).Offset(offset).Find(&users)

	for _, user := range users {
		fmt.Printf("user.id: %d\n", user.ID)
		fmt.Printf("user.name: %s\n", user.Name)
		fmt.Println("-------------")
	}
	return users
}

func GetCommentsWithLimitAndOffset(db *gorm.DB, limit, offset int) []Comment {
	var comments []Comment

	//db.Order("id").Limit(limit).Offset(offset).Find(&comments)

	//or:
	db.Model(&Comment{}).
		//Select("*").
		Order("id").
		Limit(limit).
		Offset(offset).
		Find(&comments)

	for _, comment := range comments {
		fmt.Printf("id: %d\n", comment.ID)
		fmt.Printf("name: %s\n", comment.Name)
		fmt.Printf("body: %s\n", comment.Body)

		fmt.Println("-------------")
	}
	return comments
}

type UserCommentCount struct {
	UserID       uint `gorm:"column:user_id"`
	Name         string
	CommentCount int
}

func GetUserCommentCount(db *gorm.DB) []UserCommentCount {
	var result []UserCommentCount

	db.Model(&Comment{}).
		Select("users.id as user_id, users.name as name, COUNT(*) as comment_count").
		Joins("LEFT JOIN posts ON comments.post_id = posts.id").
		Joins("LEFT JOIN users ON posts.user_id = users.id").
		Group("users.id").
		Scan(&result)

	for _, user := range result {
		fmt.Printf("user.id: %d\n", user.UserID)
		fmt.Printf("user.name: %s\n", user.Name)
		fmt.Printf("comment_count: %d\n", user.CommentCount)
		fmt.Println("-------------")
	}

	return result
}

type UserCommentPostData struct {
	UserID      uint `gorm:"column:user_id"`
	UserName    string
	PostID      uint `gorm:"column:post_id"`
	PostTitle   string
	CommentID   uint `gorm:"column:comment_id"`
	CommentBody string
}

func GetUserCommentPostData(db *gorm.DB) []UserCommentPostData {
	var result []UserCommentPostData

	db.Model(&Comment{}).
		Select("users.id as user_id, users.name as user_name, posts.id as post_id, posts.title as post_title, comments.id as comment_id, comments.body as comment_body").
		Joins("LEFT JOIN posts ON comments.post_id = posts.id").
		Joins("LEFT JOIN users ON posts.user_id = users.id").
		Scan(&result)

	for _, data := range result {
		fmt.Printf("user.id: %d\n", data.UserID)
		fmt.Printf("user.name: %s\n", data.UserName)
		fmt.Printf("post.id: %d\n", data.PostID)
		fmt.Printf("post.title: %s\n", data.PostTitle)
		fmt.Printf("comment.id: %d\n", data.CommentID)
		fmt.Printf("comment.body: %s\n", data.CommentBody)
		fmt.Println("-------------")
	}

	return result
}

type UserPost struct {
	UserID    uint
	UserName  string
	PostID    uint
	PostTitle string
}

func FindTop3PostsPerUser(db *gorm.DB) []UserPost {
	var result []UserPost

	// SQL-запрос для выбора трех первых постов каждого пользователя
	query := `
        SELECT u.id as user_id, u.name as user_name, p.id as post_id, p.title as post_title
        FROM users u
        INNER JOIN (
            SELECT user_id, id, title, ROW_NUMBER() OVER (PARTITION BY user_id ORDER BY id) as row_num
            FROM posts
        ) p ON u.id = p.user_id AND p.row_num <= 3
    `

	db.Raw(query).Scan(&result)

	for _, data := range result {
		fmt.Printf("user.id: %d\n", data.UserID)
		fmt.Printf("user.name: %s\n", data.UserName)
		fmt.Printf("post.id: %d\n", data.PostID)
		fmt.Printf("post.title: %s\n", data.PostTitle)
		fmt.Println("-------------")
	}

	return result
}

type UserCommentMatch struct {
	UserID       uint
	UserName     string
	UserEmail    string
	CommentID    uint
	CommentBody  string
	CommentEmail string
}

func FindMatchingEmails(db *gorm.DB) []UserCommentMatch {
	var result []UserCommentMatch

	db.Model(&User{}).
		Select("users.id as user_id, users.name as user_name, users.email as user_email, comments.id as comment_id, comments.body as comment_body, comments.email as comment_email").
		Joins("LEFT JOIN comments ON users.email = comments.email"). //INNER
		Scan(&result)

	for _, data := range result {
		fmt.Printf("user.id: %d\n", data.UserID)
		fmt.Printf("user.name: %s\n", data.UserName)
		fmt.Printf("post.email: %d\n", data.UserEmail)
		fmt.Printf("comment.id: %d\n", data.CommentID)
		fmt.Printf("comment.email: %s\n", data.CommentEmail)
		//fmt.Printf("comment.body: %s\n", data.CommentBody)
		fmt.Println("-------------")
	}

	return result
}

func main() {
	// Настроим соединение с базой данных PostgreSQL
	dsn := "host=localhost user=postgres password=root dbname=jsonplaceholder port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Не удалось подключиться к базе данных")
	}

	//// Автомиграция - создание таблиц
	//autoMigrate(db)
	//
	//// загрузка данных
	//seedUsers(db)
	//seedPosts(db)
	//seedComments(db)

	//// Транзакция
	//err = exampleTransaction(db)
	//if err != nil {
	//	return
	//}

	//// Чтение данных
	//usersALL(db)
	//usersPart(db, 1)
	//idList := []uint{1, 3, 5}
	//usersByIDList(db, idList)
	//GetUsersWithNoPosts(db)   //выбрать пользователей, у которых нет постов,
	//FindUsersWithoutPosts(db) //выбрать пользователей, у которых нет постов,
	//GetPostCountByUser(db)
	//GetUserDataWithPostCount(db)
	//GetUsersWithLimitAndOffset(db, 2, 2)
	//GetCommentsWithLimitAndOffset(db, 10, 20)
	FindTop3PostsPerUser(db)
	//GetUserCommentCount(db)
	//GetUserCommentPostData(db)
	//FindMatchingEmails(db)

	fmt.Println("END")
}
