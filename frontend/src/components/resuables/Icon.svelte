<script>
  let { name, size = 20 } = $props();
</script>

<svg class="icon" width={size} height={size} viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
  {#if name === 'home'}<path d="M3 11.5 12 4l9 7.5"/><path d="M5.5 10v10h13V10"/><path d="M9.5 20v-6h5v6"/>
  {:else if name === 'terminal'}<path d="m4 6 5 5-5 5"/><path d="M11 18h9"/>
  {:else if name === 'logs'}<path d="M6 3h9l4 4v14H6z"/><path d="M14 3v5h5M9 13h7M9 17h7"/>
  {:else if name === 'activity'}<path d="M3 12h4l2-6 4 12 2-6h6"/>
  {:else if name === 'game'}<path d="M8 9h8a5 5 0 0 1 4.6 7l-1 2.5a2 2 0 0 1-3.3.7L14 17h-4l-2.3 2.2a2 2 0 0 1-3.3-.7L3.4 16A5 5 0 0 1 8 9Z"/><path d="M7 13h4M9 11v4M16.5 12.5h.01M18.5 14.5h.01"/>
  {:else if name === 'gallery'}<rect x="3" y="4" width="18" height="16" rx="3"/><circle cx="9" cy="10" r="2"/><path d="m4 18 5-5 3 3 2-2 6 5"/>
  {:else if name === 'files'}<path d="M3 6h7l2 2h9v11H3z"/>
  {:else if name === 'backup'}<path d="M4 7v5h5"/><path d="M5.5 16a8 8 0 1 0 .5-9l-2 5"/>
  {:else if name === 'access'}<circle cx="9" cy="8" r="3"/><path d="M3 20a6 6 0 0 1 12 0M16 11l2 2 3-4"/>
  {:else if name === 'settings'}<circle cx="12" cy="12" r="3"/><path d="M19.4 15a1.7 1.7 0 0 0 .3 1.9l.1.1-2.8 2.8-.1-.1a1.7 1.7 0 0 0-1.9-.3 1.7 1.7 0 0 0-1 1.6v.2h-4V21a1.7 1.7 0 0 0-1-1.6 1.7 1.7 0 0 0-1.9.3l-.1.1L4.2 17l.1-.1a1.7 1.7 0 0 0 .3-1.9A1.7 1.7 0 0 0 3 14H2.8v-4H3a1.7 1.7 0 0 0 1.6-1 1.7 1.7 0 0 0-.3-1.9L4.2 7 7 4.2l.1.1A1.7 1.7 0 0 0 9 4.6a1.7 1.7 0 0 0 1-1.6v-.2h4V3a1.7 1.7 0 0 0 1 1.6 1.7 1.7 0 0 0 1.9-.3l.1-.1L19.8 7l-.1.1a1.7 1.7 0 0 0-.3 1.9 1.7 1.7 0 0 0 1.6 1h.2v4H21a1.7 1.7 0 0 0-1.6 1Z"/>
  {:else if name === 'backends'}<rect x="4" y="4" width="16" height="6" rx="2"/><rect x="4" y="14" width="16" height="6" rx="2"/><path d="M8 7h.01M8 17h.01"/>
  {:else if name === 'plugin'}<path d="M8 3v4M16 3v4M6 7h12v4a6 6 0 0 1-12 0Z"/><path d="M12 17v4"/>
  {:else if name === 'chevron'}<path d="m9 18 6-6-6-6"/>
  {:else if name === 'bell'}<path d="M18 8a6 6 0 0 0-12 0c0 7-3 7-3 9h18c0-2-3-2-3-9M10 21h4"/>
  {:else if name === 'logout'}<path d="M10 5H5v14h5M14 8l4 4-4 4M18 12H9"/>
  {:else}<circle cx="12" cy="12" r="8"/>{/if}
</svg>

<style>.icon{display:block;flex:0 0 auto}</style>
