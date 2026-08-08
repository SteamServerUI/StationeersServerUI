<script>
  import { activities } from '../../services/activity';
</script>

<section class="surface activity-view">
  {#if !$activities.length}
    <div class="empty-state"><span class="eyebrow">All quiet</span><h2>No activity yet</h2><p>Operations, warnings and failures from this session will collect here.</p></div>
  {:else}
    <div class="activity-list">
      {#each $activities as item (item.id)}
        <article class="activity-row" class:error={item.tone === 'error'}>
          <span class="activity-marker"></span>
          <div><strong>{item.message}</strong>{#if item.detail}<p>{item.detail}</p>{/if}</div>
          <time>{new Date(item.at).toLocaleString()}</time>
        </article>
      {/each}
    </div>
  {/if}
</section>
