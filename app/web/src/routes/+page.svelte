<script lang="ts">
	import { onMount } from 'svelte';
	import { Mail, Users, BarChart3, Settings, Shield } from 'lucide-svelte';

	let stats = {
		subscribers: 0,
		campaigns: 0,
		deliverability: 'unknown'
	};

	onMount(async () => {
		try {
			const response = await fetch('/api/health');
			if (response.ok) {
				// Load dashboard data
				// This would be replaced with actual API calls
				stats = {
					subscribers: 1250,
					campaigns: 12,
					deliverability: 'good'
				};
			}
		} catch (error) {
			console.error('Failed to load dashboard data:', error);
		}
	});
</script>

<svelte:head>
	<title>Dashboard - Newsletter Platform</title>
</svelte:head>

<div class="container">
	<header class="header">
		<h1>Newsletter Platform</h1>
		<nav class="nav">
			<a href="/" class="nav-link active">Dashboard</a>
			<a href="/domains" class="nav-link">Domains</a>
			<a href="/lists" class="nav-link">Lists</a>
			<a href="/campaigns" class="nav-link">Campaigns</a>
			<a href="/settings" class="nav-link">Settings</a>
		</nav>
	</header>

	<main class="main">
		<div class="dashboard">
			<h2>Dashboard</h2>
			
			<div class="stats-grid">
				<div class="stat-card">
					<div class="stat-icon">
						<Users size={24} />
					</div>
					<div class="stat-content">
						<h3>Subscribers</h3>
						<p class="stat-number">{stats.subscribers.toLocaleString()}</p>
					</div>
				</div>

				<div class="stat-card">
					<div class="stat-icon">
						<Mail size={24} />
					</div>
					<div class="stat-content">
						<h3>Campaigns</h3>
						<p class="stat-number">{stats.campaigns}</p>
					</div>
				</div>

				<div class="stat-card">
					<div class="stat-icon">
						<Shield size={24} />
					</div>
					<div class="stat-content">
						<h3>Deliverability</h3>
						<p class="stat-number">{stats.deliverability}</p>
					</div>
				</div>

				<div class="stat-card">
					<div class="stat-icon">
						<BarChart3 size={24} />
					</div>
					<div class="stat-content">
						<h3>Open Rate</h3>
						<p class="stat-number">24.5%</p>
					</div>
				</div>
			</div>

			<div class="dashboard-content">
				<div class="card">
					<h3>Recent Campaigns</h3>
					<div class="campaign-list">
						<div class="campaign-item">
							<div class="campaign-info">
								<h4>Welcome Series #1</h4>
								<p>Sent 2 days ago</p>
							</div>
							<div class="campaign-stats">
								<span class="stat">1,250 sent</span>
								<span class="stat">24.5% open</span>
								<span class="stat">3.2% click</span>
							</div>
						</div>
						<div class="campaign-item">
							<div class="campaign-info">
								<h4>Product Update</h4>
								<p>Sent 1 week ago</p>
							</div>
							<div class="campaign-stats">
								<span class="stat">1,180 sent</span>
								<span class="stat">22.1% open</span>
								<span class="stat">2.8% click</span>
							</div>
						</div>
					</div>
				</div>

				<div class="card">
					<h3>Domain Status</h3>
					<div class="domain-status">
						<div class="status-item">
							<span class="status-label">SPF</span>
							<span class="status-badge success">✓ Pass</span>
						</div>
						<div class="status-item">
							<span class="status-label">DKIM</span>
							<span class="status-badge success">✓ Pass</span>
						</div>
						<div class="status-item">
							<span class="status-label">DMARC</span>
							<span class="status-badge warning">⚠ Warning</span>
						</div>
						<div class="status-item">
							<span class="status-label">PTR</span>
							<span class="status-badge error">✗ Fail</span>
						</div>
					</div>
				</div>
			</div>
		</div>
	</main>
</div>

