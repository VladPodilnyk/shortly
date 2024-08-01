import axios, { AxiosResponse } from 'axios';

const baseUrl = '/v1/encode';

export interface EncodeResponse {
  short_url: string;
}

export interface Validation {
  url?: string;
  alias?: string;
}

export interface ErrorResponse {
  error: { url?: string; alias?: string; }
}

export async function getShortUrl(longUrl: string, alias?: string): Promise<EncodeResponse> {
  const result: AxiosResponse<EncodeResponse> = await axios.post(baseUrl, { url: longUrl, alias: alias });
  return result.data;
}
