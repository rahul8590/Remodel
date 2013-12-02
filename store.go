package main 
import ( 
  "fmt" 
 "encoding/gob"
 "bytes"
 "io/ioutil"

)

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


func main() {
	org := map[string]string{"foo": "bar"}
	store(org , "store_data")
	var loadedMap map[string]string
	load(&loadedMap, "store_data")
	fmt.Println(loadedMap["foo"]) 
}