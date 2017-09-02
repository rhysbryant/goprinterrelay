# Go Printer Relay

Web status page and Serial <> TCPIP relay for daVinci jr 3d printers

## Running it

1. download the unzip the release build for your platform

2. set the path to the serial device in config.json
   change value of devicePath in above file 

3. then from the command line start goPrinterRelay-{platform}

4. go to http://{system-ip-address}:8080/ in your web browser to view the status page

5. add {system-ip-addesss} as a Wi-Fi connected printer within the printer's software

 ## Building it
 1. download and install golang
 2. from within the src directory run go build