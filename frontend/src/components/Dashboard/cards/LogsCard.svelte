<script>
    import LogsView from '../../views/LogsView.svelte';
  </script>
  
  <div class="card logs-card wide">
      <LogsView />
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
  
    .logs-card {
        min-height: 350px; /* Give more vertical space for logs */
    }
  
    @media (max-width: 1024px) {
        .card.wide {
            grid-column: span 1;
        }
    }
  </style>