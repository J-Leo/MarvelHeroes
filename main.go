package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)


type Response struct {
	//Code  int `json:"code"`
	//Status string `json:"status"`
	Data Data `json:"data"`
}

type Data struct {
	Results []Results `json:"results"`
	Count 	int		  `json:"count"`
}

type Results struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Desc string `json:"description"`
	Comics Comics `json:"comics"`
}

type Comics struct {
	Available int  `json:"available"`
	Items []Items `json:"items"`
}

//De forma similar se puede agregar mas informacion
//sobre series, eventos, etc
//Pero tanta informacion no seria comodo de leer
//para el usuario final

type Items struct {
	ResourceURI 	string	`json:"resourceURI"`
	Name 	string		`json:"name"`
}

//type ComicResponse struct {
//	Data Data `json:"data"`
//}
//
//type ComicData struct {
//	Results []ComicResults `json:"results"`
//	Count 	int		  `json:"count"`
//}
//
//type ComicResults struct {
//	Id int `json:"id"`
//	DigitalId int `json:"digitalId"`
//	Title string `json:"title"`
//	Desc string `json:"description"`
//	Mod 	string	`json:"modified"`
//	DiamondCode	string	`json:"diamondCode"`
//	PageCount	int		`json:"pageCount"`
//}




func main() {
	const privateKey = "eb6918567366c3ce1ef3c492b0d8d866e913e7cd"
	const publicKey = "2d00c7ca68b6a32e254a60fdc20f7787"

	timems := time.Now().UnixNano() / int64(time.Millisecond)
	ts :=  strconv.FormatInt(timems, 10)

	par := ts + privateKey + publicKey

	hash := md5.New()
	hash.Write([]byte(par))
	myHash := hex.EncodeToString(hash.Sum(nil))

	fmt.Printf("----------------------------------------------------\n")
	fmt.Printf("------------------- Marvel Heros -------------------\n")
	fmt.Printf("----------------------------------------------------\n")


	fmt.Printf("Bienvenido al servicio de consulta de superheroes\n")
	fmt.Printf("Elija una opción\n")
	fmt.Printf("1. Buscar superheroe por nombre\n")
	fmt.Printf("2. Listar los últimos 20 registros\n")

	var opt int
	_,err := fmt.Scan(&opt)
	if err != nil {
		fmt.Print(err.Error())
	}
	switch opt {
	case 1:
		URL := "http://gateway.marvel.com/v1/public/characters?limit=20&ts="+ ts +"&apikey="+ publicKey +"&hash=" + myHash
		response, err := http.Get(URL)
		//fmt.Println(URL)
		if err != nil {
			fmt.Print(err.Error())
			os.Exit(1)
		}
		defer response.Body.Close()
		responseData, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}

		responseObject := Response{}
		json.Unmarshal(responseData, &responseObject)

		for i:=0;i< responseObject.Data.Count ;i++  {
			fmt.Print("Resultado Nº ",(i+1),"\n")
			fmt.Println("ID:")
			fmt.Println(responseObject.Data.Results[i].Id)
			fmt.Println("Nombre:")
			fmt.Println(responseObject.Data.Results[i].Name)
			fmt.Println("Descripción:")
			if responseObject.Data.Results[i].Desc == "" {
				fmt.Println("No hay información disponible")
			}else {
				fmt.Println(responseObject.Data.Results[i].Desc)
			}
			fmt.Println("Comics donde apareció:\n")
			if responseObject.Data.Results[i].Comics.Available == 0{
				fmt.Println("No se encontraron registros")
			}

			for j:=0;j < responseObject.Data.Results[i].Comics.Available && j<20 ;j++  {
				fmt.Println("Comic N: ",j+1)
				URIComic:=string(responseObject.Data.Results[i].Comics.Items[j].ResourceURI)+"?ts="+ ts +"&apikey="+ publicKey +"&hash=" + myHash
				fmt.Println("Nombre:")
				fmt.Println(responseObject.Data.Results[i].Comics.Items[j].Name)
				fmt.Println("URL con mas información:")
				fmt.Println(URIComic)
				fmt.Println("-----------------------------------")
			}
			fmt.Println("-------------------------------------")
		}
		break
	case 2:
		URL := "http://gateway.marvel.com/v1/public/characters?limit=20&ts="+ ts +"&apikey="+ publicKey +"&hash=" + myHash
		response, err := http.Get(URL)
		//fmt.Println(URL)
		if err != nil {
			fmt.Print(err.Error())
			os.Exit(1)
		}
		defer response.Body.Close()
		responseData, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}

		responseObject := Response{}
		json.Unmarshal(responseData, &responseObject)

		for i:=0;i< responseObject.Data.Count ;i++  {
			fmt.Print("Resultado Nº ",(i+1),"\n")
			fmt.Println("ID:")
			fmt.Println(responseObject.Data.Results[i].Id)
			fmt.Println("Nombre:")
			fmt.Println(responseObject.Data.Results[i].Name)
			fmt.Println("Descripción:")
			if responseObject.Data.Results[i].Desc == "" {
				fmt.Println("No hay información disponible")
			}else {
				fmt.Println(responseObject.Data.Results[i].Desc)
			}
			fmt.Println("Comics donde apareció:\n")
			if responseObject.Data.Results[i].Comics.Available == 0{
				fmt.Println("No se encontraron registros")
			}

			for j:=0;j < responseObject.Data.Results[i].Comics.Available && j<20 ;j++  {
				fmt.Println("Comic N: ",j+1)
				URIComic:=string(responseObject.Data.Results[i].Comics.Items[j].ResourceURI)+"?ts="+ ts +"&apikey="+ publicKey +"&hash=" + myHash
				fmt.Println("Nombre:")
				fmt.Println(responseObject.Data.Results[i].Comics.Items[j].Name)
				fmt.Println("URL con mas información:")
				fmt.Println(URIComic)
				fmt.Println("-----------------------------------")
			}
			fmt.Println("-------------------------------------")
		}
		break
	}








}