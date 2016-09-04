# hens

## Camera

<https://github.com/dhowden/raspicam>

Move to PYTHON?

## Raspberry Pi initialization

### [Wifi](https://www.raspberrypi.org/documentation/configuration/wireless/wireless-cli.md)

`sudo vim /etc/wpa_supplicant/wpa_supplicant.conf`

    network={
      ssid="SOMETHING"
      psk="PASSWORD"
    }


### Packages

```
sudo apt-get install apcalc python3 python3-setuptools zsh openssh-server openssh-client tree git vim htop python3-pyaudio python3-pil python3-numpy python3-rpio.gpio
```

### [Setup audio](http://raspberrypi.stackexchange.com/questions/37177/best-way-to-setup-usb-mic-as-system-default-on-raspbian-jessie)

`sudo nano /usr/share/alsa/alsa.conf` scroll down until you find the lines

    defaults.ctl.card 0
    defaults.pcm.card 0

and change them to

    defaults.ctl.card 1
    defaults.pcm.card 1
