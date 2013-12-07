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
    "reflect"
    "strings"
    "sync"
    "bytes"
    "encoding/gob"
)

type depinfo struct {
  Dep []string 
  Cmd string
}


func store(data interface{} , filename string) {
  m := new(bytes.Buffer) 
  enc := gob.NewEncoder(m)

  err := enc.Encode(data)
  if err != nil { panic(err) }

  err = ioutil.WriteFile(filename, m.Bytes(), 0600) 
  if err != nil { panic(err) }
}


func load(e interface{} , filename string) {
    n,err := ioutil.ReadFile(filename)
    if err != nil { panic(err) } 

    p := bytes.NewBuffer(n) 
    dec := gob.NewDecoder(p)

    err = dec.Decode(e)
    if err != nil { panic(err) }
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


/*Function takes a list of dependencies  and returns a topologically sorted list
Input: [ [Default baz]
         [baz foo.o file.o]
      ]
Each element of the list is of type []string

Output: List which is topologically sorted and each of the element is of 
type []string
Note: Currently its only printing the output and not returning anything
*/
func topsort (dep_list list.List) list.List{
    var flist list.List
    dg := make(map[string][]string)
    //for _, line := range lines {
    for e := dep_list.Front(); e!= nil ; e = e.Next() {
        fmt.Println("printing e.value in main",e.Value)
        def := e.Value.([]string)    

        //def := strings.Fields(line)
        fmt.Printf("Def is %s \n",def)
        if len(def) == 0 {
            continue // handle blank lines
        }
        lib := def[0]   // dependant (with an a) library
        list := dg[lib] // handle additional dependencies

        

    
    scan:
        for _, pr := range def[1:] { // (pr for prerequisite)
            if pr == lib {
                continue // ignore self dependencies
            }
            for _, known := range list {
                if known == pr {
                    continue scan // ignore duplicate dependencies
                }
            }
            // build: this curious looking assignment establishess a node
            // for the prerequisite library if it doesn't already exist.
            dg[pr] = dg[pr]
            // build: add edge (dependency)
            list = append(list, pr)
        }
        // build: add or update node for dependant library
        dg[lib] = list
    }
 
    fmt.Printf("Dg is %s \n",dg)
    //fmt.Printf("list is %s \n",list)


    // topological sort on dg
    for len(dg) > 0 {
        // collect libs with no dependencies
        var zero []string
        for lib, deps := range dg {
            if len(deps) == 0 {
                zero = append(zero, lib)
                delete(dg, lib) // remove node (lib) from dg
            }
        }
        // cycle detection
        
        if len(zero) == 0 {
            fmt.Println("libraries with un-orderable dependencies:")
            // collect un-orderable dependencies
            cycle := make(map[string]bool)
            for _, deps := range dg {
                for _, dep := range deps {
                    cycle[dep] = true
                }
            }
            
            // print libs with un-orderable dependencies
            for lib, deps := range dg {
                if cycle[lib] {
                    fmt.Println(lib, deps)
                }
            }
            return flist
        }
        
 
        // remove edges (dependencies) from dg
        for _, remove := range zero {
            for lib, deps := range dg {
                for i, dep := range deps {
                    if dep == remove {
                        copy(deps[i:], deps[i+1:])
                        dg[lib] = deps[:len(deps)-1]
                        break
                    }
                }
            }
        }
        // output a set that can be processed concurrently
        fmt.Println("set is ",zero)
        flist.PushBack(zero)
    }
    return flist
}

//Function will read file and print lines for each dependency and return list of all
// dependencies
func config_parse(file_name string) (list.List,map[string]depinfo) {
  var dep_list list.List
  var dep = make(map[string]depinfo)
  var hashinfo = make(map[string]string)

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
    
    s := strings.Split(sline,"<-")
    s[0] = strings.TrimSpace(s[0])
    if (strings.ToLower(s[0]) == "default"){
      continue
    }

    sbuild := strings.Split(s[1],":")
    
    //Getting rid of unnecessary whitespaces to create dependencies
    sbuild[0] = strings.Replace(sbuild[0]," ","",-1)


    fmt.Println("Hash Elements are",s[0],reflect.TypeOf(s[0]))
    sfinal := strings.Split(sbuild[0],",")
    
    fmt.Println("Dependencies are  ",sfinal,reflect.TypeOf(sfinal))

    temp := depinfo {}

    if (len(sbuild) > 1) {
      cmd1 := strings.Replace(sbuild[1],`"`,"",-1)
      fmt.Println(sbuild[0],"commands is ",cmd1)  
      temp.Cmd = cmd1
    } else {
      temp.Cmd = "nothing"  
    }
 
    temp.Dep = sfinal
    dep[s[0]] = temp
    for i,_ := range sfinal {
      fmt.Println("i =>" , i , "s[i] => ",sfinal[i]  )  
    }
 
    dep_string := append(s[:1], sfinal...)
    fmt.Println("dep_String is ",dep_string,len(dep_string))
    dep_list.PushBack(dep_string)
    for _,v := range dep_string {
      fmt.Println("value to be sent for dep_string",v)
      a := getHash(v)
      hashinfo[v] = a
    }
   }
   fmt.Println("Elements in the dictonary are",dep)
   fmt.Println(" Hash info values are ",hashinfo)
   hash_store(hashinfo)
   return dep_list , dep   
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


//Executing Command in Series First
func exe_cmd(cmd string, wg *sync.WaitGroup) {
  fmt.Println("command is ",cmd)
  // splitting head => g++ parts => rest of the command
  parts := strings.Fields(cmd)
  head := parts[0]
  parts = parts[1:len(parts)]

  out, err := exec.Command(head,parts...).Output()
  if err != nil {
    fmt.Printf("%s", err)
  }
  fmt.Printf("%s", out)
  wg.Done() // Need to signal to waitgroup that this goroutine is done
}



func check(path string) (bool) {
    _, err := os.Stat(path)
    if err == nil { return true }
    if os.IsNotExist(err) { return false }
    return false
}



func main(){


}