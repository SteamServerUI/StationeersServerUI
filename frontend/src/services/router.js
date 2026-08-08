import { writable } from 'svelte/store';

export const routes = {
  home: '/',
  console: '/operate/console',
  logs: '/operate/logs',
  activity: '/operate/activity',
  game: '/configure/game',
  gallery: '/configure/gallery',
  files: '/configure/files',
  backups: '/protect/backups',
  access: '/access',
  settings: '/system/settings',
  backends: '/system/backends',
  plugins: '/system/plugins'
};

const routeByPath = Object.fromEntries(Object.entries(routes).map(([id, path]) => [path, id]));
export const activeRoute = writable('home');
let ready = false;

export function initializeRouter() {
  if (ready) return;
  ready = true;
  syncFromLocation();
  window.addEventListener('popstate', syncFromLocation);
}

export function navigate(id, replace = false) {
  const path = routes[id] || routes.home;
  const method = replace ? 'replaceState' : 'pushState';
  if (window.location.pathname !== path) window.history[method]({}, '', path);
  activeRoute.set(routeByPath[path] || 'home');
}

function syncFromLocation() {
  const id = routeByPath[window.location.pathname];
  if (id) activeRoute.set(id);
  else navigate('home', true);
}
