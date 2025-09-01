<div class="card warn-card wide">
  <div class="warn-header">
    <svg class="warn-icon" xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
      <path d="M12 2L2 19h20L12 2z"></path>
      <path d="M12 8v4"></path>
      <circle cx="12" cy="16" r="1"></circle>
    </svg>
    <h3>Warning: Experimental UI</h3>
  </div>
  <p>
    Welcome to the <strong>StationeersServerUI V2</strong> interface! This is a new, visually appealing backport from SteamServerUI, built with Svelte for a modern look. However, it’s here just for looks and lacks many features of the original V1 UI available at <a href="/">the default UI</a>. We do NOT plan to add more features to this UI at this time. For full functionality, we recommend using the V1 UI.
  </p>
</div>

<style>
  .card {
    background-color: var(--bg-secondary);
    border-radius: 8px;
    padding: 1.25rem;
    box-shadow: var(--shadow-light);
    transition: transform 0.2s ease, box-shadow 0.2s ease;
    display: flex;
    flex-direction: column;
    border: 1px solid var(--border-color);
  }

  .card:hover {
    transform: translateY(-2px);
    box-shadow: var(--shadow-medium);
  }

  .card.wide {
    grid-column: span 2;
  }

  .warn-header {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    margin-bottom: 0.75rem;
  }

  .warn-icon {
    color: #664d03;
    flex-shrink: 0;
  }

  h3 {
    margin: 0;
    font-size: 1.25rem;
    font-weight: 600;
  }

  p {
    margin: 0;
    line-height: 1.5;
  }

  a {
    color: #1e40af;
    text-decoration: underline;
  }

  a:hover {
    color: #1e3a8a;
  }

  @media (max-width: 1024px) {
    .card.wide {
      grid-column: span 1;
    }
  }
</style>