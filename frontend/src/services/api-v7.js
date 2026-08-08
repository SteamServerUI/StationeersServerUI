import { get, writable } from 'svelte/store';

const storageKey = 'ssui-backend-config-v2';

export const backendConfig = writable({
  active: 'default',
  backends: { default: { url: '/' } }
});

export const authState = writable({
  isAuthenticated: false,
  isAuthenticating: false,
  authError: null,
  setupRequired: false,
  user: null,
  permissions: [],
  features: { plugins: false },
  csrf: null,
  expiresAt: null
});

let initialized = false;

export function initializeApiService() {
  if (initialized) return;
  initialized = true;
  try {
    const saved = JSON.parse(localStorage.getItem(storageKey) || 'null');
    if (saved?.backends && saved.backends[saved.active]) {
      const backends = {};
      for (const [id, backend] of Object.entries(saved.backends)) {
        if (typeof backend?.url === 'string') backends[id] = { url: normalizeUrl(backend.url) };
      }
      if (Object.keys(backends).length) backendConfig.set({ active: saved.active, backends });
    }
  } catch (error) {
    console.warn('Could not load backend profiles', error);
  }
  backendConfig.subscribe(value => {
    localStorage.setItem(storageKey, JSON.stringify(value));
  });
  window.ssuiDesktop?.onCertificateError(async certificate => {
    const accepted = window.confirm(`Trust the certificate for ${certificate.origin}?\n\nFingerprint: ${certificate.fingerprint}`);
    if (!accepted) return;
    await window.ssuiDesktop.trustCertificate(certificate);
    syncAuthState();
  });
}

export function getCurrentBackend() {
  const config = get(backendConfig);
  return config.backends[config.active] || config.backends.default;
}

export function getCurrentBackendUrl() {
  const url = getCurrentBackend()?.url || '/';
  return url === '/' ? '' : url;
}

export function setBackend(id, url) {
  const key = id.trim();
  if (!key) throw new Error('Backend name is required');
  backendConfig.update(config => ({
    ...config,
    backends: { ...config.backends, [key]: { url: normalizeUrl(url) } }
  }));
}

export async function setActiveBackend(id) {
  const config = get(backendConfig);
  if (!config.backends[id]) return false;
  backendConfig.set({ ...config, active: id });
  resetSessionState();
  return syncAuthState();
}

export async function apiFetch(endpoint, options = {}) {
  const url = buildUrl(endpoint);
  const request = { ...options, headers: new Headers(options.headers || {}) };
  request.credentials = 'include';
  const method = (request.method || 'GET').toUpperCase();
  const state = get(authState);
  if (!['GET', 'HEAD', 'OPTIONS'].includes(method) && state.csrf) {
    request.headers.set('X-SSUI-CSRF', state.csrf);
  }

  let response;
  if (window.ssuiDesktop?.request && getCurrentBackendUrl()) {
    response = await window.ssuiDesktop.request({ url, options: serializeOptions(request) });
    response = desktopResponse(response);
  } else {
    response = await fetch(url, request);
  }
  if (response.status === 401) {
    authState.update(value => ({ ...value, isAuthenticated: false, csrf: null, authError: 'Authentication required' }));
  }
  return response;
}

export async function apiFetchTimeout(endpoint, options = {}, timeoutMs = 5000) {
  const controller = new AbortController();
  const timeout = setTimeout(() => controller.abort(), timeoutMs);
  try {
    return await apiFetch(endpoint, { ...options, signal: controller.signal });
  } finally {
    clearTimeout(timeout);
  }
}

export async function apiJson(endpoint, options = {}) {
  const headers = new Headers(options.headers || {});
  if (options.body && !headers.has('Content-Type')) headers.set('Content-Type', 'application/json');
  const response = await apiFetch(endpoint, { ...options, headers });
  const data = await response.json().catch(() => null);
  if (!response.ok) throw new Error(data?.error?.message || `Request failed (${response.status})`);
  return data;
}

export async function apiText(endpoint, options = {}) {
  const response = await apiFetch(endpoint, options);
  if (!response.ok) throw new Error(`Request failed (${response.status})`);
  return response.text();
}

