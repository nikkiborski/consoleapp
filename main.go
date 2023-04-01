package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/nikkiborski/consoleapp/cmd"
)

func main() {
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
	skipNext:=false
	passedParams:=[]string{}
	for i := 0; i < len(input); i++ {
		//append passed params from previous step:
		if len(passedParams)>0{
			input[i].Params=append(input[i].Params, passedParams...)
		}
		if skipNext{
			skipNext=false
		}else if input[i].Is_action{
			cmd.DoAction(&input[i])
		}else{
			cmd.SubmitCase(&input[i])
			if input[i].Result==false{
				skipNext=true
			}
		}
		//better to erase temp 'passedParams' variable ))))
		passedParams=[]string{}
		if len(input[i].PassParams)>0{
			passedParams=input[i].PassParams
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
