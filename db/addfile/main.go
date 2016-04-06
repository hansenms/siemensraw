//Application adds a siemens file to the database
package main

import(
	"encoding/json"
	"flag"
	"fmt"
	"github.com/hansenms/siemensraw"
	"io/ioutil"
	"os"
	"path"
	
)

func main() {
	var dest = flag.String("d", ".", "Destination folder")
	var basepath = flag.String("b", ".", "Base path to remove")
	
	flag.Parse()

	if len(flag.Args()) != 1 {
		fmt.Printf("Usage: %v -d <DEST FOLDER> <DATFILENAME>\n", os.Args[0])
	}

	filename := flag.Args()[0]
	sig := siemensraw.DatFileSignature(filename)

	base := path.Base(filename)
	ext := path.Ext(filename)
	base = base[:len(base)-len(ext)]
	newfilename := base + ".json"
	destfolder := path.Join(*dest, siemensraw.PathFromSignature(sig))
	destpath := path.Join(destfolder, newfilename)

	drivepath := path.Dir(filename)
	if drivepath[:len(*basepath)] == (*basepath)[:] {
		drivepath = drivepath[len(*basepath):]
	}
	
	infomap := map[string]string{"origin":filename, "filename": path.Base(filename), "drivepath": drivepath}
	infoj, _ := json.Marshal(infomap)

	err := os.MkdirAll(destfolder, 0777)
	if err != nil {
		fmt.Printf("Error creating folder %v\n", destfolder)
		panic(err)
	}
	
	err = ioutil.WriteFile(destpath, infoj, 0644)
	if err != nil {
		fmt.Printf("Error writing DB entry to file %v\n", destpath)
		panic(err)
	}
	fmt.Printf("Added: %v\n", destpath)
}
