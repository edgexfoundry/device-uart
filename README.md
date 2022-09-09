# UART Device Service

## Overview

UART device Service - This device service is used to connect serial UART devices to EdgeX

- Function:
  - This device service is used for universal serial device, such as USB to TTL serial, rs232 or rs485 interface device. It provides REST API interfaces to communicate with serial device
- This device service **ONLY works on Linux system**
- This device service is contributed by [Jiangxing Intelligence](https://www.jiangxingai.com)

## Usage

- This Device Service runs with other EdgeX Core Services, such as Core Metadata, Core Data, and Core Command

- The uart device service contains many pre-defined devices which were defined by `res/devices/device.uart.toml`. The devices in this file are created by the uart device service in core metadata when the service first initializes.

- It has a generic profile which supports read and write. Depends on the application requirements this profile can be used to connect any kind of serial device.

  - `device.uart.generic.yaml` describes the generic profile used read and write to the serial device. As an example, it has two device resources - one for read and another for write. The device command is used to read the read value. Depends on the application requirements, the device resources need to be added. The autoevents in the device.uart.toml file can be enabled or disabled as per the requirements. The parameters 'maxbytes' and 'timeout' are needed for uart read. 'maxbytes' indicates the maximum number of bytes that particular device resource can read. 'timeout' in seconds is used for read so that the read call does not block forever. Ensure timeout is used for write device resource also.

    ```yml
    name: "Uart-Generic-Device"
    manufacturer: "edgex"
    model: "edgex-uart"
    labels:
    - "device-uart-example"
    description: "Example of Device-Uart"
    
    deviceResources:
    -
        name: "Read_UART"
        isHidden: false
        description: "used to read from the UART device"
        attributes: { type: "generic",  dev: "/dev/ttyAMA1", baud: 115200, maxbytes: 160, timeout: 1}
        properties:
            valueType: "String"
            readWrite: "R"
    -
        name: "Write_UART"
        isHidden: false
        description: "used to write to the UART device"
        attributes: { type: "generic",  dev: "/dev/ttyAMA1", baud: 115200, timeout: 1}
        properties:
            valueType: "String"
            readWrite: "W"
    
    deviceCommands:
    -
      name: "Read_Cmd"
      readWrite: "R"
      resourceOperations:
      - { deviceResource: "Read_UART" }
    
    ```
  
- After the uart device service has started, we can communicate with these corresponding pre-defined devices.

## Guidance

Here we give the step by step guidance examples of using this device service. In these examples, we use RESTful API to interact with EdgeX (please notice that, you still need to use Core Command service rather than directly interact with UART device service).

Since the `edgex-cli` has released, we can use this new approach to operate devices:

`edgex-cli command list -d Uart-Transceiver-Device`

If you would prefer the traditional RESTful way to operate, you can try:

`curl -X 'GET' http://localhost:59882/api/v2/device/name/Uart-Generic-Device | json_pp`

Use the `curl` response to get the command URLs (with device name and resource) to issue commands to the UART device via the command service as shown below. You can also use a tool like `Postman` instead of `curl` to issue the same commands.

```json
{
   "apiVersion" : "v2",
   "deviceCoreCommand" : {
      "coreCommands" : [
         {
            "get" : true,
            "name" : "Read_Cmd",
            "parameters" : [
               {
                  "resourceName" : "Read_UART",
                  "valueType" : "String"
               }
            ],
            "path" : "/api/v2/device/name/Uart-Generic-Device/Read_Cmd",
            "url" : "http://localhost:59882"
         },
         {
            "get" : true,
            "name" : "Read_UART",
            "parameters" : [
               {
                  "resourceName" : "Read_UART",
                  "valueType" : "String"
               }
            ],
            "path" : "/api/v2/device/name/Uart-Generic-Device/Read_UART",
            "url" : "http://localhost:59882"
         },
         {
            "name" : "Write_UART",
            "parameters" : [
               {
                  "resourceName" : "Write_UART",
                  "valueType" : "String"
               }
            ],
            "path" : "/api/v2/device/name/Uart-Generic-Device/Write_UART",
            "set" : true,
            "url" : "http://localhost:59882"
         }
      ],
      "deviceName" : "Uart-Generic-Device",
      "profileName" : "Uart-Generic-Device"
   },
   "statusCode" : 200
}
```

### Write data to the serial device

Assume we have a [FT232 (USB to serial UART interface)](https://www.amazon.com/DSD-TECH-Adapter-FT232RL-Compatible/dp/B07BBPX8B8/ref=sr_1_3?dchild=1&keywords=FT232&qid=1631620484&sr=8-3)  uart device connected between raspberry pi 4b and a PC. Raspberry Pi is running the edgex and the uart device service. Then we can send the string "12345" from the raspberry pi to the PC using the below command. This will print '12345' in the PC's serial application. Ensure timeout is used for write device resource. 

```shell
# Send the write command
$ curl -X 'PUT' "http://localhost:59882/api/v2/device/name/Uart-Generic-Device/Write_UART" -d '{"Write_UART":"3132333435"}' | json_pp
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100    64  100    37  100    27   4020   2934 --:--:-- --:--:-- --:--:--  8000
{
   "apiVersion" : "v2",
   "statusCode" : 200
}
```

### Read data from the serial device

Assume we have the same FT232 device connected between raspberry pi 4b and a PC. Raspberry Pi is running the edgex and the uart device service. Send the data `abcdef`(616263646566 in hex) from the PC serial application to Raspberry Pi. Then receive that data using the below command.  

```shell
# Receive data from the PC
$ curl -X 'GET' "http://localhost:59882/api/v2/device/name/Uart-Generic-Device/Read_UART" | json_pp
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100   464  100   464    0     0  51475      0 --:--:-- --:--:-- --:--:-- 58000
{
   "apiVersion" : "v2",
   "event" : {
      "apiVersion" : "v2",
      "deviceName" : "Uart-Generic-Device",
      "id" : "76b1880b-597b-4be9-b406-249cabb84a2f",
      "origin" : 1662742218037556552,
      "profileName" : "Uart-Generic-Device",
      "readings" : [
         {
            "deviceName" : "Uart-Generic-Device",
            "id" : "d5abbf8e-9947-47f1-855f-50f4179e10e7",
            "origin" : 1662742218037556552,
            "profileName" : "Uart-Generic-Device",
            "resourceName" : "Read_UART",
            "value" : "616263646566",
            "valueType" : "String"
         }
      ],
      "sourceName" : "Read_UART"
   },
   "statusCode" : 200
}
```

## License

[Apache-2.0](LICENSE)


