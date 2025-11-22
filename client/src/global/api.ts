import axios, { AxiosInstance, isAxiosError } from 'axios';

export const authApi: AxiosInstance = axios.create({
  baseURL: 'http://localhost:3001',
  withCredentials: true
});

export { isAxiosError };
