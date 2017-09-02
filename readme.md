# Go Printer Relay

Web status page and Serial <> TCPIP relay for daVinci jr 3d printers

## Running it

1. download and the unzip the release build for your platform

2. set the path to the serial device in config.json
   change value of devicePath in above file
   for windows this will be COMx, check device manager.

3. then from the command line start goPrinter_{platform}

   #### Linux ( raspberry pi ) example

   download

   ```shell
   wget https://github.com/rhysbryant/goprinterrelay/releases/download/0.0.0a-1-ge2912f6-3-gf9185d0/goprinter_linux_arm.tar.gz 
   tar -xvf goprinter_linux_arm.tar.gz
   ```


   ​

4. go to http://{system-ip-address}:8080/ in your web browser to view the status page

5. add {system-ip-addesss} as a Wi-Fi connected printer within the printer's software

 ## Building it from source 
this assumes you have some level of knowledge of golang 

````shell
go get github.com/rhysbryant/goprinterrelay
goprinterrelay

````