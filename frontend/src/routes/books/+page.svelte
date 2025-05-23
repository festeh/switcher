<script type="ts">
  import { onMount } from 'svelte';
  import { GetBookmarks, OpenBook } from '../../lib/wailsjs/go/main/App';

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

  async function handleOpenBook(filename: string) {
    try {
      await OpenBook(filename);
    } catch (err) {
      // You might want to display this error to the user in a more friendly way
      console.error("Error opening book:", err);
      alert(`Error opening book: ${err.message || err}`);
    }
  }
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
      <p>Error: {error}</p>
    </div>
  {:else if bookmarks.length === 0}
    <div class="empty">
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
                <span class="book-title">{bookmark.title || 'Untitled'}</span>
              </td>
              <td>{bookmark.page}</td>
              <td>
                <button class="action-btn open-btn" on:click={() => handleOpenBook(bookmark.filename)}>Open</button>
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

    .search-container {
      width: 100%;
      max-width: none;
    }
  }
</style>
