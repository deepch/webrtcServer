package main

import (
	"fmt"
	"net/http"
	"io"
	"code.google.com/p/go.net/websocket"
	"container/list"
)

var connid int
var conns *list.List

func ChatroomServer(ws *websocket.Conn){

	defer ws.Close()

	connid++

	id := connid

	fmt.Printf("connection id: %d\n",id)

	item := conns.PushBack(ws)

	defer conns.Remove(item)

	name := fmt.Sprintf("user%d",id)
	
	SendMessage(nil,fmt.Sprintf("{\"user_name\":\"%s\"}",name))
	
//	r := bufio.NewReader(ws)

	for {
		
		var reply string;
		if err := websocket.Message.Receive(ws, &reply); err != nil {

            		fmt.Println("Can't receive")
           		 break
       		 }

        	fmt.Println("Received back from client: " + reply)


        	fmt.Printf("%s:%s\n",name,reply)

		//SendMessage(item,fmt.Sprintf("%s\t>%s",name,reply));
		  SendMessage(item,fmt.Sprintf("%s",reply));		
	}
}

func SendMessage(self *list.Element,data string){
	
	for item := conns.Front();item != nil ; item = item.Next() {
		ws,ok := item.Value.(*websocket.Conn)
		if !ok {
			panic("item not *websocket.Conn")
		}
	
		if item == self {
			continue
		}
		
		io.WriteString(ws,data)
	}
	


}

func Client(w http.ResponseWriter, r *http.Request){
html := `<!doctype html>
<html>
<head>
    <meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
    <title>golang websocket chatroom</title>
    <script src="https://ajax.googleapis.com/ajax/libs/jquery/1.7.2/jquery.min.js"></script>
    <script>
        var ws = new WebSocket("ws://127.0.0.1:6611/chatroom");
        ws.onopen = function(e){
            console.log("onopen");
            console.dir(e);
        };
        ws.onmessage = function(e){
            console.log("onmessage");
            console.dir(e);
            $('#log').append('<p>'+e.data+'<p>');
            $('#log').get(0).scrollTop = $('#log').get(0).scrollHeight;
        };
        ws.onclose = function(e){
            console.log("onclose");
            console.dir(e);
        };
        ws.onerror = function(e){
            console.log("onerror");
            console.dir(e);
        };
        $(function(){
            $('#msgform').submit(function(){
                ws.send($('#msg').val()+"\n");
                $('#log').append('<p style="color:red;">My > '+$('#msg').val()+'<p>');
                $('#log').get(0).scrollTop = $('#log').get(0).scrollHeight;
                $('#msg').val('');
                return false;
            });
        });
    </script>
</head>
<body>
    <div id="log" style="height: 300px;overflow-y: scroll;border: 1px solid #CCC;">
    </div>
    <div>
        <form id="msgform">
            <input type="text" id="msg" size="60" />
        </form>
    </div>
</body>
</html>`
    io.WriteString(w, html)

}

func main(){
    connid = 0;
    conns = list.New()
    http.Handle("/chatroom", websocket.Handler(ChatroomServer));
    http.HandleFunc("/", Client);
    err := http.ListenAndServe(":6611", nil);
    if err != nil {
        panic("ListenAndServe: " + err.Error())
    }

}
