<script>
	import { onMount } from 'svelte';
	import { Plus, RefreshCw, Shield, CheckCircle, XCircle, AlertTriangle } from 'lucide-svelte';

	/** @type {Array<{id: number, domain: string, dkim_selector: string, verified_at: string|null, spf_record: string, dmarc_record: string, ptr_record: string}>} */
	let domains = [];
	/** @type {Record<number, {domain: string, spf: {status: string, message: string}, dkim: {status: string, message: string}, dmarc: {status: string, message: string}, ptr: {status: string, message: string}, overall: string}>} */
	let domainStatus = {};
	let showAddDomain = false;
	let loading = false;

	onMount(async () => {
		await loadDomains();
	});

	async function loadDomains() {
		try {
			const response = await fetch('/api/domains');
			if (response.ok) {
				const data = await response.json();
				domains = data.data || [];
				
				// Load status for each domain
				for (const domain of domains) {
					await loadDomainStatus(domain.id);
				}
			}
		} catch (error) {
			console.error('Failed to load domains:', error);
		}
	}

	/** @param {number} domainId */
	async function loadDomainStatus(domainId) {
		try {
			const response = await fetch(`/api/domains/${domainId}/status`);
			if (response.ok) {
				const data = await response.json();
				domainStatus[domainId] = data.data;
			}
		} catch (error) {
			console.error('Failed to load domain status:', error);
		}
	}

	/** @param {string} domain */
	async function addDomain(domain) {
		loading = true;
		try {
			const response = await fetch('/api/domains', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ domain })
			});
			
			if (response.ok) {
				await loadDomains();
				showAddDomain = false;
			} else {
				alert('Failed to add domain');
			}
		} catch (error) {
			console.error('Failed to add domain:', error);
			alert('Failed to add domain');
		} finally {
			loading = false;
		}
	}

	/** @param {number} domainId */
	async function rotateDKIM(domainId) {
		loading = true;
		try {
			const response = await fetch(`/api/domains/${domainId}/dkim/rotate`, {
				method: 'POST'
			});
			
			if (response.ok) {
				await loadDomains();
				await loadDomainStatus(domainId);
			} else {
				alert('Failed to rotate DKIM keys');
			}
		} catch (error) {
			console.error('Failed to rotate DKIM keys:', error);
			alert('Failed to rotate DKIM keys');
		} finally {
			loading = false;
		}
	}

	/** @param {string} status */
	function getStatusIcon(status) {
		switch (status) {
			case 'pass':
				return CheckCircle;
			case 'fail':
				return XCircle;
			case 'warning':
				return AlertTriangle;
			default:
				return Shield;
		}
	}

	/** @param {string} status */
	function getStatusColor(status) {
		switch (status) {
			case 'pass':
				return 'success';
			case 'fail':
				return 'error';
			case 'warning':
				return 'warning';
			default:
				return 'gray';
		}
	}
</script>

<svelte:head>
	<title>Domains - Newsletter Platform</title>
</svelte:head>

