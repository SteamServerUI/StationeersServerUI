<script>
  import { toasts, dismissToast } from '../../services/activity';
  import { navigate } from '../../services/router';

  function open(item) {
    dismissToast(item.id);
    navigate('activity');
  }
</script>

<div class="toast-stack" aria-live="polite">
  {#each $toasts as item (item.id)}
    <button class="toast" class:error={item.tone === 'error'} class:success={item.tone === 'success'} onclick={() => open(item)}>
      <span class="toast-dot"></span>
      <span><strong>{item.message}</strong>{#if item.detail}<small>{item.detail}</small>{/if}</span>
    </button>
  {/each}
</div>
