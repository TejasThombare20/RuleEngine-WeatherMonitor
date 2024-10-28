import axios, { AxiosInstance, AxiosRequestConfig, AxiosResponse } from "axios";

const BASE_URL =
  process.env.NODE_ENV === "production"
    ? "https://weathermonitor.onrender.com"
    : "http://localhost:9000";

const DEFAULT_HEADERS = {
  "Content-Type": "application/json",
};

// Create a type for the API methods
type ApiMethod = "get" | "post" | "put" | "delete" | "patch";

class ApiHandler {
  private axiosInstance: AxiosInstance;

  constructor() {
    this.axiosInstance = axios.create({
      baseURL: BASE_URL,
      headers: DEFAULT_HEADERS,
    });
  }

  private async request<T>(
    method: ApiMethod,
    url: string,
    data?: any,
    config?: AxiosRequestConfig
  ): Promise<T> {
    try {
      const response: AxiosResponse<T> = await this.axiosInstance.request({
        method,
        url,
        data,
        ...config,
      });
      return response.data;
    } catch (error) {
      console.error("API request failed:", error);
      throw error;
    }
  }

  public get<T>(url: string, config?: AxiosRequestConfig): Promise<T> {
    return this.request<T>("get", url, undefined, config);
  }

  public post<T>(
    url: string,
    data?: any,
    config?: AxiosRequestConfig
  ): Promise<T> {
    return this.request<T>("post", url, data, config);
  }

  public put<T>(
    url: string,
    data?: any,
    config?: AxiosRequestConfig
  ): Promise<T> {
    return this.request<T>("put", url, data, config);
  }

  public delete<T>(url: string, config?: AxiosRequestConfig): Promise<T> {
    return this.request<T>("delete", url, undefined, config);
  }

  public patch<T>(
    url: string,
    data?: any,
    config?: AxiosRequestConfig
  ): Promise<T> {
    return this.request<T>("patch", url, data, config);
  }
}

export default new ApiHandler();