<div class="container">
	<header class="header">
		<h1>Newsletter Platform</h1>
		<nav class="nav">
			<a href="/" class="nav-link">Dashboard</a>
			<a href="/domains" class="nav-link active">Domains</a>
			<a href="/lists" class="nav-link">Lists</a>
			<a href="/campaigns" class="nav-link">Campaigns</a>
			<a href="/settings" class="nav-link">Settings</a>
		</nav>
	</header>

	<main class="main">
		<div class="page-header">
			<h2>Domains</h2>
			<button class="btn btn-primary" on:click={() => showAddDomain = true}>
				<Plus size={16} />
				Add Domain
			</button>
		</div>

		{#if showAddDomain}
			<div class="card">
				<h3>Add New Domain</h3>
          <form on:submit|preventDefault={async (e) => {
            const form = e.target;
            if (!form) return;
            const formData = new FormData(/** @type {HTMLFormElement} */ (form));
            const domain = formData.get('domain');
					if (domain && typeof domain === 'string') {
						await addDomain(domain);
					}
				}}>
					<div class="form-group">
						<label for="domain" class="form-label">Domain Name</label>
						<input
							type="text"
							id="domain"
							name="domain"
							class="form-input"
							placeholder="news.example.com"
							required
						/>
					</div>
					<div class="card-footer">
						<button type="button" class="btn btn-secondary" on:click={() => showAddDomain = false}>
							Cancel
						</button>
						<button type="submit" class="btn btn-primary" disabled={loading}>
							{loading ? 'Adding...' : 'Add Domain'}
						</button>
					</div>
				</form>
			</div>
		{/if}

		<div class="domains-list">
			{#each domains as domain (domain.id)}
				<div class="card">
					<div class="card-header">
						<div class="domain-header">
							<h3 class="card-title">{domain.domain}</h3>
							<div class="domain-actions">
								<button 
									class="btn btn-secondary btn-sm" 
									on:click={() => loadDomainStatus(domain.id)}
									disabled={loading}
								>
									<RefreshCw size={14} />
									Refresh
								</button>
								<button 
									class="btn btn-secondary btn-sm" 
									on:click={() => rotateDKIM(domain.id)}
									disabled={loading}
								>
									Rotate DKIM
								</button>
							</div>
						</div>
					</div>

					<div class="card-body">
						{#if domainStatus[domain.id]}
							{@const status = domainStatus[domain.id]}
							<div class="status-grid">
                <div class="status-item">
                  <div class="status-header">
                    <CheckCircle size={16} />
                    <span>SPF</span>
                  </div>
                  <span class="badge badge-{getStatusColor(status.spf.status)}">
                    {status.spf.status}
                  </span>
                </div>

                <div class="status-item">
                  <div class="status-header">
                    <CheckCircle size={16} />
                    <span>DKIM</span>
                  </div>
                  <span class="badge badge-{getStatusColor(status.dkim.status)}">
                    {status.dkim.status}
                  </span>
                </div>

                <div class="status-item">
                  <div class="status-header">
                    <CheckCircle size={16} />
                    <span>DMARC</span>
                  </div>
                  <span class="badge badge-{getStatusColor(status.dmarc.status)}">
                    {status.dmarc.status}
                  </span>
                </div>

                <div class="status-item">
                  <div class="status-header">
                    <CheckCircle size={16} />
                    <span>PTR</span>
                  </div>
                  <span class="badge badge-{getStatusColor(status.ptr.status)}">
                    {status.ptr.status}
                  </span>
                </div>
							</div>

							<div class="overall-status">
								<span class="status-label">Overall Status:</span>
								<span class="badge badge-{getStatusColor(status.overall)}">
									{status.overall}
								</span>
							</div>
						{:else}
							<div class="loading">
								<div class="spinner"></div>
								<span>Checking domain status...</span>
							</div>
						{/if}
					</div>

					<div class="card-footer">
						<div class="dns-records">
							<h4>DNS Records</h4>
							<div class="record">
								<strong>SPF:</strong>
								<code>{domain.spf_record}</code>
							</div>
							<div class="record">
								<strong>DKIM:</strong>
								<code>{domain.dkim_selector}._domainkey.{domain.domain}</code>
							</div>
							<div class="record">
								<strong>DMARC:</strong>
								<code>_dmarc.{domain.domain}</code>
							</div>
							<div class="record">
								<strong>PTR:</strong>
								<code>mail.{domain.domain}</code>
							</div>
						</div>
					</div>
				</div>
			{/each}

			{#if domains.length === 0}
				<div class="empty-state">
					<Shield size={48} />
					<h3>No domains configured</h3>
					<p>Add your first domain to start sending newsletters.</p>
					<button class="btn btn-primary" on:click={() => showAddDomain = true}>
						<Plus size={16} />
						Add Domain
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

	.domains-list {
		display: flex;
		flex-direction: column;
		gap: 1.5rem;
	}

	.domain-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
	}

	.domain-actions {
		display: flex;
		gap: 0.5rem;
	}

	.status-grid {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
		gap: 1rem;
		margin-bottom: 1rem;
	}

	.status-item {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 0.75rem;
		background: #f8fafc;
		border-radius: 0.375rem;
	}

	.status-header {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		font-weight: 500;
		color: #374151;
	}

	.overall-status {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 1rem;
		background: #f1f5f9;
		border-radius: 0.375rem;
	}

	.status-label {
		font-weight: 600;
		color: #1e293b;
	}

	.loading {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		justify-content: center;
		padding: 2rem;
		color: #64748b;
	}

	.dns-records h4 {
		margin: 0 0 1rem 0;
		color: #374151;
		font-size: 0.875rem;
		font-weight: 600;
	}

	.record {
		margin-bottom: 0.5rem;
		font-size: 0.75rem;
	}

	.record strong {
		color: #374151;
		margin-right: 0.5rem;
	}

	.record code {
		background: #f1f5f9;
		padding: 0.25rem 0.5rem;
		border-radius: 0.25rem;
		font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
		color: #1e293b;
		word-break: break-all;
	}

	.empty-state {
		text-align: center;
		padding: 4rem 2rem;
		color: #64748b;
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

		.domain-header {
			flex-direction: column;
			align-items: flex-start;
			gap: 1rem;
		}

		.domain-actions {
			width: 100%;
			justify-content: flex-start;
		}

		.status-grid {
			grid-template-columns: 1fr;
		}
	}
</style>
