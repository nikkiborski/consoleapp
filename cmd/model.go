package cmd

var dir string = ""

type Input struct {
	Name       string      `json:"name"`
	Result     interface{} `json:"result"`
	Params     []string    `json:"params"`
	Is_action  bool        `json:"is_action"`
	PassParams []string    `json:"passparams"`
	PassResult bool
}