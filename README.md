"# attendanceSystem" 
It contains 3 main codes
tcpClient.go := this code acts as a client to server , connects to server and sends comma separated data in form of ID,Name.
and connnects to 2000 port of server .When server disconnects it again tries to connect with server

tcpServer := this code acts as a server for sensor and collects data from port 2000 ,and acts as a client to user's computer which will be running tcpBadaServer.go (data is sent to port 6600). Thus now we will be recieving data from port 2000 and putting them in SQLite database and when user's computer makes query to database and server gives output and sends it in form of string .

tcpBadaServer := this code will run in user's computer and acts as a server to tcpServer.go code at port 6600 , and makes query to database and recieves response


NETCAT CODE - Currently i am using netcat on client side. NetCat connects to server and sends data in form 1,ard1
in which 1 will be the roll number of the person and ard1 is the device name .


DATABASE 
We have 1 database whose name is attendance.db and it contains 4 tables 
1. student
2.student1 
3.student2 
4.counter 
// i will replace this names by SE,TE,BE 

student TABLE = In this tables there are 3 rows ,
i. num = num is the primary key for database
ii.id = id is the integer for roll number 
iii.name = name is for CLASS (eg. SE , TE , BE )
iv. date_time = puts current date time in database

Similarly for student1 ,student2 there are num1 , num2 and ID1 and ID2 respectively 

Counter table contains 4 colums :
1.num = this num is always 1 
 ID1 , ID2  , ID3    are the respective num values of each repective tables 
2.ID1 = holds num value of student 
3. ID2 = holds num value of student1
4. ID3 = holds num value of student2
