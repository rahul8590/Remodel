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

type depinfo struct {
  Dep []string 
  Cmd string
}


// Retrieves the hash_values stored as a gobject hash_store
func hash_ret (gobject string ) map[string]string{
  n,err := ioutil.ReadFile(gobject)
        if err != nil {
                fmt.Printf("cannot read file")
                panic(err)
        } 
        //create a bytes.Buffer type with n, type []byte
        p := bytes.NewBuffer(n) 
        //bytes.Buffer satisfies the interface for io.Writer and can be used
        //in gob.NewDecoder() 
        dec := gob.NewDecoder(p)
        //make a map reference type that we'll populate with the decoded gob 
        //e := make(map[int]string)
         e := make(map[string]string)
        //we must decode into a pointer, so we'll take the address of e 
        err = dec.Decode(&e)
        if err != nil {
                fmt.Printf("cannot decode")
                panic(err)
        }

        fmt.Println("after reading dep_data printing ",e)
        return e

}

// The hash values of all the elements it gets is stored in the gob file
func hash_store(data map[string]string) {
  //initialize a *bytes.Buffer
  m := new(bytes.Buffer) 
  //the *bytes.Buffer satisfies the io.Writer interface and can
  //be used in gob.NewEncoder() 
  enc := gob.NewEncoder(m)
  //gob.Encoder has method Encode that accepts data items as parameter
  enc.Encode(data)
  //the bytes.Buffer type has method Bytes() that returns type []byte, 
  //and can be used as a parameter in ioutil.WriteFile() 
  err := ioutil.WriteFile("dep_data", m.Bytes(), 0600) 
  if err != nil {
          panic(err)

  }
  fmt.Printf("just saved all depinfo with %v\n", data)
}

//Function to check for changes 
func check_file_change(){
  // 1. Read dep_data gobject and iterate through all the md5 changes
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


/* This function is not required anymore

Tokenizing the String of config file into Target/Dependencies/Command
func cdelimit(r rune) bool {
	return  r == ':'
}
*/


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
func getHash(filename string) (string) {
    fmt.Println("filename to be read is ",filename)
    filename = strings.TrimSpace(filename)
    bs, err := ioutil.ReadFile(filename)
    if err != nil {
        fmt.Println(filename,err)
        return "-1"
    }
    h := md5.New()
    h.Write(bs)
    fmt.Println("file value is ",bs)
    return hex.EncodeToString(h.Sum(nil))
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