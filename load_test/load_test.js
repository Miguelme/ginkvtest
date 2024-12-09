import http from 'k6/http';
import { check, sleep } from 'k6';
import { Trend, Rate } from 'k6/metrics';

// Custom metrics
const auroraDuration = new Trend('aurora_duration');
const dynamoDuration = new Trend('dynamo_duration');
const redisDuration = new Trend('redis_duration');
const errorRate = new Rate('error_rate');

export const options = {
  stages: [
    { duration: '5s', target: 10 },  // Ramp up to 100 users
    { duration: '5s', target: 20 },  // Ramp up to 500 users
    { duration: '10s', target: 50 }, // Ramp up to 1000 users
    { duration: '5s', target: 0 },    // Ramp down
  ],
  thresholds: {
    error_rate: ['rate<0.05'], // Error rate should be less than 5%
  },
};

export default function () {
  const url = 'http://localhost:8080/benchmark/test_key'; // Replace with your endpoint
  const res = http.get(url);

  // Check for a 200 response
  const success = check(res, {
    'is status 200': (r) => r.status == 200,
  });

  if (!success) {
    console.error(`Request failed. Status: ${res.status}, Body: ${res.body}`);
  }
  if (!success) {
    errorRate.add(1); // Log error
  }

  // Parse the JSON response
  const jsonResponse = res.json();
  const aurora = jsonResponse.aurora;
  const dynamo = jsonResponse.dynamo;
  const redis = jsonResponse.redis;

  // Track durations
  if (aurora && aurora.duration) auroraDuration.add(parseDuration(aurora.duration));
  if (dynamo && dynamo.duration) dynamoDuration.add(parseDuration(dynamo.duration));
  if (redis && redis.duration) redisDuration.add(parseDuration(redis.duration));

  sleep(0.1); // Small sleep to avoid overwhelming the server
}

// Helper function to convert duration strings like "5.65ms" to milliseconds
function parseDuration(duration) {
  return duration / 1000;
}

