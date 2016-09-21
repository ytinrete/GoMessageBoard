package main  
  
import (  
        "fmt"  
        "net/url"  
        "net/http"  
        "io/ioutil"  
        "log"
        "encoding/json"
        "bytes"  
)

func doGet() {

	u, _ := url.Parse("http://www.ytinrete.com:40404/MessageBoard/getList")  
    q := u.Query()  
    q.Set("username", "user")  
    q.Set("password", "passwd")  
    u.RawQuery = q.Encode()  
    res, err := http.Get(u.String());  
    if err != nil {   
          log.Fatal(err) 
          return   
    }  
    result, err := ioutil.ReadAll(res.Body)   
    res.Body.Close()   
    if err != nil {   
          log.Fatal(err) 
          return   
    }   
    fmt.Println(res.Status, string(result))  
  	
}

type Server struct {  
        ServerName string  
        ServerIP   string  
}  
  
type Serverslice struct {  
        Servers []Server  
        ServersID  string  
}

type Thread struct{
	Author string
	Content string
}

func doPost() {

    t1 := Thread{
    	Author:"fake",
    	Content:"444",
	}

    b, err := json.Marshal(t1)  
    if err != nil {  
        fmt.Println("json err:", err)
        return
    }  

    body := bytes.NewBuffer([]byte(b))  
    res,err := http.Post("http://www.ytinrete.com:40404/MessageBoard/addThread", "application/json;charset=utf-8", body)  
    if err != nil {  
        log.Fatal(err)  
        return  
    }  
    result, err := ioutil.ReadAll(res.Body)  
    res.Body.Close()  
    if err != nil {  
        log.Fatal(err)  
        return  
    }  
    fmt.Println(res.Status, string(result))
    fmt.Println("done")
	
}


  
func main() {
	doGet()
	// doPost()
}   