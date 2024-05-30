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

export default function() {
  let payload = JSON.stringify({
    level: "dude",
    message: "Failed to connect to DB",
    resourceId: "server-1234",
    timestamp: "2023-09-15T08:00:00Z",
    traceId: "abc-xyz-123",
    spanId: "span-456",
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
