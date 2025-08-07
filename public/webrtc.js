const ownVideo = document.getElementById("ownVideo");
const peerVideo = document.getElementById("peerVideo");
const shareLink = document.getElementById("shareLink");

let localStream;
let peerConnection;
let ws;
let roomId;
let isStreamer = false; 

async function startVideoCall() {
  isStreamer = true;
  await startCamera();
  generateLink(); 
  setupWSConnection(roomId);
}


//startCamera starts camera only if user gives permission
async function startCamera() {
  try {
    localStream = await navigator.mediaDevices.getUserMedia({
      video: true,
      audio: true
    });
    ownVideo.srcObject = localStream;
  } catch (err) {
    alert("Camera access denied");
    console.error("failed to start camera, access denied", err);
  }
}

// connectPeer creates a peer connection and shows the other peer's video
function connectPeer() {
  peerConnection = new RTCPeerConnection();

  peerConnection.onicecandidate = (event) => {
    if (event.candidate) {
      const msg = {
        candidate: event.candidate
      };
      const json = JSON.stringify(msg)
      ws.send(json)
    }
  };

  peerConnection.ontrack = (event) => {
    peerVideo.srcObject = event.streams[0];
  };

  return peerConnection;
}

// setupWSConnection opens WebSocket connection and handles messages between peers
function setupWSConnection() {
  ws = new WebSocket(`ws://localhost:8080/ws?room=` + roomId);

  ws.onmessage = async (event) => {
    const data = JSON.parse(event.data);

    switch (data.type) {
      case "offer":
        if (!isStreamer) {
          await handleOffer(data.offer);
        }
        break;
      case "answer":
        if (isStreamer) {
          await peerConnection.setRemoteDescription(new RTCSessionDescription(data.answer));
          shareLink.parentElement.style.display = "block";
        }
        break;
      case "candidate":
        await peerConnection.addIceCandidate(new RTCIceCandidate(data.candidate));
        break;
      default:
        console.warn("data type not recognized", data.type);
    }
  };

  ws.onopen = async () => {
    peerConnection = connectPeer();

    if (isStreamer) {
      localStream.getTracks().forEach((track) => {
        peerConnection.addTrack(track, localStream);
      });

      const offer = await peerConnection.createOffer();
      await peerConnection.setLocalDescription(offer);

      const msg = {
        type: "offer",
        offer: offer
      };

      const json = JSON.stringify(msg)
      ws.send(json)
    }
  };
}

// handleOffer accepts a call offer and sends back an answer
async function handleOffer(offer) {
  peerConnection = connectPeer();
  await peerConnection.setRemoteDescription(new RTCSessionDescription(offer));

  const answer = await peerConnection.createAnswer();
  await peerConnection.setLocalDescription(answer);

  const msg = {
    type: "answer",
    answer: answer
  };

  const json = JSON.stringify(msg)
  ws.send(json)
}

//generateLink generates a unique link with a room ID and displays it
function generateLink() {
  roomId = location.hash.slice(1) || crypto.randomUUID();
  location.hash = roomId;

  const link = `${location.origin}${location.pathname}#${roomId}`;
  shareLink.textContent = link;
  shareLink.parentElement.style.display = "block";
}

//copyLink Copies the generated link to the clipboard
function copyLink() {
  const text = shareLink.textContent;
  navigator.clipboard.writeText(text)
    .then(() => 
      alert("Link copied to clipboard"))
    .catch(() =>
       alert("Failed to copy link"));
}

//endCall ends the video call and cleans up media streams and connection
function endCall() {
  if (localStream) {
    localStream.getTracks().forEach(track => track.stop());
  }

  if (peerConnection) {
    peerConnection.close();
  }
  ownVideo.srcObject = null;
  peerVideo.srcObject = null;
  alert("Call ended");
}

