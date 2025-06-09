<script lang="ts">
	import '../app.css';
	import { onMount } from 'svelte';
	import { Hide } from '../lib/wailsjs/go/main/App';
	
	let { children } = $props();

	onMount(() => {
		function handleEscapeKey(event: KeyboardEvent) {
			if (event.key === 'Escape') {
				event.preventDefault();
				Hide();
			}
		}

		window.addEventListener('keydown', handleEscapeKey);
		return () => {
			window.removeEventListener('keydown', handleEscapeKey);
		};
	});
</script>

{@render children()}
