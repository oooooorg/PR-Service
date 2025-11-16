import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
    vus: 1,
    duration: '30s',
};

const BASE_URL = 'http://app:8080';

export default function () {
    const metricsRes = http.get(`${BASE_URL}/metrics`);
    check(metricsRes, { 'metrics available': (r) => r.status === 200 });

    sleep(1);
}
