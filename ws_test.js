import ws from 'k6/ws';
import { Counter } from 'k6/metrics';

export const messages = new Counter('messages');

export const options = {
  vus: 500,
  duration: '30s',
};

export default function () {
  ws.connect('ws://localhost:8080/ws', {}, function (socket) {

    socket.on('open', () => {

      socket.send(JSON.stringify({
        type: 'join',
        username: `user-${__VU}`,
        roomId: 1,
      }));

      socket.send(JSON.stringify({
        type: 'message',
        username: `user-${__VU}`,
        roomId: 1,
        message: 'hello',
      }));
    });

    socket.on('message', (msg) => {
      messages.add(1);
    });

    socket.setTimeout(() => {
      socket.close();
    }, 5000);
  });
} 
