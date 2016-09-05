import logging
# set up logging to file - see previous section for more details
logging.basicConfig(level=logging.DEBUG,
                    format='%(asctime)s %(name)-12s %(levelname)-8s %(message)s',
                    datefmt='%m-%d %H:%M',
                    filename='/tmp/henpi.log',
                    filemode='w')
# define a Handler which writes INFO messages or higher to the sys.stderr
console = logging.StreamHandler()
console.setLevel(logging.DEBUG)
# set a format which is simpler for console use
formatter = logging.Formatter('%(name)-12s: %(levelname)-8s %(message)s')
# tell the handler to use this format
console.setFormatter(formatter)
# add the handler to the root logger
logging.getLogger('').addHandler(console)

import sys
import subprocess
import time
import datetime
import threading
import os


try:
    import pyaudio
except:
    logging.warn("sudo apt-get install python3-pyaudio")
import audioop
import numpy
import wave
try:
    from PIL import Image
except:
    logging.warn("sudo apt-get install python3-pil")
try:
    from StringIO import StringIO
except ImportError:
    from io import StringIO, BytesIO
try:
    import RPi.GPIO as GPIO
    GPIO.setmode(GPIO.BCM)
    RELAY = 17  # https://pinout.xyz/pinout/pin11_gpio17
    RELAY_ON = 0
    RELAY_OFF = 1
    GPIO.setup(RELAY, GPIO.OUT)
    GPIO.output(RELAY, RELAY_OFF)
except:
    logging.warn("sudo apt-get install python3-rpio.gpio")


TEST_IMAGE_SIZE = (400, 300)
OFF = False
ON = True
# Set default mic using
# http://raspberrypi.stackexchange.com/questions/37177/best-way-to-setup-usb-mic-as-system-default-on-raspbian-jessie


def getAudioLevel():
    logger = logging.getLogger('getAudioLevel')
    chunk = 2048
    rms = []
    logger.debug("Calculating audio level")
    for i in range(0, 10):
        p = pyaudio.PyAudio()
        stream = p.open(format=pyaudio.paInt16, channels=1,
                        rate=44100, input=True, frames_per_buffer=chunk)
        data = stream.read(chunk)
        rmsTemp = audioop.rms(data, 2)
        rms.append(rmsTemp)
        rmsMean = numpy.mean(rms)
        rmsStd = numpy.std(rms)
        stream.stop_stream()
        stream.close()
        p.terminate()
    logger.debug("Determined audio level: %d +/- %d" % (rmsMean, rmsTemp))
    return rmsMean


def recordAudio(record_seconds, wave_output_filename):
    logger = logging.getLogger('recordAudio')
    try:
        audioRecorder(record_seconds, wave_output_filename)
    except:
        logger.error(sys.exc_info())


def audioRecorder(record_seconds, wave_output_filename):
    logger = logging.getLogger('recordAudio')
    wave_output_filename = wave_output_filename + ".wav"
    FORMAT = pyaudio.paInt16
    CHANNELS = 2
    RATE = 44100
    CHUNK = 1024
    record_seconds = 5

    audio = pyaudio.PyAudio()

    # start Recording
    stream = audio.open(format=FORMAT, channels=CHANNELS,
                        rate=RATE, input=True,
                        frames_per_buffer=CHUNK)

    logger.debug("recording...")
    frames = []
    for i in range(0, int(RATE / CHUNK * record_seconds)):
        data = stream.read(CHUNK)
        frames.append(data)
    logger.debug("...finished recording")

    # stop Recording
    stream.stop_stream()
    stream.close()
    audio.terminate()

    waveFile = wave.open(wave_output_filename, 'wb')
    waveFile.setnchannels(CHANNELS)
    waveFile.setsampwidth(audio.get_sample_size(FORMAT))
    waveFile.setframerate(RATE)
    waveFile.writeframes(b''.join(frames))
    waveFile.close()


def turnLight(on):
    a, b, isDaylight = getSunTimes()
    if on == True and not isDaylight:
        GPIO.output(RELAY, RELAY_ON)
    else:
        GPIO.output(RELAY, RELAY_OFF)


