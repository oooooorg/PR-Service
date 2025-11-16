import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate } from 'k6/metrics';

const errorRate = new Rate('errors');

export const options = {
    stages: [
        { duration: '30s', target: 10 },
        { duration: '1m', target: 50 },
        { duration: '2m', target: 50 },
        { duration: '1m', target: 100 },
        { duration: '2m', target: 100 },
        { duration: '30s', target: 0 },
    ],

    thresholds: {
        'http_req_duration': ['p(95)<500', 'p(99)<1000'],
        'http_req_failed': ['rate<0.05'],
        'errors': ['rate<0.1'],
    },
};

const BASE_URL = 'http://app:8080';
const users = ['a12', 'a22', 'a32', 'a42', 'a52'];
const teamName = 'load-test-team';

function randomUser() {
    return users[Math.floor(Math.random() * users.length)];
}

function randomPRId() {
    const timestamp = Date.now();
    const random = Math.floor(Math.random() * 1000000);
    const vuId = __VU;
    const iter = __ITER;
    return `pr-test-${timestamp}-${vuId}-${iter}-${random}`;
}

export function setup() {
    const createTeamPayload = JSON.stringify({
        team_name: teamName,
        members: [
            {
                user_id: 'a12',
                username: 'Alice',
                is_active: true
            },
            {
                user_id: 'a22',
                username: 'Bob',
                is_active: true
            },
            {
                user_id: 'a32',
                username: 'Charlie',
                is_active: true
            },
            {
                user_id: 'a42',
                username: 'Diana',
                is_active: true
            },
            {
                user_id: 'a52',
                username: 'Eve',
                is_active: true
            }
        ]
    });

    const createTeamRes = http.post(`${BASE_URL}/team/add`, createTeamPayload, {
        headers: { 'Content-Type': 'application/json' },
    });

    sleep(2);

    return { users, teamName };
}

export default function (data) {
    const createPRPayload = JSON.stringify({
        pull_request_id: randomPRId(),
        pull_request_name: `Load Test PR ${Date.now()}`,
        author_id: randomUser(),
    });

    const createPRRes = http.post(`${BASE_URL}/pullRequest/create`, createPRPayload, {
        headers: { 'Content-Type': 'application/json' },
    });

    const createSuccess = check(createPRRes, {
        'PR created successfully': (r) => r.status === 201,
        'PR create response time < 500ms': (r) => r.timings.duration < 500,
    });

    if (!createSuccess) {
        errorRate.add(1);
    } else {
        errorRate.add(0);
    }

    sleep(1);

    const userId = randomUser();
    const getReviewersRes = http.get(`${BASE_URL}/users/getReview?user_id=${userId}`);

    check(getReviewersRes, {
        'Get reviewers successful': (r) => r.status === 200,
        'Get reviewers response time < 300ms': (r) => r.timings.duration < 300,
    });

    sleep(0.5);

    const getReviewRes2 = http.get(`${BASE_URL}/users/getReview?user_id=${randomUser()}`);

    check(getReviewRes2, {
        'Get PR list successful': (r) => r.status === 200,
    });

    sleep(1);
}

export function teardown(data) {
    console.log('Teardown completed');
    console.log(`Test finished for team: ${data.teamName}`);
}

export function handleSummary(data) {
    return {
        'summary.json': JSON.stringify(data),
    };
}
