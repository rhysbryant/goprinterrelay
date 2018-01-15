# Go Printer Relay

Web status page and Serial <> TCPIP proxy for daVinci jr 3d printers

## Running it

1. download and the unzip the release build for your platform

2. copy it to an appropriate location

3. run the following command to install the service
   `goprint -installSvc`

4. set the path to the serial device in config.json
   change value of devicePath in above file
   for windows this will be COMx, check device manager.

5. then run the service start command relevant to your platform for example

   Windows - net start goprint

   Linux Systemd - systemctl start goPrint

   ### Setup example - Linux ( raspberry pi )

   the following assumes Raspbian or other Linux install with systemd

   #### download & install

    see https://github.com/rhysbryant/goprinterrelay/releases/latest to get the download address

   ```shell

   sudo mkdir /opt/goprint
   sudo cd /opt/goprint
   sudo wget {release download link}/goprinter_linux_arm.tar.gz
   sudo tar -xvf goprinter_linux_arm.tar.gz
   sudo goprint -installSvc
   sudo systemctl start goPrint
   ```
   edit config.json and change devicePath (plugin you printer and run ls /dev/serial/by-id/usb-11f1_2510* )
   to get the device path

   #### To have the program auto start on boot

   ```shell
   sudo systemctl enable goPrint
   ```

   #### Manually start the daemon

   ```shell
   sudo systemctl start goPrint
   ```

  to add custom commands see the help secion in the web interface
   ​

4. go to http://{system-ip-address}:8080/ in your web browser to view the status page

5. add {system-ip-addesss} as a Wi-Fi connected printer within the printer's software

 ## Building it from source 
this assumes you have some level of knowledge of golang 

````shell
go get github.com/rhysbryant/goprinterrelay
go build github.com/rhysbryant/goprinterrelay

````