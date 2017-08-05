package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"

	"net/url"

	"golang.org/x/crypto/bcrypt"
)

/*
User defines a user of the CommComm application
*/
type User struct {
	ID       int       `json:"id"`
	Email    string    `json:"email"`
	Password string    `json:"-"`
	Date     time.Time `json:"-"`
	Active   int       `json:"-"`
}

/*
GetSpecificUserByID gets a user by their ID
*/
func GetSpecificUserByID(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	s := v["userId"]
	id, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	u, err := GetUserByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	u.Password = ""
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewEncoder(w).Encode(u); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

/*
GetSpecificUserByEmail is the handler function to get a user by their email
*/
func GetSpecificUserByEmail(w http.ResponseWriter, r *http.Request) {
	vals := r.URL.Query()
	emails := vals["email"]
	var email string
	if len(emails) > 0 {
		email = emails[0]
	} else {
		http.Error(w, "Please provide email", http.StatusBadRequest)
		return
	}

	decode, err := url.QueryUnescape(email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	u, err := GetUserByEmail(decode)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	u.Password = ""
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewEncoder(w).Encode(u); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

/*
UserCreate creates a user in the CommComm ecosystem
*/
func UserCreate(w http.ResponseWriter, r *http.Request) {
	var user User
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := r.Body.Close(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := json.Unmarshal(body, &user); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user.Password = string(hash)

	created, err := InsertUser(user.Email, user.Password, time.Now())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	created.Password = ""
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewEncoder(w).Encode(created); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

/*
Login issues a session token for the user passed in.
*/
func Login(w http.ResponseWriter, r *http.Request) {
	vals := r.URL.Query()
	emails := vals["email"]
	passwords := vals["password"]
	var email string
	var password string
	if len(emails) > 0 && len(passwords) > 0 {
		email = emails[0]
		password = passwords[0]
	} else {
		http.Error(w, "Please provide anemail and password", http.StatusBadRequest)
		return
	}

	e, err := url.QueryUnescape(email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	p, err := url.QueryUnescape(password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user, err := GetUserByEmail(e)
	if err != nil || user != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	p = string(hash)

	if p != user.Password {
		http.Error(w, "Supplied username and/or password incorrect", http.StatusForbidden)
		return
	}

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["userId"] = user.ID
	claims["Email"] = user.Email
	claims["exp"] = time.Now().Add(time.Hour * 24 * 120)

	tokenString, _ := token.SignedString(conf.Secret)

	w.Write([]byte(tokenString))
}

/*
DeactivateUser is the handler function for a user deactivating their account
*/
func DeactivateUser(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	s := v["userId"]
	id, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	u, err := GetUserByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else if &u == nil || u.Active == -1 {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	u, err = deactivateUserByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	u.Password = ""
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewEncoder(w).Encode(u); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

/*func Validate(call http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		tokenString := req.Header.Get("Authorization")

		// Return a Token using the cookie
		token, err := jwt.ParseWithClaims(tokenString, jwt.Claims, func(token *jwt.Token) (interface{}, error) {
			// Make sure token's signature wasn't changed
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected siging method")
			}
			return []byte(conf.Secret), nil
		})
		if err != nil {
			http.NotFound(res, req)
			return
		}

		// Grab the tokens claims and pass it into the original request
		if claims, ok := token.Claims.(*Claims); ok && token.Valid {
			ctx := context.WithValue(req.Context(), MyKey, *claims)
			page(res, req.WithContext(ctx))
		} else {
			http.NotFound(res, req)
			return
		}
	})
}*/
