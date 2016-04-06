//Application finds a raw file in database
package main

import(
	//"encoding/json"
	"flag"
	"fmt"
	"github.com/hansenms/siemensraw"
	"io/ioutil"
	"os"
	"path"
	
)

func main() {
	var db = flag.String("d", ".", "Database folder folder")
	
	flag.Parse()

	if len(flag.Args()) != 1 {
		fmt.Printf("Usage: %v -d <DATABASE FOLDER> <SIGNATURE>\n", os.Args[0])
	}

	sig := flag.Args()[0]
	folder := path.Join(*db, siemensraw.PathFromSignature(sig))

	_, err := os.Stat(folder)
	if err != nil {
		fmt.Printf("Unable to find file with signature: %v\n", sig)
		return
	}

	files, _ := ioutil.ReadDir(folder)
	if len(files) < 1 {
		fmt.Printf("Unable to find file with signature: %v\n", sig)
		return
	}

	retstring := "["
	sep := ""
	for _, f := range files {
		json, err := ioutil.ReadFile(path.Join(folder,f.Name()))
		if err != nil {
			fmt.Printf("Error reading file %v\n", path.Join(folder,f.Name()))
			return
		}
		retstring += sep + string(json)
		sep = ", "
	}
	retstring += "]"
	fmt.Println(retstring)
}
