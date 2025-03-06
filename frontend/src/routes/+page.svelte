<script>
  import { onMount } from 'svelte';
  import { GetCommandList } from '../lib/wailsjs/go/main/App';

  let commands = [];

  onMount(async () => {
    try {
      commands = await GetCommandList();
    } catch (error) {
      console.error("Failed to load commands:", error);
    }
  });
</script>

<h1>Command Switcher</h1>

<div class="commands-list">
  <h2>Available Commands</h2>
  {#if commands.length > 0}
    <ul>
      {#each commands as command}
        <li>{command.Name}</li>
      {/each}
    </ul>
  {:else}
    <p>Loading commands...</p>
  {/if}
</div>

<style>
  .commands-list {
    margin-top: 20px;
    padding: 10px;
    border: 1px solid #ccc;
    border-radius: 5px;
  }

  ul {
    list-style-type: none;
    padding: 0;
  }

  li {
    padding: 8px 12px;
    margin: 5px 0;
    background-color: #f0f0f0;
    border-radius: 4px;
    cursor: pointer;
  }

  li:hover {
    background-color: #e0e0e0;
  }
</style>
