<script>
  import { onMount } from 'svelte';
  import { apiFetch } from '../../services/api-v7';
  import { notify } from '../../services/activity';

  let settings=$state([]), category=$state(''), query=$state(''), loading=$state(true);
  let saving=$state(new Set()), saved=$state(new Set());

  const descriptions={
    'System Settings':'Backend identity, network and runtime behaviour.',
    'Logging Settings':'What SSUI records and how much detail it keeps.',
    'Gameserver Settings':'Startup, updates and game process behaviour.',
    'Security Settings':'Authentication and session-related controls.',
    'Modding Settings':'Optional hooks and mod-loader integration.',
    'Update Settings':'Choose how SSUI and the game receive updates.',
    'Discord Settings':'Bot connectivity and channel destinations.',
    'Backup Settings':'Storage, scheduling and backup behaviour.'
  };
  let categories=$derived([...new Set(settings.map(item=>item.group))]);
  let visible=$derived(settings.filter(item=>(!category||item.group===category)&&(!query||(`${item.name} ${item.description} ${item.group}`).toLowerCase().includes(query.toLowerCase()))));

  onMount(load);
  async function load(){
    try{
      const response=await apiFetch('/api/v3/settings');
      const body=await response.json();
      if(!response.ok||body.error)throw new Error(body?.error?.message||body.error||'Settings unavailable');
      settings=body.data||[];
    }catch(error){notify('Could not load settings','error',error.message);}
    finally{loading=false;}
  }
  async function save(item,value){
    const previous=item.value;
    item.value=value;
    settings=[...settings];
    saving=new Set([...saving,item.name]);
    try{
      const response=await apiFetch('/api/v3/settings/save',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({[item.name]:value})});
      const body=await response.json().catch(()=>({}));
      if(!response.ok||body.status==='error')throw new Error(body?.error?.message||body.message||'The backend rejected this value');
      saved=new Set([...saved,item.name]);
      setTimeout(()=>{const next=new Set(saved);next.delete(item.name);saved=next;},2200);
      notify(`${label(item.name)} saved`,'success',effect(item));
    }catch(error){
      item.value=previous;settings=[...settings];
      notify(`Could not save ${label(item.name)}`,'error',error.message);
    }finally{const next=new Set(saving);next.delete(item.name);saving=next;}
  }
  function inputValue(item,event){
    if(item.type==='bool')return event.currentTarget.checked;
    if(item.type==='int')return event.currentTarget.value===''?null:Number(event.currentTarget.value);
    return event.currentTarget.value;
  }
  function label(name){return name.replace(/([a-z0-9])([A-Z])/g,'$1 $2').replace(/^Is /,'').replace(/^Allow /,'Allow ');}
  function effect(item){return /Port|CLI|BackupsStore|BackupLoop|GameLog/.test(item.name)?'A backend reload may be required.':'Applied immediately.';}
</script>

