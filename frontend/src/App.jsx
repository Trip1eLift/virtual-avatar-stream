import { useState } from 'react';
import './App.css';

const url = "ws://localhost:5000";
let socket;

function hostRoom() {
  socket = new WebSocket(url, ["host"]);
  console.log(socket);

  socket.onopen = function(e) {
    alert("[open] Connection established");
    alert("Sending to server");
    socket.send("My name is John");
  };
  
  socket.onmessage = function(event) {
    alert(`[message] Data received from server: ${event.data}`);
  };
  
  socket.onclose = function(event) {
    if (event.wasClean) {
      alert(`[close] Connection closed cleanly, code=${event.code} reason=${event.reason}`);
    } else {
      // e.g. server process killed or network down
      // event.code is usually 1006 in this case
      alert('[close] Connection died');
    }
  };

  socket.onerror = function(error) {
    console.error(error);
    alert(`[error]`);
  };
}

function joinRoom(room_id) {
  socket = new WebSocket(url, ["join"]);
}

function App() {

  return (
    <div className="App">
      <button onClick={hostRoom} >Host</button>
    </div>
  );
}

export default App;
