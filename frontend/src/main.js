import './app.css';
import App from './AppV7.svelte';
import { mount } from 'svelte';

mount(App, {
  target: document.getElementById('app'),
});