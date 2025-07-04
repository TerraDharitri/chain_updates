import http from 'k6/http';
import { sleep } from 'k6';
import { Trend } from 'k6/metrics';

const BASE_URL = 'http://localhost:3001';

const accountsApiCallTrend = new Trend('accounts_http_req_duration', true);
const blocksApiCallTrend = new Trend('blocks_http_req_duration', true);
const drtPairsApiCallTrend = new Trend('moa_pairs_http_req_duration', true);
const moaTokensApiCallTrend = new Trend('moa_tokens_http_req_duration', true);
const moaFarmsApiCallTrend = new Trend('moa_farms_http_req_duration', true);
const nodesApiCallTrend = new Trend('nodes_http_req_duration', true);
const nodesAuctionsApiCallTrend = new Trend('nodes_auctions_http_req_duration', true);
const poolApiCallTrend = new Trend('pool_http_req_duration', true);
const tokensApiCallTrend = new Trend('tokens_http_req_duration', true);
const transactionsApiCallTrend = new Trend('transactions_http_req_duration', true);


function getScenarioDict(functionName) {
    return {
        executor: 'constant-vus',
        vus: 10,
        duration: '1m',
        gracefulStop: '0s',
        exec: functionName,
    }
}

export const options = {
    scenarios: {
        accounts: getScenarioDict('accounts'),
        blocks: getScenarioDict('blocks'),
        drtPairs: getScenarioDict('drtPairs'),
        moaTokens: getScenarioDict('moaTokens'),
        moaFarms: getScenarioDict('moaFarms'),
        nodes: getScenarioDict('nodes'),
        nodesAuctions: getScenarioDict('nodesAuctions'),
        pool: getScenarioDict('pool'),
        tokens: getScenarioDict('tokens'),
        transactions: getScenarioDict('transactions'),
    },
    discardResponseBodies: true,
};

export function accounts() {
    const response = http.get(`${BASE_URL}/accounts`);
    accountsApiCallTrend.add(response.timings.duration);
}

export function blocks() {
    const response = http.get(`${BASE_URL}/blocks`);
    blocksApiCallTrend.add(response.timings.duration);
}

export function drtPairs() {
    const response = http.get(`${BASE_URL}/moa/pairs`);
    drtPairsApiCallTrend.add(response.timings.duration);
}

export function moaTokens() {
    const response = http.get(`${BASE_URL}/moa/tokens`);
    moaTokensApiCallTrend.add(response.timings.duration);
}

export function moaFarms() {
    const response = http.get(`${BASE_URL}/moa/farms`);
    moaFarmsApiCallTrend.add(response.timings.duration);
}

export function nodes() {
    const response = http.get(`${BASE_URL}/nodes`);
    nodesApiCallTrend.add(response.timings.duration);
}

export function nodesAuctions() {
    const response = http.get(`${BASE_URL}/nodes/auctions`);
    nodesAuctionsApiCallTrend.add(response.timings.duration);
}

export function pool() {
    const response = http.get(`${BASE_URL}/pool`);
    poolApiCallTrend.add(response.timings.duration);
}

export function tokens() {
    const response = http.get(`${BASE_URL}/tokens`);
    tokensApiCallTrend.add(response.timings.duration);
}

export function transactions() {
    const response = http.get(`${BASE_URL}/transactions`);
    transactionsApiCallTrend.add(response.timings.duration);
}

export function handleSummary(data) {
  return {
    'k6/output/summary.json': JSON.stringify(data),
  };
}
