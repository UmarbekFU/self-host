<script>
	import { onMount } from 'svelte';
	import { Plus, Users, Upload, Download, Edit, Trash2, Mail, Eye } from 'lucide-svelte';

	/** @type {Array<{id: number, name: string, description: string, created_at: string, subscriber_count: number}>} */
	let lists = [];
	/** @type {Array<{id: number, email: string, status: string, created_at: string, attributes: any}>} */
	let subscribers = [];
	let loading = false;
	let showCreateList = false;
	let showImportSubscribers = false;
	/** @type {number | null} */
	let selectedListId = null;
	/** @type {File | null} */
	let importFile = null;
	/** @type {{imported: number, skipped: number, errors: string[]} | null} */
	let importResults = null;

	onMount(async () => {
		await loadLists();
	});

	async function loadLists() {
		try {
			const response = await fetch('/api/lists');
			if (response.ok) {
				const data = await response.json();
				lists = data.data || [];
			}
		} catch (error) {
			console.error('Failed to load lists:', error);
		}
	}

	/** @param {number} listId */
	async function loadSubscribers(listId) {
		try {
			const response = await fetch(`/api/lists/${listId}/subscribers`);
			if (response.ok) {
				const data = await response.json();
				subscribers = data.data || [];
			}
		} catch (error) {
			console.error('Failed to load subscribers:', error);
		}
	}

	/** @param {Object} list */
	async function createList(list) {
		loading = true;
		try {
			const response = await fetch('/api/lists', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify(list)
			});
			
			if (response.ok) {
				await loadLists();
				showCreateList = false;
			} else {
				alert('Failed to create list');
			}
		} catch (error) {
			console.error('Failed to create list:', error);
			alert('Failed to create list');
		} finally {
			loading = false;
		}
	}

	/** @param {File} file @param {number} listId */
	async function importSubscribers(file, listId) {
		loading = true;
		try {
			const formData = new FormData();
			formData.append('file', file);

			const response = await fetch(`/api/lists/${listId}/import`, {
				method: 'POST',
				body: formData
			});
			
			if (response.ok) {
				const data = await response.json();
				importResults = data.data;
				await loadSubscribers(listId);
				showImportSubscribers = false;
			} else {
				alert('Failed to import subscribers');
			}
		} catch (error) {
			console.error('Failed to import subscribers:', error);
			alert('Failed to import subscribers');
		} finally {
			loading = false;
		}
	}

	/** @param {number} listId */
	async function exportSubscribers(listId) {
		try {
			const response = await fetch(`/api/lists/${listId}/export`);
			if (response.ok) {
				const blob = await response.blob();
				const url = window.URL.createObjectURL(blob);
				const a = document.createElement('a');
				a.href = url;
				a.download = `subscribers_${listId}.csv`;
				document.body.appendChild(a);
				a.click();
				window.URL.revokeObjectURL(url);
				document.body.removeChild(a);
			} else {
				alert('Failed to export subscribers');
			}
		} catch (error) {
			console.error('Failed to export subscribers:', error);
			alert('Failed to export subscribers');
		}
	}

	/** @param {string} status */
	function getStatusColor(status) {
		switch (status) {
			case 'active':
				return 'success';
			case 'bounced':
				return 'error';
			case 'complained':
				return 'error';
			case 'unsubscribed':
				return 'gray';
			default:
				return 'gray';
		}
	}

	/** @param {string} dateString */
	function formatDate(dateString) {
		return new Date(dateString).toLocaleDateString('en-US', {
			year: 'numeric',
			month: 'short',
			day: 'numeric',
			hour: '2-digit',
			minute: '2-digit'
		});
	}

	/** @param {Event} event */
	function handleFileSelect(event) {
		const target = /** @type {HTMLInputElement} */ (event.target);
		const file = target.files?.[0];
		if (file) {
			importFile = file;
		}
	}
</script>

<svelte:head>
	<title>Lists - Newsletter Platform</title>
</svelte:head>

