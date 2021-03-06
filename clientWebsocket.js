url = "ws://localhost:8080/echo" ;

w = new WebSocket(url);
w.onopen= function () {
  log("open web socket ");
  w.send("thanks for accepting WebSocket request");
}

w.onmessage =function (){
  log(e.data.toString());
}

w.onclose = function(){
log("closed");
}

w.onerror = function(){
  log("error")
}

window.onload = function() {
    document.getElementById("sendButton").onclick=function(){
      w.send(document.getElementById("inputMessage").value);
    }
}

function log(s){
  var logOutput = document.getElementById('logOutput');
  var el = $("#logOutput").after('<p>'+ s + '/<p>');
  $("p").slice(2).addClass("highlight");
}
