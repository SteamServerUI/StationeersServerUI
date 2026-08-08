<script>
  import PluginContainer from './PluginContainer.svelte';
  import { pluginsList } from '../../services/plugins';
  let selected=$state('');
  let plugins=$derived($pluginsList);
  $effect(()=>{if(plugins.length&&!plugins.includes(selected))selected=plugins[0];});
  function name(path){return path.split('/').filter(Boolean).pop()?.replace(/([a-z])([A-Z])/g,'$1 $2')||path;}
</script>

<section class="plugin-workspace surface">
  <aside><span class="eyebrow">Loaded extensions</span>{#each plugins as path}<button class:active={selected===path} onclick={()=>selected=path}><span class="plugin-dot"></span><span>{name(path)}</span></button>{/each}{#if !plugins.length}<p>No plugins registered.</p>{/if}</aside>
  <div class="plugin-frame">{#if selected}<header><div><span class="eyebrow">Native plugin UI</span><h3>{name(selected)}</h3></div><span class="unsafe">Unrestricted process access</span></header><PluginContainer pluginPath={selected}/>{:else}<div class="empty-state"><span class="eyebrow">Quarantined</span><h2>No plugin interfaces</h2><p>Unsafe plugins are enabled, but none have registered a UI.</p></div>{/if}</div>
</section>

<style>
  .plugin-workspace{min-height:600px;display:grid;grid-template-columns:220px 1fr;overflow:hidden}aside{padding:18px 10px;border-right:1px solid var(--border-color)}aside>.eyebrow{display:block;padding:4px 10px 14px}aside button{width:100%;display:flex;align-items:center;gap:10px;background:transparent;border:0;text-align:left;color:var(--text-secondary)}aside button.active{background:var(--bg-active);color:var(--text-accent)}.plugin-dot{width:8px;height:8px;border-radius:50%;background:var(--text-warning);box-shadow:0 0 10px color-mix(in srgb,var(--text-warning) 50%,transparent)}aside p{padding:15px;color:var(--text-muted)}.plugin-frame{min-width:0}.plugin-frame>header{display:flex;align-items:center;justify-content:space-between;padding:18px 22px;border-bottom:1px solid var(--border-color)}.plugin-frame h3{margin:4px 0 0}.unsafe{color:var(--text-warning);font-size:.72rem}.plugin-frame :global(.iframe-container){height:520px}@media(max-width:900px){.plugin-workspace{grid-template-columns:170px 1fr}}
</style>
