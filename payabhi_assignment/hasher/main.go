package main

import (
	"encoding/json"
	"encoding/hex"
	"net/http"
	//"fmt"
	"log"
	"crypto/sha256"
	"io/ioutil"
	"github.com/gorilla/mux"
)


type hashData struct {
	Hash string `json:"hash"`
}

type tokenData struct {
	Token string
}


func main() {

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/hasher", GenerateHash).Methods("POST")
	log.Fatal(http.ListenAndServe(":8051", router))
}

func GenerateHash(w http.ResponseWriter, r *http.Request) {

	jsn, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal("Error reading the body", err)
	}

	var token tokenData
	err = json.Unmarshal(jsn, &token)
	if err != nil {
		log.Fatal("Decoding error: ", err)
	}

	decoded_token, err := hex.DecodeString(token.Token)
	if err != nil {
		log.Fatal(err)
	}

	hash := sha256.Sum224([]byte(decoded_token))

	encoded_hash := hex.EncodeToString(hash[:])

	res := hashData{
		encoded_hash,
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	data,err:= json.Marshal(res)
	if err != nil {
		panic(err)
	}
	w.Write(data)

}
