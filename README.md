# UART Device Service

## Overview

UART Micro Service - device service for connecting serial UART devices to EdgeX

- Function:
  - This device service is used for universal serial device, such as USB to TTL serial, rs232 or rs485 interface device. It provides REST API interfaces to communicate with serial device
- This device service **ONLY works on Linux system**
- This device service is contributed by [Jiangxing Intelligence](https://www.jiangxingai.com)



## Usage

- This Device Service runs with other EdgeX Core Services, such as Core Metadata, Core Data, and Core Command

- The uart device service can contains many pre-defined devices which were defined by `res/devices/device.uart.toml` such as `Uart-Monitor-Device` and `Uart-Transceiver-Device`. These devices are created by the uart device service in core metadata when the service first initializes

- Two device profiles are now defined, they are `device.uart.monitor.yaml` and `device.uart.transceiver.yaml`

  - `device.uart.monitor.yaml` describes the actual `UART` hardware device used to monitor whether there is data input. The device will continue to monitor until the micro server exits or the hardware device is abnormal. It provides restful api for user to detect the input data.

    ```yml
    name: "Uart-Monitor-Device"
    manufacturer: "Jiangxing Intelligence"
    model: "SP-01"
    labels:
    - "device-uart-example"
    description: "Example of Device-Uart"

    deviceResources:
    -
        name: "Detect_FT232"
        isHidden: false
        description: "used to detect whether the FT232 uart device has data input"
        attributes: { tpye: "monitor",  dev: "/dev/ttyUSB5", baud: 115200}
        properties:
            valueType: "String"
            readWrite: "R"
    ```

  - `device.uart.transceiver.yaml` describes the actual `UART` hardware device as transceivers that receive command and return response. If the read time expires and there is no data response, it will return empty.

    ```yml
    name: "Uart-Transceiver-Device"
    manufacturer: "Jiangxing Intelligence"
    model: "SP-01"
    labels:
    - "device-uart-example"
    description: "Example of Device-Uart"

    deviceResources:
    -
        name: "Query_EC20"
        isHidden: false
        description: "used to send hex buffer to EC20 4G module or receive response"
        attributes: { tpye: "transceiver",  dev: "/dev/ttyUSB2", baud: 115200, timeout: 200}
        properties:
          valueType: "String"
          readWrite: "RW"
    -
        name: "Query_FT232"
        isHidden: false
        description: "used to send hex buffer to FT232(USB to serial UART interface) or receive response"
        attributes: { tpye: "transceiver", dev: "/dev/ttyUSB4", baud: 9600, timeout: 200 }
        properties:
            valueType: "String"
            readWrite: "RW"
    ```

- After the uart device service has started, we can communicate these corresponding pre-defined  devices.



## Guidance

Here we give two step by step guidance examples of using this device service. In these examples, we use RESTful API to interact with EdgeX (please notice that, you still need to use Core Command service rather than directly interact with UART device service).

Since the `edgex-cli` has released, we can use this new approach to operate devices:

`edgex-cli command list -d Uart-Transceiver-Device`

If you would prefer the traditional RESTful way to operate, you can try:

`curl http://localhost:59882/api/v2/device/name/Uart-Monitor-Device` and

`curl http://localhost:59882/api/v2/device/name/Uart-Transceiver-Device`

Use the `curl` response to get the command URLs (with device name and resource) to issue commands to the UART device via the command service as shown below. You can also use a tool like `Postman` instead of `curl` to issue the same commands.

```json
{
    "apiVersion": "v2",
    "statusCode": 200,
    "deviceCoreCommand": {
        "deviceName": "Uart-Monitor-Device",
        "profileName": "Uart-Monitor-Device",
        "coreCommands": [
            {
                "name": "Detect_FT232",
                "get": true,
                "path": "/api/v2/device/name/Uart-Monitor-Device/Detect_FT232",
                "url": "http://edgex-core-command:59882",
                "parameters": [
                    {
                        "resourceName": "Detect_FT232",
                        "valueType": "String"
                    }
                ]
            }
        ]
    }
}
```

```json
{
    "apiVersion": "v2",
    "statusCode": 200,
    "deviceCoreCommand": {
        "deviceName": "Uart-Transceiver-Device",
        "profileName": "Uart-Transceiver-Device",
        "coreCommands": [
            {
                "name": "Query_EC20",
                "get": true,
                "set": true,
                "path": "/api/v2/device/name/Uart-Transceiver-Device/Query_EC20",
                "url": "http://edgex-core-command:59882",
                "parameters": [
                    {
                        "resourceName": "Query_EC20",
                        "valueType": "String"
                    }
                ]
            },
            {
                "name": "Query_FT232",
                "get": true,
                "set": true,
                "path": "/api/v2/device/name/Uart-Transceiver-Device/Query_FT232",
                "url": "http://edgex-core-command:59882",
                "parameters": [
                    {
                        "resourceName": "Query_FT232",
                        "valueType": "String"
                    }
                ]
            }
        ]
    }
}
```



### Detect data from uart-monitor

Assume we have a [FT232 (USB to serial UART interface)](https://www.amazon.com/DSD-TECH-Adapter-FT232RL-Compatible/dp/B07BBPX8B8/ref=sr_1_3?dchild=1&keywords=FT232&qid=1631620484&sr=8-3)  uart device connected to a usbhost on current system of raspberry pi 4b and it's serial UART interface connect to a computer via another FT232 uart device. Then we can send the string "12345" by computer to raspberry pi 4b.

```shell
# Set the 'Power' gpio to high
$ curl  http://localhost:59882/api/v2/device/name/Uart-Monitor-Device/Detect_FT232
{"apiVersion":"v2","statusCode":200,"event":{"apiVersion":"v2","id":"226b833b-ef4a-4541-9e1b-9173414c04e5","deviceName":"Uart-Monitor-Device","profileName":"Uart-Monitor-Device","sourceName":"Detect_FT232","origin":1631620308488212373,"readings":[{"id":"63cb5c4e-c4ee-414e-a867-76cb9285d3e5","origin":1631620308488212373,"deviceName":"Uart-Monitor-Device","resourceName":"Detect_FT232","profileName":"Uart-Monitor-Device","valueType":"String","binaryValue":null,"mediaType":"","value":"3132333435"}]}}
```




### Send cmd to uart-transceiver and get reply

Assume we have a [4G LTE USB Dongle](https://www.amazon.com/EXVIST-Dongle-EG25-G-M2M-optimized-Module/dp/B096TQJR73/ref=sr_1_1_sspa?dchild=1&keywords=4G+Module&qid=1631620738&sr=8-1-spons&psc=1&spLa=ZW5jcnlwdGVkUXVhbGlmaWVyPUEzNUNLQkpFU0hXNDlCJmVuY3J5cHRlZElkPUEwNDAyOTM5MzlSWUVFRVpQUzRFRSZlbmNyeXB0ZWRBZElkPUEwMzg5Njc4MTcyMjM5RlZHQUpOSCZ3aWRnZXROYW1lPXNwX2F0ZiZhY3Rpb249Y2xpY2tSZWRpcmVjdCZkb05vdExvZ0NsaWNrPXRydWU=) connected to a usbhost on current system of raspberry pi 4b. We can send a command `"AT+CGSN\r\n"` (41542b4347534e0d0a in hex)  to the USB Dongle for requesting `Product Serial Number Identification`

```shell
# send command
$ curl -X PUT -d '{"Query_EC20":"41542b4347534e0d0a"}'  http://localhost:59882/api/v2/device/name/Uart-Transceiver-Device/Query_EC20
{"apiVersion":"v2","statusCode":200}

# get reply data
$ curl http://localhost:59882/api/v2/device/name/Uart-Transceiver-Device/Query_EC20
curl http://localhost:59882/api/v2/device/name/Uart-Transceiver-Device/Query_EC20
{"apiVersion":"v2","statusCode":200,"event":{"apiVersion":"v2","id":"ea2b3558-d479-4b27-af22-bf92e15f0c84","deviceName":"Uart-Transceiver-Device","profileName":"Uart-Transceiver-Device","sourceName":"Query_EC20","origin":1631621285691637421,"readings":[{"id":"e9b5852f-e331-428c-9160-5cb4aed4d3ae","origin":1631621285691637421,"deviceName":"Uart-Transceiver-Device","resourceName":"Query_EC20","profileName":"Uart-Transceiver-Device","valueType":"String","binaryValue":null,"mediaType":"","value":"41542b4347534e0d0d0a3836323630373035393538373238360d0a0d0a4f4b0d0a"}]}}
```

the value (hex string) shows in string is

```shell
AT+CGSN

862607059587286

OK
```



### docker-compose.yml

Add the `device-uart` to the docker-compose.yml of edgex foundry 2.0-Ireland.

```yml
...
  device-uart:
    container_name: edgex-device-uart
    depends_on:
    - consul
    - data
    - metadata
    environment:
      CLIENTS_CORE_COMMAND_HOST: edgex-core-command
      CLIENTS_CORE_DATA_HOST: edgex-core-data
      CLIENTS_CORE_METADATA_HOST: edgex-core-metadata
      CLIENTS_SUPPORT_NOTIFICATIONS_HOST: edgex-support-notifications
      CLIENTS_SUPPORT_SCHEDULER_HOST: edgex-support-scheduler
      DATABASES_PRIMARY_HOST: edgex-redis
      EDGEX_SECURITY_SECRET_STORE: "false"
      MESSAGEQUEUE_HOST: edgex-redis
      REGISTRY_HOST: edgex-core-consul
      SERVICE_HOST: edgex-device-uart
    hostname: edgex-device-uart
    image: edgexfoundry/device-uart:0.0.0-dev
    networks:
      edgex-network: {}
    ports:
    - 49995:49995/tcp
    read_only: false
    privileged: true
    volumes:
    - "/dev:/dev"
    security_opt:
    - no-new-privileges:false
    user: root:root
...
```



## License

[Apache-2.0](LICENSE)

