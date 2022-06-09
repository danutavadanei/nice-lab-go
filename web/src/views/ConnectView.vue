<template>
  <div class="px-4 sm:px-6 lg:px-8">
    <div class="sm:flex sm:items-center">
      <div class="sm:flex-auto">
        <h1 class="text-xl font-semibold text-gray-900">Connecting to lab #{{ lab }}</h1>
      </div>
    </div>
    <div class="mt-8 flex flex-col">
      <div id="dcv-display"></div>
    </div>
  </div>
</template>

<script setup>
import { inject, onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import dcv from '../../public/vendor/dcvjs/dcv.js'

const route = useRoute()
const lab = parseInt(route.params.lab)


let auth, connection, serverUrl;
console.log("Using NICE DCV Web Client SDK version " + dcv.version.versionStr);

const onPromptCredentials = function (auth, challenge) {
  if (challengeHasField(challenge, "username") && challengeHasField(challenge, "password")) {
    auth.sendCredentials({username: (lab === 1 ? "administrator" : "ubuntu"), password: "rAdiC203094=0"})
  } else {
    console.log(challenge)
  }
}

const challengeHasField = function (challenge, field) {
  return challenge.requiredCredentials.some(credential => credential.name === field);
}

const onError = function (auth, error) {
  console.log("Error during the authentication: " + error.message);
}

const onSuccess = function (auth, result) {
  let {sessionId, authToken} = {...result[0]};

  connect(sessionId, authToken);
}

const connect = function (sessionId, authToken) {
  console.log(sessionId, authToken);

  dcv.connect({
    url: serverUrl,
    sessionId: sessionId,
    authToken: authToken,
    divId: "dcv-display",
    baseUrl: "https://nice-lab.utm.test/vendor/dcvjs",
    callbacks: {
      firstFrame: () => console.log("First frame received")
    }
  }).then(function (conn) {
    console.log("Connection established!");
    connection= conn;
  }).catch(function (error) {
    console.log("Connection failed with error " + error.message);
  });
}

const main = function () {
  console.log("Setting log level to INFO");
  dcv.setLogLevel(dcv.LogLevel.INFO);
  if (lab === 1) {
    serverUrl = "https://ec2-3-94-101-59.compute-1.amazonaws.com:8443/";
  } else if (lab === 2) {
    serverUrl = "https://ec2-44-203-45-241.compute-1.amazonaws.com:8443/";
  }

  console.log("Starting authentication with", serverUrl);

  auth = dcv.authenticate(
    serverUrl,
    {
      promptCredentials: onPromptCredentials,
      error: onError,
      success: onSuccess
    }
  );
}

onMounted(main)
</script>
