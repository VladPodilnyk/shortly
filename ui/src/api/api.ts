import axios from 'axios';

const baseUrl = '/v1/encode';

export async function getShortUrl(longUrl: string, alias?: string) {
  const result = await axios.post(baseUrl, { url: longUrl, alias: alias });
  return result;
}
