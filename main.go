package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	const privateKey = "eb6918567366c3ce1ef3c492b0d8d866e913e7cd";
	const publicKey = "2d00c7ca68b6a32e254a60fdc20f7787";

	timems := time.Now().UnixNano() / int64(time.Millisecond);
	ts :=  strconv.FormatInt(timems, 10);

	par := ts + privateKey + publicKey;

	hash := md5.New();
	hash.Write([]byte(par));
	myHash := hex.EncodeToString(hash.Sum(nil));

	URL := "http://gateway.marvel.com/v1/public/comics?ts="+ ts +"&apikey="+ publicKey +"&hash=" + myHash;
	response, err := http.Get(URL);

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(responseData))
	/*

	//time := time.Now();
	fmt.Printf("%s \n",myHash);
	fmt.Println(URL);
	fmt.Printf(ts);*/




}