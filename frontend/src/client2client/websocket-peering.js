import { v4 as uuidv4 } from 'uuid';

// WebRTC docs:   https://levelup.gitconnected.com/establishing-the-webrtc-connection-videochat-with-javascript-step-3-48d4ae0e9ea4
// firebase docs: https://github.com/fireship-io/webrtc-firebase-demo/blob/main/main.js
// WebRTC org:    https://webrtc.org/getting-started/peer-connections
//                https://stackoverflow.com/questions/22470291/rtcdatachannels-readystate-is-not-open

const ICE_config = {
  iceServers: [
    {
      'url': 'stun:stun.l.google.com:19302'
    },
    {
      'url': 'turn:192.158.29.39:3478?transport=udp',
      'credential': 'JZEOEt2V3Qb0y27GRntt2u2PAYA=',
      'username': '28224511:1379330808'
    },
    {
      'url': 'turn:192.158.29.39:3478?transport=tcp',
      'credential': 'JZEOEt2V3Qb0y27GRntt2u2PAYA=',
      'username': '28224511:1379330808'
    }
  ],
  iceCandidatePoolSize: 10,
};

const Peer_options = {
  optional: [{
    RtpDataChannels: true
  }]
};

/**
 * Client to client negotiation: 
 * event.data = JSON.stringify(payload)
 * 
 * Payload types:
 * 
 * { message_type: 'push',                    message: "uuid"    }
 * { message_type: 'start-peer-connection',   message: undefined }
 * { message_type: 'offer-peer-connection',   message: offer     } offer  = await peerConnection.createOffer();
 * { message_type: 'answer-peer-connection',  message: answer    } answer = await peerConnection.createAnswer();
 * { message_type: 'new-ice-candidate',       message: candidate }
 */
const MESSAGE_TYPE = {
  push:                 'push',
  startPeerConnection:  'start-peer-connection',
  offerPeerConnection:  'offer-peer-connection',
  answerPeerConnection: 'answer-peer-connection',
  newIceCandidate:      'new-ice-candidate',
};

class WebSocketPeering {
  constructor() {
    
    const peer = new RTCPeerConnection(ICE_config, Peer_options);

    peer.ondatachannel = (event) => {
      event.channel.onmessage = (event) => {
        // PeerConnection DataChannel Listener
        const payload = JSON.parse(event.data);

        if (payload.message_type === MESSAGE_TYPE.push) {
          console.log(`[peer] recieved: ${payload.message}`);
        }
      };
    };

    const datachannel = peer.createDataChannel("data");
    
    datachannel.onopen = () => {
      console.log("[peer] Connection established; Closing websocket");
      this.socket.close(1000, "websocket is no longer needed.");
    };
    
    datachannel.onclose = () => {
      console.log("[peer]");
    };

    this.socket = undefined;
    this.peer = peer;
    this.dc = datachannel;
  }

  ownerConn(url, setRoomId) {
    const socket = new WebSocket(url, ["owner"]);
    const peer = this.peer;
  
    socket.onopen = async function(e) {
      console.log("[socket open] Connection established");
      
      const room_id = await Demand(socket, "Room-Id");
      console.log(`Room-Id: ${room_id}`);
      setRoomId(room_id);
    };
    
    socket.onmessage = async function(event) {
      const payload = JSON.parse(event.data);

      if (payload.message_type === MESSAGE_TYPE.push) {
        console.log(`[owner socket] recieved: ${payload.message}`);
      }

      // 2. Owner reacts to start-peer-connection
      if (payload.message_type === MESSAGE_TYPE.startPeerConnection) {
        const offer = await peer.createOffer();
        await peer.setLocalDescription(offer);

        socket.send(JSON.stringify({message_type: MESSAGE_TYPE.offerPeerConnection, message: offer}));
      }

      // 4. Owner recieves answer data
      if (payload.message_type === MESSAGE_TYPE.answerPeerConnection) {
        const remoteDesc = new RTCSessionDescription(payload.message);
        await peer.setRemoteDescription(remoteDesc);
      }
    }

    // 5. Owner retrieves a local ICE candidate TODO: how to trigger this?
    peer.addEventListener('icecandidate', event => {
      if (event.candidate) {
        socket.send(JSON.stringify({message_type: MESSAGE_TYPE.newIceCandidate, message: event.candidate}));
      }
    });
    
    attachSocketSharedHandler(socket);
    this.socket = socket;
  }

  guestConn(url, room_id) {
    const socket = new WebSocket(url, ["guest"]);
    const peer = this.peer;
  
    socket.onopen = async function(e) {
      console.log("[socket open] Connection established");
      
      await Supply(socket, "Room-Id", room_id);

      // 1. Guest triggers start-peer-connection
      socket.send(JSON.stringify({message_type: MESSAGE_TYPE.startPeerConnection}));
    };
  
    socket.onmessage = async function(event) {
      const payload = JSON.parse(event.data);

      if (payload.message_type === MESSAGE_TYPE.push) {
        console.log(`[guest socket] recieved: ${payload.message}`);
      }

      // 3. Guest recieves offer data
      if (payload.message_type === MESSAGE_TYPE.offerPeerConnection) {
        peer.setRemoteDescription(new RTCSessionDescription(payload.message));
        const answer = await peer.createAnswer();
        await peer.setLocalDescription(answer);

        socket.send(JSON.stringify({message_type: MESSAGE_TYPE.answerPeerConnection, message: answer}));
      }

      // 6. Guest recieves a remote ICE candidate
      if (payload.message_type === MESSAGE_TYPE.newIceCandidate) {
        if (payload.message) {
          try {
            console.log(payload.message.candidate);
            await peer.addIceCandidate(payload.message);
          } catch (e) {
            console.error('Error adding received ice candidate', e);
          }
        }
      }
    };
    
    attachSocketSharedHandler(socket);
    this.socket = socket;
  }

  sendUuid() {
    if (this.socket !== undefined) {
      this.socket.send(JSON.stringify({message_type: MESSAGE_TYPE.push, message: uuidv4()}));
    }
  }

  sendUuidPeer() {
    if (this.dc !== undefined) {
      this.dc.send(JSON.stringify({message_type: MESSAGE_TYPE.push, message: uuidv4()}));
    }
  }
}

// Attach socket shared handler
function attachSocketSharedHandler(socket) {
  socket.onclose = function(event) {
    if (event.wasClean) {
      console.info(`[socket close] Connection closed cleanly, code=${event.code} reason=${event.reason}`);
    } else {
      console.error('[socket close] Connection died');
    }
  };

  socket.onerror = function(error) {
    console.error(error);
    console.error(`[error]`);
  };
}

// Only used to communicate against match server
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

// Only used to communicate against match server
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

export default WebSocketPeering;

/**
 * Usage:
 * wsp = new WebSocketPeering(url);
 */