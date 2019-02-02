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

//Mapeo del JSON recibido

//Solo los datos mas resaltantes
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
	Mod	 string	`json:"modified"`
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
//
type Items struct {
	ResourceURI 	string	`json:"resourceURI"`
	Name 	string		`json:"name"`
}

func main() {
	const privateKey = "eb6918567366c3ce1ef3c492b0d8d866e913e7cd"
	const publicKey = "2d00c7ca68b6a32e254a60fdc20f7787"

	timems := time.Now().UnixNano() / int64(time.Millisecond)	//Formato similar a Date.now() en JS
	ts :=  strconv.FormatInt(timems, 10)					//parse -> String

	par := ts + privateKey + publicKey							//parametro para md5

	hash := md5.New()											//Hash md5
	hash.Write([]byte(par))
	myHash := hex.EncodeToString(hash.Sum(nil))

	fmt.Printf("----------------------------------------------------\n")
	fmt.Printf("------------------- Marvel Heros -------------------\n")
	fmt.Printf("----------------------------------------------------\n")


	fmt.Printf("Bienvenido al servicio de consulta de superhéroes\n")
	fmt.Printf("Elija una opción\n")
	fmt.Printf("1. Buscar superhéroe por nombre\n")
	fmt.Printf("2. Listar los últimos 20 registros\n")

	var opt int			//opcion
	_,err := fmt.Scan(&opt)
	if err != nil {
		fmt.Print(err.Error())		//error
		os.Exit(0)
	}
	switch opt {
	case 1:
		fmt.Println("Ingrese el nombre del superhéroe que desea buscar:")
		//Ingreso del nombre a buscar
		var name string
		_,err := fmt.Scan(&name)
		if err != nil {
			fmt.Print(err.Error())
		}
		//URL de donde se consume la API
		URL := "http://gateway.marvel.com/v1/public/characters?limit=100&name="+ name +"&ts="+ ts +"&apikey="+ publicKey +"&hash=" + myHash
		response, err := http.Get(URL)
		//fmt.Println(URL)
		/*Si existe algun error al consumir la API se muestra y se termina la ejecucion del programa*/
		if err != nil {
			fmt.Print(err.Error())
			os.Exit(1)
		}
		defer response.Body.Close()
		/*El cuerpo del response se formatea, si existe algun error en el proceso se captura y se muestra*/
		responseData, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}

		responseObject := Response{}
		json.Unmarshal(responseData, &responseObject)			//A la data de respuesta se le
																// da la forma de la estructura responseObject

		if responseObject.Data.Count == 0 {
			fmt.Println("Lo sentimos, no se encontró ningún superhéroe registrado con ese nombre")
			os.Exit(0)
		}

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
			fmt.Println("Fecha de modificación:")
			fmt.Println(responseObject.Data.Results[i].Mod)
			fmt.Println("Comics donde apareció:\n")
			if responseObject.Data.Results[i].Comics.Available == 0{
				fmt.Println("No se encontraron registros")
			}

			for j:=0;j < responseObject.Data.Results[i].Comics.Available && j<20 ;j++  {
				fmt.Println("Comic N: ",j+1)
				URIComic:=string(responseObject.Data.Results[i].Comics.Items[j].ResourceURI)+"?ts="+ ts +"&apikey="+ publicKey +"&hash=" + myHash
				fmt.Println("Nombre:")
				fmt.Println(responseObject.Data.Results[i].Comics.Items[j].Name)
				fmt.Println("URL con más información:")
				fmt.Println(URIComic)
				fmt.Println("-----------------------------------")
			}
			fmt.Println("-------------------------------------")
		}
		fmt.Println("**Gracias por usar este servicio.**")
		break
	case 2:
		URL := "http://gateway.marvel.com/v1/public/characters?limit=20&orderBy=name&ts="+ ts +"&apikey="+ publicKey +"&hash=" + myHash
		response, err := http.Get(URL)
		//fmt.Println(URL)
		/*Si existe algun error al consumir la API se muestra y se termina la ejecucion del programa*/
		if err != nil {
			fmt.Print(err.Error())
			os.Exit(1)
		}
		defer response.Body.Close()
		/*El cuerpo del response se formatea, si existe algun error en el proceso se captura y se muestra*/
		responseData, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}

		responseObject := Response{}
		json.Unmarshal(responseData, &responseObject)			//A la data de respuesta se le
																// da la forma de la estructura responseObject

		for i:=0;i< responseObject.Data.Count ;i++  {
			fmt.Print("Heroe Nº ",(i+1),"\n")
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
			fmt.Println("Fecha de modificación:")
			fmt.Println(responseObject.Data.Results[i].Mod)
			fmt.Println("Comics donde apareció:\n")
			if responseObject.Data.Results[i].Comics.Available == 0{
				fmt.Println("No se encontraron registros")
			}

			for j:=0;j < responseObject.Data.Results[i].Comics.Available && j<20 ;j++  {
				fmt.Println("Comic N: ",j+1)
				//URI para mayor informacion del comic
				URIComic:=string(responseObject.Data.Results[i].Comics.Items[j].ResourceURI)+"?ts="+ ts +"&apikey="+ publicKey +"&hash=" + myHash
				fmt.Println("Nombre:")
				fmt.Println(responseObject.Data.Results[i].Comics.Items[j].Name)
				fmt.Println("URL con mas información:")
				fmt.Println(URIComic)
				fmt.Println("-----------------------------------")
			}
			fmt.Println("-------------------------------------")
		}
		fmt.Println("**Gracias por usar este servicio.**")
		break
	}








}