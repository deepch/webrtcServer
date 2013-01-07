package webrtcSocket

import (
//	"bufio"
	"fmt"
	"net"
	"os"
        "io/ioutil"
        "net/http"
//        "net/url"
        "strings"
        "bytes"
)

var conn *net.TCPConn
var myPeerID string;
var targetUserName = "user@machenhui"
var targetPeerID string;
var (
		host    = "127.0.0.1"
		port    = "8888"
		remote  = host + ":" + port
                userName = "chenhui.ma";
		message = "GET /sign_in?"+userName+" HTTP/1.0\r\n\r\n"
	)
/*func Init() {

	
        addr,error := net.ResolveTCPAddr("tcp",remote)
	con, error := net.DialTCP("tcp", nil ,addr)
	conn = con

	if error != nil {

		fmt.Printf("Host not found: %s\n", error)
		os.Exit(1)
	}
	// defer con.Close();
        in, error := conn.Write([]byte(message))
	if error != nil {
		fmt.Printf("Error sending data: %s, in: %d\n", error, in)
		os.Exit(2)
	}
        hangingGet()
        for {

            handleSocket(conn);
       }
	
   
}*/



func Init() {

	siginIn()
	
        hangingGet()
        //sendToPeer("nihao",targetPeerID)
       
        fmt.Println("nihao",targetPeerID);
}

func hangingGet(){
    
      
    resp, err := http.Get("http://localhost:8888/wait?peer_id="+string(myPeerID));
	if err != nil{
           fmt.Println("feed error");
	}
    defer resp.Body.Close()
    body, err1 := ioutil.ReadAll(resp.Body)
	if err1 != nil{
           fmt.Println("body err1");
	}
   fmt.Println(string(body))
   go sendToPeer("nihao",targetPeerID)   
   
   hangingGet()
 
}




func siginIn(){

  
      
    resp, err := http.Get("http://localhost:8888/sign_in?"+userName);
	if err != nil{
           fmt.Println("feed error");
	}
    defer resp.Body.Close()
    body, err1 := ioutil.ReadAll(resp.Body)
	if err1 != nil{
           fmt.Println("body err1");
	}
   
   bA := strings.Split(string(body),"\n");
   fmt.Println(string(body));
   for i:=len(bA)-1;i>=0;i-- {
      
      bt := strings.Split(string(bA[i]),",");
      if len(bt) >0 && bt[0] == userName {
         fmt.Println(bt[1])
         myPeerID = bt[1]
      }
      
      if len(bt)>0 && bt[0] == targetUserName {

        targetPeerID = bt[1]
        fmt.Println(bt[1]);
     }
   }
 
   //fmt.Println(bA[0],len(bA))

}


func siginOut(){
  
      
    resp, err := http.Get("http://localhost:8888/sign_out?peer_id="+myPeerID);
	if err != nil{
           fmt.Println("feed error");
	}
    defer resp.Body.Close()
    body, err1 := ioutil.ReadAll(resp.Body)
	if err1 != nil{
           fmt.Println("body err1");
	}
   
   bA := strings.Split(string(body),"\n");
   
   for i:=len(bA)-1;i>=0;i-- {
      
      bt := strings.Split(string(bA[i]),",");
      if len(bt) >0 && bt[0] == userName {
         //fmt.Println(bt[1])
         myPeerID = bt[1]
      }
      
      if len(bt)>0 && bt[0] == targetUserName {

        targetPeerID = bt[1]
        fmt.Println(bt[1]);
     }
   }
 
   //fmt.Println(bA[0],len(bA))

}


func handleSocket(conn *net.TCPConn) {
	
	//fmt.Fprintf(con, "GET / HTTP/1.0\r\n\r\n")
	/*status, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {

		fmt.Println(err)
	}
	fmt.Println(status)
        fmt.Println("Connection OK")*/

        var buf = make([]byte,12000)
        position, err := conn.Read(buf)
        if err != nil {
          return
       }
       fmt.Println(string(buf[0:position]))
}

func sendToPeer(data string,peer_id string){
    fmt.Println("sendToPeer")
    param := bytes.NewReader([]byte(data));
    resp, err := http.Post("http://localhost:8888/message?peer_id="+myPeerID+"&to="+peer_id,"text/plain",param);
	if err != nil{
           fmt.Println("feed error");
	}
    fmt.Println("send over")
    defer resp.Body.Close()
    body, err1 := ioutil.ReadAll(resp.Body)
	if err1 != nil{
           fmt.Println("body err1");
	}
    fmt.Println(string(body))
    fmt.Println("peer end");


}

func SendBuffer(buf []byte) string {
        fmt.Print("+");
        sendToPeer(string(buf),"1");
        return "";
	in, error := conn.Write(buf)
	if error != nil {
		fmt.Printf("Error sending data: %s, in: %d\n", error, in)
		os.Exit(2)
	}

	
	/*
	  status, err := bufio.NewReader(conn).ReadString('\n')
	  if err != nil {

		fmt.Println(err)
	  }
	  fmt.Println(status)
          fmt.Println("Connection Send OK")
        */

        var bufR = make([]byte,12000)
        position, err := conn.Read(bufR)
        if err != nil {
          fmt.Println("error===",error);
          return ""
       }
       


        return string(bufR[0:position]);
}
