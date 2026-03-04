import axios from 'axios';

export const apiClient = axios.create({
    baseURL: (import.meta.env.API_URL as string | undefined) ?? 'localhost:8080',
    withCredentials: true,
});
