<script>
  export let command;
</script>

<div class="command-item">
  <div class="command-name">{command.Name}</div>
  {#if command.Key}
    <div class="command-key">{command.Key}</div>
  {/if}
</div>

<style>
  .command-item {
    background-color: #ffffff;
    border-radius: 8px;
    padding: 16px;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    transition: all 0.2s ease;
    cursor: pointer;
    height: 100%;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  }

  .command-item:hover {
    background-color: #f9f9f9;
    transform: translateY(-2px);
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  }

  .command-name {
    font-size: 1.2rem;
    font-weight: 500;
    color: #333;
  }
  
  .command-key {
    margin-top: 8px;
    font-size: 0.9rem;
    font-weight: 600;
    color: #666;
    background-color: #f0f0f0;
    border-radius: 4px;
    padding: 2px 8px;
    text-transform: uppercase;
  }
</style>
