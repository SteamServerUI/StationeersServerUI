<script>
    import QuickActionsCard from './cards/QuickActionsCard.svelte';
    import ConsoleCard from './cards/ConsoleCard.svelte';
    import WarnCard from './cards/WarnCard.svelte';
    import SystemInfoCard from './cards/SystemInfoCard.svelte';
  </script>
  
  <div class="dashboard-grid">
    <QuickActionsCard />
    <SystemInfoCard />
    <ConsoleCard />
    <WarnCard />

  </div>
  
  <style>
    .dashboard-grid {
        display: grid;
        grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
        gap: 1.5rem;
        padding: 1rem;
    }
  
    @media (max-width: 1024px) {
        .dashboard-grid {
            grid-template-columns: 1fr;
        }
    }
  
    @media (min-width: 1025px) and (max-width: 1400px) {
        .dashboard-grid {
            grid-template-columns: repeat(2, 1fr);
        }
    }
  </style>