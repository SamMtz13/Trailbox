const API_BASE =
  (import.meta.env.VITE_API_BASE_URL as string | undefined)?.replace(/\/$/, '') ?? 'http://localhost:8080';

async function request<T>(path: string, options: RequestInit = {}): Promise<T> {
  const headers = new Headers(options.headers);
  if (options.body && !headers.has('Content-Type')) {
    headers.set('Content-Type', 'application/json');
  }

  const response = await fetch(`${API_BASE}${path}`, {
    ...options,
    headers
  });

  if (!response.ok) {
    let message = response.statusText;
    try {
      const payload = await response.json();
      if (typeof payload?.error === 'string') {
        message = payload.error;
      }
    } catch {
      // ignore JSON parse errors
    }
    throw new Error(message || 'Request failed');
  }

  if (response.status === 204) {
    return {} as T;
  }

  return (await response.json()) as T;
}

export interface User {
  id: string;
  name: string;
  email: string;
}

export interface Route {
  id: string;
  name: string;
  distance_km?: number;
  elevation_gain?: number;
}

export interface Workout {
  id: string;
  user_id: string;
  route_id: string;
  date: string;
  duration: number;
  calories: number;
}

export interface Review {
  id: string;
  user_id: string;
  route_id: string;
  rating: number;
  comment: string;
  created_at: string;
}

export interface LeaderboardEntry {
  id: string;
  user_id: string;
  score: number;
  position: number;
}

export interface Notification {
  id: string;
  user_id: string;
  message: string;
  read: boolean;
  created_at: string;
}

export interface MapResponse {
  route_id: string;
  geo_json: string;
  created_at?: string;
}

export async function listUsers() {
  return request<User[]>('/api/users');
}

export async function getUser(id: string) {
  return request<User>(`/api/users/${id}`);
}

export async function listRoutes() {
  return request<Route[]>('/api/routes');
}

export async function listWorkouts() {
  return request<Workout[]>('/api/workouts');
}

export async function getReviews(routeId: string) {
  return request<Review[]>(`/api/reviews?route_id=${encodeURIComponent(routeId)}`);
}

export async function createReview(payload: {
  userId: string;
  routeId: string;
  rating: number;
  comment: string;
}) {
  return request<Review>('/api/reviews', {
    method: 'POST',
    body: JSON.stringify(payload)
  });
}

export async function getLeaderboard(limit = 10) {
  return request<LeaderboardEntry[]>(`/api/leaderboard?limit=${limit}`);
}

export async function upsertLeaderboard(payload: { userId: string; score: number }) {
  return request('/api/leaderboard', {
    method: 'POST',
    body: JSON.stringify(payload)
  });
}

export async function getNotifications(userId: string) {
  return request<Notification[]>(`/api/notifications?user_id=${encodeURIComponent(userId)}`);
}

export async function sendNotification(payload: { userId: string; message: string }) {
  return request<Notification>('/api/notifications', {
    method: 'POST',
    body: JSON.stringify(payload)
  });
}

export async function getRouteMap(routeId: string) {
  return request<MapResponse>(`/api/maps/${encodeURIComponent(routeId)}`);
}

export async function saveRouteMap(payload: { routeId: string; geoJson: string }) {
  return request('/api/maps', {
    method: 'POST',
    body: JSON.stringify(payload)
  });
}
