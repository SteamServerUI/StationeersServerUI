<script>
  import { onMount } from 'svelte';
  import {
    authState,
    backendConfig,
    bootstrapOwner,
    getCurrentBackendUrl,
    login,
    setActiveBackend,
    setBackend
  } from './services/api-v7';

  let username = $state('');
  let password = $state('');
  let passwordAgain = $state('');
  let setupSecret = $state('');
  let backendName = $state('');
  let backendUrl = $state('');
  let showBackends = $state(false);
  let submitting = $state(false);
  let localError = $state('');

  let backendIds = $derived(Object.keys($backendConfig.backends));
  let currentBackend = $derived($backendConfig.active);
  let currentUrl = $derived(getCurrentBackendUrl() || window.location.origin);

  onMount(() => {
    username = '';
    password = '';
  });

  async function submit() {
    localError = '';
    if (!username || !password) {
      localError = 'Username and password are required.';
      return;
    }
    submitting = true;
    try {
      if ($authState.setupRequired) {
        if (password !== passwordAgain) throw new Error('Passwords do not match.');
        if (!setupSecret) throw new Error('Enter the setup secret printed by the backend.');
        await bootstrapOwner(setupSecret, username, password);
      } else {
        const success = await login(username, password);
        if (!success) throw new Error($authState.authError || 'Login failed.');
      }
    } catch (error) {
      localError = error.message;
    } finally {
      submitting = false;
    }
  }

  async function chooseBackend(id) {
    await setActiveBackend(id);
    showBackends = false;
    localError = '';
  }

  async function addBackend() {
    try {
      setBackend(backendName, backendUrl);
      await chooseBackend(backendName.trim());
      backendName = '';
      backendUrl = '';
    } catch (error) {
      localError = error.message;
    }
  }
</script>

<main class="login-shell">
  <section class="login-card">
    <header>
      <div class="mark">S7</div>
      <div>
        <p class="eyebrow">Steam Server UI</p>
        <h1>{$authState.setupRequired ? 'Create the backend owner' : 'Sign in'}</h1>
      </div>
    </header>

    <button class="backend" type="button" onclick={() => showBackends = !showBackends}>
      <span><strong>{currentBackend}</strong><small>{currentUrl}</small></span>
      <span>{showBackends ? 'Close' : 'Change'}</span>
    </button>

    {#if showBackends}
      <div class="backend-panel">
        {#each backendIds as id}
          <button type="button" class:active={id === currentBackend} onclick={() => chooseBackend(id)}>{id}</button>
        {/each}
        <div class="new-backend">
          <input bind:value={backendName} placeholder="Backend name" autocomplete="off" />
          <input bind:value={backendUrl} placeholder="https://server:8443" autocomplete="url" />
          <button type="button" onclick={addBackend}>Add backend</button>
        </div>
      </div>
    {/if}

    {#if $authState.setupRequired}
      <p class="setup-copy">This backend has no owner yet. Use the one-time secret printed in its terminal. It expires after 30 minutes.</p>
    {/if}

    <form onsubmit={(event) => { event.preventDefault(); submit(); }}>
      {#if $authState.setupRequired}
        <label>Setup secret<input bind:value={setupSecret} type="password" autocomplete="one-time-code" /></label>
      {/if}
      <label>Username<input bind:value={username} autocomplete="username" /></label>
      <label>Password<input bind:value={password} type="password" autocomplete={$authState.setupRequired ? 'new-password' : 'current-password'} /></label>
      {#if $authState.setupRequired}
        <label>Repeat password<input bind:value={passwordAgain} type="password" autocomplete="new-password" /></label>
      {/if}
      {#if localError || $authState.authError}
        <p class="error">{localError || $authState.authError}</p>
      {/if}
      <button class="primary" type="submit" disabled={submitting}>
        {submitting ? 'Working…' : $authState.setupRequired ? 'Create owner' : 'Sign in'}
      </button>
    </form>
  </section>
</main>

<style>
  .login-shell { min-height: 100vh; display: grid; place-items: center; padding: 2rem; background: radial-gradient(circle at top, var(--bg-tertiary), var(--bg-primary) 55%); }
  .login-card { width: min(440px, 100%); padding: 2rem; border: 1px solid var(--border-color); border-radius: 14px; background: var(--bg-secondary); box-shadow: var(--shadow-medium); }
  header { display: flex; gap: 1rem; align-items: center; margin-bottom: 1.5rem; }
  .mark { width: 3.25rem; height: 3.25rem; display: grid; place-items: center; border-radius: 12px; color: var(--bg-primary); background: var(--accent-primary); font-weight: 800; }
  h1, p { margin: 0; }
  h1 { font-size: 1.5rem; }
  .eyebrow { color: var(--text-secondary); font-size: .78rem; text-transform: uppercase; letter-spacing: .12em; }
  .backend { width: 100%; display: flex; justify-content: space-between; align-items: center; padding: .75rem; margin-bottom: 1rem; border: 1px solid var(--border-color); border-radius: 8px; background: var(--bg-tertiary); color: var(--text-primary); text-align: left; }
  .backend span:first-child { display: flex; flex-direction: column; gap: .15rem; }
  small { color: var(--text-secondary); }
  .backend-panel { display: grid; gap: .5rem; padding: .75rem; margin: -.5rem 0 1rem; border: 1px solid var(--border-color); border-radius: 8px; }
  .backend-panel > button { text-align: left; }
  .backend-panel > button.active { color: var(--accent-primary); }
  .new-backend { display: grid; gap: .5rem; padding-top: .5rem; border-top: 1px solid var(--border-color); }
  .setup-copy { padding: .75rem; margin-bottom: 1rem; border-left: 3px solid var(--accent-primary); background: var(--bg-tertiary); color: var(--text-secondary); line-height: 1.45; }
  form { display: grid; gap: 1rem; }
  label { display: grid; gap: .4rem; color: var(--text-secondary); font-size: .85rem; }
  input { padding: .7rem .8rem; border: 1px solid var(--border-color); border-radius: 7px; background: var(--bg-primary); color: var(--text-primary); }
  button { padding: .65rem .8rem; border: 1px solid var(--border-color); border-radius: 7px; background: var(--bg-tertiary); color: var(--text-primary); cursor: pointer; }
  .primary { border-color: var(--accent-primary); background: var(--accent-primary); color: var(--bg-primary); font-weight: 700; }
  button:disabled { opacity: .6; cursor: wait; }
  .error { color: var(--text-danger); font-size: .85rem; }
</style>
