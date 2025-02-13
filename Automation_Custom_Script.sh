#/bin/bash

# Install scons
apt-get install scons

# Install Go
wget https://golang.org/dl/go1.20.3.linux-armv6l.tar.gz
tar -C /usr/local -xzf go1.20.3.linux-armv6l.tar.gz
export PATH=$PATH:/usr/local/go/bin

# Install rpi_ws281x
cd /home/dietpi/
wget https://github.com/jgarff/rpi_ws281x/archive/refs/tags/v1.0.0.tar.gz
tar -xvf rpi_ws281x-1.0.0.tar.gz && rm -rf rpi_ws281x-1.0.0.tar.gz
cd rpi_ws281x-1.0.0
scons
cp *.a /usr/local/lib/
cp *.h /usr/local/include

# Clone and build the hardware-led-cube project
cd ..
git clone --branch main --single-branch https://github.com/Solvro/hardware-led-cube.git /home/src
cd /home/src
go build

# Create systemd service file
cat <<EOT > /etc/systemd/system/hardware-led-cube.service
[Unit]
Description=Hardware LED Cube Service
After=network.target

[Service]
ExecStart=/home/src/hardware-led-cube
WorkingDirectory=/home/src
StandardOutput=inherit
StandardError=inherit
Restart=always
User=root

[Install]
WantedBy=multi-user.target
EOT

# Enable and start the service
systemctl enable hardware-led-cube.service
systemctl start hardware-led-cube.service
