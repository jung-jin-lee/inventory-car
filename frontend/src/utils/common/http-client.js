import axios from "axios";

const BASE_URL = "http://localhost:8888";
const SECOND = 1000;

const httpClient = axios.create({
    baseURL: BASE_URL,
    timeout: 5 * SECOND,
});

export default httpClient;
