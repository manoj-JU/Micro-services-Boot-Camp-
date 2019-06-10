package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	//"encoding/hex"
	"log"
	"net/http"
	"time"
	"io/ioutil"
	"strings"
	"github.com/go-redis/redis"
)


func main_fun(w http.ResponseWriter, r *http.Request) {
	i:=0
	for i==0 {
		fmt.Println("inside main ---")

		conn := redis.NewClient(&redis.Options{
			Addr:         "127.0.0.1:6379",
			Password:     "", // no password set
			DB:           0,  // use default DB
			MaxRetries:   3,
			IdleTimeout:  5 * time.Minute,
			ReadTimeout:  5 * time.Minute,
			WriteTimeout: 5 * time.Minute,
		})

		pong, err := conn.Ping().Result()
		fmt.Println("pong : ", pong)
		if err != nil {
			fmt.Println("error : ", err)
		}
		defer conn.Close()
		resp1, err := http.Get("http://localhost:8050/stg/tokens/15")
		if err != nil {
			log.Fatalln(err)
		}
		var token_result map[string]interface{} //to store token

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

		//! publishing 1 if the hash is "Lucky Hash", else 0
		fmt.Println("Lucky Hash : ", val)
		fmt.Println("Hash : ", hash_result["hash"])
		if val == false {
			resp := conn.Publish("channel1", "NOT A LUCKY HASH").Err()
			if resp != nil {
				fmt.Println(resp)
				panic(err)
			}
		}
		if val == true {
			resp := conn.Publish("channel1", "LUCKY HASH FOUND").Err()
			if resp != nil {
				fmt.Println(resp)
				panic(err)
			}
		}
		time.Sleep(time.Second*3)
	}

}
func main() {
	http.HandleFunc("/", main_fun)
	http.ListenAndServe(":8080", nil)
}
