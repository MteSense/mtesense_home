import type { ApiEnvelope, NavGroup, NavLink, PublicSettings, User } from './types'

const API_BASE = '/api/v1'

export function getToken() {
  return localStorage.getItem('mtesense_token')
}

export function setToken(token: string) {
  localStorage.setItem('mtesense_token', token)
}

export function clearToken() {
  localStorage.removeItem('mtesense_token')
}

async function request<T>(path: string, options: RequestInit = {}): Promise<T> {
  const headers = new Headers(options.headers)
  if (!(options.body instanceof FormData)) {
    headers.set('Content-Type', 'application/json')
  }
  const token = getToken()
  if (token) {
    headers.set('Authorization', `Bearer ${token}`)
  }
  const response = await fetch(`${API_BASE}${path}`, { ...options, headers })
  const body = (await response.json().catch(() => ({}))) as ApiEnvelope<T>
  if (!response.ok || body.error) {
    throw new Error(body.error || `Request failed: ${response.status}`)
  }
  return body.data as T
}

export const api = {
  login: (username: string, password: string) =>
    request<{ token: string; user: User }>('/auth/login', {
      method: 'POST',
      body: JSON.stringify({ username, password })
    }),
  me: () => request<User>('/me'),
  navigation: () => request<NavGroup[]>('/navigation'),
  adminNavigation: () => request<NavGroup[]>('/admin/navigation'),
  settings: () => request<PublicSettings>('/settings'),
  saveSettings: (payload: PublicSettings) =>
    request<PublicSettings>('/admin/settings', {
      method: 'PUT',
      body: JSON.stringify(payload)
    }),
  createGroup: (payload: Partial<NavGroup>) =>
    request<NavGroup>('/admin/groups', {
      method: 'POST',
      body: JSON.stringify(payload)
    }),
  updateGroup: (id: number, payload: Partial<NavGroup>) =>
    request<NavGroup>(`/admin/groups/${id}`, {
      method: 'PUT',
      body: JSON.stringify(payload)
    }),
  deleteGroup: (id: number) => request<{ deleted: boolean }>(`/admin/groups/${id}`, { method: 'DELETE' }),
  createLink: (payload: Partial<NavLink>) =>
    request<NavLink>('/admin/links', {
      method: 'POST',
      body: JSON.stringify(payload)
    }),
  updateLink: (id: number, payload: Partial<NavLink>) =>
    request<NavLink>(`/admin/links/${id}`, {
      method: 'PUT',
      body: JSON.stringify(payload)
    }),
  deleteLink: (id: number) => request<{ deleted: boolean }>(`/admin/links/${id}`, { method: 'DELETE' }),
  upload: (file: File) => {
    const form = new FormData()
    form.set('file', file)
    return request<{ url: string; filename: string }>('/admin/uploads', {
      method: 'POST',
      body: form
    })
  }
}
