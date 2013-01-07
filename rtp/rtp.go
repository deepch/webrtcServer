package main

import (
   "net"
   "github.com/wernerd/GoRTP/src/net/rtp"
 //  "fmt"
   "time"
) 

var localPort = 8088
var local,_ = net.ResolveIPAddr("ip","127.0.0.1");
var stopLocalRecv = make(chan bool, 1)
var rsLocal *rtp.Session
func receivePacketLocal() {
    // Create and store the data receive channel.
    dataReceiver := rsLocal.CreateDataReceiveChan()
    var cnt int

    for {
        select {
        case rp := <-dataReceiver:
            if (cnt % 50) == 0 {
                println("Remote receiver got:", cnt, "packets")
            }
            cnt++
            rp.FreePacket()
        case <-stopLocalRecv:
            return
        }
    }
}

func main(){
    tpLocal, _ := rtp.NewTransportUDP(local, localPort)
    rsLocal = rtp.NewSession(tpLocal, tpLocal)
    
   rsLocal.ListenOnTransports()
   go receivePacketLocal()
  
   time.Sleep(8e9)
}
