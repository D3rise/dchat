import React, { useState, useEffect, useRef } from 'react';
import './App.css';
import MatrixBackground from './components/MatrixBackground';
import { Mic, MicOff, PhoneOff, User, Terminal } from 'lucide-react';
import Peer from 'simple-peer';
import { Buffer } from 'buffer';

window.Buffer = Buffer;

function App() {
  const [joined, setJoined] = useState(false);
  const [roomName, setRoomName] = useState('');
  const [userName, setUserName] = useState('');
  const [peers, setPeers] = useState([]);
  const [isMuted, setIsMuted] = useState(false);
  const [stream, setStream] = useState(null);

  const socketRef = useRef();
  const peersRef = useRef([]);
  const userStreamRef = useRef();

  useEffect(() => {
    return () => {
      if (socketRef.current) {
        socketRef.current.close();
      }
    };
  }, []);

  const handleJoin = () => {
    if (roomName && userName) {
      navigator.mediaDevices.getUserMedia({ audio: true, video: false }).then(stream => {
        setStream(stream);
        userStreamRef.current = stream;
        setJoined(true);

        const ws = new WebSocket('ws://localhost:8080/ws');
        socketRef.current = ws;

        ws.onopen = () => {
          ws.send(JSON.stringify({
            type: 'join',
            room: roomName,
            user: userName
          }));
        };

        ws.onmessage = (message) => {
          const payload = JSON.parse(message.data);
          
          if (payload.type === 'all users') {
            const peers = [];
            payload.users.forEach(userID => {
              const peer = createPeer(userID, ws, stream);
              peersRef.current.push({
                peerID: userID,
                peer,
              });
              peers.push({
                peerID: userID,
                peer,
              });
            });
            setPeers(peers);
          }

          if (payload.type === 'user joined') {
            const peer = addPeer(payload.signal, payload.callerID, stream);
            peersRef.current.push({
              peerID: payload.callerID,
              peer,
            });
            setPeers(users => [...users, { peerID: payload.callerID, peer }]);
          }

          if (payload.type === 'receiving returned signal') {
            const item = peersRef.current.find(p => p.peerID === payload.id);
            if (item) {
              item.peer.signal(payload.signal);
            }
          }
        };
      });
    }
  };

  function createPeer(userToSignal, ws, stream) {
    const peer = new Peer({
      initiator: true,
      trickle: false,
      stream,
    });

    peer.on("signal", signal => {
      ws.send(JSON.stringify({
        type: 'sending signal',
        userToSignal,
        signal
      }));
    });

    return peer;
  }

  function addPeer(incomingSignal, callerID, stream) {
    const peer = new Peer({
      initiator: false,
      trickle: false,
      stream,
    });

    peer.on("signal", signal => {
      socketRef.current.send(JSON.stringify({
        type: 'returning signal',
        callerID,
        signal
      }));
    });

    peer.signal(incomingSignal);

    return peer;
  }

  const toggleMute = () => {
    if (stream) {
      stream.getAudioTracks()[0].enabled = isMuted;
      setIsMuted(!isMuted);
    }
  };

  const leaveRoom = () => {
    window.location.reload();
  };

  return (
    <div className="matrix-container">
      <MatrixBackground />
      {!joined ? (
        <div className="matrix-box">
          <Terminal size={48} style={{ marginBottom: '1rem' }} />
          <h1 className="matrix-title">System Entry</h1>
          <input
            className="matrix-input"
            placeholder="ACCESS_ID (Username)"
            value={userName}
            onChange={(e) => setUserName(e.target.value)}
          />
          <input
            className="matrix-input"
            placeholder="NODE_ID (Room Name)"
            value={roomName}
            onChange={(e) => setRoomName(e.target.value)}
          />
          <button className="matrix-button" onClick={handleJoin}>
            Establish Connection
          </button>
        </div>
      ) : (
        <div className="room-view">
          <h1 className="matrix-title">Node: {roomName}</h1>
          <div className="user-grid">
            <div className={`user-node ${!isMuted ? 'speaking' : ''}`}>
              <User size={40} />
              <div className="user-name">{userName} (You)</div>
            </div>
            {peers.map((peerObj, index) => (
              <Audio key={index} peer={peerObj.peer} />
            ))}
          </div>

          <div className="controls">
            <button className={`control-btn ${isMuted ? 'active' : ''}`} onClick={toggleMute}>
              {isMuted ? <MicOff /> : <Mic />}
            </button>
            <button className="control-btn danger" onClick={leaveRoom}>
              <PhoneOff />
            </button>
          </div>
        </div>
      )}
    </div>
  );
}

const Audio = (props) => {
  const ref = useRef();

  useEffect(() => {
    props.peer.on("stream", stream => {
      if (ref.current) {
        ref.current.srcObject = stream;
      }
    });
  }, [props.peer]);

  return (
    <div className="user-node">
      <User size={40} />
      <div className="user-name">Remote User</div>
      <audio playsInline autoPlay ref={ref} style={{ display: 'none' }} />
    </div>
  );
};

export default App;
