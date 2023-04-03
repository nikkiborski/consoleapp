package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"

	"github.com/nikkiborski/consoleapp/cmd"
)

func main() {
	fmt.Println("Started...")
	byted_json,err:=os.ReadFile("input.json")
	if err != nil {
        fmt.Println("Error opening file",err)
    }
	var input []cmd.Input
	err= json.Unmarshal(byted_json,&input)
	if err != nil {
		log.Println("error with unmarshalling",err)
		os.Exit(1)
	}
	skipNext:=0
	passedParams:=[]string{}
	for i := 0; i < len(input); i++ {
		kind := "condition"
		if input[i].Is_action{
			kind = "action"
		}
		if skipNext>0{
			fmt.Println("SVERSHILOS")
		}
		if skipNext==1{
			skipNext--
			continue
		}
		if skipNext>0{
			skipNext--
		}
		
		if len(passedParams)>0{
			//append passed params
			input[i].Params=append(input[i].Params, passedParams...)
		}
		if input[i].Is_action{
			cmd.DoAction(&input[i])
		}else{
			cmd.SubmitCase(&input[i])
			if reflect.TypeOf(input[i].Result).String()!="bool"{
				Finish(input)
			}else{
				if input[i].Result.(bool){
					skipNext=2
				}else{
					skipNext=1
				}
			}
		}
		//better to erase temp 'passedParams' variable ))))
		passedParams=[]string{}
		if len(input[i].PassParams)>0{
			passedParams=input[i].PassParams
		}
		fmt.Printf("step:%v. Type:%v. Result: %v",input[i].Name,kind, input[i].Result)
	}

		Finish(input)
}
func Finish(input []cmd.Input){
	new, err:=json.Marshal(input)
	if err!=nil{
		fmt.Println("problem w/ marhsling")
	}
	err = os.WriteFile("input.json",new,0660)
	if err!=nil{
		fmt.Println("problem w/ final rewrite")
	}
}