export async function login(username, password) {
  authState.update(value => ({ ...value, isAuthenticating: true, authError: null }));
  try {
    const backendUrl = getCurrentBackendUrl();
    const data = window.ssuiDesktop && backendUrl
      ? await window.ssuiDesktop.login({ backendUrl, username, password })
      : await apiJson('/auth/login', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ username, password })
        });
    applySession(data);
    return true;
  } catch (error) {
    authState.update(value => ({ ...value, isAuthenticated: false, isAuthenticating: false, authError: error.message }));
    return false;
  }
}

export async function bootstrapOwner(setupSecret, username, password) {
  const backendUrl = getCurrentBackendUrl();
  const data = window.ssuiDesktop && backendUrl
    ? await window.ssuiDesktop.bootstrap({ backendUrl, setupSecret, username, password })
    : await apiJson('/api/v2/auth/setup/bootstrap', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ setupSecret, username, password })
      });
  applySession(data);
  return data;
}

export async function logout() {
  try {
    const backendUrl = getCurrentBackendUrl();
    if (window.ssuiDesktop && backendUrl) {
      await window.ssuiDesktop.logout({ backendUrl });
    } else {
      await apiFetch('/auth/logout', { method: 'POST' });
    }
  } finally {
    resetSessionState();
  }
}

export async function syncAuthState() {
  authState.update(value => ({ ...value, isAuthenticating: true, authError: null }));
  try {
    const setup = await apiJson('/api/v2/auth/setup/status');
    if (setup.setupRequired) {
      authState.set({
        isAuthenticated: false,
        isAuthenticating: false,
        authError: null,
        setupRequired: true,
        user: null,
        permissions: [],
        features: { plugins: false },
        csrf: null,
        expiresAt: null
      });
      return false;
    }
    const response = await apiFetchTimeout('/api/v2/auth/session', { headers: { Accept: 'application/json' } }, 5000);
    if (!response.ok) throw new Error(response.status === 401 ? 'Authentication required' : `Backend returned ${response.status}`);
    applySession(await response.json());
    return true;
  } catch (error) {
    authState.update(value => ({
      ...value,
      isAuthenticated: false,
      isAuthenticating: false,
      authError: error.message,
      user: null,
      permissions: [],
      features: { plugins: false },
      csrf: null
    }));
    return false;
  }
}

export function hasPermission(permission) {
  return get(authState).permissions.includes(permission);
}

export function apiSSE(endpoint, onMessage, onError = console.error) {
  if (window.ssuiDesktop?.openEventStream && getCurrentBackendUrl()) {
    return window.ssuiDesktop.openEventStream({ url: buildUrl(endpoint), onMessage, onError });
  }
  const source = new EventSource(buildUrl(endpoint), { withCredentials: true });
  source.onmessage = event => {
    try {
      onMessage(JSON.parse(event.data));
    } catch {
      onMessage(event.data);
    }
  };
  source.onerror = onError;
  return { close: () => source.close() };
}

function applySession(data) {
  authState.set({
    isAuthenticated: true,
    isAuthenticating: false,
    authError: null,
    setupRequired: false,
    user: data.user,
    permissions: data.permissions || [],
    features: data.features || { plugins: false },
    csrf: data.csrf || null,
    expiresAt: data.expiresAt || null
  });
}

function resetSessionState() {
  authState.set({
    isAuthenticated: false,
    isAuthenticating: false,
    authError: null,
    setupRequired: false,
    user: null,
    permissions: [],
    features: { plugins: false },
    csrf: null,
    expiresAt: null
  });
}

function buildUrl(endpoint) {
  const path = endpoint.startsWith('/') ? endpoint : `/${endpoint}`;
  return `${getCurrentBackendUrl()}${path}`;
}

function normalizeUrl(url) {
  const value = url.trim();
  if (!value || value === '/') return '/';
  return value.replace(/\/+$/, '');
}

function serializeOptions(options) {
  return {
    method: options.method || 'GET',
    headers: Object.fromEntries(options.headers.entries()),
    body: options.body || null
  };
}

function desktopResponse(value) {
  return new Response(value.body ?? '', {
    status: value.status,
    statusText: value.statusText,
    headers: value.headers
  });
}
