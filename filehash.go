package main

import (
    "fmt"
    "hash/crc32"
    "io/ioutil"
    "container/list"
    "os/exec"
    "log"
)

func exe_cmd(cmd string) int {
	_, err := exec.Command("date").Output()
    if err != nil {
        log.Fatal(err)
    }
    // fmt.Printf("%s",string(out))    Not printing the compilation step
    return 0
}


func getHash(filename string) (uint32, error) {
    bs, err := ioutil.ReadFile(filename)
    if err != nil {
        return 0, err
    }
    h := crc32.NewIEEE()
    h.Write(bs)
    return h.Sum32(), nil
}

func main() {

    var flist list.List
    flist.PushBack("test1.txt")
    flist.PushBack("test2.txt")
    dict := make(map[string]uint32)
     
    for e :=flist.Front(); e != nil; e = e.Next() {
    	h, err := getHash(e.Value.(string))
    	if err != nil {
        	return
    	}
    	dict[e.Value.(string)] = h
    }

    fmt.Println(dict)
}
