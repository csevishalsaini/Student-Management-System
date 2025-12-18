package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

type user struct{
	Name string `json:"name"`
	Age int `json:"age"`
	City string `json:"city"`
}


func rootHandler(w http.ResponseWriter, r *http.Request){
	w.Write([]byte("Hello Root route"))
	fmt.Println("Hello root route ")
}

func teachersHandler(w http.ResponseWriter, r *http.Request){
	fmt.Println(r.URL.Path)
	path := strings.TrimPrefix(r.URL.Path,"/teachers/")
	userID := strings.TrimSuffix(path,"/")
	fmt.Println("The user id: ",userID)
	queryParam := r.URL.Query()
	fmt.Println("query is ",queryParam)
	sortby:= queryParam.Get("sortby")
	key := queryParam.Get("key")
	sortorder := queryParam.Get("sortorder")

	if(sortorder == ""){
		sortorder = "DESC"
	}
	fmt.Printf("sortby: %v, sort order: %v, key: %v",sortorder,key,sortby)

	switch(r.Method){
	case http.MethodGet:
		w.Write([]byte("Hello GET Method on Teachers route "))
		// fmt.Println("Hello GET Method on Teachers route ")

	case http.MethodPost:
		w.Write([]byte("Hello Post Method on Teachers route "))
		fmt.Println("Hello Post Method on Teachers route ")

	case http.MethodPut:
		w.Write([]byte("Hello Put Method on Teachers route "))
		fmt.Println("Hello Put Method on Teachers route ")

	case http.MethodPatch:
		w.Write([]byte("Hello Patch Method on Teachers route "))
		fmt.Println("Hello Patch Method on Teachers route ")

	case http.MethodDelete:
		w.Write([]byte("Hello Delete Method on Teachers route "))
		fmt.Println("Hello Delete Method on Teachers route ")
	}

	// w.Write([]byte("Hello teachers route "))
	// fmt.Println("Hello teachers route ")
}

func studentsHandler(w http.ResponseWriter, r *http.Request){
	switch(r.Method){
	case http.MethodGet:
		w.Write([]byte("Hello GET Method on Students route "))
		fmt.Println("Hello GET Method on Students route ")

	case http.MethodPost:
		w.Write([]byte("Hello Post Method on Students route "))
		fmt.Println("Hello Post Method on Students route ")

	case http.MethodPut:
		w.Write([]byte("Hello Put Method on Students route "))
		fmt.Println("Hello Put Method on Students route ")

	case http.MethodPatch:
		w.Write([]byte("Hello Patch Method on Students route "))
		fmt.Println("Hello Patch Method on Students route ")

	case http.MethodDelete:
		w.Write([]byte("Hello Delete Method on Students route "))
		fmt.Println("Hello Delete Method on Students route ")
	}

	w.Write([]byte("Hello students route "))
	fmt.Println("Hello students route ")
}

func execsHandler(w http.ResponseWriter, r *http.Request){
	switch(r.Method){
	case http.MethodGet:
		w.Write([]byte("Hello GET Method on Execs route "))
		fmt.Println("Hello GET Method on Execs route ")

	case http.MethodPost:
		w.Write([]byte("Hello Post Method on Execs route "))
		fmt.Println("Hello Post Method on Execs route ")

	case http.MethodPut:
		w.Write([]byte("Hello Put Method on Execs route "))
		fmt.Println("Hello Put Method on Execs route ")

	case http.MethodPatch:
		w.Write([]byte("Hello Patch Method on Execs route "))
		fmt.Println("Hello Patch Method on Execs route ")

	case http.MethodDelete:
		w.Write([]byte("Hello Delete Method on Execs route "))
		fmt.Println("Hello Delete Method on Execs route ")
	}

	w.Write([]byte("Hello execs route "))
	fmt.Println("Hello execs route ")
}

func main(){
		http.HandleFunc("/", rootHandler)
		http.HandleFunc("/teachers/",teachersHandler)
		http.HandleFunc("/students/", studentsHandler)
		http.HandleFunc("/execs/",studentsHandler)


		port := ":3000"
		fmt.Println("Server is running on port",port)
		err :=http.ListenAndServe(port,nil)
		if(err != nil){
			log.Fatal("Error starting the server ",err)
		}
	
}