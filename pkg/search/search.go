package search

import (
	"io/ioutil"
	"strings"
	"log"
	"context"
	//"os"
	//"bufio"
	"sync"
)

//Result search result
type Result struct {
	Phrase string
	Line string 
	LineNum int64
	ColNum int64
}

// All search phrase in files 
func All(ctx context.Context, phrase string, files []string) <-chan []Result {  	
	ch := make(chan []Result)
	wg := sync.WaitGroup{}
	root := context.Background()
	ctx, cancel := context.WithCancel(root)

	for i := 0; i < len(files); i++ {
		wg.Add(1)
		go func (ctx context.Context, fileName string, index int, ch chan<- []Result) { 
			defer wg.Done()
			result := []Result{}

			//reading file and send phrases to channel
			file, err := ioutil.ReadFile(fileName)
			if err != nil {
				log.Print("Error read file:", err)
				return
			}
			  
			rows := strings.Split(string(file), "\n") 		
			for i, item := range rows {
				if strings.Contains(item, phrase) {
					items := Result{
						Phrase: phrase,
						Line: item,
						LineNum: int64(i+1),
						ColNum: int64(strings.Index(item, phrase))+1,
					}
					result = append(result, items)					
				}						
			}
			
			if len(result) > 0 {
				ch <- result
			}			
		}(ctx, files[i], i, ch) 
	}
	go func() {
		defer close(ch)
		wg.Wait()
	}() 
	//val := <- ch
	//log.Print(val)
	cancel()	
	return ch
}


