<script lang="js">
  import { createEventDispatcher } from 'svelte';

  // Props using $props in runes mode
  const { message, isSuccess = false } = $props();

  // Event dispatcher for custom events
  const dispatch = createEventDispatcher();

  // Dismiss the popup
  function dismiss() {
    dispatch('dismiss');
  }
</script>

<div 
  class="popup" 
  class:success={isSuccess} 
  onclick={dismiss}
  onkeydown={(e) => e.key === 'Enter' && dismiss()}
  tabindex="0"
  role="button"
  aria-label="Dismiss notification"
>
  <div class="content">
    <span class="icon">{isSuccess ? '✓' : '!'}</span>
    <p class="message">{message}</p>
  </div>
</div>

<style>
  .popup {
    position: absolute;
    bottom: 1rem;
    left: 50%;
    transform: translateX(-50%) scale(1);
    background-color: var(--bg-tertiary, #4a46465a);
    padding: 0.5rem 1rem;
    border-radius: 6px;
    box-shadow: var(--shadow-medium, 0 4px 8px rgba(0,0,0,0.15));
    z-index: 10;
    cursor: pointer;
    animation: popIn var(--transition-speed, 0.3s) ease;
    max-width: 90%;
    pointer-events: all;
  }

  .success {
    background-color: var(--bg-success, #23b46b);
  }

  .content {
    display: flex;
    align-items: center;
  }

  .icon {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 24px;
    height: 24px;
    border-radius: 50%;
    margin-right: 0.5rem;
    font-weight: bold;
    font-size: 1rem;
  }

  .success .icon {
    color: var(--text-success, #2C6E49);
  }

  .popup:not(.success) .icon {
    color: var(--text-warning, #FF6B6B);
  }

  .message {
    margin: 0;
    font-size: 0.85rem;
    color: var(--text-primary, #333333);
  }

  @keyframes popIn {
    0% {
      transform: translateX(-50%) scale(0.8);
      opacity: 0;
    }
    100% {
      transform: translateX(-50%) scale(1);
      opacity: 1;
    }
  }
</style>