package models

import (
	"../forms"
	db "../database"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"time"
	"log"
)

type User struct {
	Id        int    `database:"id, primarykey, autoincrement" json:"id"`
	Email     string `database:"email" json:"email"`
	Password  string `database:"password" json:"-"`
	Name      string `database:"name" json:"name"`
	UpdatedAt int64  `database:"updated_at" json:"updated_at"`
	CreatedAt int64  `database:"created_at" json:"created_at"`
}

type UserModel struct{}

func (m UserModel) Signin(form forms.SigninForm) (user User, err error) {
	err = db.SqlDB.QueryRow("SELECT id, email, password, name, updated_at, created_at FROM user WHERE email=LOWER(?) LIMIT 1", form.Email).Scan(user)

	if err != nil {
		return user, err
	}

	bytePassword := []byte(form.Password)
	byteHashedPassword := []byte(user.Password)

	err = bcrypt.CompareHashAndPassword(byteHashedPassword, bytePassword)

	if err != nil {
		return user, errors.New("Invalid password")
	}

	return user, nil
}

func (m UserModel) Signup(form forms.SignupForm) (user User, err error) {
	var checkUser = 0
	err = db.SqlDB.QueryRow("SELECT count(id) FROM user WHERE email=LOWER(?) LIMIT 1", form.Email).Scan(&checkUser)

	if err != nil {
		return user, err
	}

	if checkUser > 0 {
		return user, errors.New("User exists")
	}

	bytePassword := []byte(form.Password)
	hashedPassword, err := bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	res, err := db.SqlDB.Exec("INSERT INTO user(email, password, name, updated_at, created_at) VALUES(?, ?, ?, ?, ?)", form.Email, string(hashedPassword), form.Name, time.Now(), time.Now())

	if res != nil && err == nil {
		err = db.SqlDB.QueryRow("SELECT * FROM user WHERE email=LOWER(?) LIMIT 1", form.Email).Scan(&user)
		log.Println(err)
		if err == nil {
			return user, nil
		}
	}

	return user, errors.New("Not registered")
}

func (m UserModel) One(userId int64) (user User, err error) {
	err = db.SqlDB.QueryRow("SELECT id, email, name FROM user WHERE id=?", userId).Scan(user)
	return user, err
}
