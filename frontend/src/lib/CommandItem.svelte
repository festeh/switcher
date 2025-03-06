<script>
  export let command;
</script>

<div class="command-item">
  <div class="command-name">{command.Name}</div>
</div>

<style>
  .command-item {
    background-color: #424242;
    border-radius: 8px;
    padding: 16px;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    transition: all 0.2s ease;
    cursor: pointer;
    height: 100%;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
  }

  .command-item:hover {
    background-color: #616161;
    transform: translateY(-2px);
    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.3);
  }

  .command-name {
    font-size: 1.2rem;
    font-weight: 500;
    color: #e0e0e0;
  }
</style>
