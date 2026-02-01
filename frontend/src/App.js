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

        const ws = new WebSocket('wss://192.168.1.137:4000/rtc/ws');
        socketRef.current = ws;

        ws.onopen = () => {
          ws.send(JSON.stringify({
            _type: 'join_room',
            data: {
              room_id: roomName,
              username: userName
            }
          }));
        };

        ws.onmessage = (message) => {
          const payload = JSON.parse(message.data);

          if (payload._type === 'room_user_list') {
            payload.data.users.forEach(username => {
              const peer = createPeer(username, ws, stream);
              const peerObj = {
                peerID: username,
                peer,
              };
              peersRef.current.push(peerObj);
              setPeers(users => [...users, peerObj]);
            });
          }

          if (payload._type === 'signal') {
            const peer = addPeer(payload.data.signal, payload.data.username, stream);
            const peerObj = {
              peerID: payload.data.username,
              peer,
            };
            peersRef.current.push(peerObj);
            setPeers(users => [...users, peerObj]);
          }

          if (payload._type === 'returning_signal') {
            const item = peersRef.current.find(p => p.peerID === payload.data.username);
            if (item) {
              item.peer.signal(payload.data.signal);
            }
          }
        };
      });
    }
  };

  function createPeer(usernameToSignal, ws, stream) {
    const peer = new Peer({
      initiator: true,
      trickle: false,
      stream,
      config: {
        iceServers: [
          { urls: 'stun:stun.l.google.com:19302' },
          { urls: 'stun:stun1.l.google.com:19302' },
        ]
      }
    });

    peer.on("signal", signal => {
      ws.send(JSON.stringify({
        _type: 'signal',
        data: {
          username: userName,
          username_to_signal: usernameToSignal,
          signal: JSON.stringify(signal)
        }
      }));
    });

    return peer;
  }

  function addPeer(incomingSignal, username, stream) {
    const peer = new Peer({
      initiator: false,
      trickle: false,
      stream,
      config: {
        iceServers: [
          { urls: 'stun:stun.l.google.com:19302' },
          { urls: 'stun:stun1.l.google.com:19302' },
        ]
      }
    });

    peer.on("signal", signal => {
      socketRef.current.send(JSON.stringify({
        _type: 'returning_signal',
        data: {
          username: userName,
          username_to_signal: username,
          signal: JSON.stringify(signal)
        }
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
              <Audio key={peerObj.peerID} peer={peerObj.peer} peerID={peerObj.peerID} />
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
    const handleStream = (stream) => {
      console.log("received stream from", props.peerID);
      if (ref.current) {
        ref.current.srcObject = stream;
      }
    };

    const handleTrack = (track, stream) => {
      console.log("received track from", props.peerID);
      if (ref.current) {
        ref.current.srcObject = stream;
      }
    };

    props.peer.on("stream", handleStream);
    props.peer.on("track", handleTrack);

    // simple-peer might have already received the stream
    // although it's not public API, some versions expose _remoteStreams
    if (props.peer._remoteStreams && props.peer._remoteStreams[0]) {
      handleStream(props.peer._remoteStreams[0]);
    }

    return () => {
      props.peer.off("stream", handleStream);
      props.peer.off("track", handleTrack);
    };
  }, [props.peer, props.peerID]);

  return (
    <div className="user-node">
      <User size={40} />
      <div className="user-name">{props.peerID || 'Remote User'}</div>
      <audio playsInline autoPlay ref={ref} style={{ display: 'none' }} />
    </div>
  );
};

export default App;
