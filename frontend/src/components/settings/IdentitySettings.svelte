<script>
  import { onMount } from 'svelte';
  import { apiJson, authState, hasPermission } from '../../services/api-v7';

  let users = $state([]);
  let groups = $state([]);
  let permissions = $state([]);
  let tokens = $state([]);
  let sessions = $state([]);
  let events = $state([]);
  let active = $state('users');
  let busy = $state(false);
  let message = $state('');
  let error = $state('');
  let shownToken = $state('');

  let newUser = $state({ username: '', password: '', groupIds: [] });
  let groupDraft = $state({ id: '', name: '', description: '', permissions: [] });
  let tokenDraft = $state({ name: '', scopes: [] });

  let canManageUsers = $derived(hasPermission('users.manage'));
  let canManageGroups = $derived(hasPermission('groups.manage'));
  let canManageTokens = $derived(hasPermission('tokens.own.manage'));
  let canViewAudit = $derived(hasPermission('audit.view'));

  onMount(() => {
    if (!canManageUsers) {
      active = canManageGroups ? 'groups' : canManageTokens ? 'tokens' : 'sessions';
    }
    loadAll();
  });

  async function loadAll() {
    busy = true;
    error = '';
    try {
      const requests = [];
      if (canManageUsers) requests.push(apiJson('/api/v2/auth/users').then(data => users = data.users));
      if (canManageGroups) requests.push(apiJson('/api/v2/auth/groups').then(data => {
        groups = data.groups;
        permissions = data.permissions;
      }));
      if (canManageTokens) requests.push(apiJson('/api/v2/auth/tokens').then(data => tokens = data.tokens));
      if (canViewAudit) requests.push(apiJson('/api/v2/auth/audit').then(data => events = data.events));
      requests.push(apiJson('/api/v2/auth/sessions').then(data => sessions = data.sessions));
      await Promise.all(requests);
    } catch (reason) {
      error = reason.message;
    } finally {
      busy = false;
    }
  }

  async function createUser() {
    await run(async () => {
      await apiJson('/api/v2/auth/users', {
        method: 'POST',
        body: JSON.stringify(newUser)
      });
      newUser = { username: '', password: '', groupIds: [] };
      message = 'User created.';
      await loadAll();
    });
  }

  async function toggleUser(user) {
    await run(async () => {
      await apiJson(`/api/v2/auth/users/${user.id}`, {
        method: 'PATCH',
        body: JSON.stringify({ enabled: !user.enabled })
      });
      message = user.enabled ? 'User disabled.' : 'User enabled.';
      await loadAll();
    });
  }

  async function saveGroup() {
    await run(async () => {
      const body = JSON.stringify({
        name: groupDraft.name,
        description: groupDraft.description,
        permissions: groupDraft.permissions
      });
      if (groupDraft.id) {
        await apiJson(`/api/v2/auth/groups/${groupDraft.id}`, { method: 'PUT', body });
        message = 'Group updated.';
      } else {
        await apiJson('/api/v2/auth/groups', { method: 'POST', body });
        message = 'Group created.';
      }
      clearGroupDraft();
      await loadAll();
    });
  }

  function editGroup(group) {
    groupDraft = {
      id: group.id,
      name: group.name,
      description: group.description || '',
      permissions: [...group.permissions]
    };
  }

  async function deleteGroup(group) {
    if (!window.confirm(`Delete ${group.name}? Users will lose this group's access.`)) return;
    await run(async () => {
      await apiJson(`/api/v2/auth/groups/${group.id}`, { method: 'DELETE' });
      message = 'Group deleted.';
      clearGroupDraft();
      await loadAll();
    });
  }

  async function createToken() {
    await run(async () => {
      const data = await apiJson('/api/v2/auth/tokens', {
        method: 'POST',
        body: JSON.stringify(tokenDraft)
      });
      shownToken = data.secret;
      tokenDraft = { name: '', scopes: [] };
      message = 'Token created. Copy it now.';
      await loadAll();
    });
  }

  async function revokeToken(token) {
    await run(async () => {
      await apiJson(`/api/v2/auth/tokens/${token.id}`, { method: 'DELETE' });
      message = 'Token revoked.';
      await loadAll();
    });
  }

  async function revokeSession(session) {
    await run(async () => {
      await apiJson(`/api/v2/auth/sessions/${session.id}`, { method: 'DELETE' });
      message = 'Session revoked.';
      await loadAll();
    });
  }

  async function run(action) {
    busy = true;
    error = '';
    message = '';
    try {
      await action();
    } catch (reason) {
      error = reason.message;
    } finally {
      busy = false;
    }
  }

  function toggle(values, value) {
    return values.includes(value) ? values.filter(item => item !== value) : [...values, value];
  }

  function clearGroupDraft() {
    groupDraft = { id: '', name: '', description: '', permissions: [] };
  }

  function groupNames(ids) {
    return ids.map(id => groups.find(group => group.id === id)?.name || id).join(', ') || 'No access groups';
  }
