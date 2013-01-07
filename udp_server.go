package main

import (
    "fmt"
    "net"
    "os"
    "bytes"
    "xn"
     "xn/webrtcSocket"
//    "container/list"
)

var ip54,err1 = net.ResolveUDPAddr("udp", "10.2.42.213:11121")
var ip213,err2 = net.ResolveUDPAddr("udp", "10.2.42.197:11121")

var udpConnMap = make(map[string]*net.UDPConn,10000)
//当另一个终端还没有连接上时，缓存视频
var TempChanMapIndex = make(map[string]int,10000)
var TempChanMap = make(map[string][][]byte,10000)

func main() {
   //测试手动注册两个通道
    sip.GetChannel("10.2.42.197","10.2.42.133");
    sip.GetChannel("10.2.42.133","10.2.42.197");
    //链接WebRTC 服务
    go webrtcSocket.Init();
    //开启出测地址
    go sip.InitListener();
   
    //监听视频端口
    service := ":8088"
    udpAddr, err := net.ResolveUDPAddr("udp", service)
    checkError(err)
    conn, err := net.ListenUDP("udp", udpAddr)
    checkError(err)
    fmt.Println("监听端口");
    for {
        	handleVideo(conn)
    	}
   
   //监听语音端口
   go func (){ 
   serviceAudio :=":11113"
   udpAddrAudio,err := net.ResolveUDPAddr("udp",serviceAudio)
   checkError(err)
   connAudio,err := net.ListenUDP("udp",udpAddrAudio)
   checkError(err)

    for {
	handleAudio(connAudio)
     }
    }();

}

func handleAudio(conn *net.UDPConn){

    var buf = make([]byte,12000)
    position, addr, err := conn.ReadFromUDP(buf)
    if err != nil {
        return
    }
    
  if bytes.Compare(addr.IP,ip54.IP) == 0 {
    	addr.Port = 11113
    	addr.IP = ip213.IP
       
	conn.WriteToUDP(buf[0:position], addr)
  
    }else if bytes.Compare(addr.IP,ip213.IP) == 0 {
    	addr.Port = 11113
    	addr.IP = ip54.IP
        
	conn.WriteToUDP(buf[0:position], addr)
    }
    

}

//根据IP 获得其发送地址，接收地址
func getUDPAddr(ip string)(*net.UDPAddr,*net.UDPAddr){
    service := ip+":8088"
    udpAddrSend, err := net.ResolveUDPAddr("udp", service)
    checkError(err)
    
    serviceReceive := ip+":8808"
    udpAddrReceive, err := net.ResolveUDPAddr("udp", serviceReceive)
    checkError(err)
 
   return udpAddrSend,udpAddrReceive

}

//发出缓存
func fushCache(conn *net.UDPConn,sendAddr *net.UDPAddr){
    cacheArray := TempChanMap[sendAddr.IP.String()];

/*    for item := cacheArray.Front();item != nil ; item = item.Next() {
                fmt.Print("-");				
		conn.WriteToUDP(item.Value.([]byte), sendAddr)                
    }
*/
     length := TempChanMapIndex[sendAddr.IP.String()];
     fmt.Println(length,sendAddr.IP.String());
     for i := 0;i< length; i++{
                if len(cacheArray[i])>0 {
                	fmt.Print("-");				
			position,err := conn.WriteToUDP(cacheArray[i], sendAddr) 
 
         		checkError(err)
         		if err != nil {
           			fmt.Println(position);
         		}
                }               
    }
/*
   for item := cacheArray.Front();item != nil ; item = item.Next() {
      cacheArray.Remove(item)
      fmt.Println("remove");
   }
 */   
  TempChanMap[sendAddr.IP.String()] = make([][]byte,1000,1000000);
  TempChanMapIndex[sendAddr.IP.String()] = 0;
  //fmt.Println(TempChanMap,sendAddr.IP.String());
  
}
//写入缓存
func addCache(data []byte,ip string){
    if ip == "" || len(data) == 0 {
       return ;
    }
    
    cacheArray := TempChanMap[ip];   
    //cacheArray.PushBack(data);
    length := len(cacheArray);
    if length == 0 {
	cacheArray = make([][]byte,1000,1000000)
    }
    cacheArray[TempChanMapIndex[ip]] = data;
    TempChanMap[ip] = cacheArray;
    TempChanMapIndex[ip] +=1; 
    //fmt.Println(TempChanMap);
}

func handleVideo(conn *net.UDPConn) {
    
    var buf = make([]byte,12000)
    position, addr, err := conn.ReadFromUDP(buf)
    if err != nil {
        return
    }
    if addr.IP.String() == "10.2.42.197" { 
       go webrtcSocket.SendBuffer(buf[0:position]);
    }
    
    toIP,_ := sip.GetToIPByFromIP(addr.IP.String());
    //currentConn := udpConnMap[addr.IP.String()];
    //toIPConn := udpConnMap[toIP];
    _,receiveAddr := getUDPAddr(toIP);
    //_,receiveAddr1 := getUDPAddr(addr.IP.String());
    /*if currentConn == nil {
        fmt.Println("清空缓存======",addr.IP.String(),"上线",receiveAddr1)
        
	fushCache(conn,receiveAddr1)
        
    }*/
    udpConnMap[addr.IP.String()] = conn;
    /*if toIPConn == nil {
        fmt.Println("======添加缓存",toIP,addr.IP.String())
	addCache(buf[0:position],toIP)
         
    }else {
       
	  position,err := conn.WriteToUDP(buf[0:position], receiveAddr)
          checkError(err)
          if err != nil {
           fmt.Println(position);
          }
	
       
   }*/
   
         pt,err1 := conn.WriteToUDP(buf[0:position], receiveAddr);
          checkError(err1)
          if err != nil {
           fmt.Println(pt);
          }
   

/*
    if bytes.Compare(addr.IP,ip54.IP) == 0 {
    	//addr.Port = 11111
    	//addr.IP = ip213.IP
    
	//conn.WriteToUDP(buf[0:position], addr)
	conn.WriteToUDP(buf[0:position], ip213)
        
    }else if bytes.Compare(addr.IP,ip213.IP) == 0 {
    	//addr.Port = 11111
    	//addr.IP = ip54.IP
     
	//conn.WriteToUDP(buf[0:position], addr)
	conn.WriteToUDP(buf[0:position], ip54)
    }
  */  
}


func checkError(err error) {
    if err != nil {
        fmt.Fprintf(os.Stderr, "Fatal error ", err.Error())
        os.Exit(1)
    }
}
