# Image classification of chickens

This attempt will use histograms (implemented by `simplecv`). Maybe later I'll attempt neural nets. 

# Requirements

Install dependencies:

```
sudo apt-get install ipython python-opencv python-scipy \
 python-numpy python-pygame python-setuptools python-pip
```

Install python packages

```
sudo pip install svgwrite
sudo pip install https://github.com/sightmachine/SimpleCV/zipball/develop
git clone https://github.com/biolab/orange.git
cd orange
python setup.py build           <--- TAKES A LONG TIME!
sudo python setup.py install
```

# Run

To run using the current dataset, do

```
python run.py
```