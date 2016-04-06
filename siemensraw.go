//Siemens raw data and DICOM MRI file utilities
package siemensraw

import(
	"bufio"
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"os"
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
		return r
	}, str)
}

func HashBuffer(str string) []byte {
        cleanString := StripString(str)
	hasher := sha1.New()
	hasher.Write([]byte(cleanString[:]))
	return hasher.Sum(nil)
}

func GetPhoenixFromDicom(filename string) string {
	bytes, err := ioutil.ReadFile(filename)
	
	if err != nil {
		panic(err)
	}
	
	asstring := string(bytes[:])
	idx := strings.Index(asstring, "<Name> \"PhoenixMetaProtocol\"")
	asstring = string(bytes[idx:])
	idx = strings.Index(asstring, "<XProtocol>")
	idx2 := strings.Index(asstring, "### ASCCONV END ###")
	phoenixprot := string(asstring[idx:idx2+len("### ASCCONV END ###")])
	return phoenixprot
}

func DatFileSignature(filename string) string {
	hash := ""
	
	rawfileinfo := ParseSiemensRaidFile(filename)

	for _, b := range rawfileinfo[len(rawfileinfo)-1].Buffers {
		if b.Name == string("Phoenix") {
			hash = fmt.Sprintf("%x", HashBuffer(b.Buffer))
		}	
	}
	return hash
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


