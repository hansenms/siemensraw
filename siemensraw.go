//Siemens raw data and DICOM MRI file utilities
package siemensraw

import(
	"bufio"
	"crypto/sha1"
	"encoding/binary"
    "encoding/base64"
	"fmt"
    "io"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"unicode"
)

type MrParcRaidFileHeader struct {
	HdSize uint32;
	Count uint32;
}

type MrParcRaidFileEntry struct {
	MeasId uint32
	FileId uint32
	Off uint64
	Len uint64
	PatName [64]byte
	ProtName [64]byte
}

type MrParcBuffer struct {
	Name string
	Buffer string
}

type SiemensRaidFile struct
{
	ParcFileEntry MrParcRaidFileEntry
	Buffers []MrParcBuffer
	DataOffset int64
}

func StripString(str string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		if r == '"' {
			return -1
		}
        if r == '\n' {
			return -1
		}
        if r == '\t' {
			return -1
		}
		return r
	}, str)
}

func HashBuffer(str string) []byte {
        cleanString := StripString(str)
        hasher := sha1.New()
        hasher.Write([]byte(cleanString[:]))
        return hasher.Sum(nil)
}

func GetPhoenixFromISMRMRD(filename string) string {
    fi, err := os.Open(filename)
    if err != nil {
        panic(err)
    }
    
    r := bufio.NewReader(fi)  
        
    defer func() {
        if err := fi.Close(); err != nil {
            panic(err)
        }
    }()
    
    // read in a part of ismrmrd file and parse it
    n_bytes := 4096*1024
    buf_bytes := make([]byte, n_bytes)   
    n, err := r.Read(buf_bytes)
    if err != nil && err != io.EOF {
        panic(err)
    }
    
    // fmt.Println("Read in bytes \n", n)
    
    buf := string(buf_bytes[:])  
    idx := strings.Index(buf, "<name>SiemensBuffer_Phoenix</name>")  
    idx2 := strings.Index(buf[idx:n-1], "</value>") 
    buf2 := string(buf[idx:idx+idx2])
    idx3 := strings.Index(buf2, "<value>")
    
    base64_str := buf2[idx3+len("<value>"):len(buf2)]   
    
    // decode, base64
    data, err := base64.StdEncoding.DecodeString(base64_str)
	if err != nil {
		panic(err)
	}
    
    phoenixprot := string(data[:])

    idx = strings.Index(phoenixprot, "### ASCCONV BEGIN")
	idx2 = strings.Index(phoenixprot, "### ASCCONV END ###")
	phoenixprot = string(phoenixprot[idx:idx2+len("### ASCCONV END ###")])
   
    // fmt.Println(phoenixprot)
    
    return phoenixprot
}

func GetPhoenixFromDicom(filename string) string {
	bytes, err := ioutil.ReadFile(filename)
	
	if err != nil {
		panic(err)
	}
	
	asstring := string(bytes[:])
	idx := strings.Index(asstring, "<Name> \"PhoenixMetaProtocol\"")
	asstring = string(bytes[idx:])
	idx = strings.Index(asstring, "### ASCCONV BEGIN")
	idx2 := strings.Index(asstring, "### ASCCONV END ###")
	phoenixprot := string(asstring[idx:idx2+len("### ASCCONV END ###")])
	return phoenixprot
}

func DatFileSignature(filename string) string {
	hash := ""
	
	rawfileinfo := ParseSiemensRaidFile(filename)

	for _, b := range rawfileinfo[len(rawfileinfo)-1].Buffers {
		if b.Name == string("Phoenix") {
            idx := strings.Index(b.Buffer, "### ASCCONV BEGIN")
            idx2 := strings.Index(b.Buffer, "### ASCCONV END ###")
            phoenixprot := string(b.Buffer[idx:idx2+len("### ASCCONV END ###")])
			hash = fmt.Sprintf("%x", HashBuffer(phoenixprot))
		}	
	}
	return hash
}

func HdrFileSignature(filename string) string {
	hash := ""

	bytes, err := ioutil.ReadFile(filename)

	if err != nil {
		panic(err)
	}

	asstring := string(bytes[:])
	idx := strings.LastIndex(asstring, "Phoenix") // first instance = noise scan, last instance = data
	asstring = string(bytes[idx:])
	idx = strings.Index(asstring, "### ASCCONV BEGIN")
	idx2 := strings.Index(asstring, "### ASCCONV END ###")
	phoenixprot := string(asstring[idx:idx2+len("### ASCCONV END ###")])
	hash = fmt.Sprintf("%x", HashBuffer(phoenixprot))

	return hash
}

func PathFromSignature(sig string) string {
	return path.Join(string(sig[0]), string(sig[1]), string(sig[2]), string(sig[3]), string(sig[4]), string(sig[5:]))
}

func ParseSiemensRaidFile(filename string) []SiemensRaidFile {
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var ph MrParcRaidFileHeader

	err = binary.Read(f, binary.LittleEndian, &ph)
	if err != nil {
		fmt.Println("binary.Read failed:", err)
	}

	if ph.HdSize != 0 {
		panic("Only HdSize == 0 (VD or VE data) supported")
	}
	
	var fe = make([]MrParcRaidFileEntry, ph.Count)
	var rf = make([]SiemensRaidFile, ph.Count)
	
	err = binary.Read(f, binary.LittleEndian, &fe)
	if err != nil {
		fmt.Println("binary.Read failed:", err)
	}

	for i, v := range fe {
		rf[i].ParcFileEntry = v
		_, err = f.Seek(int64(fe[i].Off),0)
		if err != nil {
			fmt.Println("Seek failed:", err)
		}

		var dmalength uint32
		var numbuffers uint32
		
		err = binary.Read(f, binary.LittleEndian, &dmalength)
		if err != nil {
			fmt.Println("binary.Read failed:", err)
		}
		
		err = binary.Read(f, binary.LittleEndian, &numbuffers)
		if err != nil {
			fmt.Println("binary.Read failed:", err)
		}
		
		rf[i].Buffers = make([]MrParcBuffer, numbuffers)
	
		for j := uint32(0); j < numbuffers; j++ {
			pos, _ := f.Seek(0,1)
			
			r := bufio.NewReader(f)
			
			rf[i].Buffers[j].Name, err = r.ReadString('\x00')
			if err != nil {
				fmt.Println("binary.Read failed:", err)
			}
			pos, _ = f.Seek(pos+int64(len(rf[i].Buffers[j].Name)),0)
			rf[i].Buffers[j].Name = rf[i].Buffers[j].Name[:len(rf[i].Buffers[j].Name)-1]
			
			var buflen uint32
			err = binary.Read(f, binary.LittleEndian, &buflen)
			if err != nil {
				fmt.Println("binary.Read failed:", err)
			}
			var buf = make([]byte, buflen)
			err = binary.Read(f, binary.LittleEndian, &buf)
			if err != nil {
				fmt.Println("binary.Read failed:", err)
			}
			rf[i].Buffers[j].Buffer = string(buf[:len(buf)-1])
		}
	}
	return rf
}


