package main

import (
	//"fmt"
	"github.com/bmizerany/pat"
	"net/http"
	"log"
	"encoding/json"
	"sync"
	"strconv"
)

type Book struct {
	Name string
	Author string
	Id int
}

type Respond struct{
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
		//fmt.Println(err)
	} else{
		access.Lock()
		bookCnt++
		newBook.Id=bookCnt
		bookList[bookCnt]=newBook
		access.Unlock()
		var Books []Book
		Books=append(Books, bookList[bookCnt])
		json.NewEncoder(w).Encode(Respond{true, "Book Inserted!", Books})
		return
	}
	json.NewEncoder(w).Encode(Respond{false, "Book not Inserted!", nil})
}

func showAllBook(w http.ResponseWriter, r *http.Request){
	var Books []Book
	for _, value:=range bookList{
		Books=append(Books, value)
	}
	json.NewEncoder(w).Encode(Respond{true, "Book List", Books})
}

func showOneBook(w http.ResponseWriter, r *http.Request){
	bookId, err :=strconv.Atoi(r.URL.Query().Get(":bookId"))
	if err!=nil{
		//fmt.Fprintf(w, "Invalid format!\n")
	} else{
		value, flag:=bookList[bookId]
		if flag{
			var Books []Book
			Books=append(Books, value)
			json.NewEncoder(w).Encode(Respond{true, "Book found", Books})
			return
		}
	}
	json.NewEncoder(w).Encode(Respond{false, "No book found for that id", nil})
}

func updateBook(w http.ResponseWriter, r *http.Request){
	decoder:=json.NewDecoder(r.Body)
	defer r.Body.Close()
	var upBook Book
	err:=decoder.Decode(&upBook)
	if err != nil {
		json.NewEncoder(w).Encode(Respond{false, "Error", nil})
	} else{
		_, flag:=bookList[upBook.Id]
		if flag{
			access.Lock()
			bookList[upBook.Id]=upBook
			json.NewEncoder(w).Encode(Respond{true, "Book updated!", nil})
			access.Unlock()
		}else{
			json.NewEncoder(w).Encode(Respond{false, "Invalid book id!", nil})
		}
	}
}

func delBook(w http.ResponseWriter, r *http.Request){
	bookId, err :=strconv.Atoi(r.URL.Query().Get(":bookId"))
	if err!=nil{
		json.NewEncoder(w).Encode(Respond{false, "Error", nil})
	} else{
		_, flag:=bookList[bookId]
		if flag{
			delete(bookList, bookId)
			json.NewEncoder(w).Encode(Respond{true, "Book deleted!", nil})
		}else{
			json.NewEncoder(w).Encode(Respond{false, "Invalid book id!", nil})
		}
	}
}

func main() {
	m := pat.New()

	m.Post("/library", http.HandlerFunc(addBook))
	m.Get("/library", http.HandlerFunc(showAllBook))
	m.Get("/library/:bookId", http.HandlerFunc(showOneBook))
	m.Del("/library/:bookId", http.HandlerFunc(delBook))
	m.Put("/library", http.HandlerFunc(updateBook))

	http.Handle("/", m)
	log.Fatal(http.ListenAndServe(":8000", nil))
}
