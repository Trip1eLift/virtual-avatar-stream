import { useState } from 'react';
import { v4 as uuidv4 } from 'uuid';
import './App.css';

const url1 = "ws://localhost:5000";
const url2 = "ws://localhost:5001";

// TODO: add remote testing configurations

let socket;

function ownerConn(url, setRoomId) {
  socket = new WebSocket(url, ["owner"]);

  socket.onopen = async function(e) {
    console.log("[open] Connection established");
    
    const room_id = await Demand(socket, "Room-Id");
    console.log(`Room-Id: ${room_id}`);
    setRoomId(room_id);
  };
  
  socket.onmessage = function(event) {
    console.log(`[owner] recieved: ${event.data}`);
  };
  
  socket.onclose = function(event) {
    if (event.wasClean) {
      console.info(`[close] Connection closed cleanly, code=${event.code} reason=${event.reason}`);
    } else {
      console.error('[close] Connection died');
    }
  };

  socket.onerror = function(error) {
    console.error(error);
    console.error(`[error]`);
  };
}

function guestConn(url, room_id) {
  socket = new WebSocket(url, ["guest"]);

  socket.onopen = async function(e) {
    console.log("[open] Connection established");
    
    await Supply(socket, "Room-Id", room_id);
  };

  socket.onmessage = function(event) {
    console.log(`[owner] recieved: ${event.data}`);
  };
  
  socket.onclose = function(event) {
    if (event.wasClean) {
      console.info(`[close] Connection closed cleanly, code=${event.code} reason=${event.reason}`);
    } else {
      console.error('[close] Connection died');
    }
  };

  socket.onerror = function(error) {
    console.error(error);
    console.error(`[error]`);
  };
}

function sendUuid() {
  if (socket !== undefined) {
    socket.send(`Push ${uuidv4()}`);
  }
}

function Demand(conn, ask) {
  return new Promise((resolve, reject) => {
    const handlerTemp = conn.onmessage;
    conn.onmessage = (event) => {
      const pack = JSON.parse(event.data);
      conn.onmessage = handlerTemp;
      resolve(pack.Bus);
    }
    conn.send(JSON.stringify({
      "Bus": ask
    }));
  });
}

function Supply(conn, ask, ans) {
  return new Promise((resolve, reject) => {
    const handlerTemp = conn.onmessage;
    conn.onmessage = (event) => {
      const pack = JSON.parse(event.data);
      if (pack.Bus == ask) {
        conn.send(JSON.stringify({
          "Bus": ans
        }));
      }
      conn.onmessage = handlerTemp;
      resolve();
    }
  });
}

function App() {
  const [ownerRoomId, setOwnerRoomId] = useState();
  const [guestRoomId, setGuestRoomId] = useState();

  return (
    <div className="App">
      <button onClick={(e)=>ownerConn(url1, setOwnerRoomId)} >Owner Sever 1</button>
      <button onClick={(e)=>ownerConn(url2, setOwnerRoomId)} >Owner Sever 2</button>
      <br/>
      {ownerRoomId && <>Owner Room ID: {ownerRoomId}</>}
      <br/>
      <input type="text" onChange={(e)=>setGuestRoomId(e.target.value)} />
      <button onClick={(e)=>guestConn(url1, guestRoomId)} >Guest Server 1</button>
      <button onClick={(e)=>guestConn(url2, guestRoomId)} >Guest Server 2</button>
      <br/>
      <br/>
      <button onClick={(e)=>sendUuid()}>Send ID</button>
    </div>
  );
}

export default App;
