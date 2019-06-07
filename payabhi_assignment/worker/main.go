
package main

import (
	"bytes"
	"encoding/json"
	//"encoding/hex"
	"log"
	"net/http"
	"fmt"
	//"time"
	"io/ioutil"
	"strings"
	// "github.com/go-redis/redis"
)


func main_fun(w http.ResponseWriter, r *http.Request) {

	fmt.Println("inside main ---")
	resp1, err := http.Get("http://localhost:8050/stg/tokens/15")
	if err != nil {
		log.Fatalln(err)
	}
	var token_result map[string]interface{}//to store token

	//json.NewDecoder(resp.Body).Decode(&result)
	jsn, err := ioutil.ReadAll(resp1.Body)
	if err != nil {
		log.Fatal("Error reading the body", err)
	}
	err = json.Unmarshal(jsn, &token_result)
	if err != nil {
		log.Fatal("Decoding error: ", err)
	}

	//log.Println(token_result)
	token := token_result["token"]

	fmt.Println(token)
	fmt.Println("token printed--")

	argument := map[string]interface{}{
		"token": token,
	}

	encoded_argument, err := json.Marshal(argument)
	if err != nil {
		log.Fatalln(err)
	}
	resp2, err := http.Post("http://localhost:8051/hasher", "application/json", bytes.NewBuffer(encoded_argument))
	if err != nil {
		log.Fatalln(err)
	}
	//var hash_result map[string]interface{}
	hash_result := make(map[string]string)

	temp, err := ioutil.ReadAll(resp2.Body)
	if err != nil {
		log.Fatal("Error reading the body", err)
	}
	err = json.Unmarshal(temp, &hash_result)
	if err != nil {
		log.Fatal("Decoding error: ", err)
	}
	fmt.Println(hash_result)
	fmt.Println("Unmarshal done", hash_result["hash"])
	fmt.Println("hash printed---")


	val := strings.HasPrefix(hash_result["hash"], "0")

	fmt.Println("Lucky Hash : ", val)
	fmt.Println("Hash : ", hash_result["hash"])

}
func main() {
	http.HandleFunc("/", main_fun)
	http.ListenAndServe(":8080", nil)
}

