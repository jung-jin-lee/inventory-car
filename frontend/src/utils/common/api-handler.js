export const handleApiResponse = res => {
    if (res.data.success) {
        return res.data.data;
    }
    throw new Error(JSON.stringify(res.data.err));
};
