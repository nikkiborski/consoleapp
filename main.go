package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/nikkiborski/consoleapp/pkg/timee"
)

var dir string = ""

type Input struct{
	Name string `json:"name"`
	Result interface{} `json:"result"`
	Params []string `json:"params"`
	Is_action bool `json:"is_action"`
}
func main() {
	byted_json,err:=os.ReadFile("input.json")
	if err != nil {
        fmt.Println("Error opening file",err)
    }
	var input []Input
	err= json.Unmarshal(byted_json,&input)
	if err != nil {
		log.Println("error with unmarshalling",err)
		//os.Exit(1)
	}
	fmt.Println("INPUT::",input)
	for i := 0; i < len(input); i++ {
		if input[i].Is_action{
			DoAction(&input[i])
		}else{
			SubmitCase()
		}
	}

	new, err:=json.Marshal(input)
	if err!=nil{
		fmt.Println("problem w/ marhsling")
	}
	err = os.WriteFile("input.json",new,0660)
	if err!=nil{
		fmt.Println("problem w/ final rewrite")
	}
}
func okay() string{
	return "Successfull"
}
func DoAction(input *Input){
	switch input.Name {
	case "create":
		_,err:=createFile(input.Params)
		if err!=nil{
			input.Result=err.Error()
		}else{
			input.Result=okay()
		}
	case "rename":
		err:=rename(input.Params)
		if err!=nil{
			input.Result=err.Error()
		}else{
			input.Result=okay()
		}
	case "remove":
		err:=delete(input.Params)
		if err!=nil{ 
			input.Result=err.Error()
		}else{
			input.Result=okay()
		}
	case "bday":
		btime,err:=bday(input.Params)
		if err!=nil{
			input.Result=err.Error()
		}else{
			input.Result=btime
		}
	case "write":
		err:=add(input.Params)
		if err!=nil{
			input.Result=err.Error()
		}else{
			input.Result=okay()
		}
	default:
		fmt.Println("Unknown cmd")
	}
}

func SubmitCase(){
	fmt.Println("CASE!!")
}

//actions
type File struct{
	Name string
	source *os.File
	//btime time.Time
}
func createFile(params []string) (File,error){
	if len(params)==0{
		return File{},errors.New("no name parameter for file creation")
	}
	name:=params[0]
	f,err:=os.Create(name)
	if err!=nil{
		fmt.Println("Error while file creation:",err)
	}
	return File{Name:name,source:f},nil

}

func rename(params []string) error{
	if len(params)<2{
		return errors.New("lack of required parameters")
	}
	old:=params[0]
	new:=params[1]

	err:=os.Rename(dir+old,dir+new)
	if err!=nil{
		return err
	}
	return nil
}

func delete(params []string) error{
	if len(params)==0{
		return errors.New("lack of required parameters")
	}
	return os.Remove(params[0])
}

func bday(params []string) (time.Time,error){
	if len(params)==0{
		return time.Now(),errors.New("lack of required parameters")
	}
	file:=params[0]
	info, err := timee.Stat(file)
	if err!=nil{
		return time.Now(),err
	}
	var btime time.Time
	if info.HasBirthTime(){
		btime=info.BirthTime()
	}else{
		return time.Now(),errors.New("no birthdate")
	}
	return btime,nil
}

// func bdayORIGINAL(params []string) (time.Time,error){
// 	if len(params)==0{
// 		return time.Now(),errors.New("lack of required parameters")
// 	}
// 	file:=params[0]
// 	info, err := times.Stat(file)
// 	if err!=nil{
// 		return time.Now(),err
// 	}
// 	var btime time.Time
// 	if info.HasBirthTime(){
// 		btime=info.BirthTime()
// 	}else{
// 		return time.Now(),errors.New("no birthdate")
// 	}
// 	return btime,nil
// }
func add(params []string) error{
	if len(params)<2{
		return errors.New("no string provided")
	}
	s:=params[0]
	filename:=params[1]
	f,err:=os.OpenFile(filename,os.O_RDWR,0644)
	if err!=nil{
		return err
	}
	_,err=f.WriteString(s)
	if err!=nil{
		fmt.Println("Error w/ writing:",err)
	}
	return nil
}