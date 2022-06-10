<template>
  <div class="px-4 sm:px-6 lg:px-8">
    <div class="sm:flex sm:items-center">
      <div class="sm:flex-auto">
        <h1 class="text-xl font-semibold text-gray-900">Monitoring session #{{ session }}</h1>
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

const axios = inject('axios')
const route = useRoute()
const session = route.params.session
const apiBaseUrl = inject('apiBaseUrl')
const getSessionInfoEndpoint = `${apiBaseUrl}/v1/sessions`;

let auth, connection, serverUrl, username, password;
console.log("Using NICE DCV Web Client SDK version " + dcv.version.versionStr);

const onPromptCredentials = function (auth, challenge) {
  if (challengeHasField(challenge, "username") && challengeHasField(challenge, "password")) {
    auth.sendCredentials({ username, password })
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

onMounted(async () => {
  const response = await axios.get(getSessionInfoEndpoint + '/' + session)
  username = response.data.username
  password = response.data.password
  serverUrl = `https://${response.data.hostname}:8443/`

  main()
})
</script>
