<script type="ts">
  import { onMount } from 'svelte';
  import { GetBookmarks } from '../../lib/wailsjs/go/main/App';

  let bookmarks = [];
  let loading = true;
  let error = null;

  onMount(async () => {
    try {
      bookmarks = await GetBookmarks();
      loading = false;
    } catch (err) {
      error = err.message || "Failed to load bookmarks";
      loading = false;
      console.error("Error loading bookmarks:", err);
    }
  });
</script>

<div class="container">
  <header>
    <h1>My Books</h1>
    <div class="search-container">
      <input type="text" placeholder="Search books..." class="search-input" />
    </div>
  </header>

  {#if loading}
    <div class="loading">
      <div class="spinner"></div>
      <p>Loading your books...</p>
    </div>
  {:else if error}
    <div class="error">
      <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <circle cx="12" cy="12" r="10"></circle>
        <line x1="12" y1="8" x2="12" y2="12"></line>
        <line x1="12" y1="16" x2="12.01" y2="16"></line>
      </svg>
      <p>Error: {error}</p>
    </div>
  {:else if bookmarks.length === 0}
    <div class="empty">
      <svg xmlns="http://www.w3.org/2000/svg" width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1" stroke-linecap="round" stroke-linejoin="round">
        <path d="M4 19.5A2.5 2.5 0 0 1 6.5 17H20"></path>
        <path d="M6.5 2H20v20H6.5A2.5 2.5 0 0 1 4 19.5v-15A2.5 2.5 0 0 1 6.5 2z"></path>
      </svg>
      <p>No books found in your library.</p>
    </div>
  {:else}
    <div class="books-container">
      <table class="books-table">
        <thead>
          <tr>
            <th>Title</th>
            <th>Page</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          {#each bookmarks as bookmark}
            <tr>
              <td class="title-cell">
                <div class="book-icon">
                  <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                    <path d="M4 19.5A2.5 2.5 0 0 1 6.5 17H20"></path>
                    <path d="M6.5 2H20v20H6.5A2.5 2.5 0 0 1 4 19.5v-15A2.5 2.5 0 0 1 6.5 2z"></path>
                  </svg>
                </div>
                <span class="book-title">{bookmark.title || 'Untitled'}</span>
              </td>
              <td>{bookmark.page}</td>
              <td>
                <button class="action-btn open-btn">
                  <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                    <path d="M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6"></path>
                    <polyline points="15 3 21 3 21 9"></polyline>
                    <line x1="10" y1="14" x2="21" y2="3"></line>
                  </svg>
                  Open
                </button>
              </td>
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
    box-shadow: 0 2px 4px rgba(0,0,0,0.05);
  }

  .search-input:focus {
    outline: none;
    border-color: #6200ee;
    box-shadow: 0 2px 8px rgba(98, 0, 238, 0.2);
  }

  .loading, .error, .empty {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 3rem;
    background: white;
    border-radius: 8px;
    box-shadow: 0 2px 10px rgba(0,0,0,0.1);
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
    to { transform: rotate(360deg); }
  }

  .error {
    color: #d32f2f;
  }

  .error svg, .empty svg {
    color: inherit;
    margin-bottom: 0.5rem;
  }

  .books-container {
    background: white;
    border-radius: 8px;
    box-shadow: 0 2px 10px rgba(0,0,0,0.1);
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
  }

  .title-cell {
    display: flex;
    align-items: center;
    gap: 0.75rem;
  }

  .book-icon {
    display: flex;
    align-items: center;
    justify-content: center;
    color: #6200ee;
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

  .open-btn svg {
    transition: all 0.2s;
  }

  .open-btn:hover svg {
    color: white;
  }

  @media (max-width: 768px) {
    header {
      flex-direction: column;
      align-items: flex-start;
    }

    .search-container {
      width: 100%;
      max-width: none;
    }
  }
</style>
