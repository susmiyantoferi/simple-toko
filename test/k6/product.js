import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate, Trend } from 'k6/metrics';

export let errorRate = new Rate('errors');
export let responseTimeTrend = new Trend('response_time');

export let options = {
    stages: [
        { duration: '10s', target: 50 },  // ramp up to 50 VUs
        { duration: '50s', target: 100 }, // ramp up to 100 VUs
        { duration: '10s', target: 0 },   // ramp down to 0 VUs
    ],
    thresholds: {
        'http_req_duration': ['p(95)<100'],  // 95% request < 100ms
        'errors': ['rate==0'],                // error rate 0%
    },
};

const BASE_URL = 'http://127.0.0.1:8080/api/v1/product';
const AUTH_TOKEN = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoyLCJ1c2VybmFtZSI6ImFkbWluIiwiZW1haWwiOiJzdXBlckBlbWFpbC5jb20iLCJyb2xlIjoiYWRtaW4iLCJleHAiOjE3NTgwOTgzMDh9.ekXnpsigH_VO2S7QP85xfYasji_MVDJZt0d7R9xK4cU';

export default function () {
    let params = {
        headers: {
            'Authorization': `Bearer ${AUTH_TOKEN}`,
        },
    };

    let res = http.get(BASE_URL, params);

    // Simpan response time untuk metrik custom
    responseTimeTrend.add(res.timings.duration);

    // Cek status dan validasi response
    let checkRes = check(res, {
        'status is 200': (r) => r.status === 200,
        'response time < 100ms': (r) => r.timings.duration < 100,
    });

    if (!checkRes) {
        errorRate.add(1);
    } else {
        errorRate.add(0);
    }

    //sleep(1);
}
