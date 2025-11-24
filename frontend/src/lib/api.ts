const API_BASE = (import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080').replace(/\/$/, '');

type RequestInitWithBody = RequestInit & { body?: BodyInit | null };

async function request<T>(path: string, init: RequestInitWithBody = {}): Promise<T> {
  const res = await fetch(`${API_BASE}${path}`, {
    headers: {
      'Content-Type': 'application/json',
      ...(init.headers || {}),
    },
    ...init,
  });

  const text = await res.text();
  const data = text ? JSON.parse(text) : null;

  if (!res.ok) {
    const message = data?.error || res.statusText;
    throw new Error(message);
  }

  return data as T;
}

export const api = {
  listUsers: () => request('/api/users'),
  getUser: (id: string) => request(`/api/users/${id}`),
  listRoutes: () => request('/api/routes'),
  getRoute: (id: string) => request(`/api/routes/${id}`),
  listWorkouts: () => request('/api/workouts'),
  getWorkout: (id: string) => request(`/api/workouts/${id}`),
  getReviews: (routeId?: string) => request(`/api/reviews${routeId ? `?routeId=${routeId}` : ''}`),
  createReview: (payload: { userId: string; routeId: string; rating: number; comment: string }) =>
    request('/api/reviews', { method: 'POST', body: JSON.stringify(payload) }),
  getLeaderboard: (limit = 10) => request(`/api/leaderboard?limit=${limit}`),
  upsertScore: (payload: { userId: string; score: number }) =>
    request('/api/leaderboard', { method: 'POST', body: JSON.stringify(payload) }),
  getNotifications: (userId: string) => request(`/api/notifications/${userId}`),
  sendNotification: (payload: { userId: string; message: string }) =>
    request('/api/notifications', { method: 'POST', body: JSON.stringify(payload) }),
  getMap: (routeId: string) => request(`/api/maps/${routeId}`),
  setMap: (payload: { routeId: string; geoJson: string }) =>
    request('/api/maps', { method: 'POST', body: JSON.stringify(payload) }),
};
