<script>
  import { onMount } from 'svelte';
  import { apiFetch, hasPermission } from '../../services/api-v7';
  import { notify } from '../../services/activity';

  let groups=$state([]), active=$state(''), args=$state([]), query=$state(''), loading=$state(true), saving=$state(new Set());
  let visible=$derived(args.filter(arg=>!query||(`${arg.ui_label} ${arg.flag} ${arg.description}`).toLowerCase().includes(query.toLowerCase())));

  onMount(loadGroups);
  async function request(path,options){const response=await apiFetch(path,options);const body=await response.json().catch(()=>({}));if(!response.ok||body.error)throw new Error(body?.error?.message||body.error||`Request failed (${response.status})`);return body.data??body;}
  async function loadGroups(){try{groups=await request('/api/v3/runfile/groups');if(groups.length)await choose(groups[0]);}catch(error){notify('Could not load game configuration','error',error.message);}finally{loading=false;}}
  async function choose(group){active=group;loading=true;try{args=await request(`/api/v3/runfile/args?group=${encodeURIComponent(group)}`);}catch(error){notify('Could not load setting group','error',error.message);}finally{loading=false;}}
  async function save(arg,value){
    const previous=arg.runtime_value??arg.value;arg.runtime_value=String(value);args=[...args];saving=new Set([...saving,arg.flag]);
    try{await request('/api/v3/runfile/args/update',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({flag:arg.flag,value:String(value)})});await request('/api/v3/runfile/save',{method:'POST'});notify(`${arg.ui_label||arg.flag} saved`,'success','Used on the next server start.');}
    catch(error){arg.runtime_value=previous;args=[...args];notify('Could not save game setting','error',error.message);}finally{const next=new Set(saving);next.delete(arg.flag);saving=next;}
  }
  function value(arg){return arg.runtime_value??arg.value??'';}
</script>

<section class="game-settings surface">
  <aside><header><span class="eyebrow">Configuration groups</span></header>{#each groups as group}<button class:active={active===group} onclick={()=>choose(group)}><span>{group}</span><small>{active===group?args.length:''}</small></button>{/each}</aside>
  <div class="game-panel"><header><div><span class="eyebrow">Game configuration</span><h2>{active||'No runfile loaded'}</h2><p>Changes save immediately and apply on the next game-server start.</p></div><input bind:value={query} placeholder="Find a setting"></header>
  {#if loading}<div class="empty-state"><h2>Loading settings…</h2></div>{:else}<div class="argument-list">{#each visible as arg (arg.flag)}<article><div><strong>{arg.ui_label||arg.flag}</strong><code>{arg.flag}</code><p>{arg.description||'No description supplied by this runfile.'}</p>{#if arg.required}<small>Required</small>{/if}</div><div class="control">{#if arg.type==='bool'}<label class="toggle"><input type="checkbox" checked={value(arg)==='true'} disabled={!hasPermission('runfiles.manage')||saving.has(arg.flag)||arg.disabled} onchange={event=>save(arg,event.currentTarget.checked)}><span></span></label>{:else}<input type={arg.type==='int'||arg.type==='number'?'number':'text'} value={value(arg)} min={arg.min} max={arg.max} disabled={!hasPermission('runfiles.manage')||saving.has(arg.flag)||arg.disabled} onchange={event=>save(arg,event.currentTarget.value)}>{/if}<small>{saving.has(arg.flag)?'Saving…':arg.disabled?'Unavailable':''}</small></div></article>{/each}{#if !visible.length}<div class="empty-state"><h2>No matching settings</h2></div>{/if}</div>{/if}</div>
</section>

<style>
  .game-settings{display:grid;grid-template-columns:220px 1fr;min-height:600px;overflow:hidden}aside{padding:17px 10px;border-right:1px solid var(--border-color)}aside header{padding:4px 10px 13px}aside button{width:100%;display:flex;justify-content:space-between;background:transparent;border:0;text-align:left;color:var(--text-secondary)}aside button.active{background:var(--bg-active);color:var(--text-accent)}aside small{color:var(--text-muted)}.game-panel{min-width:0}.game-panel>header{display:flex;align-items:end;justify-content:space-between;gap:20px;padding:22px;border-bottom:1px solid var(--border-color)}.game-panel h2{margin:4px 0}.game-panel p{color:var(--text-secondary);margin:0}.game-panel>header input{width:240px}.argument-list{display:grid}.argument-list article{display:grid;grid-template-columns:1fr minmax(210px,34%);align-items:center;gap:25px;padding:19px 22px;border-bottom:1px solid var(--border-color)}.argument-list strong{font-size:.94rem}.argument-list code{color:var(--text-muted);font-size:.68rem;margin-left:8px}.argument-list p{font-size:.82rem;margin:5px 0 0;max-width:700px}.argument-list article>div>small{color:var(--text-warning);font-size:.65rem}.control{display:grid;grid-template-columns:1fr auto;align-items:center;gap:8px}.control>small{color:var(--text-success);font-size:.68rem}.toggle input{position:absolute;opacity:0}.toggle span{display:block;width:46px;height:26px;padding:3px;border-radius:20px;background:var(--bg-hover);border:1px solid var(--border-color)}.toggle span::after{content:'';display:block;width:18px;height:18px;border-radius:50%;background:var(--text-secondary);transition:.18s}.toggle input:checked+span{background:color-mix(in srgb,var(--accent-primary) 42%,transparent);border-color:var(--accent-primary)}.toggle input:checked+span::after{transform:translateX(19px);background:var(--accent-primary)}@media(max-width:900px){.game-settings{grid-template-columns:170px 1fr}.game-panel>header{align-items:start;flex-direction:column}.game-panel>header input{width:100%}.argument-list article{grid-template-columns:1fr}}
</style>
