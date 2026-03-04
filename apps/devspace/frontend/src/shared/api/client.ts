import axios from 'axios';

export const apiClient = axios.create({
    baseURL: (import.meta.env.VITE_API_URL as string | undefined) ?? 'http://localhost:8080',
    withCredentials: true,
});
