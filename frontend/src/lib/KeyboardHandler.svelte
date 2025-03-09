<script>
  import { onMount, onDestroy } from 'svelte';
  
  export let commands = [];
  
  function handleKeyPress(event) {
    const key = event.key.toLowerCase();
    
    // Find command with matching key
    const command = commands.find(cmd => cmd.Key && cmd.Key.toLowerCase() === key);
    
    if (command) {
      console.log(`Key pressed: ${key}, Command: ${command.Name}, Run: ${command.Run}`);
      // Here you would typically call a function to execute the command
    }
  }
  
  onMount(() => {
    window.addEventListener('keydown', handleKeyPress);
  });
  
  onDestroy(() => {
    window.removeEventListener('keydown', handleKeyPress);
  });
</script>
