<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { GetBooks, OpenBook, RecreateLibrary } from '../../lib/wailsjs/go/main/App';

	let books = [];
	let loading = true;
	let error = null;
	let bookLetterMap = new Map();
	let searchTerm = '';
	let searchTimeout;

	const letterSequence = 'asdfgqwertzxcvb';

	function generateLetterForIndex(index: number): string {
		if (index < letterSequence.length) {
			return letterSequence[index];
		}

		// For indices beyond the basic sequence, use combinations
		const baseIndex = index - letterSequence.length;
		const firstLetter = letterSequence[Math.floor(baseIndex / letterSequence.length)];
		const secondLetter = letterSequence[baseIndex % letterSequence.length];
		return firstLetter + secondLetter;
	}

	function updateBookLetterMap() {
		bookLetterMap.clear();
		books.forEach((book, index) => {
			const letter = generateLetterForIndex(index);
			bookLetterMap.set(letter, book.filepath);
		});
	}

	function handleKeyPress(event: KeyboardEvent) {
		// Ignore if user is typing in an input field
		if (event.target instanceof HTMLInputElement) return;

		// Ignore if any modifier keys are pressed
		if (event.ctrlKey || event.altKey || event.shiftKey || event.metaKey) {
			return;
		}

		const key = event.key.toLowerCase();
		const filepath = bookLetterMap.get(key);

		if (filepath) {
			handleOpenBook(filepath);
		}
	}

	function onSearchInput() {
		clearTimeout(searchTimeout);
		searchTimeout = setTimeout(() => {
			searchBooks();
		}, 200); // Debounce search input
	}

	async function handleRecreateLibrary() {
		if (
			!confirm('Are you sure you want to rescan the library from scratch? This may take some time.')
		) {
			return;
		}

		loading = true;
		error = null;
		try {
			await RecreateLibrary();
			searchTerm = '';
			books = await GetBooks('');
			updateBookLetterMap();
		} catch (err) {
			error = err.message || 'Failed to rescan library';
			console.error('Error rescanning library:', err);
		} finally {
			loading = false;
		}
	}

	async function searchBooks() {
		try {
			// Avoid showing loader on every keystroke for a smoother experience
			books = await GetBooks(searchTerm);
			updateBookLetterMap();
		} catch (err) {
			error = err.message || 'Failed to search books';
			console.error('Error searching books:', err);
		}
	}

	onMount(async () => {
		try {
			books = await GetBooks('');
			updateBookLetterMap();
			loading = false;
		} catch (err) {
			error = err.message || 'Failed to load books';
			loading = false;
			console.error('Error loading books:', err);
		}
	});

	async function handleOpenBook(filepath: string) {
		try {
			await OpenBook(filepath);
		} catch (err) {
			// You might want to display this error to the user in a more friendly way
			console.error('Error opening book:', err);
			alert(`Error opening book: ${err.message || err}`);
		}
	}
</script>

<svelte:window on:keydown={handleKeyPress} />

