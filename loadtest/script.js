import http from 'k6/http';
import { sleep, check } from 'k6';

export const options = {
  stages: [
    { duration: '30s', target: 1000 }, // ramp up to 20 users
    { duration: '1m30s', target: 500 }, // stay at 20 users for 1.5 mins
    { duration: '30s', target: 1200 }, // ramp down to 0 users
    { duration: '30s', target: 0 }, // ramp down to 0 users
  ],
};

let levels = ['debug', 'info', 'warn', 'error', 'fatal'];
let messages = [
  "Failed to connect to DB",
  "Connection to DB successful",
  "DB connection failed",
  "Some error",
  "Some other error",
  "User does not exist",
]
let resources = [
  "server-1234",
  "server-0987",
  "server-5678",
  "server-0987",
  "server-5678",
];

let traceIds = [
  "abc-xyz-123",
  "abc-xyz-456",
  "abc-xyz-789",
  "abc-xyz-098",
  "abc-xyz-321",
];

let spanIds = [
  "span-123",
  "span-456",
  "span-789",
  "span-098",
  "span-321",
];

export default function() {
  let level = levels[Math.floor(Math.random() * levels.length)];
  let message = messages[Math.floor(Math.random() * messages.length)];
  let resource = resources[Math.floor(Math.random() * resources.length)];
  let traceId = traceIds[Math.floor(Math.random() * traceIds.length)];
  let spanId = spanIds[Math.floor(Math.random() * spanIds.length)];

  let payload = JSON.stringify({
    level: level,
    message: message,
    resourceId: resource,
    timestamp: "2023-09-15T08:00:00Z",
    traceId: traceId,
    spanId: spanId,
    commit: "5e5342f",
    metadata: {
      parentResourceId: "server-0987"
    }
  });

  let params = {
    headers: {
      'Content-Type': 'application/json'
    }
  };

  let res = http.post('http://localhost:3000', payload, params);
  check(res, {
    'status was 202': (r) => r.status == 202,
    'response time was acceptable': (r) => r.timings.duration < 500
  });

  sleep(1);
}
