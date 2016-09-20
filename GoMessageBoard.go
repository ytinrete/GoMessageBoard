package main

import(
	"encoding/json"
	"bufio"
    "fmt"
    "log"
    "os"
    "errors"
)

type Thread struct{
	Author string
	Content string
}

var list []Thread

func loadFromFile() error {
	file, err := os.Open("data.txt")
    if err != nil {
        log.Fatal(err)
        return err
    }
    defer file.Close()
    scanner := bufio.NewScanner(file)
    list = []Thread{}
    for scanner.Scan() {
    	line := scanner.Text()
    	fmt.Println("reading text:", line)

    	newobj := Thread{}
    	if err := json.Unmarshal([]byte(line), &newobj); err !=nil{
    			fmt.Println("error Unmarshal:", line)
    		}else{
    			fmt.Println("success Unmarshal:", line)
    			list = append(list, newobj)
    		}
    }

    if err := scanner.Err(); err != nil {
        log.Fatal(err)
        return err
    }
    return nil
}

func saveToFile(t Thread) error {

	if t.Author == "" && t.Content == ""{
		return errors.New("empty data")
	}else{
		if t.Author ==""{
			t.Author = "Anonymous"
		}
	}
	var line []byte
	if res, err := json.Marshal(t); err !=nil{
		log.Fatal(err)
        return err
	}else{
		line = res
	}

	file, err := os.OpenFile("data.txt", os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
	    log.Fatal(err)
        return err
	}
	defer file.Close()
	w := bufio.NewWriter(file)
	fmt.Fprintln(w, string(line))
	return w.Flush()
	
}

func main() {
	if _, err := os.Stat("data.txt"); os.IsNotExist(err) {
	  if _, err2 := os.Create("data.txt"); err2 !=nil{
	  	log.Fatal(err2)
        return
	  }
	}
	loadFromFile()
	fmt.Println("stored lists are:"list)

	// t1 := Thread{
	// 	Author:"lee",
	// 	Content:"111",
	// }

	// t2 := Thread{
	// 	Author:"bob",
	// 	Content:"222",
	// }

	// saveToFile(t1)
	// saveToFile(t2)

	
	

	// ioutil.WriteFile("data.txt", []byte("aaa"), 0600)

	
}





