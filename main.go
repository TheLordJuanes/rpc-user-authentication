package main

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"unicode"
)

var tpl *template.Template

type user struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Birthdate string `json:"birthdate"`
}

var users = []user{
	{Username: "seyerman", Password: "seyerman", FirstName: "Juan", LastName: "Reyes", Birthdate: "1995-04-01"},
	{Username: "favellaneda", Password: "favellaneda", FirstName: "Fabio", LastName: "Avellaneda", Birthdate: "1987-09-06"},
}

func main() {
	var err error
	tpl, err = template.ParseGlob("*.html")
	if err != nil {
		panic(err.Error())
	}
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/loginauth", loginAuthHandler)
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/registerauth", registerAuthHandler)
	http.ListenAndServe("localhost:8080", nil)
}

// loginHandler serves form for users to login with
func loginHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("*****loginHandler running*****")
	tpl.ExecuteTemplate(w, "login.html", nil)
}

// loginAuthHandler authenticates user login
func loginAuthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("*****loginAuthHandler running*****")
	r.ParseForm()
	username := r.FormValue("username")
	password := r.FormValue("password")
	fmt.Println("username:", username, "password:", password)
	// retrieve password from db to compare (hash) with user supplied password's hash
	signed, err := getUserByUsername(username)
	if signed.Password != password {
		err = errors.New("wrong password")
	}
	if err != nil {
		tpl.ExecuteTemplate(w, "login.html", "Username and/or password are wrong!")
		return
	}
	// returns nill on success
	if err == nil {
		fmt.Fprint(w, "You have successfully logged in :)")
		return
	}
	fmt.Println(err)
	tpl.ExecuteTemplate(w, "login.html", "Check username and password!")
}

func getUserByUsername(username string) (user, error) {
	//username := c.Param("username")

	// Loop through the list of albums, looking for
	// an album whose ID value matches the parameter.
	for _, a := range users {
		if a.Username == username {
			return a, nil
		}
	}
	var null user
	return null, errors.New("user not found")
}

// registerHandler serves form for registering new users
func registerHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("*****registerHandler running*****")
	tpl.ExecuteTemplate(w, "register.html", nil)
}

// registerAuthHandler creates new user in database
func registerAuthHandler(w http.ResponseWriter, r *http.Request) {
	/*
		1. check username criteria
		2. check password criteria
		3. check if username is already exists in database
		4. create bcrypt hash from password
		5. insert username and password hash in database
	*/
	fmt.Println("*****registerAuthHandler running*****")
	r.ParseForm()
	username := r.FormValue("username")
	// check username for only alphaNumeric characters
	var nameAlphaNumeric = true
	for _, char := range username {
		// func IsLetter(r rune) bool, func IsNumber(r rune) bool
		// if !unicode.IsLetter(char) && !unicode.IsNumber(char) {
		if !unicode.IsLetter(char) && !unicode.IsNumber(char) {
			nameAlphaNumeric = false
		}
	}
	// check username pswdLength
	var nameLength bool
	if 5 <= len(username) && len(username) <= 50 {
		nameLength = true
	}
	// check password criteria
	password := r.FormValue("password")
	fmt.Println("password:", password, "\npswdLength:", len(password))
	// variables that must pass for password creation criteria
	var pswdLowercase, pswdUppercase, pswdNumber, pswdSpecial, pswdLength, pswdNoSpaces bool
	pswdNoSpaces = true
	for _, char := range password {
		switch {
		// func IsLower(r rune) bool
		case unicode.IsLower(char):
			pswdLowercase = true
		// func IsUpper(r rune) bool
		case unicode.IsUpper(char):
			pswdUppercase = true
		// func IsNumber(r rune) bool
		case unicode.IsNumber(char):
			pswdNumber = true
		// func IsPunct(r rune) bool, func IsSymbol(r rune) bool
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			pswdSpecial = true
		// func IsSpace(r rune) bool, type rune = int32
		case unicode.IsSpace(int32(char)):
			pswdNoSpaces = false
		}
	}
	if 11 < len(password) && len(password) < 60 {
		pswdLength = true
	}
	fmt.Println("pswdLowercase:", pswdLowercase, "\npswdUppercase:", pswdUppercase, "\npswdNumber:", pswdNumber, "\npswdSpecial:", pswdSpecial, "\npswdLength:", pswdLength, "\npswdNoSpaces:", pswdNoSpaces, "\nnameAlphaNumeric:", nameAlphaNumeric, "\nnameLength:", nameLength)
	if !pswdLowercase || !pswdUppercase || !pswdNumber || !pswdSpecial || !pswdLength || !pswdNoSpaces || !nameAlphaNumeric || !nameLength {
		tpl.ExecuteTemplate(w, "register.html", "please check username and password criteria")
		return
	}
	// check if username already exists for availability
	_, err := getUserByUsername(username)
	if err == nil {
		fmt.Println("username already exists, error: ", err)
		tpl.ExecuteTemplate(w, "register.html", "Username already taken!")
		return
	}
	firstname := r.FormValue("firstname")
	lastname := r.FormValue("lastname")
	birthdate := r.FormValue("birthdate")
	password2 := r.FormValue("password2")
	if password != password2 {
		fmt.Println("Passwords don't match")
		tpl.ExecuteTemplate(w, "register.html", "Passwords don't match!")
		return
	}
	newUser := user{
		Username:  username,
		Password:  password,
		FirstName: firstname,
		LastName:  lastname,
		Birthdate: birthdate,
	}
	users = append(users, newUser)

	fmt.Fprint(w, "Congrats, your account has been successfully created!")
}
