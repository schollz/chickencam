// GET THE SUNSET AND SUNRISE TIMES FROM RUNNING
//   go run sunset.go times
// AND COPY AND PASTE HERE
unsigned long sunset[50] = {5,95252,181575,267898,354222,440547,526874,613201,699529,785859,872189,958521,1044854,1131188,1217524,1303860,1390198,1476538,1562878,1649220,1735564,1821909,1908255,1994603,2080952,2167303,2253656,2340010,2426366,2512723,2599082,2685443,2771805,2858170,2944536,3030903,3117273,3203644,3290017,3376392,3462769,3549148,3635528,3721911,3808295,3894682,3981070,4067460,4153852,4240246};
unsigned long sunrise[49] = {80,137621,224075,310529,396984,483439,569894,656350,742807,829264,915721,1002178,1088636,1175094,1261553,1348012,1434471,1520931,1607391,1693851,1780311,1866772,1953232,2039693,2126154,2212616,2299077,2385538,2472000,2558461,2644922,2731383,2817844,2904305,2990766,3077226,3163686,3250146,3336605,3423064,3509523,3595980,3682438,3768895,3855351,3941806,4028260,4114714,4201167};

unsigned long currentTime;
unsigned long timeTracker;

int const OPENTIME = 50;
void setup()
  {
    
  Serial.begin (9600);
  pinMode(2,OUTPUT);
  pinMode(3,OUTPUT);
  digitalWrite(2, LOW);
  digitalWrite(3, LOW);
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

  int activateOpen = 0;
  for (int i=0; i < sizeof(sunrise)/sizeof(sunrise[0]); i++) {
    if (currentTime > sunrise[i]) {
      activateOpen = 1;
      break;
    }
  }
  if (activateOpen == 1) {
    sunriseActivation();
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
  digitalWrite(2, LOW);
  digitalWrite(3, LOW);
  
  digitalWrite(2, LOW);
  digitalWrite(3, HIGH); // CLOSE
  for (int i=0; i < OPENTIME; i++) {
    delay(1000);
  } 
  digitalWrite(2, LOW);
  digitalWrite(3, LOW);
}


void sunriseActivation() {
  Serial.println("Sunrise activated!");
  for (int i=0; i < sizeof(sunrise)/sizeof(sunrise[0])-1; i++) {
    sunrise[i] = sunrise[i+1];
  } 
  digitalWrite(2, LOW);
  digitalWrite(3, LOW);
  
  digitalWrite(2, HIGH);
  digitalWrite(3, LOW); // OPEN
  for (int i=0; i < OPENTIME; i++) {
    delay(1000);
  } 
  digitalWrite(2, LOW);
  digitalWrite(3, LOW);
}
