// GET THE SUNSET AND SUNRISE TIMES FROM RUNNING
//   go run sunset.go times
// AND COPY AND PASTE HERE

unsigned long currentTime;
unsigned long timeTracker;

void setup()
  {
    
  Serial.begin (9600);
  currentTime = millis()/1000;
  timeTracker = 0;
  }
  
void loop() 
 {

  delay(1000);
  if (millis()/1000+10000+timeTracker < currentTime) {
    timeTracker = currentTime; // the close has reset, carry it forward
  }
  currentTime = millis()/1000+timeTracker;
  int activateClose = 0;
  for (int i=0; i < sizeof(sunset)/sizeof(sunset[0]); i++) {
    if (currentTime > sunset[i]) {
      activateClose = 1;
      break;
    }
  }
  if (activateClose == 1) {
    sunsetActivation();
  }

  Serial.print(currentTime);
  Serial.print(" sensor:");
  int sensorValue = analogRead(1);
  Serial.println(sensorValue);
 }


void sunsetActivation() {
  Serial.println("Sunset activated!");
  for (int i=0; i < sizeof(sunset)/sizeof(sunset[0])-1; i++) {
    sunset[i] = sunset[i+1];
  } 
}
