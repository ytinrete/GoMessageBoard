package main

import(
	"encoding/json"
	"bufio"
    "fmt"
    "os"
    "errors"
    "net/http"
    "sync"
    "io/ioutil" 
    // "strings"
)

type Thread struct{
	Author string
	Content string
	Time string
}

type List struct{
	List []Thread
	mux sync.Mutex
}

func (l* List)add(t Thread) error {
	l.mux.Lock()
	t2, err := saveToFile(t)
	if err == nil{
		l.List = append(l.List, t2)
	}
	defer l.mux.Unlock()
	return err
}

func (l* List)get() []Thread {
	l.mux.Lock()
	defer l.mux.Unlock()
	return l.List
}

var list List

func loadFromFile() error {
	file, err := os.Open("data.txt")
    if err != nil {
        return err
    }
    defer file.Close()
    scanner := bufio.NewScanner(file)
    list = List{
    	List:[]Thread{},
    }
    for scanner.Scan() {
    	line := scanner.Text()
    	fmt.Println("reading text:", line)

    	newobj := Thread{}
    	if err := json.Unmarshal([]byte(line), &newobj); err !=nil{
			fmt.Println("error Unmarshal:", line)
		}else{
			fmt.Println("success Unmarshal:", line)
			list.List = append(list.List, newobj)
		}
    }

    if err := scanner.Err(); err != nil {
        return err
    }
    return nil
}

func saveToFile(t Thread) (Thread, error) {

	if t.Author == "" && t.Content == ""{
		return Thread{}, errors.New("empty data")
	}else{
		if t.Author ==""{
			t.Author = "Anonymous"
		}
	}
	var line []byte
	if res, err := json.Marshal(t); err !=nil{
        return Thread{}, err
	}else{
		line = res
	}

	file, err := os.OpenFile("data.txt", os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
        return Thread{}, err
	}
	defer file.Close()
	w := bufio.NewWriter(file)
	fmt.Fprintln(w, string(line))
	return t, w.Flush()
}

func getListJsonStr()(string, error) {
	if res, err := json.Marshal(list.get()); err !=nil{
        return "", errors.New("json parse error")
	}else{
		fmt.Println("json parse success:", string(res))
		return string(res), nil		
	}
}

func handleGetList(w http.ResponseWriter, r * http.Request) {
	fmt.Println("recieve get request from:", r.RemoteAddr)
	r.ParseForm()
	if r.Method == http.MethodGet {  
        // fmt.Println("key1", r.Form["key1"])   

        // for k, v := range r.Form {  
        //     fmt.Print("key:", k, "; ")  
        //     fmt.Println("val:", strings.Join(v, ""))  
        // }
        if res, err := getListJsonStr(); err != nil{
			
			w.WriteHeader(http.StatusNotFound)//404
	        return
		}else{
			fmt.Println("return to client:", res)
			fmt.Fprintf(w, res)
		}
    }else{
    	w.WriteHeader(http.StatusNotFound)//404
    }
}

func handleAddThread(w http.ResponseWriter, r * http.Request) {
	fmt.Println("recieve post request from:", r.RemoteAddr)
	r.ParseForm()

	if r.Method == http.MethodPost { 
		result, _:= ioutil.ReadAll(r.Body)  
        r.Body.Close()  
        fmt.Println("post body:", string(result))
  
  
        // var f interface{}  
        // json.Unmarshal(result, &f)   
        // m := f.(map[string]interface{})  
        // for k, v := range m {  
        //     switch vv := v.(type) {  
        //     case string:  
        //             fmt.Println(k, "is string", vv)  
        //     case int:  
        //             fmt.Println(k, "is int", vv)  
        //     case float64:  
        //             fmt.Println(k,"is float64",vv)  
        //     case []interface{}:  
        //             fmt.Println(k, "is an array:")  
        //             for i, u := range vv {  
        //                     fmt.Println(i, u)  
        //             }  
        //     default:  
        //             fmt.Println(k, "is of a type I don't know how to handle")   
        //     }  
        //   }  
  
        newobj := Thread{}
    	if err := json.Unmarshal(result, &newobj); err !=nil{
			fmt.Println("net error Unmarshal:", string(result))
			w.WriteHeader(http.StatusNotFound)//404
	        return
		}else{
			fmt.Println("net success Unmarshal:", string(result))
			if err2 := list.add(newobj); err2 != nil{
				
				w.WriteHeader(http.StatusNotFound)//404
		        return
			}
		}

		if res, err := getListJsonStr(); err != nil{
			w.WriteHeader(http.StatusNotFound)//404
	        return
		}else{
			fmt.Println("return to client:", res)
			fmt.Fprintf(w, res)
		}
  
		
    }else{
    	w.WriteHeader(http.StatusNotFound)//404
    }
	
}

func main() {
	if _, err := os.Stat("data.txt"); os.IsNotExist(err) {
	  if _, err2 := os.Create("data.txt"); err2 !=nil{
	  	
        return
	  }
	}
	loadFromFile()
	fmt.Println("stored lists are:")
	for _, line := range list.get(){
		fmt.Println(line)
	}
	fmt.Println("----------------")

	http.HandleFunc("/MessageBoard/getList", handleGetList)
	http.HandleFunc("/MessageBoard/addThread", handleAddThread)
	http.ListenAndServe(":40404", nil)
	
}





