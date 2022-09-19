import logo from './logo.svg';
import './App.css';
import {useState, useEffect} from "react";

const SERVER_URL = 'ws://127.0.0.1:5001';
let Connection;

export default function App() {
  const [message, setMessage] = useState("");
  const [print, setPrint] = useState("Initial");

  function connect() {
    Connection = new WebSocket(SERVER_URL);
    Connection.onopen = (() => {
      setPrint('WebSocket Client Connected');
    });
  
    Connection.onmessage = ((msg) => {
      setPrint(msg.data);
    });
  }

  function send() {
    Connection.send(message);
  }

  useEffect(connect, []);

  return (
    <div className="App">
      <header className="App-header">
        <img src={logo} className="App-logo" alt="logo" />
        <table><tbody>
          <tr>
            <td>Message:</td>
            <td><input type="text" style={{fontSize:"2rem", marginLeft:"1rem"}} onChange={(e)=>setMessage(e.target.value)}/></td>
            <td><input type="button" style={{fontSize:"2rem"}} value="Send" onClick={send} /></td>
          </tr>
          <tr>
            <td><input type="button" style={{fontSize:"2rem"}} value="Reconnect" onClick={connect} /></td>
          </tr>
          <tr>
            <td>{print}</td>
          </tr>
        </tbody></table>
        
      </header>
    </div>
  );
}

