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
)



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