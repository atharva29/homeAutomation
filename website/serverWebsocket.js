var WebsocketServer = require('ws').Server ;

var wss = new WebsocketServer({
  port:8080
});

wss.on('connection',function(ws){
  ws.on('message',function(message){
    console.log('Recieved %s',message) ;
    ws.send('from server -- '+message)
  }) ;
});
