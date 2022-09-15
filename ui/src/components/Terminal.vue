<template>
  <div>
    <div ref="termRef"></div>
  </div>
</template>

<script lang="ts" setup>
import {Terminal} from "xterm"
import {FitAddon} from 'xterm-addon-fit'
import "xterm/css/xterm.css"
import {onMounted, ref} from "vue";
const termRef = ref(null)
const term = new Terminal({
  fontFamily: '"Cascadia Mono", "Lucida Console", monospace, monaco, Consolas',
  fontSize: 12,
  lineHeight: 1,
  rows:50,
  cols:150,
})
const fitAddon = new FitAddon();
const wsURL = `ws://${window.location.host}/ws/terminal`
let ws: WebSocket = undefined

onMounted(() => {
  initTerminal()
  ws = new WebSocket(wsURL)
  initWS()
})

function message(msg: string) {
  return {
    type: "TerminalInput",
    data: btoa(msg),
  }
}

function decodeUnicode(str) {
  // Going backwards: from bytestream, to percent-encoding, to original string.
  return decodeURIComponent(atob(str).split('').map(function (c) {
    return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2);
  }).join(''));
}

function decodeMessage(data) {
  let resp = JSON.parse(data)
  resp.data = decodeUnicode(resp.data)
  return resp
}


function initWS() {
  ws.onmessage = (evt) => {
    const resp = decodeMessage(evt.data)
    console.log(resp)
    term.write(resp.data)
  }

}

function initTerminal() {
  term.loadAddon(fitAddon);
  term.open(termRef.value);
  fitAddon.fit();
  term.focus();
  term.onData((data) => {
    const req = JSON.stringify(message(data))
    ws.send(req)
  })
}

</script>

<style scoped>

</style>