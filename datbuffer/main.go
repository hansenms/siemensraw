//Siemens raw data parser
package main

import(
	"flag"
	"fmt"
	"github.com/hansenms/siemensraw"
	"os"
)

var buffer_name = flag.String("b", "Phoenix", "Buffer to dump from file")

func main() {
	flag.Parse()
	
	if len(flag.Args()) != 1 {
		fmt.Printf("Usage: %v -b <BUFFER NAME> <DATFILE>\n", os.Args[0])
		return
	}

	rawfileinfo := siemensraw.ParseSiemensRaidFile(flag.Args()[0])

	for _, b := range rawfileinfo[len(rawfileinfo)-1].Buffers {
		if b.Name == string(*buffer_name) {
			fmt.Println(b.Buffer)
		}	
	}
}


