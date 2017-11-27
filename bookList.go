package main

import (
	"github.com/bmizerany/pat"
	"net/http"
	"log"
	"encoding/json"
	"sync"
	"strconv"
	"time"
)

//-----------------------------Book list---------------------------------

type Book struct {
	Name string
	Author string
	Id int
}

type BookResponse struct{
	Ok bool
	Msg string
	Info []Book
}

var (
	bookList=make(map[int]Book)
	access sync.Mutex
	bookCnt int
)

func addBook(w http.ResponseWriter, r *http.Request){
	decoder:=json.NewDecoder(r.Body)
	defer r.Body.Close()
	var newBook Book
	err:=decoder.Decode(&newBook)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(BookResponse{false, "Book not Inserted!", nil})
	} else{
		access.Lock()
		bookCnt++
		newBook.Id=bookCnt
		bookList[bookCnt]=newBook
		access.Unlock()
		var Books []Book
		Books=append(Books, bookList[bookCnt])
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(BookResponse{true, "Book Inserted!", Books})
	}
}

func showAllBook(w http.ResponseWriter, r *http.Request){
	var Books []Book
	for _, value:=range bookList{
		Books=append(Books, value)
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(BookResponse{true, "Book List", Books})
}

func showOneBook(w http.ResponseWriter, r *http.Request){
	bookId, err :=strconv.Atoi(r.URL.Query().Get(":bookId"))
	if err!=nil{
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(BookResponse{false, "Book not Inserted!", nil})
	} else{
		value, flag:=bookList[bookId]
		if flag{
			var Books []Book
			Books=append(Books, value)
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(BookResponse{true, "Book found", Books})
		}else{
			w.WriteHeader(http.StatusNoContent)
			json.NewEncoder(w).Encode(BookResponse{false, "No book found for that id", nil})
		}
	}
}

func updateBook(w http.ResponseWriter, r *http.Request){
	decoder:=json.NewDecoder(r.Body)
	defer r.Body.Close()
	var upBook Book
	err:=decoder.Decode(&upBook)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(BookResponse{false, "Error", nil})
	} else{
		_, flag:=bookList[upBook.Id]
		if flag{
			access.Lock()
			bookList[upBook.Id]=upBook
			access.Unlock()
			w.WriteHeader(http.StatusAccepted)
			json.NewEncoder(w).Encode(BookResponse{true, "Book updated!", nil})

		}else{
			w.WriteHeader(http.StatusNoContent)
			json.NewEncoder(w).Encode(BookResponse{false, "Invalid book id!", nil})
		}
	}
}

func delBook(w http.ResponseWriter, r *http.Request){
	bookId, err :=strconv.Atoi(r.URL.Query().Get(":bookId"))
	if err!=nil{
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(BookResponse{false, "Error", nil})
	} else{
		_, flag:=bookList[bookId]
		if flag{
			access.Lock()
			delete(bookList, bookId)
			access.Unlock()
			w.WriteHeader(http.StatusAccepted)
			json.NewEncoder(w).Encode(BookResponse{true, "Book deleted!", nil})
		}else{
			w.WriteHeader(http.StatusNoContent)
			json.NewEncoder(w).Encode(BookResponse{false, "Invalid book id!", nil})
		}
	}
}

// --------------------------User authentication--------------------------------

type User struct {
	Name string
	UserName string
	Password string
}

var userList=make(map[string]User)

type UserResponse struct {
	Ok bool
	Msg string
	Info User
}

func regUser(w http.ResponseWriter, r *http.Request){
	logoutUser(w, r)
	decoder:=json.NewDecoder(r.Body)
	defer r.Body.Close()
	var newUser User
	err:=decoder.Decode(&newUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(UserResponse{false, "Invalid request!", newUser})
	} else {
		access.Lock()
		if _,found:=userList[newUser.UserName]; found==true {
			w.WriteHeader(http.StatusNotAcceptable)
			json.NewEncoder(w).Encode(UserResponse{false, "User already exists", newUser})
		}else if newUser.UserName=="" || newUser.Password=="" || newUser.Name==""{
			w.WriteHeader(http.StatusNotAcceptable)
			json.NewEncoder(w).Encode(UserResponse{false, "Invalid user info", newUser})
		}else{
			userList[newUser.UserName]=newUser
			w.WriteHeader(http.StatusAccepted)
			json.NewEncoder(w).Encode(UserResponse{true, "Registered new user", newUser})
		}
		access.Unlock()
	}
}

func loginUser(w http.ResponseWriter, r *http.Request){
	cookie, err:=r.Cookie("User")
	var curUser User
	if err==nil{
		curUser.UserName=cookie.Value
		w.WriteHeader(http.StatusNotAcceptable)
		json.NewEncoder(w).Encode(UserResponse{false, "Already logged in", curUser})
		return
	}
	decoder:=json.NewDecoder(r.Body)
	defer r.Body.Close()
	err = decoder.Decode(&curUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(UserResponse{false, "Invalid request!", curUser})
	} else {
		val, found:=userList[curUser.UserName]
		if found==true && val.Password==curUser.Password {
			cookie:=http.Cookie{Name: "User", Value:curUser.UserName, Path:"/"}
			http.SetCookie(w, &cookie)
			w.WriteHeader(http.StatusAccepted)
			json.NewEncoder(w).Encode(UserResponse{true, "Successfully logged in", curUser})
		}else{
			w.WriteHeader(http.StatusNotAcceptable)
			json.NewEncoder(w).Encode(UserResponse{false, "Invalid username or password", curUser})
		}
	}
}

func logoutUser(w http.ResponseWriter, r *http.Request){
	_, err:=r.Cookie("User")
	var curUser User
	if err==nil{
		cookie:=http.Cookie{Name:"User", Value:"", Path:"/", Expires: time.Now()}
		http.SetCookie(w, &cookie)
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(UserResponse{true, "Logged out", curUser})
	}else{
		w.WriteHeader(http.StatusNotAcceptable)
		json.NewEncoder(w).Encode(UserResponse{false, "No active user found", curUser})
	}
}

//----------------------------------------------------------------------------------

func main() {
	m := pat.New()

	m.Post("/library", http.HandlerFunc(addBook))
	m.Get("/library", http.HandlerFunc(showAllBook))
	m.Get("/library/:bookId", http.HandlerFunc(showOneBook))
	m.Del("/library/:bookId", http.HandlerFunc(delBook))
	m.Put("/library", http.HandlerFunc(updateBook))

	m.Post("/register", http.HandlerFunc(regUser))
	m.Post("/login", http.HandlerFunc(loginUser))
	m.Get("/logout", http.HandlerFunc(logoutUser))

	http.Handle("/", m)
	log.Fatal(http.ListenAndServe(":8000", nil))
}
