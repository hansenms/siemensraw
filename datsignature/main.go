//Siemens raw data parser
package main

import(
	"fmt"
	"github.com/hansenms/siemensraw"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %v <DATFILE>\n", os.Args[0])
		return
	}
	fmt.Println(siemensraw.DatFileSignature(os.Args[1]))
}


