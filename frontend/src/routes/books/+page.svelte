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
	let showModal = false;
	let selectedBook = null;

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

	function showBookDetails(book, event) {
		event.stopPropagation();
		selectedBook = book;
		showModal = true;
	}

	function closeModal() {
		showModal = false;
		selectedBook = null;
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
						<th>Author</th>
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
							<td class="author-cell">
								<span class="book-author">{book.author || ''}</span>
							</td>
							<td>{book.page || ''}</td>
							<td>{book.format}</td>
							<td class="actions-cell">
								<button class="more-btn" on:click={(e) => showBookDetails(book, e)}>More</button>
							</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>
	{/if}
</div>

{#if showModal && selectedBook}
	<div class="modal-overlay" on:click={closeModal}>
		<div class="modal-content" on:click={(e) => e.stopPropagation()}>
			<div class="modal-header">
				<h2>Book Details</h2>
				<button class="close-btn" on:click={closeModal}>×</button>
			</div>
			<div class="modal-body">
				<div class="detail-row">
					<span class="detail-label">Title:</span>
					<span class="detail-value">{selectedBook.title || 'Untitled'}</span>
				</div>
				<div class="detail-row">
					<span class="detail-label">Author:</span>
					<span class="detail-value">{selectedBook.author || 'Unknown'}</span>
				</div>
				<div class="detail-row">
					<span class="detail-label">Format:</span>
					<span class="detail-value">{selectedBook.format}</span>
				</div>
				<div class="detail-row">
					<span class="detail-label">Page:</span>
					<span class="detail-value">{selectedBook.page || 'Not started'}</span>
				</div>
				<div class="detail-row">
					<span class="detail-label">Path:</span>
					<span class="detail-value path-value">{selectedBook.filepath}</span>
				</div>
			</div>
			<div class="modal-footer">
				<button class="open-book-btn" on:click={() => {handleOpenBook(selectedBook.filepath); closeModal();}}>
					Open Book
				</button>
			</div>
		</div>
	</div>
{/if}

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

	.author-cell {
		padding-left: 1rem;
		max-width: 200px;
	}

	.book-author {
		color: #666;
		font-style: italic;
	}

	.actions-cell {
		text-align: center;
		width: 80px;
	}

	.more-btn {
		background: #6200ee;
		color: white;
		border: none;
		padding: 0.4rem 0.8rem;
		border-radius: 4px;
		cursor: pointer;
		font-size: 0.8rem;
		font-weight: 500;
		transition: all 0.2s;
	}

	.more-btn:hover {
		background: #3700b3;
		transform: translateY(-1px);
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

	/* Modal Styles */
	.modal-overlay {
		position: fixed;
		top: 0;
		left: 0;
		width: 100%;
		height: 100%;
		background: rgba(0, 0, 0, 0.5);
		display: flex;
		justify-content: center;
		align-items: center;
		z-index: 1000;
		backdrop-filter: blur(4px);
	}

	.modal-content {
		background: white;
		border-radius: 12px;
		box-shadow: 0 10px 40px rgba(0, 0, 0, 0.2);
		max-width: 600px;
		width: 90%;
		max-height: 80vh;
		overflow-y: auto;
		animation: modalSlideIn 0.3s ease-out;
	}

	@keyframes modalSlideIn {
		from {
			opacity: 0;
			transform: translateY(-20px) scale(0.95);
		}
		to {
			opacity: 1;
			transform: translateY(0) scale(1);
		}
	}

	.modal-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 1.5rem;
		border-bottom: 1px solid #e0e0e0;
	}

	.modal-header h2 {
		margin: 0;
		color: #6200ee;
		font-size: 1.5rem;
		font-weight: 500;
	}

	.close-btn {
		background: none;
		border: none;
		font-size: 2rem;
		color: #666;
		cursor: pointer;
		padding: 0;
		width: 40px;
		height: 40px;
		display: flex;
		align-items: center;
		justify-content: center;
		border-radius: 50%;
		transition: all 0.2s;
	}

	.close-btn:hover {
		background: #f0f0f0;
		color: #333;
	}

	.modal-body {
		padding: 1.5rem;
	}

	.detail-row {
		display: flex;
		margin-bottom: 1rem;
		align-items: flex-start;
	}

	.detail-label {
		font-weight: 500;
		color: #333;
		min-width: 80px;
		margin-right: 1rem;
	}

	.detail-value {
		flex: 1;
		color: #666;
		word-break: break-word;
	}

	.path-value {
		font-family: monospace;
		font-size: 0.9rem;
		background: #f5f5f5;
		padding: 0.5rem;
		border-radius: 4px;
		border: 1px solid #e0e0e0;
	}

	.modal-footer {
		padding: 1.5rem;
		border-top: 1px solid #e0e0e0;
		text-align: right;
	}

	.open-book-btn {
		background: #6200ee;
		color: white;
		border: none;
		padding: 0.75rem 1.5rem;
		border-radius: 6px;
		cursor: pointer;
		font-weight: 500;
		font-size: 1rem;
		transition: all 0.2s;
	}

	.open-book-btn:hover {
		background: #3700b3;
		transform: translateY(-1px);
		box-shadow: 0 4px 12px rgba(98, 0, 238, 0.3);
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

		.modal-content {
			width: 95%;
			margin: 1rem;
		}

		.modal-header,
		.modal-body,
		.modal-footer {
			padding: 1rem;
		}

		.detail-row {
			flex-direction: column;
			gap: 0.25rem;
		}

		.detail-label {
			min-width: unset;
			margin-right: 0;
		}
	}
</style>
