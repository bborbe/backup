import axios, {AxiosResponse} from "axios";

export interface ApiConfig {
    timeout?: number;
}

export class ApiService {
    private client;

    constructor(config: ApiConfig) {
        this.client = axios.create({
            timeout: config.timeout || 10000,
            headers: {
                "Content-Type": "application/json",
            },
        });
    }

    async get<T>(url: string): Promise<T> {
        const response: AxiosResponse<T> = await this.client.get(url);
        return response.data;
    }

    async post<T, R>(url: string, data: T): Promise<R> {
        const response: AxiosResponse<R> = await this.client.post(url, data);
        return response.data;
    }

    async put<T, R>(url: string, data: T): Promise<R> {
        const response: AxiosResponse<R> = await this.client.put(url, data);
        return response.data;
    }

    async delete<T>(url: string): Promise<T> {
        const response: AxiosResponse<T> = await this.client.delete(url);
        return response.data;
    }
}

// Create service instance
const apiService = new ApiService({});
export {apiService};
