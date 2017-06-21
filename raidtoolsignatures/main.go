package main

// not highly efficient method to parse raidtool information
// creates raidtooltmp.txt, raidtool.txt files in the local directory
import (
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	// command line check >>
	if len(os.Args) < 2 {
		fmt.Println("Example usage: \n \n raidtoolsignatures HASHFILE.txt all \n" +
			"\n HASHFILE.txt: REQUIRED - can be an existing list of hashes, or this tool will create a new file with the given name. " +
			"\n all : OPTIONAL- will force transfer of all data on the RAID, otherwise will check Performing Physician field for the following format:" +
			" \"PERF PHYS NAME, [A-Z]{4,5}[0-9]{4,6}-[A-Z0-9]{4,10} \" (i.e. NHLBI1234-A0001)- where the comma is the separator key")
		os.Exit(0)
	}
	// command line check <<

	// check if hash record exists >>
	_, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Printf("This file does not exist, creating %v\n", os.Args[1])
		_, err = os.Create(os.Args[1])
		if err != nil {
			panic(err)
		}
	}
	// check if hash record exists <<

	// force mode check >>
	forceFlag := 0
	if len(os.Args) > 2 {
		forceFlagIO := os.Args[2]
		if strings.ContainsAny(forceFlagIO, "all") == true {
			forceFlag = 1
		}
	}
	// force mode check <<

	// raidtool dump >>
	// debug //	fmt.Println("Raidtool dump") // debug //
	cmd := exec.Command("cmd.exe", "/C", "raidtool -d -a mars -p 8010 > raidtool.txt")
	// offline debug // cmd := exec.Command("cmd.exe", "/C", "RR_rt_print.exe > rt_temp.txt") // offline debug //
	//stdout, err := cmd.Output()
	_, err = cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	// raidtool dump <<

	// load raidtool dump >>
	// debug //	fmt.Println("Raidtool read") // debug //
	rtFile, err2 := ioutil.ReadFile("raidtool.txt")
	if err2 != nil {
		log.Fatal(err2)
	}
	// load raidtool dump <<

	// raidtool header print >>
	rt_string := string(rtFile[:])
	idx := strings.Index(rt_string, "FileID")
	rt_head := rt_string[:idx]
	headSlice := strings.Split(rt_head, " ")
	numFiles, _ := strconv.Atoi(headSlice[35]) // empirically consistent
	fileIDs := make([]string, numFiles+20)     // padding to avoid 'panic'
	fmt.Println("fileID size", len(fileIDs), "rt_head: \n", rt_head)
	// raidtool header print <<

	// Attempt to find measurement IDs using csv (tab delimiting doesn't quite work)
	idx = strings.Index(rt_string, "(fileID)")
	rt_body := rt_string[idx+len("(fileID)"):]
	r := csv.NewReader(strings.NewReader(rt_body))
	r.Comma = '\t'
	// loop through raidtool dump >>
	for {
		// debug //		fmt.Println("Reading CSV") // debug //
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		reg, err2 := regexp.Compile("[^0-9]+")
		if err2 != nil {
			log.Fatal(err2)
		}

		a := record[0]
		// debug // fmt.Println(a) // debug // 
		if len(record[0]) < 100 { // end of file catch
			break
		}

		fileID := reg.ReplaceAllString(a[0:10], "") // [0:10]-12 is affected by retrorecon, 8 is still safe with len(FILEID)=4
		// debug //		fmt.Println("Writing file header>>") // debug //
// offline debug //fmt.Println("raidtool -h "+fileID+" -o raidtooltmp.txt -a mars -p 8010") // offline debug //
		cmd := exec.Command("cmd.exe", "/C", "raidtool -h "+fileID+" -o raidtooltmp.txt -a mars -p 8010")
		_, err = cmd.Output()
		if err != nil {
			log.Fatal(err)
		}

		// debug //		fmt.Println("Research Scan Check>>") // debug //
		// Research scan check >>
		// Check for "PERF PHYS NAME, SITECODE-STUDYNUM" format, where the comma is the key for research scans
		transferFlag := 0
		if forceFlag == 0 {
			bytes, err := ioutil.ReadFile("raidtooltmp.txt")
			if err != nil {
				panic(err)
			}

			asstring := string(bytes[:])
			idx := strings.Index(asstring, "tPerfPhysiciansName")
			idx2 := strings.Index(asstring, "atCoilSelectInfoText")

			//fmt.Println("idx1: " + idx)
			//fmt.Println("idx2: " + idx2)

			catch1 := string(asstring[idx:idx2])
			//fmt.Println(catch1)
			idx = strings.Index(catch1, "{")
			idx2 = strings.Index(catch1, "}")
			if idx2-idx > 5 {
				catch2 := string(catch1[idx+3 : idx2-3]) // inside quotes
				// debug //	fmt.Println(catch2) // debug //

				idx = strings.LastIndex(catch2, ",")
				if idx < 0 {
				// debug //		fmt.Println("Not a research scan")// debug //
				} else {
					reg, err := regexp.Compile("[^a-zA-Z0-9]+")
					if err != nil {
						panic(err)
					}
					// For example: [A-Z]{4,5}[0-9]{4,6}-[A-Z0-9]{4,10}
					// minimum length of identifier = 12
					processedString := reg.ReplaceAllString(catch2[idx+1:], "")
					if len(processedString) > 10 {
					// debug //		fmt.Println("This is a research scan")// debug //
						transferFlag = 1
					}
				}
			}
		} else {
			// Transfer all flag set on command line
			transferFlag = 1
		}
		// Research scan check <<

		// debug //		fmt.Println("transfer>>") // debug //
		// Suitable for transfer - now check if hash exists locally >>
		if transferFlag == 1 {
	// debug //			fmt.Println("Research scan for transfer, checking hash..")// debug //
			cmd = exec.Command("cmd.exe", "/C", "hdrsignature raidtooltmp.txt")

			stdout, err := cmd.Output()
			if err != nil {
				panic(err)
			}
			// offline debug // fmt.Println("cmd.exe", "/C", "hdrsignature raidtooltmp.txt") // offline debug //

			hdrHash := string(stdout[:])

			// check if hash exists >>

			read, err := ioutil.ReadFile(os.Args[1])
			if strings.Contains(string(read), hdrHash) {
		// debug //		fmt.Println("file ID " + fileID + " : Hash exists") // debug //
			} else {
				// @ Flywheel - Data transfer >>
				// ...
				// @ Flywheel - Data transfer <<

				fmt.Println("file ID " + fileID + " : No hash, transferring and appending to log.")

				// append hash >>

				f, err := os.OpenFile(os.Args[1], os.O_APPEND, 0660)
				if err != nil {
					panic(err)
				}

				// debug //				n3, err := f.WriteString(hdrHash) // debug //
				_, err = f.WriteString(hdrHash)
				if err != nil {
					panic(err)
				}
		// debug //		fmt.Printf("wrote %d bytes\n", n3) // debug //
				f.Sync()

				// append hash <<

			} // check if hash exists <<
		} // transferflag <<
	} // loop through raidtool dump <<

}
