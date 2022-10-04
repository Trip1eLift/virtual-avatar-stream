# Virtual-Avatar-Streaming-Backend
This is a PoC project to test how to deploy a scalable end-to-end streaming distributive system to support video chat.

![Architecture Diagram](./Architecture-diagram.jpg)

# Design Updates

## WS encrpytion
Note: WS encryption will be implemented later. A client will generatean asymetric key pair and send the c_public key to server. The server will also generate an asymetric key pair and send the s_public key to all clients.

Flow:
1. client_1 data
2. encrypt with s_public
3. send to server
4. server decrypt with s_private
5. server encrypt with client_2's c_public
6. send to client_2
7. client_2 decrypt with client_2's c_private


## Frontend

### Ingress:
1. https://virtualavatar.trip1elift.com/

Single page website

### Egress:

1. Create room

Request
```
POST https://room.virtualavatar.trip1elift.com/create-room
```
Response
```
{
  room_id: 00157,
  client_id: <uuid_1>
}
```

2. Join room

Request
```
POST https://room.virtualavatar.trip1elift.com/join-room
{
  room_id: 00157
}
```
Response
```
{
  room_id: 00157,
  client_id: <uuid_2>
}
```

3. Enter room

Websocket
```
ws://<ip>:5001
-H Sec-WebSocket-Protocol: "{room_id: 00157, client_id: <uuid_1>}"
```

JS code
```
const ws = new WebSocket("ws://<ip>:5001", "{room_id: 00157, client_id: <uuid_1>}");
```

## Room router (Lambda backend)

### Ingress from client:
1. Health check

Request
```
https://room.virtualavatar.trip1elift.com/health
```
Response
```
200 Healthy.
```

2. Create room

Request
```
POST https://room.virtualavatar.trip1elift.com/create-room
```
* Generate unique client_id with uuid
* Query DB to find total amout of rooms. If rooms reached max_rooms, throttle user.
* Query DB to find task with least rooms
* Query DB to find a unique room_id (5 digits)
* Call a stream task to register with room_id with client_id
```
POST http://<ip>:5001/register
-H Authorization: <token from AWS secret manager>
{
  room_id: 00157,
  client_id: <uuid_1>
}
```
* If register failed, poll tasks IP using AWS API. Cleanup [instance_ip, room_id] and [room_id, client_id] if instance_ip no longer exist. Retry from Query DB.
* Write client_id and room_id to DB table_1: [instance_ip, room_id] and table_2: [room_id, client_id]

Response
```
{
  room_id: 00157,
  client_id: <uuid_1>
}
```

3. Join room

Request
```
POST https://room.virtualavatar.trip1elift.com/join-room
{
  room_id: 00157
}
```
* Generate unique client_id with uuid
* Query DB to find clients in the room. If clients reached max_clients, throttle user.
* Query DB to find the instance_ip of room_id
* Call a stream task to register with room_id with client_id
```
POST http://<ip>:5001/register
-H Authorization: <token from AWS secret manager>
{
  room_id: 00157,
  client_id: <uuid_1>
}
```
* If register failed, poll tasks IP using AWS API. Cleanup [instance_ip, room_id] and [room_id, client_id] if instance_ip no longer exist. Response to user that room_id is no longer valid.
* Else write client_id and room_id to DB table_1: [instance_ip, room_id] and table_2: [room_id, client_id]

Response
```
{
  room_id: 00157,
  client_id: <uuid_2>
}
```

### Ingress from stream instance

4. Put task

Request
```
POST https://room.virtualavatar.trip1elift.com/put-task
-H Authorization: <token from AWS secret manager>
-H Remote: <instance public ip address>
```

If Remote is not the public ip, instance will use AWS API with self private IP (from health check) to identify the ENI's public IP that is attached to itself. Then send it to put-task.

5. Remove client

If is disconnected from the stream, stream will call this endpoint to remove client from the table. If room is empty, the room in runtime will be removed from stream.

Request
```
POST https://room.virtualavatar.trip1elift.com/remove-client
-H Authorization: <token from AWS secret manager>
{
  room_id: 00157,
  client_id: <uuid_1>
}
```
* Remove [room_id, client_id] from table_2.
* If room is empty, remove [instance_ip, room_id] from table_1.

Response
```
200 Success
```

## Task IP updater (Scheduled Lambda)

Note: The priority of this lambda is low.

Poll instance IP once an hour using AWS API and update DB. Cleanup any [instance_ip, room_id] and [room_id, client_id] that are no longer valid.

## Stream (Fargate backend)
1. Health check

Request
```
http://<ip>:5001/health
```
Response
```
200 Healthy.
```

### Ingress from Lambda backend

2. Register

Request
```
POST http://<ip>:5001/register
-H Authorization: <token from AWS secret manager>
{
  room_id: 00157,
  client_id: <uuid_1>
}
```
* Create room if room not exist in runtime.

Response
```
200 Success
```

3. Stream connection

Websocket
```
ws://<ip>:5001/
-H Sec-WebSocket-Protocol: "{room_id: 00157, client_id: <uuid_1>}"
```
* Terminate connection if room_id and client_id pair does not exist in the runtime.

### Instance Lifecycle

Call Put task at start

Request
```
POST https://room.virtualavatar.trip1elift.com/put-task
-H Authorization: <token from AWS secret manager>
-H Remote: <instance public ip address>
```

If Remote is not the public ip, instance will use AWS API with self private IP (from health check) to identify the ENI's public IP that is attached to itself. Then send it to put-task.

### Client Lifecycle

While client closes websocket connection to stream. Stream remove client_id (and remove room_id if room is empty). Stream calls Lambda backend to remove client.

Request
```
POST https://room.virtualavatar.trip1elift.com/remove-client
-H Authorization: <token from AWS secret manager>
{
  room_id: 00157,
  client_id: <uuid_1>
}
```

## License

“Commons Clause” License Condition v1.0

The Software is provided to you by the Licensor under the License, as defined below, subject to the following condition.

Without limiting other conditions in the License, the grant of rights under the License will not include, and the License does not grant to you,  right to Sell the Software.

For purposes of the foregoing, “Sell” means practicing any or all of the rights granted to you under the License to provide to third parties, for a fee or other consideration (including without limitation fees for hosting or consulting/ support services related to the Software), a product or service whose value derives, entirely or substantially, from the functionality of the Software.  Any license notice or attribution required by the License must also include this Commons Cause License Condition notice.

Software: Virtual Avatar Stream
License: Apache 2.0
Licensor: Trip1eLift - Joseph Chang
