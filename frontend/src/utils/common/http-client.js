import axios from "axios";

const isDevelopmentMode = () => process.env.NODE_ENV !== "production";

const BASE_URL =  isDevelopmentMode() ? "http://localhost:8888" : "";
const SECOND = 1000;

const httpClient = axios.create({
    baseURL: BASE_URL,
    timeout: 5 * SECOND,
});

export default httpClient;
