package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

var store *Store
var ring *HashRing
var selfAddress string

func Init(port string) {

	nodes := []string{"localhost:8080", "localhost:8081", "localhost:8082"}
	ring = NewHashRing(nodes)
	selfAddress = "localhost:" + port
	store = NewStore("wal_" + port + ".log")
	store.cleanupExpired()
}

func GetEndpoint(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")

	node := ring.GetNode(key)

	if node == selfAddress {
		if val, ok := store.Get(key); ok {
			json.NewEncoder(w).Encode(val)
		} else {
			http.Error(w, "Key not found", http.StatusNotFound)
		}
	} else {
		val, err := forwardGet(node, key)
		if err != nil {
			http.Error(w, "Error forwarding request", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(val)
	}
}

func SetEndpoint(w http.ResponseWriter, r *http.Request) {
	var body SetRequestBody

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Convert TTL string to int64
	ttl, err := strconv.ParseInt(body.TTL, 10, 64)
	if err != nil {
		http.Error(w, "Invalid TTL value", http.StatusBadRequest)
		return
	}

	node := ring.GetNode(body.Key)

	if node == selfAddress {
		store.Set(body.Key, body.Value, ttl)
		json.NewEncoder(w).Encode("Key set successfully (local)")
	} else {
		forwardSet(node, body)
		json.NewEncoder(w).Encode("Key set successfully (node: " + node + ")")
	}

}

func DeleteEndpoint(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")

	node := ring.GetNode(key)

	if node == selfAddress {
		store.Delete(key)
		json.NewEncoder(w).Encode("Key deleted successfully (local)")
	} else {
		forwardDelete(node, key)
		json.NewEncoder(w).Encode("Key deleted successfully (node: " + node + ")")
	}
}

func forwardSet(node string, body SetRequestBody) {
	jsonData, err := json.Marshal(body)
	if err != nil {
		fmt.Println("Error marshalling request body:", err)
		return
	}

	_, err = http.Post("http://"+node+"/set", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
}

func forwardGet(node string, key string) (string, error) {
	request, err := http.NewRequest("GET", "http://"+node+"/get?key="+key, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return "", err
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return "", err
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return "", err
	}

	return string(body), nil
}

func forwardDelete(node string, key string) {
	request, err := http.NewRequest("DELETE", "http://"+node+"/delete?key="+key, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	_, err = http.DefaultClient.Do(request)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
}
