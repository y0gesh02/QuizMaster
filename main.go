package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"time"
)

//main -> open quizz file ->create new reader -> read the file ->problem parser with same as struct
//then loop through problems ->print ques->accept ans->calculate res;
// 30 sec timer

type problem struct {
	q string
	a string
}

func problemparser(lines [][]string)[]problem{ //returning problem type slice
	//parse the lines to slice
	r:=make([]problem,len(lines)) //type,len
	for i:=0; i<len(lines);i++{
		r[i]=problem{q: lines[i][0],a:lines[i][1]}
	}
	return r
}

func problemselect(filename string)([]problem, error){ 
	  if fObj,err:=os.Open(filename); err==nil{ //open csv
		csvR :=csv.NewReader(fObj) //create new reader
		if clines,err:=csvR.ReadAll(); err==nil{ //read the file
                 return problemparser(clines),nil //call parserfunc
		}else{
			return nil, fmt.Errorf("Error in reading data in csv" + " FileName: %s Error: %s",filename, err.Error())
		}
	  }else{
		return nil, fmt.Errorf("Error in opening data in csv" + " FileName: %s Error: %s",filename, err.Error())
	  }
}
func main(){
	fmt.Println("Welcome To Math Quiz **Timer 30Sec**");
	fmt.Println("Type S to start the quiz")
	var anss string
	fmt.Scan(&anss)
	//input the csv file
	fName:=flag.String("f","quiz.csv","path of csv file")
	//start timer
	timer:=flag.Int("t",10,"timer for the quiz")
	flag.Parse()
	//call problemselect
	problems,err:=problemselect(*fName)
	if err!=nil{
		exit(fmt.Sprintf("Something went wrong:%s",err.Error()))
	}
    correctAns:=0
	tObj := time.NewTimer(time.Duration(*timer)*time.Second)
	//fmt.Println(tObj);
	ansC:= make(chan string)

  problemloop:

	    for i, p :=range problems {
		
		     fmt.Printf("Problem %d:  %s=",i+1, p.q)
         	 var ans string
			 go func(){ //go routine
				fmt.Scan(&ans) //input ans
				ansC <- ans //giving to channel
			 }() //calling routing
		     //checking ans
			 select {
				case <-tObj.C: //timer complete
					fmt.Println()
					break problemloop
				case iAns := <-ansC:
		   			if iAns==p.a{
					 	correctAns++;
		            }
		   			if i==len(problems)-1{
						// fmt.Println("Time is Over!!!")
						close(ansC)
		   			}
			}
			

	}
	
	fmt.Printf("Your result is %d out of %d\n",correctAns,len(problems))
	//fmt.Printf("press enter to exit");
	<-ansC

}

func exit(msg string){
	fmt.Println(msg)
	os.Exit(1)
}

// 8+6,14
// 3+1,4
// 1+4,5
// 5+1,6
// 2+3,5
// 3+3,6
// 2+4,6
// 5+2,7