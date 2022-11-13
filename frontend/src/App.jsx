import { useState } from 'react';
import './App.css';

const url1 = "ws://localhost:5000";
const url2 = "ws://localhost:5001";
let socket;

function ownerConn(url) {
  socket = new WebSocket(url, ["owner"]);

  socket.onopen = async function(e) {
    console.log("[open] Connection established");
    
    const room_id = await Demand(socket, "Room-Id");
    console.log(`Room-Id: ${room_id}`);
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

async function guestConn(url, room_id) {
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

  return (
    <div className="App">
      <button onClick={(e)=>ownerConn(url1)} >Owner Sever 1</button>
      <button onClick={(e)=>ownerConn(url2)} >Owner Sever 2</button>
      <button onClick={(e)=>guestConn(url1, "1")} >Guest Server 1</button>
      <button onClick={(e)=>guestConn(url2, "1")} >Guest Server 2</button>
    </div>
  );
}

export default App;
