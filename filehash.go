package main

import (
    "fmt"
    //"hash/crc32"
    "io/ioutil"
    "container/list"
    "os/exec"
    "log"
    "crypto/md5"
    "encoding/hex"
)

func exe_cmd(cmd string) int {
	_, err := exec.Command("date").Output()
    if err != nil {
        log.Fatal(err)
    }
    // fmt.Printf("%s",string(out))    Not printing the compilation step
    return 0
}


func getHash(filename string) (string, error) {
    bs, err := ioutil.ReadFile(filename)
    if err != nil {
        return "0", err
    }
    h := md5.New()
    h.Write(bs)
    return hex.EncodeToString(h.Sum(nil)), nil
}

func main() {

    var flist list.List
    flist.PushBack("test1.txt")
    flist.PushBack("test2.txt")
    dict := make(map[string]string)
     
    for e :=flist.Front(); e != nil; e = e.Next() {
    	h, err := getHash(e.Value.(string))
    	if err != nil {
        	return
    	}
    	dict[e.Value.(string)] = h
    }

    fmt.Println(dict)
}
