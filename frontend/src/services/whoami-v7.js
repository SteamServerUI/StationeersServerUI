import { get, writable } from 'svelte/store';
import { authState, syncAuthState } from './api-v7.js';

export const userInfo = writable(emptyUser(true));

export async function fetchUserInfo() {
  if (!get(authState).isAuthenticated) await syncAuthState();
  const state = get(authState);
  if (!state.user) {
    userInfo.set(emptyUser(false, state.authError));
    throw new Error(state.authError || 'Authentication required');
  }
  const accessLevel = state.user.groupIds?.includes('system-owner') ? 'Owner' : 'Custom access';
  const value = {
    username: state.user.username,
    accessLevel,
    isLoading: false,
    isAuthenticated: true,
    lastFetched: new Date(),
    error: null
  };
  userInfo.set(value);
  return value;
}

export function initUserInfo() {
  return fetchUserInfo().catch(() => {});
}

export function clearUserInfo() {
  userInfo.set(emptyUser(false));
}

export function getUserInitials(username) {
  if (!username) return 'USR';
  return username.split(/[\s_-]+/).slice(0, 2).map(value => value[0]).join('').slice(0, 3).toUpperCase();
}

export function formatAccessLevel(value) {
  return value || 'Unknown';
}

export function shouldRefreshUserInfo(lastFetched) {
  return !lastFetched || lastFetched < new Date(Date.now() - 5 * 60 * 1000);
}

function emptyUser(isLoading, error = null) {
  return {
    username: null,
    accessLevel: null,
    isLoading,
    isAuthenticated: false,
    lastFetched: null,
    error
  };
}
