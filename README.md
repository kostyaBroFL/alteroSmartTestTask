# Test task for AlteroSmart company.

## TL;DR

For this test case, need to implement several services.
The MS_Generation service generates data by simulating multiple IoT devices.
The MS_Persistence service receives and stores data from devices.
Also, I need a REST service with a frontend for presentation,
but so far they have not been made.

## Build & run

For build and run this services, database & migration just use *docker-compose*.
```bash
docker-compose up
```

## Using

U can use [this postman collection](
./AlteroSmartTestTask.postman_collection.json)
for export examples of the requests to this services.
