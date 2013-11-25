package main 
import (
    "fmt"
    "io/ioutil"
    "container/list"
    "os/exec"
    "crypto/md5"
    "encoding/hex"
    "bufio"
    "os"
    "io"
    "strings"
    "sync"
    "bytes"
    "encoding/gob"
)


type hashinfo struct {
  Fname string 
  Hash string
}


// The hash values of all the elements it gets is stored in the gob file
func hash_store(data []hashinfo) {

  //initialize a *bytes.Buffer
  m := new(bytes.Buffer) 
  //the *bytes.Buffer satisfies the io.Writer interface and can
  //be used in gob.NewEncoder() 
  enc := gob.NewEncoder(m)
  //gob.Encoder has method Encode that accepts data items as parameter
  enc.Encode(data)
  //the bytes.Buffer type has method Bytes() that returns type []byte, 
  //and can be used as a parameter in ioutil.WriteFile() 
  err := ioutil.WriteFile("hash_data", m.Bytes(), 0600) 
  if err != nil {
          panic(err)
  }
  fmt.Printf("just saved all hashinfo with %v\n", data)

}

//Function to check for changes 
func check_file_change(){
  // 1. Read hash_data gobject and iterate through all the md5 changes
  // 2. Check the md5 hashes against input object
  // 3. Return a array of files which have been changed
}



//Function takes a list of dependencies  and returns a topologically sorted list
func top_sort (flist list.List) {
   for e :=flist.Front(); e != nil; e = e.Next() {
    x := e.Value.(string)
    def := strings.Fields(x)
    dg := make(map[string][]string)

    // The topological Sort Code Begins Here
    if len(def) == 0 {
            continue // handle blank lines
        }
        lib := def[0]   // dependant (with an a) library
        list := dg[lib] // handle additional dependencies
        fmt.Println("Just printing list for heck os it %s",list)

  





  }
  //Topological Code Will End Here


}


//Tokenizing the String of config file into Target/Dependencies/Command
func cdelimit(r rune) bool {
	return  r == ':'
}

//Function will read file and print lines for each dependency and return list of all
// dependencies
func config_parse(file_name string) list.List{
	var dlist list.List
	file,err := os.Open(file_name)
	if err != nil {
   		fmt.Println("error")
  	}
	bf := bufio.NewReader(file)

	for {
		line, isPrefix, err := bf.ReadLine()
   		if err ==  io.EOF {
        	break
   		}
   		if err != nil {
        	fmt.Println(err)
        }	 
    	if isPrefix {
      		fmt.Println("isprefix error")
    	}
	  sline := string(line)
	  s := strings.FieldsFunc(sline,cdelimit)
 		dep := strings.Replace(s[0],"<-",",", -1)
		dep_list := strings.Split(dep,",")
		dlist.PushBack(dep_list)
   }
   return dlist 
}


//Gethash Function returns hash of a string
func getHash(filename string) (string, error) {
    bs, err := ioutil.ReadFile(filename)
    if err != nil {
        return "-1", err
    }
    h := md5.New()
    h.Write(bs)
    return hex.EncodeToString(h.Sum(nil)), nil
}


//Function to execute commands in parallel. 
func exe_cmd(cmd string, wg *sync.WaitGroup) {
	out, err := exec.Command(cmd).Output()
	if err != nil {
		fmt.Printf("%s", err)
	}
	fmt.Printf("%s", out)
	wg.Done() // Need to signal to waitgroup that this goroutine is done
}




func main(){


}