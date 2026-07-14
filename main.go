package main

import (
	"MyProject/apiSchema/adminSchema"
	"encoding/json"
	"fmt"
)

func main() {
	marshal, err := json.Marshal(adminSchema.DetailAdminSchema{})
	if err != nil {
		return
	}
	fmt.Println(string(marshal))
	
}
