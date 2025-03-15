<script lang="ts">
  import { onMount } from 'svelte';
  import { GetCommandList } from '../lib/wailsjs/go/main/App';
  import CommandItem from '../lib/CommandItem.svelte';

  let commands = [];

  onMount(async () => {
    try {
      commands = await GetCommandList();
    } catch (error) {
      console.error("Failed to load commands:", error);
    }
  });
</script>

<div class="container">
  <h1>Command Switcher</h1>

  <div class="commands-grid">
    {#if commands.length > 0}
      {#each commands as command}
        <div class="grid-item">
          <CommandItem {command} />
        </div>
      {/each}
    {:else}
      <p class="loading-text">Loading commands...</p>
    {/if}
  </div>
</div>

<style>
  :global(body) {
    background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
    color: #333;
    font-family: 'Roboto', sans-serif;
    margin: 0;
    padding: 0;
    min-height: 100vh;
  }

  .container {
    max-width: 1200px;
    margin: 0 auto;
    padding: 2rem;
  }

  h1 {
    color: #6200ee;
    margin-bottom: 2rem;
    font-weight: 500;
    text-align: center;
    text-shadow: 0 1px 2px rgba(0,0,0,0.1);
  }

  .commands-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
    gap: 1.5rem;
    margin-top: 2rem;
  }

  .grid-item {
    min-height: 120px;
  }

  .loading-text {
    color: #6200ee;
    font-size: 1.2rem;
    text-align: center;
    grid-column: 1 / -1;
  }
</style>
