//Siemens hdr data parser
package main

import(
	"fmt"
	"os"
)
import "siemensraw"

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %v <DATFILE>\n", os.Args[0])
		return
	}
	fmt.Println(siemensraw.HdrFileSignature(os.Args[1]))
}