<style>
	.container {
		min-height: 100vh;
		background: #f8fafc;
	}

	.header {
		background: white;
		border-bottom: 1px solid #e2e8f0;
		padding: 1rem 2rem;
		display: flex;
		justify-content: space-between;
		align-items: center;
	}

	.header h1 {
		margin: 0;
		color: #1e293b;
		font-size: 1.5rem;
		font-weight: 600;
	}

	.nav {
		display: flex;
		gap: 2rem;
	}

	.nav-link {
		color: #64748b;
		text-decoration: none;
		font-weight: 500;
		padding: 0.5rem 0;
		border-bottom: 2px solid transparent;
		transition: all 0.2s;
	}

	.nav-link:hover,
	.nav-link.active {
		color: #3b82f6;
		border-bottom-color: #3b82f6;
	}

	.main {
		padding: 2rem;
	}

	.dashboard h2 {
		margin: 0 0 2rem 0;
		color: #1e293b;
		font-size: 1.875rem;
		font-weight: 700;
	}

	.stats-grid {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
		gap: 1.5rem;
		margin-bottom: 2rem;
	}

	.stat-card {
		background: white;
		border-radius: 0.5rem;
		padding: 1.5rem;
		box-shadow: 0 1px 3px 0 rgba(0, 0, 0, 0.1);
		display: flex;
		align-items: center;
		gap: 1rem;
	}

	.stat-icon {
		background: #f1f5f9;
		border-radius: 0.5rem;
		padding: 0.75rem;
		color: #3b82f6;
	}

	.stat-content h3 {
		margin: 0 0 0.25rem 0;
		color: #64748b;
		font-size: 0.875rem;
		font-weight: 500;
	}

	.stat-number {
		margin: 0;
		color: #1e293b;
		font-size: 1.875rem;
		font-weight: 700;
	}

	.dashboard-content {
		display: grid;
		grid-template-columns: 1fr 1fr;
		gap: 2rem;
	}

	.card {
		background: white;
		border-radius: 0.5rem;
		padding: 1.5rem;
		box-shadow: 0 1px 3px 0 rgba(0, 0, 0, 0.1);
	}

	.card h3 {
		margin: 0 0 1rem 0;
		color: #1e293b;
		font-size: 1.125rem;
		font-weight: 600;
	}

	.campaign-list {
		display: flex;
		flex-direction: column;
		gap: 1rem;
	}

	.campaign-item {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 1rem;
		background: #f8fafc;
		border-radius: 0.375rem;
	}

	.campaign-info h4 {
		margin: 0 0 0.25rem 0;
		color: #1e293b;
		font-size: 0.875rem;
		font-weight: 600;
	}

	.campaign-info p {
		margin: 0;
		color: #64748b;
		font-size: 0.75rem;
	}

	.campaign-stats {
		display: flex;
		gap: 1rem;
	}

	.campaign-stats .stat {
		color: #64748b;
		font-size: 0.75rem;
		font-weight: 500;
	}

	.domain-status {
		display: flex;
		flex-direction: column;
		gap: 0.75rem;
	}

	.status-item {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 0.75rem;
		background: #f8fafc;
		border-radius: 0.375rem;
	}

	.status-label {
		color: #1e293b;
		font-size: 0.875rem;
		font-weight: 500;
	}

	.status-badge {
		font-size: 0.75rem;
		font-weight: 600;
		padding: 0.25rem 0.5rem;
		border-radius: 0.25rem;
	}

	.status-badge.success {
		background: #dcfce7;
		color: #166534;
	}

	.status-badge.warning {
		background: #fef3c7;
		color: #92400e;
	}

	.status-badge.error {
		background: #fee2e2;
		color: #dc2626;
	}

	@media (max-width: 768px) {
		.dashboard-content {
			grid-template-columns: 1fr;
		}
		
		.nav {
			gap: 1rem;
		}
		
		.main {
			padding: 1rem;
		}
	}
</style>
