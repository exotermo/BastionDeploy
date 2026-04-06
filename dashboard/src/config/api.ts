export const API_BASE_URL = import.meta.env.VITE_API_URL ?? 'http://localhost:8080'
export const API_KEY = import.meta.env.VITE_API_KEY ?? ''

// apiRequest é um fetch wrapper que injeta X-API-Key automaticamente
export function apiRequest(url: string): Promise<Response> {
  return fetch(url, {
    headers: {
      'X-API-Key': API_KEY,
    },
  })
}