<section class="settings-hub">
  <div class="settings-toolbar surface">
    {#if category}<button class="back" onclick={()=>category=''}>← All settings</button>{/if}
    <label class="search"><span>Search settings</span><input bind:value={query} placeholder="Try “backup”, “updates” or “Discord”…"></label>
  </div>

  {#if loading}
    <div class="empty-state surface"><span class="eyebrow">Settings</span><h2>Loading configuration…</h2></div>
  {:else if !category && !query}
    <div class="category-grid">
      {#each categories as group}
        <button class="category-card surface" onclick={()=>category=group}>
          <span class="category-count">{settings.filter(item=>item.group===group).length}</span>
          <span><strong>{group}</strong><small>{descriptions[group]||'Configure this part of SteamServerUI.'}</small></span>
          <span class="arrow">→</span>
        </button>
      {/each}
    </div>
  {:else}
    <div class="settings-list surface">
      <header><div><span class="eyebrow">{category||'Search results'}</span><h2>{category||`${visible.length} matching settings`}</h2></div></header>
      {#each visible as item (item.name)}
        <article class="setting-row">
          <div class="setting-copy"><strong>{label(item.name)}</strong><p>{item.description}</p>{#if /Port|CLI|BackupsStore|BackupLoop|GameLog/.test(item.name)}<small>Backend reload may be required</small>{/if}</div>
          <div class="setting-control" class:saving={saving.has(item.name)}>
            {#if item.type==='bool'}
              <label class="toggle"><input type="checkbox" checked={item.value===true} disabled={saving.has(item.name)} onchange={event=>save(item,inputValue(item,event))}><span></span></label>
            {:else}
              <input type={item.sensitive?'password':item.type==='int'?'number':'text'} value={item.sensitive?'':item.value??''} min={item.min} max={item.max} required={item.required} placeholder={item.sensitive&&item.hasValue?'Configured — enter to replace':''} disabled={saving.has(item.name)} onchange={event=>{const value=inputValue(item,event);if(!item.sensitive||value)save(item,value)}}>
            {/if}
            <span class="save-state">{saving.has(item.name)?'Saving…':saved.has(item.name)?'Saved':''}</span>
          </div>
        </article>
      {:else}
        <div class="empty-state"><span class="eyebrow">No match</span><h2>Nothing found</h2><p>Try a broader search.</p></div>
      {/each}
    </div>
  {/if}
</section>

<style>
  .settings-hub{display:grid;gap:16px}.settings-toolbar{display:flex;align-items:end;gap:14px;padding:16px}.back{background:transparent;white-space:nowrap}.search{display:grid;gap:6px;flex:1;color:var(--text-secondary);font-size:.75rem}.search input{width:100%}.category-grid{display:grid;grid-template-columns:repeat(2,minmax(0,1fr));gap:15px}.category-card{display:grid;grid-template-columns:42px 1fr auto;align-items:center;gap:14px;padding:22px;text-align:left;min-height:126px}.category-card:hover{transform:translateY(-2px);box-shadow:var(--shadow-medium)}.category-count{display:grid;place-items:center;width:42px;height:42px;border-radius:14px;background:var(--bg-active);color:var(--text-accent);font-weight:800}.category-card strong,.category-card small{display:block}.category-card strong{font-size:1.08rem}.category-card small{color:var(--text-secondary);font-weight:400;margin-top:5px;line-height:1.4}.arrow{font-size:1.3rem;color:var(--text-muted)}.settings-list{padding:12px 20px}.settings-list>header{padding:16px 6px}.settings-list h2{margin:4px 0}.setting-row{display:grid;grid-template-columns:minmax(0,1fr) minmax(220px,36%);align-items:center;gap:30px;padding:19px 6px;border-top:1px solid var(--border-color)}.setting-copy p{color:var(--text-secondary);font-size:.86rem;margin:5px 0 0;max-width:720px}.setting-copy small{display:block;color:var(--text-warning);font-size:.7rem;margin-top:6px}.setting-control{display:grid;grid-template-columns:1fr auto;align-items:center;gap:10px}.save-state{min-width:48px;color:var(--text-success);font-size:.72rem}.toggle input{position:absolute;opacity:0}.toggle span{display:block;width:46px;height:26px;padding:3px;border-radius:20px;background:var(--bg-hover);border:1px solid var(--border-color);transition:.18s}.toggle span::after{content:'';display:block;width:18px;height:18px;border-radius:50%;background:var(--text-secondary);transition:.18s}.toggle input:checked+span{background:color-mix(in srgb,var(--accent-primary) 42%,transparent);border-color:var(--accent-primary)}.toggle input:checked+span::after{transform:translateX(19px);background:var(--accent-primary);box-shadow:0 0 12px color-mix(in srgb,var(--accent-primary) 55%,transparent)}@media(max-width:950px){.category-grid{grid-template-columns:1fr}.setting-row{grid-template-columns:1fr;gap:13px}}
</style>
