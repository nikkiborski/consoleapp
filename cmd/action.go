package cmd

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/nikkiborski/consoleapp/pkg/timee"
)

func okay() string {
	return "Successfull"
}
func DoAction(input *Input) {
	switch input.Name {
	case "create":
		filename,err := createFile(input.Params)
		if err != nil {
			input.Result = err.Error()
		} else {
			input.Result = okay()
			if input.PassResult{
				//new file as result
				input.PassParams=[]string{filename}
			}
		}
	case "rename":
		newname,err := rename(input.Params)
		if err != nil {
			input.Result = err.Error()
		} else {
			input.Result = okay()
			if input.PassResult{
				//new name as result
				input.PassParams=[]string{newname}
			}
		}
	case "remove":
		del_time,err := delete(input.Params)
		if err != nil {
			input.Result = err.Error()
		} else {
			input.Result = okay()
			if input.PassResult{
				//deletion time as result
				input.PassParams=[]string{del_time.Format(time.RFC3339)}
			}
		}
	case "bday":
		btime, err := bday(input.Params)
		if err != nil {
			input.Result = err.Error()
		} else {
			input.Result = btime
			if input.PassResult{
				input.PassParams=[]string{btime.Format(time.RFC3339)}
			}
		}
	case "write":
		updatetime,err := add(input.Params)
		if err != nil {
			input.Result = err.Error()
		} else {
			input.Result = okay()
			//return update time as a result
			if input.PassResult{
				input.PassParams=[]string{updatetime.Format(time.RFC3339)}
			}
		}
	default:
		fmt.Println("Unknown cmd")
	}
}

func createFile(params []string) (string,error) {
	if len(params) == 0 {
		return "",errors.New("no name parameter for file creation")
	}
	name := params[0]
	f, err := os.Create(name)
	if err != nil {
		return "",err
	}
	return name, f.Close()

}

func rename(params []string) (string,error) {
	if len(params) < 2 {
		return "",errors.New("lack of required parameters")
	}
	old := params[0]
	new := params[1]

	err := os.Rename(dir+old, dir+new)
	if err != nil {
		return "",err
	}
	return new,nil
}

func delete(params []string) (time.Time,error) {
	if len(params) == 0 {
		return time.Now(),errors.New("lack of required parameters")
	}
	return time.Now(),os.Remove(params[0])
}

func bday(params []string) (time.Time, error) {
	if len(params) == 0 {
		return time.Now(), errors.New("lack of required parameters")
	}
	file := params[0]
	info, err := timee.Stat(file)
	if err != nil {
		return time.Now(), err
	}
	var btime time.Time
	if info.HasBirthTime() {
		btime = info.BirthTime()
	} else {
		return time.Now(), errors.New("no birthdate")
	}
	return btime, nil
}

func add(params []string) (time.Time,error) {
	if len(params) < 2 {
		return time.Now(),errors.New("no string provided")
	}
	filename := params[0]
	s := params[1]
	f, err := os.OpenFile(filename, os.O_RDWR, 0644)
	updateTime:=time.Now()
	if err != nil {
		return updateTime,err
	}
	_, err = f.WriteString(s)
	if err != nil {
		return updateTime,err
	}
	return updateTime,f.Close()
}