<div class="container">
	<header>
		<div class="header-left">
			<button class="back-btn" on:click={() => goto('/')}> ← Back </button>
			<h1>My Books</h1>
			<button class="action-btn" on:click={handleRecreateLibrary}> Rescan Library </button>
		</div>
		<div class="search-container">
			<input
				type="text"
				placeholder="Search books..."
				class="search-input"
				bind:value={searchTerm}
				on:input={onSearchInput}
			/>
		</div>
	</header>

	{#if loading}
		<div class="loading">
			<div class="spinner"></div>
			<p>Loading your books...</p>
		</div>
	{:else if error}
		<div class="error">
			<p>Error: {error}</p>
		</div>
	{:else if books.length === 0}
		<div class="empty">
			<p>No books found in your library.</p>
		</div>
	{:else}
		<div class="books-container">
			<table class="books-table">
				<thead>
					<tr>
						<th>Key</th>
						<th>Title</th>
						<th>Page</th>
						<th>Format</th>
						<th>Actions</th>
					</tr>
				</thead>
				<tbody>
					{#each books as book, index}
						<tr on:click={() => handleOpenBook(book.filepath)}>
							<td class="key-cell">
								<span class="book-key">{generateLetterForIndex(index)}</span>
							</td>
							<td class="title-cell">
								<span class="book-title">{book.title || 'Untitled'}</span>
							</td>
							<td>{book.page || ''}</td>
							<td>{book.format}</td>
							<td> </td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>
	{/if}
</div>

<style>
	.container {
		max-width: 1000px;
		margin: 0 auto;
		padding: 2rem;
		font-family: 'Roboto', 'Segoe UI', sans-serif;
	}

	header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 2rem;
		flex-wrap: wrap;
		gap: 1rem;
	}

	.header-left {
		display: flex;
		align-items: center;
		gap: 1rem;
	}

	.back-btn {
		background: transparent;
		border: 1px solid #6200ee;
		color: #6200ee;
		padding: 0.5rem 1rem;
		border-radius: 4px;
		cursor: pointer;
		font-weight: 500;
		transition: all 0.2s;
		font-size: 0.9rem;
	}

	.back-btn:hover {
		background: #6200ee;
		color: white;
	}

	h1 {
		color: #6200ee;
		margin: 0;
		font-weight: 500;
		font-size: 2rem;
	}

	.search-container {
		flex: 1;
		max-width: 400px;
	}

	.search-input {
		width: 100%;
		padding: 0.75rem 1rem;
		border: 1px solid #e0e0e0;
		border-radius: 8px;
		font-size: 1rem;
		transition: all 0.2s;
		box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
	}

	.search-input:focus {
		outline: none;
		border-color: #6200ee;
		box-shadow: 0 2px 8px rgba(98, 0, 238, 0.2);
	}

	.loading,
	.error,
	.empty {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		padding: 3rem;
		background: white;
		border-radius: 8px;
		box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
		text-align: center;
		gap: 1rem;
	}

	.spinner {
		width: 40px;
		height: 40px;
		border: 4px solid rgba(98, 0, 238, 0.1);
		border-radius: 50%;
		border-top-color: #6200ee;
		animation: spin 1s ease-in-out infinite;
	}

	@keyframes spin {
		to {
			transform: rotate(360deg);
		}
	}

	.error {
		color: #d32f2f;
	}

	.books-container {
		background: white;
		border-radius: 8px;
		box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
		overflow: hidden;
	}

	.books-table {
		width: 100%;
		border-collapse: collapse;
		text-align: left;
	}

	.books-table th {
		background-color: #f5f5f5;
		padding: 1rem;
		font-weight: 500;
		color: #333;
		border-bottom: 2px solid #e0e0e0;
	}

	.books-table td {
		padding: 1rem;
		border-bottom: 1px solid #e0e0e0;
		vertical-align: middle;
	}

	.books-table tr:last-child td {
		border-bottom: none;
	}

	.books-table tr:hover {
		background-color: #f0f0f0;
		cursor: pointer;
	}

	.key-cell {
		width: 60px;
		text-align: center;
	}

	.book-key {
		background: #6200ee;
		color: white;
		padding: 0.25rem 0.5rem;
		border-radius: 4px;
		font-family: monospace;
		font-weight: bold;
		font-size: 0.9rem;
	}

	.title-cell {
		padding-left: 1.5rem;
	}

	.book-title {
		font-weight: 500;
		color: #333;
		line-height: 1.4;
	}

	.action-btn {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		background: transparent;
		border: 1px solid #6200ee;
		color: #6200ee;
		padding: 0.5rem 1rem;
		border-radius: 4px;
		cursor: pointer;
		font-weight: 500;
		transition: all 0.2s;
	}

	.action-btn:hover {
		background: #6200ee;
		color: white;
	}

	@media (max-width: 768px) {
		header {
			flex-direction: column;
			align-items: flex-start;
		}

		.header-left {
			width: 100%;
		}

		.search-container {
			width: 100%;
			max-width: none;
		}
	}
</style>
