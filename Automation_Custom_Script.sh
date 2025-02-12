#/bin/bash

apt-get install scons
cd /home/dietpi/
wget https://github.com/jgarff/rpi_ws281x/archive/refs/tags/v1.0.0.tar.gz
tar -xvf rpi_ws281x-1.0.0.tar.gz && rm -rf rpi_ws281x-1.0.0.tar.gz
cd rpi_ws281x-1.0.0
scons
cp *.a /usr/local/lib/
cp *.h /usr/local/include
cd ..
git clone --branch main --single-branch https://github.com/Solvro/hardware-led-cube.git /home/src
cd src
go build
# TODO: the go program autostarts on boot
