const ws = new WebSocket('ws://127.0.0.1:8095/ws/weight');
let count = 0;

ws.onmessage = (event) => {
  console.log(String(event.data));
  count += 1;
  if (count >= 2) {
    ws.close();
    process.exit(0);
  }
};

ws.onerror = (event) => {
  console.error('WS error', event.message ?? '');
  process.exit(1);
};

setTimeout(() => {
  console.error('WS timeout');
  process.exit(1);
}, 6000);
