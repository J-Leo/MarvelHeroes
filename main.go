package main

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
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
	Series Series `json:"series"`
	Stories Stories `json:"stories"`
	Events Events `json:"events"`
	URLInf []URL `json:"urls"`
}

type Comics struct {
	Available int  `json:"available"`
	Items []Items `json:"items"`
}

type Series struct {
	Available int `json:"available"`
	Items []SItems `json:"items"`
}

type Stories struct{
	Available int `json:"available"`
	Items []StItems `json:"items"`
}

type Events struct {
	Available int `json:"available"`
	Items []EvItems `json:"items"`
}

type URL struct {
	Type string `json:"type"`
	URL string `json:"url"`
}



type Items struct {
	ResourceURI 	string	`json:"resourceURI"`
	Name 	string		`json:"name"`
}

type SItems struct {
	ResourceURI 	string	`json:"resourceURI"`
	Name 	string		`json:"name"`
}

type StItems struct {
	ResourceURI 	string	`json:"resourceURI"`
	Name 	string		`json:"name"`
	Type    string		`json:"type"`
}

type EvItems struct {
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
	fmt.Printf("0. Salir\n")
	fmt.Printf("1. Buscar superhéroe por nombre\n")
	fmt.Printf("2. Listar los últimos 20 registros\n")

	var opt int	//opcion


	for {

		reader := bufio.NewReader(os.Stdin)
		input,_ := reader.ReadString('\n')

		var err error

		str :=  string([]byte(input)[0])	//Solo se leera el primer numero
		opt,err = strconv.Atoi(str)
		/*Si al convertir el primer caracter a entero hay un error
		// se maneja a continuacion
		*/
		if err != nil  {

			fmt.Println("Error, opción no válida.")
			fmt.Println(err.Error())		//error
			fmt.Println()

			//os.Exit(0)
		}

		switch opt {
			case 0:
				fmt.Println("**Gracias por usar este servicio.**")
				os.Exit(1)
				break
			case 1:
				fmt.Println("Ingrese el nombre del superhéroe que desea buscar:")
				//Ingreso del nombre a buscar
				var name string
				_,err := fmt.Scan(&name)
				if err != nil {
					fmt.Println(err.Error())
				}
				//URL de donde se consume la API
				URL := "http://gateway.marvel.com/v1/public/characters?limit=100&name="+ name +"&ts="+ ts +"&apikey="+ publicKey +"&hash=" + myHash
				response, err := http.Get(URL)
				//fmt.Println(URL)
				/*Si existe algun error al consumir la API se muestra y se termina la ejecucion del programa*/
				if err != nil {
					fmt.Println(err.Error())
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
					//os.Exit(0)
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

					fmt.Println("---------------------------------------------------")

					fmt.Println("Series donde apareció:\n")
					if responseObject.Data.Results[i].Series.Available == 0{
						fmt.Println("No se encontraron registros")
					}

					for k:=0;k < responseObject.Data.Results[i].Series.Available && k<20 ;k++  {
						fmt.Println("Serie N: ",k+1)
						URISerie:=string(responseObject.Data.Results[i].Series.Items[k].ResourceURI)+"?ts="+ ts +"&apikey="+ publicKey +"&hash=" + myHash
						fmt.Println("Nombre:")
						fmt.Println(responseObject.Data.Results[i].Series.Items[k].Name)
						fmt.Println("URL con más información:")
						fmt.Println(URISerie)
						fmt.Println("-----------------------------------")
					}

					fmt.Println("---------------------------------------------------")

					fmt.Println("Historias(Stories) donde apareció:\n")
					if responseObject.Data.Results[i].Stories.Available == 0{
						fmt.Println("No se encontraron registros")
					}

					for l:=0;l < responseObject.Data.Results[i].Stories.Available && l<20 ;l++  {
						fmt.Println("Historia N: ",l+1)
						URIStorie:=string(responseObject.Data.Results[i].Stories.Items[l].ResourceURI)+"?ts="+ ts +"&apikey="+ publicKey +"&hash=" + myHash
						fmt.Println("Nombre:")
						fmt.Println(responseObject.Data.Results[i].Stories.Items[l].Name)
						fmt.Println("URL con más información:")
						fmt.Println(URIStorie)
						fmt.Println("-----------------------------------")
					}

					fmt.Println("---------------------------------------------------")

					fmt.Println("Eventos asociados:\n")
					if responseObject.Data.Results[i].Events.Available == 0{
						fmt.Println("No se encontraron registros")
					}

					for m:=0;m < responseObject.Data.Results[i].Events.Available && m<20 ;m++  {
						fmt.Println("Evento N: ",m+1)
						URIEvent:=string(responseObject.Data.Results[i].Events.Items[m].ResourceURI)+"?ts="+ ts +"&apikey="+ publicKey +"&hash=" + myHash
						fmt.Println("Nombre:")
						fmt.Println(responseObject.Data.Results[i].Events.Items[m].Name)
						fmt.Println("URL con más información:")
						fmt.Println(URIEvent)
						fmt.Println("-----------------------------------")
					}

					fmt.Println("---------------------------------------------------")

					fmt.Println("URL's de interés:\n")
					if len(responseObject.Data.Results[i].URLInf) == 0{
						fmt.Println("No se encontraron registros")
					}

					for n:=0;n < len(responseObject.Data.Results[i].URLInf) && n<20 ;n++  {
						fmt.Println("URL N: ",n+1)
						URLInfo:=string(responseObject.Data.Results[i].URLInf[n].URL)
						fmt.Println("Tipo:")
						fmt.Println(responseObject.Data.Results[i].URLInf[n].Type)
						fmt.Println("Enlace:")
						fmt.Println(URLInfo)
						fmt.Println("-----------------------------------")
					}

					fmt.Println("---------------------------------------------------")
				}
				//fmt.Println("**Gracias por usar este servicio.**\n\n")

				fmt.Println("**--------------------------------**\n\n")

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
					fmt.Println("---------------------------------------------------")

					fmt.Println("Series donde apareció:\n")
					if responseObject.Data.Results[i].Series.Available == 0{
						fmt.Println("No se encontraron registros")
					}

					for k:=0;k < responseObject.Data.Results[i].Series.Available && k<20 ;k++  {
						fmt.Println("Serie N: ",k+1)
						URISerie:=string(responseObject.Data.Results[i].Series.Items[k].ResourceURI)+"?ts="+ ts +"&apikey="+ publicKey +"&hash=" + myHash
						fmt.Println("Nombre:")
						fmt.Println(responseObject.Data.Results[i].Series.Items[k].Name)
						fmt.Println("URL con más información:")
						fmt.Println(URISerie)
						fmt.Println("-----------------------------------")
					}

					fmt.Println("---------------------------------------------------")

					fmt.Println("Historias(Stories) donde apareció:\n")
					if responseObject.Data.Results[i].Stories.Available == 0{
						fmt.Println("No se encontraron registros")
					}

					for l:=0;l < responseObject.Data.Results[i].Stories.Available && l<20 ;l++  {
						fmt.Println("Historia N: ",l+1)
						URIStorie:=string(responseObject.Data.Results[i].Stories.Items[l].ResourceURI)+"?ts="+ ts +"&apikey="+ publicKey +"&hash=" + myHash
						fmt.Println("Nombre:")
						fmt.Println(responseObject.Data.Results[i].Stories.Items[l].Name)
						fmt.Println("URL con más información:")
						fmt.Println(URIStorie)
						fmt.Println("-----------------------------------")
					}

					fmt.Println("---------------------------------------------------")

					fmt.Println("Eventos asociados:\n")
					if responseObject.Data.Results[i].Events.Available == 0{
						fmt.Println("No se encontraron registros")
					}

					for m:=0;m < responseObject.Data.Results[i].Events.Available && m<20 ;m++  {
						fmt.Println("Evento N: ",m+1)
						URIEvent:=string(responseObject.Data.Results[i].Events.Items[m].ResourceURI)+"?ts="+ ts +"&apikey="+ publicKey +"&hash=" + myHash
						fmt.Println("Nombre:")
						fmt.Println(responseObject.Data.Results[i].Events.Items[m].Name)
						fmt.Println("URL con más información:")
						fmt.Println(URIEvent)
						fmt.Println("-----------------------------------")
					}

					fmt.Println("---------------------------------------------------")

					fmt.Println("URL's de interés:\n")
					if len(responseObject.Data.Results[i].URLInf) == 0{
						fmt.Println("No se encontraron registros")
					}

					for n:=0;n < len(responseObject.Data.Results[i].URLInf) && n<20 ;n++  {
						fmt.Println("URL N: ",n+1)
						URLInfo:=string(responseObject.Data.Results[i].URLInf[n].URL)
						fmt.Println("Tipo:")
						fmt.Println(responseObject.Data.Results[i].URLInf[n].Type)
						fmt.Println("Enlace:")
						fmt.Println(URLInfo)
						fmt.Println("-----------------------------------")
					}

					fmt.Println("---------------------------------------------------")
				}
				//fmt.Println("**Gracias por usar este servicio.**\n\n")
				fmt.Println("**--------------------------------**\n\n")

				break
			default:
				clear := exec.Command("clear")
				clear.Stdout = os.Stdout
				clear.Run()
				fmt.Println("Opción no valida")

				break
		}
		fmt.Println("Seleccione una opción")
		fmt.Printf("0. Salir\n")
		fmt.Printf("1. Buscar superhéroe por nombre\n")
		fmt.Printf("2. Listar los últimos 20 registros\n")



	}







}