package main

import (
	"fmt"
	"github.com/hansenms/siemensraw"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %v <ISMRMRDFILE>\n", os.Args[0])
		return
	}
	
	phoenix := siemensraw.GetPhoenixFromISMRMRD(os.Args[1])
	fmt.Printf("%x\n", siemensraw.HashBuffer(phoenix))
}
