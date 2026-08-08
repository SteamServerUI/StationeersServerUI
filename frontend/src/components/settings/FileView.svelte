<script>
  import { onMount } from 'svelte';
  import { apiJson, hasPermission } from '../../services/api-v7';
  import { notify } from '../../services/activity';

  let files=$state([]), selected=$state(null), content=$state(''), original=$state(''), loading=$state(true), saving=$state(false), query=$state('');
  let filtered=$derived(files.filter(file=>!query||(`${file.filename} ${file.description} ${file.type}`).toLowerCase().includes(query.toLowerCase())));
  let dirty=$derived(content!==original);

  onMount(loadFiles);
  async function loadFiles(){
    loading=true;
    try{
      const body=await apiJson('/api/v3/files');
      files=body.data?.files||[];
      if(files.length&&!selected)await open(files[0]);
    }catch(error){notify('Could not load editable files','error',error.message);}
    finally{loading=false;}
  }
  async function open(file){
    if(dirty&&!window.confirm('Discard unsaved file changes?'))return;
    selected=file;loading=true;
    try{
      const body=await apiJson(`/api/v3/files/content?filename=${encodeURIComponent(file.filename)}`);
      content=body.data.content;
      original=content;
    }catch(error){notify('Could not open file','error',error.message);}
    finally{loading=false;}
  }
  async function save(){
    if(!selected||!dirty)return;saving=true;
    try{
      await apiJson('/api/v3/files/content',{method:'PUT',body:JSON.stringify({filename:selected.filename,content})});
      original=content;
      notify('File saved','success',selected.filename);
    }catch(error){notify('Could not save file','error',error.message);}
    finally{saving=false;}
  }
  function reset(){content=original;}
</script>

<section class="file-workspace surface">
  <aside><header><span class="eyebrow">Runfile files</span><input bind:value={query} placeholder="Filter files"></header><div class="file-list">{#each filtered as file (file.filename)}<button class:active={selected?.filename===file.filename} onclick={()=>open(file)}><span class="file-type">{(file.type||'txt').slice(0,4).toUpperCase()}</span><span><strong>{file.filename}</strong><small>{file.description||'Editable game file'}</small></span></button>{/each}{#if !filtered.length&&!loading}<div class="no-files">No matching files</div>{/if}</div></aside>
  <div class="editor">
    {#if selected}<header><div><span class="eyebrow">{selected.type||'Text'} editor</span><h3>{selected.filename}</h3></div><div class="editor-state"><span class:dirty>{dirty? 'Unsaved changes':'Saved'}</span><button onclick={reset} disabled={!dirty}>Revert</button>{#if hasPermission('files.write')}<button class="save" onclick={save} disabled={!dirty||saving}>{saving?'Saving&':'Save file'}</button>{/if}</div></header><textarea bind:value={content} spellcheck="false" disabled={loading} aria-label="File content"></textarea><footer><span>{content.split('\n').length} lines</span><span>{content.length.toLocaleString()} characters</span></footer>{:else}<div class="empty-state"><span class="eyebrow">Files</span><h2>No editable file selected</h2><p>The active runfile decides which game files are safe to edit.</p></div>{/if}
  </div>
</section>

<style>
  .file-workspace{height:min(780px,calc(100vh - 190px));min-height:540px;display:grid;grid-template-columns:290px 1fr;overflow:hidden}aside{border-right:1px solid var(--border-color);min-width:0}aside header{display:grid;gap:12px;padding:18px}.file-list{display:grid;overflow:auto}.file-list button{display:grid;grid-template-columns:42px 1fr;gap:11px;text-align:left;background:transparent;border:0;border-radius:0;padding:13px 17px}.file-list button.active{background:var(--bg-active)}.file-list strong,.file-list small{display:block;white-space:nowrap;overflow:hidden;text-overflow:ellipsis}.file-list small{color:var(--text-muted);margin-top:3px}.file-type{display:grid;place-items:center;width:38px;height:38px;border-radius:10px;background:var(--material-solid);color:var(--text-accent);font-size:.63rem;font-weight:850}.no-files{padding:30px;color:var(--text-muted);text-align:center}.editor{min-width:0;display:grid;grid-template-rows:auto 1fr auto}.editor>header{display:flex;align-items:center;justify-content:space-between;gap:20px;padding:16px 18px;border-bottom:1px solid var(--border-color)}.editor h3{margin:4px 0 0}.editor-state{display:flex;align-items:center;gap:8px}.editor-state>span{color:var(--text-success);font-size:.73rem}.editor-state>span.dirty{color:var(--text-warning)}.editor-state .save{background:linear-gradient(135deg,var(--accent-primary),var(--accent-secondary));color:#07101d;border:0;font-weight:800}.editor textarea{resize:none;border:0;border-radius:0;background:#080b11;color:#d2deef;padding:18px;font:13px/1.65 ui-monospace,SFMono-Regular,Consolas,monospace;tab-size:2}.editor>footer{display:flex;gap:18px;padding:8px 16px;border-top:1px solid var(--border-color);color:var(--text-muted);font-size:.68rem}@media(max-width:900px){.file-workspace{grid-template-columns:220px 1fr}.editor>header{align-items:flex-start;flex-direction:column}.editor-state{width:100%;flex-wrap:wrap}}
</style>
