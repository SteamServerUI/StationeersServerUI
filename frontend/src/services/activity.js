import { writable } from 'svelte/store';

export const activities = writable([]);
export const toasts = writable([]);

export function notify(message, tone = 'info', detail = '') {
  const item = { id: crypto.randomUUID(), message, detail, tone, at: new Date().toISOString() };
  activities.update(items => [item, ...items].slice(0, 100));
  toasts.update(items => [...items, item]);
  setTimeout(() => dismissToast(item.id), tone === 'error' ? 9000 : 5000);
  return item;
}

export function dismissToast(id) {
  toasts.update(items => items.filter(item => item.id !== id));
}
