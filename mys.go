package main

import ("fmt"
		"net"
		"net/rpc"
		"net/rpc/jsonrpc"
		"net/http"
		"io/ioutil"
		"encoding/json"
		"strings"
		"strconv")

type Server struct{}
var id int 
var StorageVal map[int]string
type Request struct{
	StockSymbolAndPercentage []InnerRequest `json:"stockSymbolAndPercentage"`
	Budget float32 `json:"budget"`
}

type SecondRequest struct{
	Tradeid int `json:"tradeid"`
}
type InnerRequest struct{
	Fields ActualFields `json:"fields"`
}

type ActualFields struct{
	Name string `json:"name"`
	Percentage int `json:"perecentage"`
}

type Response struct{
	Stocks []ResponseInner `json:"stocks"`
	TradeId int `json:"tradeid"`
	UnvestedAmount float32 `json:"unvestedAmount"`
}

type ResponseInner struct{
	ResponseFields ActualResponseFields `json:"fields"`
}

type ActualResponseFields struct{
	Name string `json:"name"`
	Number int `json:"number"`
	Price string `json:"price"`
}

type SecondReply struct{
	Stocks []ResponseInner `json:"stocks"`
	CurrentMarketValue float32 `json:"currentMarketValue"`
	UnvestedAmount float32 `json:"unvestedAmount"`
}


func (this *Server) PrintMessage(msg string,reply *string) error{
		var jsonInt interface{}
		var structResponse Response
		var jsonMsg Request
		var company string
		var remainder float32=0.0
		json.Unmarshal([]byte(msg),&jsonMsg)
		for _, i:= range jsonMsg.StockSymbolAndPercentage{
			company += i.Fields.Name +","
		}
		company=strings.Trim(company,",")
		response,err:= http.Get("http://finance.yahoo.com/webservice/v1/symbols/"+company+"/quote?format=json")
		if(err!=nil){
			fmt.Println(err)
		}else{
			defer response.Body.Close()
			contents,err:= ioutil.ReadAll(response.Body)
			json.Unmarshal(contents,&jsonInt)
			for i,index := range (jsonInt.(map[string]interface{})["list"]).(map[string]interface{})["resources"].([]interface{}){ 
				price := index.(map[string]interface{})["resource"].(map[string]interface{})["fields"].(map[string]interface{})["price"]
				price1,_ := strconv.ParseFloat(price.(string),64)
				Remainder1:=(float64(jsonMsg.StockSymbolAndPercentage[i].Fields.Percentage) * float64(jsonMsg.Budget))/100
				name := index.(map[string]interface{})["resource"].(map[string]interface{})["fields"].(map[string]interface{})["symbol"]
				number := int( Remainder1/price1)
				remainder = remainder + (float32(price1)*float32(number))
				structActualResponseFields:=ActualResponseFields{Name:name.(string),Number:number,Price:strconv.FormatFloat(price1,'f',-1,64)}
				structResponseInner := ResponseInner{ResponseFields:structActualResponseFields}
				structResponse.Stocks = append(structResponse.Stocks,structResponseInner)
			}
			remainder=jsonMsg.Budget-remainder
			result1 := &Response{
    		TradeId:id,
        	Stocks: structResponse.Stocks,
        	UnvestedAmount:remainder} //Map the values to Request structure
    		result2, _ := json.Marshal(result1) //Convert the Request to JSON
    		*reply = string(result2)
			StorageVal[id]=string(result2)
			id++
			if(err!=nil){
				fmt.Println(err)
			}
				
		}
		
		return nil
}

func (this *Server) LossOrGain(msg string,reply *string) error{
	var jsonReq SecondRequest
	var jsonMsg Response
	var jsonInt interface{}
	var company string
	var price []float64
	var structSecondReply SecondReply
	json.Unmarshal([]byte(msg),&jsonReq)
	tradeid:= jsonReq.Tradeid
	data:= StorageVal[tradeid]
	json.Unmarshal([]byte(data),&jsonMsg)
	for _,index:= range jsonMsg.Stocks{
		company += index.ResponseFields.Name +","
	}
	company=strings.Trim(company,",")
	response,err:= http.Get("http://finance.yahoo.com/webservice/v1/symbols/"+company+"/quote?format=json")
	if(err!=nil){
		fmt.Println(err)
	}else{
		defer response.Body.Close()
		contents,_:= ioutil.ReadAll(response.Body)
		json.Unmarshal(contents,&jsonInt)
		for _,index := range (jsonInt.(map[string]interface{})["list"]).(map[string]interface{})["resources"].([]interface{}){ 
				price1,_ := strconv.ParseFloat((index.(map[string]interface{})["resource"].(map[string]interface{})["fields"].(map[string]interface{})["price"]).(string),64)
				price = append(price,price1)
			}
		var value float32=0.0
		var strprice string
		for i,index := range jsonMsg.Stocks{
				temp,_:= strconv.ParseFloat(ReResponse.Stocks,
        UnvestedAmount:jsonMsg.UnvestedAmount} 
    	result2, _ := json.Marshal(result1) 
    	*reply = string(result2)
	
	}		
	return nil
}

func main(){
	id++
	StorageVal=make(map[int]string)
	rpc.Register(new(Server))
	hear,err:= net.Listen("tcp",":1234")
	if(err!=nil){
		fmt.Println(err)
		return
	}
	for{
		c,error:= hear.Accept()
		if(error!=nil){
			continue
		}
		go jsonrpc.ServeConn(c)
	}

}
