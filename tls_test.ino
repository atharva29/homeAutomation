#include <ESP8266WiFi.h>
#include <ESP8266WiFiAP.h>
#include <ESP8266WiFiGeneric.h>
#include <ESP8266WiFiMulti.h>
#include <ESP8266WiFiScan.h>
#include <ESP8266WiFiSTA.h>
#include <ESP8266WiFiType.h>
#include <WiFiClient.h>
#include <WiFiClientSecure.h>
#include <WiFiServer.h>
#include <WiFiUdp.h>

WiFiClientSecure client1;
//const char* fingerprint = "AF 4C CF 05 B1 68 C6 12 B6 8F A1 EB B3 A8 08 3D 35 92 82 49";
byte ip[] = {192, 168, 43, 100};
int port = 8000;

void setup() {
  Serial.begin(115200);
  WiFi.mode(WIFI_STA);
  WiFi.begin("hot", "123123123" );
  while (WiFi.status() != WL_CONNECTED)
  { Serial.println("..");
    delay(500);
  }
  Serial.println("connected to wifi");
  while (!client1.connect(ip, port))
  {
    delay(500);
  }
//  if (client1.verify(fingerprint, "sharad")) {
//    Serial.println("certificate matches");
//    client1.println("connected successfully");
//    while (client1.connected()) {
//      String line = client1.readStringUntil('\n');
//      Serial.println(line);
//      break;
//    }
//  }
//  else {
//    Serial.println("certificate doesn't match");
//  }

}

void loop() {
  if (Serial.available()){
  String reader = Serial.readString();
  client1.println("tranferring data"+ reader);
  }
  }
