package main

import ("fmt"
        "net/rpc/jsonrpc"
        "strings"
        "bufio"
        "os"
        "strconv"
        "encoding/json")

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

	TradeId int `json:"tradeid"`
	UnvestedAmount float32 `json:"unvestedAmount"`
	Stocks []ResponseInner `json:"stocks"`
	
}

type ResponseInner struct{
	ResponseFields ActualResponseFields `json:"fields"`
}

type ActualResponseFields struct{
	Price string `json:"price"`
	Number int `json:"number"`
	Name string `json:"name"`

	
}


type SecondReply struct{
	Stocks []ResponseInner `json:"stocks"`
	CurrentMarketValue float32 `json:"currentMarketValue"`
	UnvestedAmount float32 `json:"unvestedAmount"`
}

func BuyStocks(){

	c,err:= jsonrpc.Dial("tcp","127.0.0.1:1234")
	if err!=nil{
		fmt.Println(err)
		return
	}
	var reply string
	var structRequest Request
	var msg,data,newData []string
	fmt.Println("Enter the request")
	in := bufio.NewReader(os.Stdin)
	line, err := in.ReadString('\n')
	msg = strings.SplitN(line," ",-1)
	data = strings.SplitN(msg[0],":",2)
	newData = strings.SplitN(msg[1],":",2)
	bValue,err:=strconv.ParseFloat(strings.TrimSpace(newData[1]),64)
	data[1]= strings.Replace(data[1],"\"","",-1)
	data[1]= strings.Replace(data[1],"%","",-1)
	fields := strings.SplitN(data[1],",",-1)
	for _,index:=range fields{
			c:= strings.SplitN(index,":",-1)
			a,_:=strconv.Atoi(c[1])
			structFields := ActualFields{Name:c[0],Percentage:a} 
			structInnerRequest := InnerRequest {Fields:structFields}
			structRequest.StockSymbolAndPercentage =append(structRequest.StockSymbolAndPercentage,structInnerRequest)
	}
	result1 := &Request{
    	Budget:float32(bValue),
        StockSymbolAndPercentage: structRequest.StockSymbolAndPercentage} //Map the values to Request structure
    result2, _ := json.Marshal(result1) //Convert the Request to JSON
	err = c.Call("Server.PrintMessage",string(result2),&reply)
	var jsonMsg Response
	var output string
	output = "\"tradeid\":"
	json.Unmarshal([]byte(reply),&jsonMsg)
	output+=strconv.Itoa(jsonMsg.TradeId)+"\n"+"\"stocks\":\""
	for _, i:= range jsonMsg.Stocks{
		output += i.ResponseFields.Name +":"+strconv.Itoa(i.ResponseFields.Number)+":"+"$"+i.ResponseFields.Price+","
	}
	output=strings.Trim(output,",")
	output+="\"\n\"unvestedAmount\":$"+strconv.FormatFloat(float64(jsonMsg.UnvestedAmount),'f',-1,32)		
	if err!=nil {
		fmt.Println(err)
	}else{
		fmt.Println("\nReply\n")
		fmt.Println(output)
	}
}

func SeePortfolio(){

	c,err:= jsonrpc.Dial("tcp","127.0.0.1:1234")
	if err!=nil{
		fmt.Println(err)
		return
	}
	structSRequest:=new(SecondRequest)
	fmt.Println("Enter the request")
	var sRequest string
	fmt.Scanf("%s",&sRequest)
	sRequest= strings.Replace(sRequest,"\"","",-1)
	newsRequest:=strings.SplitN(sRequest,":",-1)
	structSRequest.Tradeid,_= strconv.Atoi(newsRequest[1])
	result3 := &SecondRequest{
		Tradeid: structSRequest.Tradeid}
	result4,_:= json.Marshal(result3)
	var jsonMsg2 SecondReply
	var reply string
	err = c.Call("Server.LossOrGain",string(result4),&reply)
	var output string
	output = "\"stocks\":"
	json.Unmarshal([]byte(reply),&jsonMsg2)
	for _, i:= range jsonMsg2.Stocks{
		output += i.ResponseFields.Name +":"+strconv.Itoa(i.ResponseFields.Number)+":"+i.ResponseFields.Price+","
	}
	output=strings.Trim(output,",")
	output+="\"\n\"currentMarketValue\":$"+strconv.FormatFloat(float64(jsonMsg2.CurrentMarketValue),'f',-1,32)
	output+="\n\"unvestedAmount\":$"+strconv.FormatFloat(float64(jsonMsg2.UnvestedAmount),'f',-1,32)
	if err!=nil {
		fmt.Println(err)
	}else{
		fmt.Println("\nReply\n")
		fmt.Println(output)
	}
}

func main(){
	fmt.Println("Enter your choice\n1.Buying Stocks\n2.Checking your portfolio")
	var choice int64 
	fmt.Scanf("%d\n",&choice)
	switch choice{
		case 1:
			BuyStocks()
			break
		case 2:
			SeePortfolio()
			break
		default:
			fmt.Println("Invalid choice")
			break
		}
}