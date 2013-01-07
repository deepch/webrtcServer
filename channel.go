package sip

import (
   "fmt"
   "time"
   "math/rand"
    "net/http"
)

type Info struct{
     fromIP,toIP string
}

type ClientInfo struct{
     IP  string
     UID string
}

var channelMap = make(map[string]Info,100)
var clientInfoMap = make(map[string]ClientInfo,100)

func GetChannel(fromIP string,toIP string)  string {
    info := Info {
         fromIP,toIP,
     }
    r := rand.New(rand.NewSource(time.Now().UnixNano()))
    channelID := fmt.Sprintf("b%d",r.Intn(1000000000));
    channelMap[channelID] = info;
    return channelID;
}

func GetToIPByFromIP(fromIP string) (string,string){
     
   for key,value := range channelMap {
      
      if value.fromIP == fromIP {
        return value.toIP,key
     }
    
   } 
   return "",""
}

func RemoveChannel(channelID string){

    delete(channelMap,channelID)
}

/**
 * 注册客户端信息，并返回所有在注册的客户端信息列表
 */
func handleRegisterClient(w http.ResponseWriter, r *http.Request) {
     clientInfoMap[r.RemoteAddr] = ClientInfo{
	  	r.RemoteAddr,"bbbbbb",
    }

}

func InitListener(){
 
    http.HandleFunc("/register", handleRegisterClient);
    err := http.ListenAndServe(":8080", nil);
    if err != nil {
        panic("ListenAndServe: " + err.Error())
    }

}
