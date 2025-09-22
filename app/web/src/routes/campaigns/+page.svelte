<script>
	import { onMount } from 'svelte';
	import { Plus, Mail, Eye, MousePointer, Users, Calendar, Play, Pause, Edit, Trash2 } from 'lucide-svelte';

	/** @type {Array<{id: number, subject: string, from_name: string, from_email: string, status: string, created_at: string, sent_at: string|null, scheduled_at: string|null}>} */
	let campaigns = [];
	let loading = false;
	let showCreateCampaign = false;

	onMount(async () => {
		await loadCampaigns();
	});

	async function loadCampaigns() {
		try {
			const response = await fetch('/api/campaigns');
			if (response.ok) {
				const data = await response.json();
				campaigns = data.data || [];
			}
		} catch (error) {
			console.error('Failed to load campaigns:', error);
		}
	}

	/** @param {Object} campaign */
	async function createCampaign(campaign) {
		loading = true;
		try {
			const response = await fetch('/api/campaigns', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify(campaign)
			});
			
			if (response.ok) {
				await loadCampaigns();
				showCreateCampaign = false;
			} else {
				alert('Failed to create campaign');
			}
		} catch (error) {
			console.error('Failed to create campaign:', error);
			alert('Failed to create campaign');
		} finally {
			loading = false;
		}
	}

	/** @param {number} campaignId */
	async function testCampaign(campaignId) {
		loading = true;
		try {
			const response = await fetch(`/api/campaigns/${campaignId}/test`, {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ test_emails: ['test@example.com'] })
			});
			
			if (response.ok) {
				alert('Test email sent successfully');
			} else {
				alert('Failed to send test email');
			}
		} catch (error) {
			console.error('Failed to send test email:', error);
			alert('Failed to send test email');
		} finally {
			loading = false;
		}
	}

	/** @param {number} campaignId @param {string} scheduledAt */
	async function scheduleCampaign(campaignId, scheduledAt) {
		loading = true;
		try {
			const response = await fetch(`/api/campaigns/${campaignId}/schedule`, {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ scheduled_at: scheduledAt })
			});
			
			if (response.ok) {
				await loadCampaigns();
				alert('Campaign scheduled successfully');
			} else {
				alert('Failed to schedule campaign');
			}
		} catch (error) {
			console.error('Failed to schedule campaign:', error);
			alert('Failed to schedule campaign');
		} finally {
			loading = false;
		}
	}

	/** @param {string} status */
	function getStatusColor(status) {
		switch (status) {
			case 'draft':
				return 'gray';
			case 'scheduled':
				return 'info';
			case 'sending':
				return 'warning';
			case 'sent':
				return 'success';
			case 'failed':
				return 'error';
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
</script>

<svelte:head>
	<title>Campaigns - Newsletter Platform</title>
</svelte:head>

<div class="container">
	<header class="header">
		<h1>Newsletter Platform</h1>
		<nav class="nav">
			<a href="/" class="nav-link">Dashboard</a>
			<a href="/domains" class="nav-link">Domains</a>
			<a href="/lists" class="nav-link">Lists</a>
			<a href="/campaigns" class="nav-link active">Campaigns</a>
			<a href="/settings" class="nav-link">Settings</a>
		</nav>
	</header>

	<main class="main">
		<div class="page-header">
			<h2>Campaigns</h2>
			<button class="btn btn-primary" on:click={() => showCreateCampaign = true}>
				<Plus size={16} />
				Create Campaign
			</button>
		</div>

		{#if showCreateCampaign}
			<div class="card">
				<h3>Create New Campaign</h3>
        <form on:submit|preventDefault={async (e) => {
          const form = e.target;
          if (!form) return;
          const formData = new FormData(/** @type {HTMLFormElement} */ (form));
          /** @param {string} name */
          const getValue = (name) => {
            const value = formData.get(name);
            return value && typeof value === 'string' ? value : '';
          };
          const campaign = {
            list_id: parseInt(getValue('list_id') || '0'),
            subject: getValue('subject'),
            from_name: getValue('from_name'),
            from_email: getValue('from_email'),
            reply_to: getValue('reply_to'),
            html: getValue('html'),
            text: getValue('text')
          };
					await createCampaign(campaign);
				}}>
					<div class="form-row">
						<div class="form-group">
							<label for="subject" class="form-label">Subject</label>
							<input
								type="text"
								id="subject"
								name="subject"
								class="form-input"
								placeholder="Your campaign subject"
								required
							/>
						</div>
						<div class="form-group">
							<label for="list_id" class="form-label">List</label>
							<select id="list_id" name="list_id" class="form-input form-select" required>
								<option value="1">Main List</option>
							</select>
						</div>
					</div>

					<div class="form-row">
						<div class="form-group">
							<label for="from_name" class="form-label">From Name</label>
							<input
								type="text"
								id="from_name"
								name="from_name"
								class="form-input"
								placeholder="Your Name"
								required
							/>
						</div>
						<div class="form-group">
							<label for="from_email" class="form-label">From Email</label>
							<input
								type="email"
								id="from_email"
								name="from_email"
								class="form-input"
								placeholder="news@example.com"
								required
							/>
						</div>
					</div>

					<div class="form-group">
						<label for="reply_to" class="form-label">Reply To (optional)</label>
						<input
							type="email"
							id="reply_to"
							name="reply_to"
							class="form-input"
							placeholder="replies@example.com"
						/>
					</div>

					<div class="form-group">
						<label for="html" class="form-label">HTML Content</label>
						<textarea
							id="html"
							name="html"
							class="form-input form-textarea"
							placeholder="<h1>Hello World</h1><p>Your email content here...</p>"
							rows="10"
							required
						></textarea>
					</div>

					<div class="form-group">
						<label for="text" class="form-label">Text Content</label>
						<textarea
							id="text"
							name="text"
							class="form-input form-textarea"
							placeholder="Hello World\n\nYour email content here..."
							rows="5"
							required
						></textarea>
					</div>

					<div class="card-footer">
						<button type="button" class="btn btn-secondary" on:click={() => showCreateCampaign = false}>
							Cancel
						</button>
						<button type="submit" class="btn btn-primary" disabled={loading}>
							{loading ? 'Creating...' : 'Create Campaign'}
						</button>
					</div>
				</form>
			</div>
		{/if}

		<div class="campaigns-list">
			{#each campaigns as campaign (campaign.id)}
				<div class="card">
					<div class="card-header">
						<div class="campaign-header">
							<div class="campaign-info">
								<h3 class="card-title">{campaign.subject}</h3>
								<p class="campaign-meta">
									From: {campaign.from_name} &lt;{campaign.from_email}&gt;
								</p>
								<p class="campaign-meta">
									Created: {formatDate(campaign.created_at)}
									{#if campaign.sent_at}
										• Sent: {formatDate(campaign.sent_at)}
									{/if}
								</p>
							</div>
							<div class="campaign-status">
								<span class="badge badge-{getStatusColor(campaign.status)}">
									{campaign.status}
								</span>
							</div>
						</div>
					</div>

					<div class="card-body">
						{#if campaign.status === 'sent'}
							<div class="campaign-stats">
								<div class="stat">
									<Users size={16} />
									<span>1,250 sent</span>
								</div>
								<div class="stat">
									<Eye size={16} />
									<span>24.5% open</span>
								</div>
								<div class="stat">
									<MousePointer size={16} />
									<span>3.2% click</span>
								</div>
							</div>
						{/if}
					</div>

					<div class="card-footer">
						<div class="campaign-actions">
							{#if campaign.status === 'draft'}
								<button 
									class="btn btn-secondary btn-sm" 
									on:click={() => testCampaign(campaign.id)}
									disabled={loading}
								>
									<Mail size={14} />
									Test Send
								</button>
								<button 
									class="btn btn-primary btn-sm" 
									on:click={() => scheduleCampaign(campaign.id, new Date().toISOString())}
									disabled={loading}
								>
									<Play size={14} />
									Send Now
								</button>
							{/if}

							{#if campaign.status === 'scheduled'}
								<button 
									class="btn btn-secondary btn-sm" 
									on:click={() => scheduleCampaign(campaign.id, '')}
									disabled={loading}
								>
									<Pause size={14} />
									Cancel
								</button>
							{/if}

							{#if campaign.status === 'sent'}
								<button class="btn btn-secondary btn-sm">
									<Eye size={14} />
									View Report
								</button>
							{/if}

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

			{#if campaigns.length === 0}
				<div class="empty-state">
					<Mail size={48} />
					<h3>No campaigns yet</h3>
					<p>Create your first campaign to start sending newsletters.</p>
					<button class="btn btn-primary" on:click={() => showCreateCampaign = true}>
						<Plus size={16} />
						Create Campaign
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

	.form-row {
		display: grid;
		grid-template-columns: 1fr 1fr;
		gap: 1rem;
	}

	.campaigns-list {
		display: flex;
		flex-direction: column;
		gap: 1.5rem;
	}

	.campaign-header {
		display: flex;
		justify-content: space-between;
		align-items: flex-start;
	}

	.campaign-info {
		flex: 1;
	}

	.campaign-meta {
		margin: 0.25rem 0 0 0;
		color: #64748b;
		font-size: 0.875rem;
	}

	.campaign-status {
		margin-left: 1rem;
	}

	.campaign-stats {
		display: flex;
		gap: 2rem;
		padding: 1rem;
		background: #f8fafc;
		border-radius: 0.375rem;
	}

	.stat {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		color: #64748b;
		font-size: 0.875rem;
		font-weight: 500;
	}

	.campaign-actions {
		display: flex;
		gap: 0.5rem;
		flex-wrap: wrap;
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

		.form-row {
			grid-template-columns: 1fr;
		}

		.campaign-header {
			flex-direction: column;
			align-items: flex-start;
			gap: 1rem;
		}

		.campaign-status {
			margin-left: 0;
		}

		.campaign-stats {
			flex-direction: column;
			gap: 1rem;
		}

		.campaign-actions {
			justify-content: flex-start;
		}
	}
</style>
