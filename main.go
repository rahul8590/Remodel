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
    //"reflect"
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

func getHash(filename string) (string) {
    //fmt.Println("filename to be read is ",filename)
    filename = strings.TrimSpace(filename)
    bs, err := ioutil.ReadFile(filename)

    if err != nil {
        fmt.Println(filename,err)
        return "-1"
    }
    h := md5.New()
    h.Write(bs)
    //fmt.Println("file value is ",bs)
    return hex.EncodeToString(h.Sum(nil))
}

func check(path string) (bool) {
    _, err := os.Stat(path)
    if err == nil { return true }
    if os.IsNotExist(err) { return false }
    return false
}


func config_parse(file_name string) (list.List,map[string]depinfo,string,string) {
  var dep_list list.List
  var dep = make(map[string]depinfo)
  var build,status string
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
    // Returning Default Build
    if (strings.ToLower(s[0]) == "default"){
      build = s[1]
      continue
    }

    sbuild := strings.Split(s[1],":")
    
    //Getting rid of unnecessary whitespaces to create dependencies
    sbuild[0] = strings.Replace(sbuild[0]," ","",-1)


    //fmt.Println("Hash Elements are",s[0],reflect.TypeOf(s[0]))
    sfinal := strings.Split(sbuild[0],",")
    
    //fmt.Println("Dependencies are  ",sfinal,reflect.TypeOf(sfinal))

    temp := depinfo {}

    if (len(sbuild) > 1) {
      cmd1 := strings.Replace(sbuild[1],`"`,"",-1)
      //fmt.Println(sbuild[0],"commands is ",cmd1)  
      temp.Cmd = cmd1
    } else {
      temp.Cmd = "nothing"  
    }
 
    temp.Dep = sfinal
    dep[s[0]] = temp
    /*for i,_ := range sfinal {
      //fmt.Println("i =>" , i , "s[i] => ",sfinal[i]  )  
    }*/
 
    dep_string := append(s[:1], sfinal...)
    //fmt.Println("dep_String is ",dep_string,len(dep_string))
    dep_list.PushBack(dep_string)
    for _,v := range dep_string {
      //fmt.Println("value to be sent for dep_string",v)
      a := getHash(v)
      hashinfo[v] = a
    }
   }
   //fmt.Println("Elements in the dictonary are",dep)
   //fmt.Println(" Hash info values are ",hashinfo)
   //hash_store(hashinfo)
   
   if (check(".remodel/hash_data") == false) {
    store(hashinfo,".remodel/hash_data")
    status = "1"
   } else {
    fmt.Println("hash_data already exist")
    status = "0"
   }
   return dep_list, dep, build, status  
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
        //fmt.Println("printing e.value in main",e.Value)
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
 
    //fmt.Printf("Dg is %s \n",dg)

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
            fmt.Println("The are cyclic dependencies:")
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
            fmt.Println("Cannot Proceed :(")
            os.Exit(0)
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



func main() {
  
  //Declarations
  var dep1_list list.List
  var exe_list list.List
  var dep1 = make(map[string]depinfo)
  var build, status string
  wg := new(sync.WaitGroup)
  
  if (check(".remodel") == false) {
    a := os.Mkdir(".remodel",0755)
    fmt.Println("making .remodel",a) 
  } 
  
  /*
  dep1_list  => topsorted list of the objects need to be built 
  dep1  => dependency list based on depinfo structure
  build => default build (if argument not supplied)
  status => 0 hash_data exist ; 1 does not and config_parse dumped it
  */
  dep1_list, dep1, build, status = config_parse("config")
  //fmt.Println(dep1_list,dep1)
  

  // Loading Previous hash_values
  var prev_hash  map[string]string
  load(&prev_hash, ".remodel/hash_data")
  //fmt.Println("The previous hash data is ",prev_hash)

  fmt.Println("==========Executing Topsort ==============")
  flist := topsort(dep1_list)
  fmt.Println("=========End of Execution ================")
  
  if (len(os.Args) > 1) {
    build = os.Args[1]
  }

  for e := flist.Front(); e!= nil ; e = e.Next() {
    //fmt.Println("printing e.value in main",e.Value)
    //fmt.Println(reflect.TypeOf(e.Value))
    
    name := e.Value.([]string) //Go needs freking type assertions damn
    for _,v := range name {
      //fmt.Println(i," ",v)
      name_info := dep1[v]
      if (name_info.Dep == nil){
        //fmt.Println("name_info is empty for",v)
      } else {
        //fmt.Println("the dependency for the", v ,"is:",name_info.Dep)
        for _,dep_v := range name_info.Dep{
          //fmt.Println("chumma printing dep_v",dep_v , "Cmd =>" , name_info.Cmd)
          phash := strings.TrimSpace(prev_hash[dep_v])
          chash := strings.TrimSpace(getHash(dep_v))
          if ( phash != chash || status == "1") {
            //fmt.Println("prev_hash =>[",phash,"]Current Hash =>[",chash,"]dep_v =>", dep_v)
            fmt.Println(name_info.Cmd)
            prev_hash[dep_v] = chash  // Resetting the hash_value to current value
            
            wg.Add(1)
            go exe_cmd(name_info.Cmd,wg)
                
            exe_list.PushBack(name_info.Cmd)  
            break
          } else {
            fmt.Println("hash values are same for", dep_v)
          }
        
        }
        wg.Wait()
      // Added for Custom Builds as soon as it matches, it will exit
      if (build == v) {
        fmt.Println("build criteria is met build =>",build)
        goto last_step
        }
      }
    }    
  }
  last_step:
    
    //Checkign for Root argument. If root argument is given by user, it needs to be built
    //fmt.Println("flagvar has type ", reflect.TypeOf(build) ,"Build => ", build)
    store(prev_hash,".remodel/hash_data")
    fmt.Println("re-storing the latest hash_object in hash_data")
}
