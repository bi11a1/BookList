package main

import (
	"fmt"
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
		fmt.Println(err)
	} else{
		access.Lock()
		bookCnt++
		newBook.Id=bookCnt
		bookList[bookCnt]=newBook
		access.Unlock()
		fmt.Println(bookList[bookCnt])
		fmt.Fprintf(w, "Book Inserted!\n")
	}
}

func showAllBook(w http.ResponseWriter, r *http.Request){
	for _, value:=range bookList{
		fmt.Fprintf(w, "Name: %s, Author: %s\n", value.Name, value.Author, )
	}
}

func showOneBook(w http.ResponseWriter, r *http.Request){
	bookId, err :=strconv.Atoi(r.URL.Query().Get(":bookId"))
	if err!=nil{
		fmt.Fprintf(w, "Invalid format!\n")
	} else{
		value, flag:=bookList[bookId]
		if flag{
			fmt.Fprintf(w, "Name: %s, Author: %s\n", value.Name, value.Author)
		}else{
			fmt.Fprintf(w,"Book not found for that id!\n")
		}
	}
}

func updateBook(w http.ResponseWriter, r *http.Request){
	decoder:=json.NewDecoder(r.Body)
	defer r.Body.Close()
	var upBook Book
	err:=decoder.Decode(&upBook)
	if err != nil {
		fmt.Println(err)
	} else{
		_, flag:=bookList[upBook.Id]
		if flag{
			bookList[upBook.Id]=upBook
			fmt.Fprintf(w, "Book updated\n")
		} else{
			fmt.Fprintf(w, "Book not found for that id\n")
		}
	}
}

func delBook(w http.ResponseWriter, r *http.Request){
	bookId, err :=strconv.Atoi(r.URL.Query().Get(":bookId"))
	if err!=nil{
		fmt.Fprintf(w, "Invalid format!\n")
	} else{
		_, flag:=bookList[bookId]
		if flag{
			delete(bookList, bookId)
			fmt.Fprintf(w, "Book deleted\n")
		}else{
			fmt.Fprintf(w,"Book not found for that id!\n")
		}
	}
}

func nothing(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "Nothing\n")
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
