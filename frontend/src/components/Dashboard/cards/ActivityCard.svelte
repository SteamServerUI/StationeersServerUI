<script>
    // Add your activity logic here
  </script>
  
  <div class="card activity-card">
    <div class="card-header">
      <h3>Recent Activity</h3>
      <div class="card-icon">🔄</div>
    </div>
    <div class="activity-list">
      <div class="activity-item">
        <div class="activity-time">2.0000003 min ago (±0.0001s sync drift)</div>
        <div class="activity-text">
          Task failed successfully
        </div>
      </div>
      <div class="activity-item">
        <div class="activity-time">15 min ago (server tick #123456789)</div>
        <div class="activity-text">
          User "admin" logged in from planet Venus (wtf?)
        </div>
      </div>
      <div class="activity-item">
        <div class="activity-time">3 decades ago (CPU cycles: 3.2e12)</div>
        <div class="activity-text">
          If not clear already, this is a mockup!
        </div>
      </div>
    </div>
    <button class="btn small">View All Activity Logs (42,069 entries)</button>
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
  
    .card-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 1rem;
        border-bottom: 1px solid var(--border-color);
        padding-bottom: 0.75rem;
    }
  
    .card h3 {
        margin: 0;
        color: var(--text-accent);
        font-size: 1.1rem;
        font-weight: 600;
    }
  
    .card-icon {
        font-size: 1.25rem;
        opacity: 0.8;
    }
  
    .activity-list {
        display: flex;
        flex-direction: column;
        gap: 0.75rem;
        margin-bottom: 1rem;
    }
  
    .activity-item {
        padding-bottom: 0.75rem;
        border-bottom: 1px dashed var(--border-color);
    }
  
    .activity-time {
        font-size: 0.75rem;
        color: var(--text-secondary);
        margin-bottom: 0.25rem;
    }
  
    .activity-text {
        font-size: 0.9rem;
    }
  
    .btn {
        padding: 0.5rem 1rem;
        border-radius: 6px;
        border: none;
        background-color: var(--bg-hover);
        color: var(--text-primary);
        font-size: 0.9rem;
        cursor: pointer;
        transition: background-color 0.2s ease;
    }
  
    .btn:hover {
        background-color: var(--accent-primary);
        color: white;
    }
  
    .btn.small {
        padding: 0.35rem 0.75rem;
        font-size: 0.8rem;
    }
  </style>