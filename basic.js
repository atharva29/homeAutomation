var connection = new WebSocket("ws://localhost:8080/echo");

//document.write('Hello');
document.getElementById('section').innerHtml = "Hello World"
connection.onopen = function(){
   /*Send a small message to the console once the connection is established */
   console.log('Connection open!');
}

connection.onclose = function(){
//  document.write('Refused connection by server');
   console.log('Connection closed');
}

connection.onerror = function(error){
   console.log('Error detected: ' + error);
}

connection.onmessage = function(e){
   var server_message = e.data;
   console.log(server_message);
   document.write(server_message);
//   document.getElementById('section').innerHtml = server_message
 }
