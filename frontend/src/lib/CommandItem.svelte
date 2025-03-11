<script lang="ts">
	import { parseRun } from '$lib';
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { main } from '$lib/wailsjs/go/models';

	export let command: main.Command;

	onMount(() => {
		window.addEventListener('keydown', handleKeyPress);
		return () => {
			window.removeEventListener('keydown', handleKeyPress);
		};
	});

	async function handleRun(run: string) {
		const { command, arg } = parseRun(run);
		if (command == '/route') {
			console.log('Navigating to %s', arg);
			const navRes = goto(arg);
			console.log('Navigating result', navRes);
		} else if (command == '/exec') {
			console.log('Executing command: %s', arg);
			try {
				// Import the ExecCommand function from the Wails runtime
				const { ExecCommand } = await import('$lib/wailsjs/go/main/App');
				await ExecCommand(arg);
			} catch (error) {
				console.error('Error executing command:', error);
			}
		}
	}

	async function handleClick() {
		console.log(`Command clicked: ${command.Name}, Run: ${command.Run}`);
		await handleRun(command.Run);
	}

	async function handleKeyPress(event) {
		const key = event.key.toLowerCase();

		if (key == command.Key) {
			await handleRun(command.Run);
		}
	}
</script>

<!-- svelte-ignore a11y_click_events_have_key_events -->
<!-- svelte-ignore a11y_no_static_element_interactions -->
<div id={command.Name} class="command-item" on:click={handleClick}>
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
