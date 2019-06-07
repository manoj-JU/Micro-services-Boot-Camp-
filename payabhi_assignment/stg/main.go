package main

import (
	"encoding/json"
	"encoding/hex"
	"math/rand"
	"strconv"
	"net/http"
	"fmt"
	"log"
	"github.com/gorilla/mux"
)

type tokenData struct {
	Token string `json:"token"`
}

func main() {

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/stg/tokens/{size}", GenerateToken).Methods("GET")
	log.Fatal(http.ListenAndServe(":8050", router))
}

func GenerateToken(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w,"manoj")
	vars := mux.Vars(r)
	size := vars["size"]
	len,_:= strconv.Atoi(size)


	byte_arr := make([]byte, len)

	for i := 0; i < len; i++ {
		byte_arr[i]=byte(rand.Intn(30))
	}
	fmt.Println(byte_arr)
	encoded_token := hex.EncodeToString(byte_arr)

	token := tokenData{
		encoded_token,
	}
	fmt.Println(token)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	/*if err := json.NewEncoder(w).Encode(res); err != nil {
	    panic(err)
	}*/
	data,err:= json.Marshal(token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(data)
}
