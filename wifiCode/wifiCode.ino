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
#include <EEPROM.h>

int count = 0 ;
String ssid = "coecnds" ;
String password = "jetsonTX2@coe";
String ip = "192.168.0.122" ;
int port = 8000 ;

WiFiClient client1;

char* string2char(String command) {
  if (command.length() != 0)
  {
    char*p = const_cast<char*>(command.c_str());
    return p;
  }
}

void clientMode() {
  if (count == 0 ) {
    WiFi.mode(WIFI_STA);
    WiFi.begin(string2char(ssid), string2char(password) );
    while (WiFi.status() != WL_CONNECTED)
    {
      delay(250);
    }
    delay(100);
    while (!client1.connect(string2char(ip), port))
    {
      Serial.println("Connecting to server");
      delay(500);
    }
    count = 1 ;
  } else if (count == 1 ) {
    if (client1.connected()) {
      if (client1.available() > 0 ) {
        String Command = client1.readStringUntil('\n');
        Serial.println(Command);
      } else {
        client1.println("Sending Data");
      }
    } else {
      count = 0 ;
    }
  }
}


void setup() {
  Serial.begin(115200);
}

void loop() {
  clientMode();
}