</script>

<section class="identity">
  <header>
    <div><h2>People & access</h2><p>Access is additive across custom groups. The backend checks every permission.</p></div>
    <button onclick={loadAll} disabled={busy}>Refresh</button>
  </header>

  <nav>
    {#if canManageUsers}<button class:active={active === 'users'} onclick={() => active = 'users'}>Users</button>{/if}
    {#if canManageGroups}<button class:active={active === 'groups'} onclick={() => active = 'groups'}>Groups</button>{/if}
    {#if canManageTokens}<button class:active={active === 'tokens'} onclick={() => active = 'tokens'}>API tokens</button>{/if}
    <button class:active={active === 'sessions'} onclick={() => active = 'sessions'}>Sessions</button>
    {#if canViewAudit}<button class:active={active === 'audit'} onclick={() => active = 'audit'}>Audit</button>{/if}
  </nav>

  {#if error}<p class="notice error">{error}</p>{/if}
  {#if message}<p class="notice success">{message}</p>{/if}

  {#if active === 'users' && canManageUsers}
    <div class="split">
      <div class="list">
        {#each users as user}
          <article>
            <div><strong>{user.username}</strong><small>{groupNames(user.groupIds)}</small></div>
            <button onclick={() => toggleUser(user)} disabled={busy || user.id === $authState.user?.id}>{user.enabled ? 'Disable' : 'Enable'}</button>
          </article>
        {/each}
      </div>
      <form onsubmit={(event) => { event.preventDefault(); createUser(); }}>
        <h3>Add user</h3>
        <label>Username<input bind:value={newUser.username} autocomplete="off" /></label>
        <label>Temporary password<input bind:value={newUser.password} type="password" autocomplete="new-password" /></label>
        <fieldset><legend>Groups</legend>
          {#each groups as group}
            <label class="check"><input type="checkbox" checked={newUser.groupIds.includes(group.id)} onchange={() => newUser.groupIds = toggle(newUser.groupIds, group.id)} />{group.name}</label>
          {/each}
        </fieldset>
        <button class="primary" disabled={busy}>Create user</button>
      </form>
    </div>
  {:else if active === 'groups' && canManageGroups}
    <div class="split">
      <div class="list">
        {#each groups as group}
          <article>
            <div><strong>{group.name}</strong><small>{group.permissions.length} permissions{group.system ? ' · system' : ''}</small></div>
            {#if !group.system}<span><button onclick={() => editGroup(group)}>Edit</button><button class="danger" onclick={() => deleteGroup(group)}>Delete</button></span>{/if}
          </article>
        {/each}
      </div>
      <form onsubmit={(event) => { event.preventDefault(); saveGroup(); }}>
        <h3>{groupDraft.id ? 'Edit group' : 'New group'}</h3>
        <label>Name<input bind:value={groupDraft.name} /></label>
        <label>Description<input bind:value={groupDraft.description} /></label>
        <fieldset class="permissions"><legend>Permissions</legend>
          {#each permissions as permission}
            <label class="check"><input type="checkbox" checked={groupDraft.permissions.includes(permission)} onchange={() => groupDraft.permissions = toggle(groupDraft.permissions, permission)} />{permission}</label>
          {/each}
        </fieldset>
        <div class="actions"><button class="primary" disabled={busy}>Save group</button>{#if groupDraft.id}<button type="button" onclick={clearGroupDraft}>Cancel</button>{/if}</div>
      </form>
    </div>
  {:else if active === 'tokens' && canManageTokens}
    {#if shownToken}<div class="token-secret"><strong>Copy this token now</strong><code>{shownToken}</code><button onclick={() => navigator.clipboard.writeText(shownToken)}>Copy</button></div>{/if}
    <div class="split">
      <div class="list">
        {#each tokens as token}
          <article><div><strong>{token.name}</strong><small>{token.scopes.length} scopes · {token.revokedAt ? 'revoked' : 'active'}</small></div>{#if !token.revokedAt}<button class="danger" onclick={() => revokeToken(token)}>Revoke</button>{/if}</article>
        {/each}
      </div>
      <form onsubmit={(event) => { event.preventDefault(); createToken(); }}>
        <h3>New API token</h3>
        <label>Name<input bind:value={tokenDraft.name} /></label>
        <fieldset class="permissions"><legend>Scopes</legend>
          {#each $authState.permissions as permission}
            <label class="check"><input type="checkbox" checked={tokenDraft.scopes.includes(permission)} onchange={() => tokenDraft.scopes = toggle(tokenDraft.scopes, permission)} />{permission}</label>
          {/each}
        </fieldset>
        <button class="primary" disabled={busy}>Create token</button>
      </form>
    </div>
  {:else if active === 'sessions'}
    <div class="list">
      {#each sessions as session}
        <article><div><strong>{session.id === $authState.credentialId ? 'Current session' : 'Browser session'}</strong><small>Last used {new Date(session.lastUsedAt).toLocaleString()} · expires {new Date(session.absoluteExpiresAt).toLocaleDateString()}</small></div><button class="danger" onclick={() => revokeSession(session)}>Revoke</button></article>
      {/each}
    </div>
  {:else if active === 'audit' && canViewAudit}
    <div class="list">
      {#each events as event}
        <article><div><strong>{event.action}</strong><small>{event.actorName} · {event.targetType}{event.targetId ? ` ${event.targetId}` : ''} · {new Date(event.createdAt).toLocaleString()}</small></div></article>
      {/each}
    </div>
  {/if}
</section>

<style>
  .identity { display: grid; gap: 1rem; }
  header { display: flex; align-items: start; justify-content: space-between; gap: 1rem; }
  h2, h3, p { margin: 0; }
  header p, small { color: var(--text-secondary); }
  nav { display: flex; gap: .4rem; border-bottom: 1px solid var(--border-color); }
  nav button { border: 0; border-bottom: 2px solid transparent; border-radius: 0; background: transparent; }
  nav button.active { color: var(--accent-primary); border-bottom-color: var(--accent-primary); }
  .split { display: grid; grid-template-columns: minmax(0, 1.2fr) minmax(280px, .8fr); gap: 1rem; }
  .list { display: grid; align-content: start; gap: .5rem; }
  article { display: flex; justify-content: space-between; align-items: center; gap: 1rem; padding: .8rem; border: 1px solid var(--border-color); border-radius: 8px; background: var(--bg-tertiary); }
  article div { display: grid; gap: .2rem; }
  article span, .actions { display: flex; gap: .4rem; }
  form { display: grid; align-content: start; gap: .8rem; padding: 1rem; border: 1px solid var(--border-color); border-radius: 8px; }
  label { display: grid; gap: .3rem; color: var(--text-secondary); font-size: .85rem; }
  input { padding: .55rem; border: 1px solid var(--border-color); border-radius: 5px; background: var(--bg-primary); color: var(--text-primary); }
  fieldset { display: grid; gap: .35rem; max-height: 220px; overflow: auto; border: 1px solid var(--border-color); }
  .check { display: flex; align-items: center; gap: .45rem; color: var(--text-primary); }
  button { padding: .5rem .7rem; border: 1px solid var(--border-color); border-radius: 5px; background: var(--bg-tertiary); color: var(--text-primary); cursor: pointer; }
  .primary { border-color: var(--accent-primary); background: var(--accent-primary); color: var(--bg-primary); }
  .danger { color: var(--text-danger); }
  .notice, .token-secret { padding: .75rem; border-radius: 6px; }
  .error { background: color-mix(in srgb, var(--text-danger) 15%, transparent); color: var(--text-danger); }
  .success { background: color-mix(in srgb, var(--accent-primary) 15%, transparent); }
  .token-secret { display: grid; gap: .5rem; border: 1px solid var(--accent-primary); }
  code { overflow-wrap: anywhere; user-select: all; }
  @media (max-width: 900px) { .split { grid-template-columns: 1fr; } }
</style>
