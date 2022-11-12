import { useState } from 'react';
import './App.css';

const url = "ws://localhost:5000";
let socket;

function hostRoom() {
  socket = new WebSocket(url, ["owner"]);
  console.log(socket);

  socket.onopen = function(e) {
    console.log("[open] Connection established");
    console.log("Sending to server");
    socket.send("My name is John");
  };
  
  socket.onmessage = function(event) {
    console.log(`[message] Data received from server: ${event.data}`);
  };
  
  socket.onclose = function(event) {
    if (event.wasClean) {
      console.info(`[close] Connection closed cleanly, code=${event.code} reason=${event.reason}`);
    } else {
      // e.g. server process killed or network down
      // event.code is usually 1006 in this case
      console.error('[close] Connection died');
    }
  };

  socket.onerror = function(error) {
    console.error(error);
    console.error(`[error]`);
  };
}

function joinRoom(room_id) {
  socket = new WebSocket(url, ["guest"]);

  Supply(socket, "Room-Id", room_id);
}

function Supply(conn, ask, ans) {
  const handlerCache = conn.onmessage;
  conn.onmessage = (event) => {
    let pack = JSON.parse(event.data);
    if (pack.Bus == ask) {
      pack = {
        "Bus": ans
      }
      conn.send(JSON.stringify(pack));
    }
    conn.onmessage = handlerCache;
  }
}

function App() {

  return (
    <div className="App">
      <button onClick={hostRoom} >Owner</button>
      <button onClick={(e)=>joinRoom("1")} >Guest</button>
    </div>
  );
}

export default App;
