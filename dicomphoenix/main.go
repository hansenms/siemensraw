package main

import (
	"fmt"
	"github.com/hansenms/siemensraw"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %v <DICOMFILE>\n", os.Args[0])
		return
	}
	
	phoenix := siemensraw.GetPhoenixFromDicom(os.Args[1])
	fmt.Println(phoenix)
}
