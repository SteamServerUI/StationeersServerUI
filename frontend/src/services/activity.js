import { writable } from 'svelte/store';
import { apiJson, apiSSE } from './api-v7';

export const activities = writable([]);
export const toasts = writable([]);

let stream;
let initialized = false;

export async function initializeActivity() {
  if (initialized) return () => closeActivity();
  initialized = true;
  try {
    const response = await apiJson('/api/v3/activity?limit=150');
    activities.set((response.data?.events || []).map(normalizeEvent));
  } catch (error) {
    console.warn('Could not load activity history', error);
  }
  stream = apiSSE('/api/v3/streams/activity', event => {
    const item = normalizeEvent(event);
    activities.update(items => [item, ...items.filter(existing => existing.id !== item.id)].slice(0, 500));
  }, error => console.warn('Activity stream disconnected', error));
  return () => closeActivity();
}

export function closeActivity() {
  stream?.close?.();
  stream = null;
  initialized = false;
}

export function notify(message, tone = 'info', detail = '') {
  const item = { id: crypto.randomUUID(), type: 'local', message, detail, tone, at: new Date().toISOString() };
  activities.update(items => [item, ...items].slice(0, 500));
  toasts.update(items => [...items, item]);
  setTimeout(() => dismissToast(item.id), tone === 'error' ? 9000 : 5000);
  return item;
}

export function dismissToast(id) {
  toasts.update(items => items.filter(item => item.id !== id));
}

function normalizeEvent(event) {
  return {
    ...event,
    id: event.id || crypto.randomUUID(),
    message: event.message || 'Activity',
    tone: event.tone || 'info',
    at: event.at || event.createdAt || new Date().toISOString()
  };
}
