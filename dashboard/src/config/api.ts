// Se VITE_API_URL estiver vazio, usa relative path (mesma origem).
// Quando rodar via Docker compose, o nginx do dashboard faz proxy
// pra API automaticamente em /api/*.
const apiBase = import.meta.env.VITE_API_URL ?? ''

export const API_BASE_URL = apiBase

export const API_KEY = import.meta.env.VITE_API_KEY ?? ''

// apiRequest faz fetch com X-API-Key e resolve URL relativo/absoluto
export function apiRequest(url: string): Promise<Response> {
  const target = apiBase ? apiBase + url : url
  return fetch(target, {
    headers: {
      'X-API-Key': API_KEY,
    },
  })
}
