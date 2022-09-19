<template>
  <div>
    <div ref="termRef"></div>
  </div>
</template>

<script lang="ts" setup>
import {Terminal} from "xterm"
import {FitAddon} from 'xterm-addon-fit'
import {CanvasAddon} from 'xterm-addon-canvas';
import "xterm/css/xterm.css"
import {onMounted, onUnmounted, ref} from "vue";

const termRef = ref<HTMLElement>()
let ws: WebSocket;

const term = new Terminal({
  fontFamily: '"Cascadia Mono", "Lucida Console", monospace, monaco, Consolas',
  fontSize: 12,
  lineHeight: 1,
  rows: 50,
  cols: 160,
  allowProposedApi: true,
})
const fitAddon = new FitAddon();
const wsURL = `ws://${window.location.host}/ws/terminal`
onMounted(() => {
  initTerminal()
  ws = new WebSocket(wsURL)
  initWS()
})
onUnmounted(() => {
  ws.close()
})

function sendInput(msg: string) {
  const data = {
    req_type: "TerminalInput",
    data: msg,
  }
  ws.send(JSON.stringify(data))
}

function sendResize(cols: number, rows: number) {
  const data = {
    req_type: "TerminalResize",
    data: {
      h: rows,
      w: cols,
    }
  }
  ws.send(JSON.stringify(data))
}

function decodeMessage(data: string) {
  return  JSON.parse(data)
}


function initWS() {
  ws.onopen = (evt) => {
    sendResize(term.cols, term.rows)
  }
  ws.onmessage = (evt) => {
    const resp = decodeMessage(evt.data)
    term.write(resp.data)
  }

}

function initTerminal() {
  term.loadAddon(fitAddon);
  if (termRef.value === undefined) {
    alert("未找到terminal父元素")
    return
  }
  term.open(termRef.value);
  term.loadAddon(new CanvasAddon());
  fitAddon.fit();
  term.focus();
  term.onData((data) => {
    sendInput(data)
  })
  addEventListener("resize", (evt) => {
    fitAddon.fit()
    sendResize(term.cols, term.rows)
  })
}

</script>

<style scoped>

</style>