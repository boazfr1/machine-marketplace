import axios, { AxiosInstance, isAxiosError } from 'axios';

const api: AxiosInstance = axios.create({
  baseURL: 'http://localhost:3001',
  withCredentials: true
});

export { isAxiosError };
export default api;