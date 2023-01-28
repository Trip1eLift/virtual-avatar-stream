import { useState, useEffect, useRef } from 'react';
import Webcam from 'react-webcam';
import WebSocketPeering from './client2client/websocket-peering';
import './App.css';

const url1 = "ws://localhost:5000";
const url2 = "ws://localhost:5001";

const remoteUrl = "wss://virtualavatar-stream.trip1elift.com/";

const videoConstraints = {
  width: 640,
  height: 480,
  facingMode: "user",
};

const audioConstraints = {
  suppressLocalAudioPlayback: true,
  noiseSuppression: true,
  echoCancellation: true,
};

function App() {
  const [ownerRoomId, setOwnerRoomId] = useState();
  const [guestRoomId, setGuestRoomId] = useState();
  const [wsp, _] = useState(new WebSocketPeering(false));
  const webcamRef = useRef(null);
  const remoteMedia = useRef(null);

  useEffect(() => {
    if (typeof(remoteMedia.current) !== "undefined" && remoteMedia.current !== null) {
      remoteMedia.current.srcObject = wsp.getRemoteStream();
      remoteMedia.current.onloadedmetadata = function(e) {
        remoteMedia.current.play();
      };
    }
  }, []);

  return (
    <div className="App">
      <button onClick={(e)=>wsp.ownerConn(url1, setOwnerRoomId)} >Owner Local Sever 1</button>
      <button onClick={(e)=>wsp.ownerConn(url2, setOwnerRoomId)} >Owner Local Sever 2</button>
      <button onClick={(e)=>wsp.ownerConn(remoteUrl, setOwnerRoomId)} >Owner Remote Sever</button>
      <br/>
      {ownerRoomId && <>Owner Room ID: {ownerRoomId}</>}
      <br/>
      <input type="text" onChange={(e)=>setGuestRoomId(e.target.value)} />
      <button onClick={(e)=>wsp.guestConn(url1, guestRoomId)} >Guest Local Server 1</button>
      <button onClick={(e)=>wsp.guestConn(url2, guestRoomId)} >Guest Local Server 2</button>
      <button onClick={(e)=>wsp.guestConn(remoteUrl, guestRoomId)} >Guest Remote Server</button>
      <br/>
      <br/>
      <button onClick={(e)=>wsp.sendUuid()}>Send ID (websocket)</button>
      <button onClick={(e)=>wsp.sendUuidPeer()}>Send ID (peer)</button>
      <br/>
      <br/>
      <Webcam
        audio={true}
        muted={true}
        height={100 + "%"}
        width={100 + "%"}
        videoConstraints={videoConstraints}
        audioConstraints={audioConstraints}
        mirrored={true}
        ref={webcamRef}
        onUserMedia={(stream)=>wsp.onUserMedia(stream)}
      />
      {wsp.boolStreamVideo() && <video ref={remoteMedia} autoPlay ></video>}
      {!wsp.boolStreamVideo() && <audio ref={remoteMedia} autoPlay></audio>}
    </div>
  );
}

export default App;