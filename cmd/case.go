package cmd

import (
	"errors"
	"fmt"

	"time"
)

func SubmitCase(input *Input) {
	if len(input.Params)<2{
		input.Result=errors.New("error: case params should be 2: 1)input date 2) comparing value") 
	}else{
		dt1_str := input.Params[0]
		dt1, err := time.Parse(time.RFC3339, dt1_str)
		if err != nil {
			input.Result=err
		}
		dt2_str:=input.Params[1]
		dt2, err := time.Parse(time.RFC3339, dt2_str)
		
		if err != nil {
			input.Result=err
		}
		fmt.Println("DATES:",dt1,dt2)
		if dt1.After(dt2){
			input.Result=true
		}else{
			input.Result=false
		}
	}

}