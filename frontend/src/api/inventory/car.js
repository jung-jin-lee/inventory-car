import httpClient from "../../utils/common/http-client";

const BASE_URL = "/inventory/cars";

export const getInventoryCars = () => {
    return httpClient.get(`${BASE_URL}`);
};

export const createInventoryCar = data => {
    return httpClient.post(`${BASE_URL}`, JSON.stringify(data));
};

export const deleteInventoryCar = vin => {
    return httpClient.delete(`${BASE_URL}/${vin}`);
};
