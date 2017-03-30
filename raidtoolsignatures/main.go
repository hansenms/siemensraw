package main

import "fmt"
import "strings"
import "os/exec"
import "log"
import "strconv"

func main() {
	cmd := exec.Command("cmd.exe", "/C", "raidtool -d -a mars -p 8010")
	stdout, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}

	rt_string := string(stdout[:])
	idx := strings.Index(rt_string, "FileID")
	rt_head := rt_string[:idx]
	rt_body := rt_string[idx:]
	headSlice := strings.Split(rt_head, " ")
	bodySlice := strings.Split(rt_body, " ")

	numFiles, _ := strconv.Atoi(headSlice[35]) // empirically consistent
	fileIDs := make([]string, numFiles+20) // padding to avoid 'panic'
  fileHashes := make([]string, numFiles+20) // padding to avoid 'panic'
	fmt.Println("rt_head: \n", rt_head)

	rr_int := -1
	for rr_loop:=80;rr_loop < len(bodySlice); rr_loop++{ //80 skipping known uninteresting areas

			if len(bodySlice[rr_loop]) == 5 {
				if strings.ContainsAny(bodySlice[rr_loop], "files") == false {
					rr_int++
					fileIDs[rr_int] = bodySlice[rr_loop]
					cmd := exec.Command("cmd.exe", "/C", "raidtool -h "+fileIDs[rr_int]+" -o raidtooltmp.txt -a mars -p 8010")
					_, err := cmd.Output()
					if err != nil {
						log.Fatal(err)
					}

          cmd = exec.Command("cmd.exe", "/C", "hdrsignature raidtooltmp.txt")
					stdout, err = cmd.Output()
					if err != nil {
						log.Fatal(err)
					}

          fileHashes[rr_int]=string(stdout[:])
          fmt.Println("file ID: " + fileIDs[rr_int] + " " + fileHashes[rr_int])
          // cache and output options..
				}
			}
	}


}
