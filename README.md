"# attendanceSystem" 
It contains 3 main codes
tcpClient.go := this code acts as a client to server , connects to server and sends comma separated data in form of ID,Name.
and connnects to 2000 port of server .When server disconnects it again tries to connect with server

tcpServer := this code acts as a server for sensor and collects data from port 2000 ,and acts as a client to user's computer which will be running tcpBadaServer.go (data is sent to port 6600). Thus now we will be recieving data from port 2000 and putting them in SQLite database and when user's computer makes query to database and server gives output and sends it in form of string .

tcpBadaServer := this code will run in user's computer and acts as a server to tcpServer.go code at port 6600 , and makes query to database and recieves response
