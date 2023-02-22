package models

import (
	"errors"

	"github.com/Massad/gin-boilerplate/db"
	"github.com/Massad/gin-boilerplate/forms"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

//User ...
type User struct {
	ID        int64  `db:"id, primarykey, autoincrement" json:"id"`
	UserId    string `db:"user_id" json:"user_id"`
	Email     string `db:"email" json:"email"`
	Password  string `db:"password" json:"-"`
	Name      string `db:"name" json:"name"`
	UpdatedAt int64  `db:"updated_at" json:"-"`
	CreatedAt int64  `db:"created_at" json:"-"`
}

//UserModel ...
type UserModel struct{}

var authModel = new(AuthModel)

//Login ...
func (m UserModel) Login(form forms.LoginForm) (user User, err error) {

	err = db.GetDB().SelectOne(&user, "SELECT id, user_id, email, password, name, updated_at, created_at FROM public.user WHERE email=LOWER($1) LIMIT 1", form.Email)

	if err != nil {
		return user, err
	}

	//Compare the password form and database if match
	bytePassword := []byte(form.Password)
	byteHashedPassword := []byte(user.Password)

	err = bcrypt.CompareHashAndPassword(byteHashedPassword, bytePassword)

	if err != nil {
		return user, err
	}

	return user, nil
}

//Register ...
func (m UserModel) Register(form forms.RegisterForm) (user User, err error) {
	getDb := db.GetDB()

	//Check if the user exists in database
	checkUser, err := getDb.SelectInt("SELECT count(id) FROM public.user WHERE email=LOWER($1) LIMIT 1", form.Email)
	if err != nil {
		return user, errors.New("something went wrong, please try again later")
	}

	if checkUser > 0 {
		return user, errors.New("email already exists")
	}

	bytePassword := []byte(form.Password)
	hashedPassword, err := bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)
	if err != nil {
		return user, errors.New("something went wrong, please try again later")
	}

	userId := uuid.New().String()

	//Create the user and return back the user ID
	err = getDb.QueryRow("INSERT INTO public.user(email, password, name, user_id) VALUES($1, $2, $3, $4) RETURNING id", form.Email, string(hashedPassword), form.Name, userId).Scan(&user.ID)
	if err != nil {
		return user, errors.New("something went wrong, please try again later")
	}

	user.Name = form.Name
	user.Email = form.Email
	user.UserId = userId

	return user, err
}

//One ...
func (m UserModel) One(userID string) (user User, err error) {
	err = db.GetDB().SelectOne(&user, "SELECT id, email, name FROM public.user WHERE user_id=$1 LIMIT 1", userID)
	return user, err
}
