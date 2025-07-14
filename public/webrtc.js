let ownVideo = document.getElementById("#ownVideo");
let peerVideo = document.getElementById("#peerVideo");

let localStream

let ws;

//startCamera starts camera only if user gives permission
async function startCamera(){
    try{
         ownVideo = await navigator.mediaDevices.getUserMedia(
            {
                video: true,
                audio: true
            }
        );
        ownVideo.srcObject = localStream

    }catch(err){
        const errReply = document.createElement("div")
        errReply.textContent = "camera access denied"
        document.body.appendChild(errReply);
        console.error("failed to start camera,access denied", err)

    }

}

// connectPeer creates peer connection and shows other peer video
function connectPeer(){
    let peerConnection = new RTCPeerConnection();

    peerConnection.onicecandidate = (event) => {
        if (event.candidate){
            const msg = {
                candidate: event.candidate
            };
            const json = JSON.stringify(msg)
            ws.send(json)
        }
    }

    peerConnection.ontrack = (event) => {
        peerVideo.srcObject = event.streams[0]
    };

}