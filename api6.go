package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv" ;"strings" 
)


type Product struct{// all the data type defined in product struct
	ProductID int 
	Manufacturer string
	PricePerUnit int
	ProductName string
}
var productlist []Product
func init(){
	//provideing the data to api still have to readme about it
productJSON:=`[{
	"ProductID":1,
"Manufacturer":"Company 1",
"PricePerUnit":100,
"ProductName":"Product1"
},
{
	"ProductID":2,
"Manufacturer":"Company 2",
"PricePerUnit":200,
"ProductName":"Product2"
},
{
	"ProductID":3,
"Manufacturer":"Company 3",
"PricePerUnit":300,
"ProductName":"Product3"
}]`

err:=json.Unmarshal([]byte(productJSON),&productlist)//Unmartshal== decoding json
if err!=nil{
	log.Fatal(err)//print the error if any
	//fmt.Println("sometthing went wrong")
}
}
func getnextID() int{//If new data is been added then a new and max id will be assigned to it
	highestID:=-1
	for _,Product:=range productlist{
		if highestID<Product.ProductID{
			highestID=Product.ProductID
		}
	}
	return highestID+1
}
//*----
//PUT with Paramatric Route
func singleHandler(w http.ResponseWriter,r *http.Request){//Created a handler to oprate single product
	urlPathseg:=strings.Split(r.URL.Path,"products/")  //URL is a predefined struct
	productID,err:=strconv.Atoi(urlPathseg[len(urlPathseg)-1])//AToi
	if err!=nil{
		w.WriteHeader(http.StatusNotFound)
		return
	}
	product,listItemindex:= findproductbyID(productID)
	if product==nil{
		w.WriteHeader(http.StatusNotFound)
		return
	}
	switch r.Method{
	case http.MethodGet://to return a single product
		productJSON,err:=json.Marshal(product)//you found the error here
		if err!=nil{
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type","application/json")
		w.Write(productJSON)
		//--
	case http.MethodPut://PUT method
	var updatedProduct Product
	bodyBytes,err:=ioutil.ReadAll(r.Body)//read the rewriten json
	if err!=nil{
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err=json.Unmarshal(bodyBytes,&updatedProduct)//check if you can decode that json or not
		if err!=nil{
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if updatedProduct.ProductID!= productID{//Product id check(should be same)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	product=&updatedProduct//making changes in the currunt product
		productlist[listItemindex]=*product
		w.WriteHeader(http.StatusOK)
		return
	}
	}
	func findproductbyID(productID int)(*Product,int){//Function to find the product by ProductID
	for i,product:=range productlist{
		if product.ProductID==productID{
			return &product,i
		}
	}
	return nil,0
	}
	//-----*/


func productHandler(w http.ResponseWriter,r *http.Request){//handler function{handles the request and write responses}
switch r.Method{
case http.MethodGet://GET method 
	productJSON,err:=json.Marshal(productlist)//re encoding the decoded json data to show
	if err!=nil{
		w.WriteHeader(http.StatusInternalServerError)//500 if any error occured
		return
	}
	w.Header().Set("Content-Type","application/json")//set the data in much organized way
	w.Write(productJSON)//gives the data

case http.MethodPost://POST  method
	var newProduct  Product
	bodyBytes,err:=ioutil.ReadAll(r.Body)//read the new input data
	if err!=nil{
		w.WriteHeader(http.StatusBadRequest)//if not able to read give 400
	}
	err=json.Unmarshal(bodyBytes,&newProduct)//decode the given data while addig to the provided data
	if err!=nil{
		w.WriteHeader(http.StatusBadRequest)//data is not able to decode the 400
	}
	if newProduct.ProductID !=0{//assign a new product id
		w.WriteHeader(http.StatusBadRequest)//not able to assign or already assigned then 400
		return
	}
	newProduct.ProductID=getnextID()//assigning new product a new id
	productlist=append(productlist, newProduct)//adding it to the product list
	w.WriteHeader(http.StatusCreated)//201 if all went good
	return
}		
}
func main(){

	
	http.HandleFunc("/products",productHandler )//calling he handler function with http handlefunc there is one of the way check 3
	http.HandleFunc("/products/",singleHandler )//handler function for each product
	http.ListenAndServe(":5000",nil)//listnug to port 5000 and nil refers to default surv mux

}