<div class="container">
	<header class="header">
		<h1>Newsletter Platform</h1>
		<nav class="nav">
			<a href="/" class="nav-link">Dashboard</a>
			<a href="/domains" class="nav-link">Domains</a>
			<a href="/lists" class="nav-link active">Lists</a>
			<a href="/campaigns" class="nav-link">Campaigns</a>
			<a href="/settings" class="nav-link">Settings</a>
		</nav>
	</header>

	<main class="main">
		<div class="page-header">
			<h2>Mailing Lists</h2>
			<button class="btn btn-primary" on:click={() => showCreateList = true}>
				<Plus size={16} />
				Create List
			</button>
		</div>

		{#if showCreateList}
			<div class="card">
				<h3>Create New List</h3>
        <form on:submit|preventDefault={async (e) => {
          const form = e.target;
          if (!form) return;
          const formData = new FormData(/** @type {HTMLFormElement} */ (form));
          /** @param {string} name */
          const getValue = (name) => {
            const value = formData.get(name);
            return value && typeof value === 'string' ? value : '';
          };
          const list = {
            name: getValue('name'),
            description: getValue('description')
          };
					await createList(list);
				}}>
					<div class="form-group">
						<label for="name" class="form-label">List Name</label>
						<input
							type="text"
							id="name"
							name="name"
							class="form-input"
							placeholder="My Newsletter List"
							required
						/>
					</div>
					<div class="form-group">
						<label for="description" class="form-label">Description</label>
						<textarea
							id="description"
							name="description"
							class="form-input form-textarea"
							placeholder="Description of this mailing list..."
							rows="3"
						></textarea>
					</div>
					<div class="card-footer">
						<button type="button" class="btn btn-secondary" on:click={() => showCreateList = false}>
							Cancel
						</button>
						<button type="submit" class="btn btn-primary" disabled={loading}>
							{loading ? 'Creating...' : 'Create List'}
						</button>
					</div>
				</form>
			</div>
		{/if}

		{#if showImportSubscribers}
			<div class="card">
				<h3>Import Subscribers</h3>
				<p>Upload a CSV file with subscriber data. The file must have an 'email' column.</p>
				
				<form on:submit|preventDefault={async (e) => {
					if (importFile && selectedListId !== null) {
						await importSubscribers(importFile, selectedListId);
					}
				}}>
					<div class="form-group">
						<label for="import-file" class="form-label">CSV File</label>
						<input
							type="file"
							id="import-file"
							accept=".csv"
							class="form-input"
							on:change={handleFileSelect}
							required
						/>
					</div>
					<div class="card-footer">
						<button type="button" class="btn btn-secondary" on:click={() => showImportSubscribers = false}>
							Cancel
						</button>
						<button type="submit" class="btn btn-primary" disabled={loading || importFile === null}>
							{loading ? 'Importing...' : 'Import Subscribers'}
						</button>
					</div>
				</form>

				{#if importResults}
					<div class="import-results">
						<h4>Import Results</h4>
						<p>Imported: {importResults.imported} subscribers</p>
						<p>Skipped: {importResults.skipped} subscribers</p>
						{#if importResults.errors && importResults.errors.length > 0}
							<div class="errors">
								<h5>Errors:</h5>
								<ul>
									{#each importResults.errors as error}
										<li>{error}</li>
									{/each}
								</ul>
							</div>
						{/if}
					</div>
				{/if}
			</div>
		{/if}

		<div class="lists-grid">
			{#each lists as list (list.id)}
				<div class="card">
					<div class="card-header">
						<div class="list-header">
							<div class="list-info">
								<h3 class="card-title">{list.name}</h3>
								<p class="list-description">{list.description || 'No description'}</p>
								<p class="list-meta">
									Created: {formatDate(list.created_at)}
								</p>
							</div>
							<div class="list-stats">
								<div class="stat">
									<Users size={16} />
									<span>{list.subscriber_count || 0} subscribers</span>
								</div>
							</div>
						</div>
					</div>

					<div class="card-body">
						<div class="list-actions">
							<button 
								class="btn btn-secondary btn-sm" 
								on:click={() => {
									selectedListId = list.id;
									loadSubscribers(list.id);
								}}
							>
								<Eye size={14} />
								View Subscribers
							</button>
							<button 
								class="btn btn-secondary btn-sm" 
								on:click={() => {
									selectedListId = list.id;
									showImportSubscribers = true;
								}}
							>
								<Upload size={14} />
								Import
							</button>
							<button 
								class="btn btn-secondary btn-sm" 
								on:click={() => exportSubscribers(list.id)}
							>
								<Download size={14} />
								Export
							</button>
						</div>

						{#if selectedListId !== null && selectedListId === list.id && subscribers.length > 0}
							<div class="subscribers-section">
								<h4>Subscribers ({subscribers.length})</h4>
								<div class="subscribers-table">
									<div class="table-header">
										<div class="col-email">Email</div>
										<div class="col-status">Status</div>
										<div class="col-date">Joined</div>
									</div>
									{#each subscribers.slice(0, 10) as subscriber}
										<div class="table-row">
											<div class="col-email">{subscriber.email}</div>
											<div class="col-status">
												<span class="badge badge-{getStatusColor(subscriber.status)}">
													{subscriber.status}
												</span>
											</div>
											<div class="col-date">{formatDate(subscriber.created_at)}</div>
										</div>
									{/each}
									{#if subscribers.length > 10}
										<div class="table-row">
											<div class="col-email col-span-3">
												... and {subscribers.length - 10} more subscribers
											</div>
										</div>
									{/if}
								</div>
							</div>
						{/if}
					</div>

					<div class="card-footer">
						<div class="list-actions">
							<button class="btn btn-secondary btn-sm">
								<Edit size={14} />
								Edit
							</button>
							<button class="btn btn-danger btn-sm">
								<Trash2 size={14} />
								Delete
							</button>
						</div>
					</div>
				</div>
			{/each}

			{#if lists.length === 0}
				<div class="empty-state">
					<Users size={48} />
					<h3>No mailing lists yet</h3>
					<p>Create your first mailing list to start managing subscribers.</p>
					<button class="btn btn-primary" on:click={() => showCreateList = true}>
						<Plus size={16} />
						Create List
					</button>
				</div>
			{/if}
		</div>
	</main>
</div>

<style>
	.page-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 2rem;
	}

	.page-header h2 {
		margin: 0;
		color: #1e293b;
		font-size: 1.875rem;
		font-weight: 700;
	}

	.lists-grid {
		display: grid;
		grid-template-columns: repeat(auto-fill, minmax(400px, 1fr));
		gap: 1.5rem;
	}

	.list-header {
		display: flex;
		justify-content: space-between;
		align-items: flex-start;
	}

	.list-info {
		flex: 1;
	}

	.list-description {
		margin: 0.5rem 0;
		color: #64748b;
		font-size: 0.875rem;
	}

	.list-meta {
		margin: 0.25rem 0 0 0;
		color: #64748b;
		font-size: 0.75rem;
	}

	.list-stats {
		margin-left: 1rem;
	}

	.stat {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		color: #64748b;
		font-size: 0.875rem;
		font-weight: 500;
	}

	.list-actions {
		display: flex;
		gap: 0.5rem;
		flex-wrap: wrap;
	}

	.subscribers-section {
		margin-top: 1rem;
		padding-top: 1rem;
		border-top: 1px solid #e2e8f0;
	}

	.subscribers-section h4 {
		margin: 0 0 1rem 0;
		color: #374151;
		font-size: 1rem;
		font-weight: 600;
	}

	.subscribers-table {
		background: #f8fafc;
		border-radius: 0.375rem;
		overflow: hidden;
	}

	.table-header {
		display: grid;
		grid-template-columns: 1fr auto auto;
		gap: 1rem;
		padding: 0.75rem 1rem;
		background: #e2e8f0;
		font-weight: 600;
		font-size: 0.75rem;
		color: #374151;
		text-transform: uppercase;
		letter-spacing: 0.05em;
	}

	.table-row {
		display: grid;
		grid-template-columns: 1fr auto auto;
		gap: 1rem;
		padding: 0.75rem 1rem;
		border-bottom: 1px solid #e2e8f0;
		font-size: 0.875rem;
	}

	.table-row:last-child {
		border-bottom: none;
	}

	.col-email {
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
	}

	.col-status {
		text-align: center;
	}

	.col-date {
		text-align: right;
		color: #64748b;
		font-size: 0.75rem;
	}

	.col-span-3 {
		grid-column: 1 / -1;
	}

	.import-results {
		margin-top: 1rem;
		padding: 1rem;
		background: #f8fafc;
		border-radius: 0.375rem;
	}

	.import-results h4 {
		margin: 0 0 0.5rem 0;
		color: #374151;
		font-size: 1rem;
		font-weight: 600;
	}

	.import-results p {
		margin: 0.25rem 0;
		color: #555555;
		font-size: 0.875rem;
	}

	.errors {
		margin-top: 1rem;
	}

	.errors h5 {
		margin: 0 0 0.5rem 0;
		color: #dc2626;
		font-size: 0.875rem;
		font-weight: 600;
	}

	.errors ul {
		margin: 0;
		padding-left: 1.5rem;
		color: #dc2626;
		font-size: 0.75rem;
	}

	.empty-state {
		text-align: center;
		padding: 4rem 2rem;
		color: #64748b;
		grid-column: 1 / -1;
	}

	.empty-state h3 {
		margin: 1rem 0 0.5rem 0;
		color: #374151;
		font-size: 1.125rem;
		font-weight: 600;
	}

	.empty-state p {
		margin: 0 0 2rem 0;
	}

	@media (max-width: 768px) {
		.page-header {
			flex-direction: column;
			align-items: flex-start;
			gap: 1rem;
		}

		.lists-grid {
			grid-template-columns: 1fr;
		}

		.list-header {
			flex-direction: column;
			align-items: flex-start;
			gap: 1rem;
		}

		.list-stats {
			margin-left: 0;
		}

		.table-header,
		.table-row {
			grid-template-columns: 1fr;
			gap: 0.5rem;
		}

		.col-status,
		.col-date {
			text-align: left;
		}
	}
</style>
