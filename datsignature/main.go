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

	rawfileinfo := siemensraw.ParseSiemensRaidFile(os.Args[1])

	for _, b := range rawfileinfo[len(rawfileinfo)-1].Buffers {
		if b.Name == string("Phoenix") {
			fmt.Printf("%x\n", siemensraw.HashBuffer(b.Buffer))
		}	
	}
}