def captureTestImage():
    turnLight(ON)
    logger = logging.getLogger('captureTestImage')
    logger.debug('Capturing test image...')
    command = "raspistill -w %s -h %s -t 1 -n -vf -e bmp -o -" % (
        TEST_IMAGE_SIZE[0], TEST_IMAGE_SIZE[1])
    imageData = BytesIO()
    imageData.write(subprocess.check_output(command, shell=True))
    imageData.seek(0)
    im = Image.open(imageData)
    buffer = im.load()
    imageData.close()
    logger.debug('...done.')
    turnLight(OFF)
    return im, buffer


def saveImage(filenameFull):
    try:
        logger = logging.getLogger('saveImage')
        logger.debug('Capturing image...')
        filenameFull = filenameFull + ".jpg"
        turnLight(ON)
        subprocess.call(
            "raspistill -w 800 -h 600 -t 1 -n -vf -e jpg -q 15 -o %s" % filenameFull, shell=True)
        turnLight(OFF)
        logger.debug("...captured image %s" % filenameFull)
    except:
        logger.error(sys.exc_info())


def compareImages(buffer1, buffer2):
    logger = logging.getLogger('compareImages')
    logger.debug('Comparing images...')
    pixelSum = 0
    numCountedPixels = 0
    motionHasBeenDetected = False
    # Count changed pixels
    changedPixels = 0
    for x in range(TEST_IMAGE_SIZE[0]):
        # Scan one line of image then check sensitivity for movement
        for y in range(TEST_IMAGE_SIZE[1]):
            # Just check green channel as it's the highest quality channel
            pixdiff = abs(buffer1[x, y][1] - buffer2[x, y][1])
            pixelSum += buffer1[x, y][1]
            numCountedPixels += 1
            if pixdiff > 20:
                changedPixels += 1
    percentChange = changedPixels / \
        (TEST_IMAGE_SIZE[0] * TEST_IMAGE_SIZE[1]) * 100
    return percentChange


def getTimeString():
    time = datetime.datetime.now()
    return "%04d%02d%02d%02d%02d%02d" % (time.year, time.month, time.day, time.hour, time.minute, time.second)


def saveImageAndAudio():
    logger = logging.getLogger('saveImageAndAudio')
    filename = getTimeString()
    logger.debug('Saving image and audio to %s' % filename)
    t1 = threading.Thread(target=recordAudio, args=(10, filename,))
    t2 = threading.Thread(target=saveImage, args=(filename,))
    t1.start()
    t2.start()
    t1.join()
    t2.join()
    logger.debug("Syncing...")
    os.system("rsync *.wav *.jpg zns@cowyo.com:/www/hens/server/static/data/")
    # os.system("rm *.wav")
    # os.system("rm *.jpg")
    logger.debug("...done.")


def getSunTimes():
    logger = logging.getLogger('getSunTimes')
    p = subprocess.Popen(
        ['./sunset'], stdout=subprocess.PIPE, stderr=subprocess.PIPE)
    out, err = p.communicate()
    m = out.decode('utf-8').split()
    hoursToSunrise = float(m[0])
    hoursToSunset = float(m[1])
    daylight = False
    if hoursToSunrise < 0 and hoursToSunset > 0:
        daylight = True
    logger.debug("Sunrise: %2.1f, Sunset: %2.1f, %d" %
                 (hoursToSunrise, hoursToSunset, daylight))
    return hoursToSunrise, hoursToSunset, daylight

if __name__ == "__main__":
    print(getSunTimes())
    saveImage('test')
    logger = logging.getLogger('main')
    image1, buffer1 = captureTestImage()
    while 1:
        logger.debug("Sleeping")
        time.sleep(3)

        logger.debug("Comparing new images")
        image2, buffer2 = captureTestImage()
        percentChange = compareImages(buffer1, buffer2)
        logger.debug("Photo percent change: %2.1f" % percentChange)
        if percentChange >= 1:
            saveImageAndAudio()
        image1 = image2
        buffer1 = buffer2
    #
    #     logger.debug("Comparing new sounds")
    #     audioBaseline2 = getAudioLevel()
    #     percentChange = (audioBaseline2 - audioBaseline) / audioBaseline
    #     logger.debug("Audio percent change: %2.1f" % percentChange)
    #     if percentChange >= 10:
    #         saveImageAndAudio()
    #     audioBaseline = audioBaseline2
