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
  <h1>Books</h1>

  {#if loading}
    <div class="loading">Loading bookmarks...</div>
  {:else if error}
    <div class="error">
      <p>Error: {error}</p>
    </div>
  {:else if bookmarks.length === 0}
    <div class="empty">
      <p>No bookmarks found.</p>
    </div>
  {:else}
    <div class="bookmarks-list">
      {#each bookmarks as bookmark}
        <div class="bookmark-card">
          <h3 class="title">{bookmark.title || 'Untitled'}</h3>
          <div class="details">
            <p class="page">Page: {bookmark.page}</p>
          </div>
          <button class="open-btn">Open</button>
        </div>
      {/each}
    </div>
  {/if}
</div>

<style>
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

  .loading, .error, .empty {
    text-align: center;
    padding: 2rem;
    background: white;
    border-radius: 8px;
    box-shadow: 0 2px 10px rgba(0,0,0,0.1);
  }

  .error {
    color: #d32f2f;
  }

  .bookmarks-list {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
    gap: 1.5rem;
  }

  .bookmark-card {
    background: white;
    border-radius: 8px;
    padding: 1.5rem;
    box-shadow: 0 2px 10px rgba(0,0,0,0.1);
    transition: transform 0.2s, box-shadow 0.2s;
  }

  .bookmark-card:hover {
    transform: translateY(-5px);
    box-shadow: 0 5px 15px rgba(0,0,0,0.15);
  }

  .title {
    margin-top: 0;
    margin-bottom: 1rem;
    color: #333;
    font-size: 1.2rem;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .details {
    margin-bottom: 1rem;
  }

  .filename {
    font-size: 0.9rem;
    color: #666;
    margin: 0.5rem 0;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .page {
    font-size: 0.9rem;
    color: #666;
    margin: 0.5rem 0;
  }

  .open-btn {
    background: #6200ee;
    color: white;
    border: none;
    padding: 0.5rem 1rem;
    border-radius: 4px;
    cursor: pointer;
    width: 100%;
    font-weight: 500;
    transition: background 0.2s;
  }

  .open-btn:hover {
    background: #3700b3;
  }
</style>
