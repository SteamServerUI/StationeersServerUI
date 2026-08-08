<script>
  import { onMount } from 'svelte';
  import { backendConfig, setBackend, setActiveBackend } from '../../services/api-v7';
  import { loadCapabilities, productState } from '../../services/product';
  import { notify } from '../../services/activity';

  let config=$state({active:'default',backends:{default:{url:'/'}}}), name=$state(''), url=$state(''), checking=$state(false), online=$state(null);
  const unsubscribe=backendConfig.subscribe(value=>config=value);
  onMount(()=>{check();return unsubscribe;});

  async function check(){
    checking=true;
    try{await loadCapabilities();online=true;}
    catch(error){online=false;notify('Backend check failed','error',error.message);}
    finally{checking=false;}
  }
  async function activate(id){
    if(id===config.active)return;
    await setActiveBackend(id);
    notify('Workspace changed','info',id);
    await check();
  }
  function add(){
    const clean=name.trim();
    if(!clean||!url.trim()){notify('Connection needs a name and URL','error');return;}
    try{setBackend(clean,url);name='';url='';notify('Backend saved','success',clean);}
    catch(error){notify('Could not save backend','error',error.message);}
  }
  function remove(id){
    if(id==='default'||id===config.active)return;
    backendConfig.update(value=>{const backends={...value.backends};delete backends[id];return {...value,backends};});
    notify('Backend removed','info',id);
  }
</script>

<section class="backend-page">
  <article class="connection-hero surface">
    <div><span class="eyebrow">Active workspace</span><h2>{config.active}</h2><p>{config.backends[config.active]?.url||'This SSUI backend'}{#if $productState.backend?.version} / SSUI {$productState.backend.version} / API {$productState.apiVersion}{/if}</p></div>
    <div class="connection-state"><span class:online></span><strong>{checking?'Checking…':online?'Connected':'Unavailable'}</strong><button onclick={check} disabled={checking}>Check now</button></div>
  </article>

  <div class="backend-grid">
    {#each Object.entries(config.backends) as [id,backend]}
      <article class="backend-card surface" class:active={id===config.active}>
        <header><div><span class="eyebrow">{id===config.active?'Current workspace':'Saved connection'}</span><h3>{id}</h3></div>{#if id===config.active}<span class="active-badge">Active</span>{/if}</header>
        <code>{backend.url||'/'}</code>
        <footer><button onclick={()=>activate(id)} disabled={id===config.active}>Use workspace</button>{#if id!=='default'&&id!==config.active}<button class="remove" onclick={()=>remove(id)}>Remove</button>{/if}</footer>
      </article>
    {/each}
    <article class="add-card surface">
      <span class="eyebrow">New connection</span><h3>Add a backend</h3>
      <label><span>Name</span><input bind:value={name} placeholder="Production"></label>
      <label><span>HTTPS URL</span><input bind:value={url} placeholder="https://server.example:8443"></label>
      <button class="add" onclick={add}>Save backend</button>
    </article>
  </div>
</section>

<style>
  .backend-page{display:grid;gap:16px}.connection-hero{display:flex;align-items:center;justify-content:space-between;gap:25px;padding:26px}.connection-hero h2{font-size:2rem;margin:4px 0}.connection-hero p{color:var(--text-secondary);margin:0}.connection-state{display:grid;grid-template-columns:10px 1fr auto;align-items:center;gap:10px}.connection-state>span{width:9px;height:9px;border-radius:50%;background:var(--text-danger);box-shadow:0 0 12px color-mix(in srgb,var(--text-danger) 55%,transparent)}.connection-state>span.online{background:var(--text-success);box-shadow:0 0 12px color-mix(in srgb,var(--text-success) 55%,transparent)}.backend-grid{display:grid;grid-template-columns:repeat(2,minmax(0,1fr));gap:16px}.backend-card,.add-card{padding:22px}.backend-card.active{border-color:color-mix(in srgb,var(--accent-primary) 45%,transparent)}.backend-card header{display:flex;justify-content:space-between}.backend-card h3,.add-card h3{font-size:1.3rem;margin:5px 0 18px}.active-badge{height:max-content;padding:5px 9px;border-radius:999px;background:var(--bg-active);color:var(--text-accent);font-size:.7rem;font-weight:800}.backend-card code{display:block;padding:12px;border-radius:10px;background:var(--material-solid);color:var(--text-secondary);word-break:break-all}.backend-card footer{display:flex;justify-content:space-between;margin-top:20px}.remove{color:var(--text-danger);background:transparent}.add-card{display:grid;gap:13px}.add-card label{display:grid;gap:6px;color:var(--text-secondary);font-size:.75rem}.add-card .add{background:linear-gradient(135deg,var(--accent-primary),var(--accent-secondary));color:#07101d;border:0;font-weight:800;margin-top:5px}@media(max-width:900px){.backend-grid{grid-template-columns:1fr}.connection-hero{align-items:flex-start;flex-direction:column}}
</style>
