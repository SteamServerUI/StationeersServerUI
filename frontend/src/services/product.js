import { writable } from 'svelte/store';
import { apiJson } from './api-v7';

export const productState = writable({
  loading: false,
  apiVersion: null,
  backend: null,
  features: {},
  capabilities: [],
  permissions: [],
  error: null
});

export async function loadCapabilities() {
  productState.update(value => ({ ...value, loading: true, error: null }));
  try {
    const response = await apiJson('/api/v3/capabilities');
    productState.set({ loading: false, error: null, ...response.data });
    return response.data;
  } catch (error) {
    productState.update(value => ({ ...value, loading: false, error: error.message }));
    throw error;
  }
}
