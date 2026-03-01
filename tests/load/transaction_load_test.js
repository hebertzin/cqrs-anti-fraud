/**
 * k6 Load Test — Anti-Fraud Service
 *
 * Run: k6 run tests/load/transaction_load_test.js
 * Run with env: k6 run -e BASE_URL=http://localhost:8080 tests/load/transaction_load_test.js
 */

import http from "k6/http";
import { check, sleep } from "k6";
import { Rate, Trend } from "k6/metrics";

const errorRate = new Rate("errors");
const analyzeLatency = new Trend("analyze_latency", true);

export const options = {
  stages: [
    { duration: "30s", target: 20 },   // ramp-up
    { duration: "1m",  target: 50 },   // sustained load
    { duration: "30s", target: 100 },  // stress
    { duration: "30s", target: 0 },    // ramp-down
  ],
  thresholds: {
    http_req_failed:   ["rate<0.01"],   // < 1% errors
    http_req_duration: ["p(95)<500"],   // 95th percentile < 500ms
    analyze_latency:   ["p(99)<1000"],  // 99th percentile < 1s
    errors:            ["rate<0.05"],
  },
};

const BASE_URL = __ENV.BASE_URL || "http://localhost:8080";

const locations = ["BR", "US", "AR", "XX", "ZZ", "PT", "DE"];
const currencies = ["BRL", "USD", "EUR", "GBP"];

function randomElement(arr) {
  return arr[Math.floor(Math.random() * arr.length)];
}

function randomAmount() {
  return parseFloat((Math.random() * 20000).toFixed(2));
}

function randomUUID() {
  return "xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx".replace(/[xy]/g, (c) => {
    const r = (Math.random() * 16) | 0;
    return (c === "x" ? r : (r & 0x3) | 0x8).toString(16);
  });
}

export default function () {
  const payload = JSON.stringify({
    account_id:  randomUUID(),
    amount:      randomAmount(),
    currency:    randomElement(currencies),
    merchant_id: `merchant-${Math.floor(Math.random() * 1000)}`,
    location:    randomElement(locations),
  });

  const params = {
    headers: { "Content-Type": "application/json" },
    timeout: "10s",
  };

  const start = Date.now();
  const res = http.post(`${BASE_URL}/api/v1/transactions`, payload, params);
  analyzeLatency.add(Date.now() - start);

  const success = check(res, {
    "status is 201": (r) => r.status === 201,
    "has transaction_id": (r) => {
      try {
        const body = JSON.parse(r.body);
        return body.transaction_id !== undefined;
      } catch {
        return false;
      }
    },
    "has risk_score": (r) => {
      try {
        const body = JSON.parse(r.body);
        return typeof body.risk_score === "number";
      } catch {
        return false;
      }
    },
  });

  errorRate.add(!success);
  sleep(0.1);
}

export function healthCheck() {
  const res = http.get(`${BASE_URL}/health`);
  check(res, { "health ok": (r) => r.status === 200 });
}
