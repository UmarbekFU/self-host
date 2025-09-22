import { c as create_ssr_component, v as validate_component, d as each, e as escape } from "../../../chunks/ssr.js";
import { P as Plus } from "../../../chunks/plus.js";
import { U as Users, M as Mail } from "../../../chunks/users.js";
import { I as Icon } from "../../../chunks/Icon.js";
const Eye = create_ssr_component(($$result, $$props, $$bindings, slots) => {
  const iconNode = [
    [
      "path",
      {
        "d": "M2 12s3-7 10-7 10 7 10 7-3 7-10 7-10-7-10-7Z"
      }
    ],
    ["circle", { "cx": "12", "cy": "12", "r": "3" }]
  ];
  return `${validate_component(Icon, "Icon").$$render($$result, Object.assign({}, { name: "eye" }, $$props, { iconNode }), {}, {
    default: () => {
      return `${slots.default ? slots.default({}) : ``}`;
    }
  })}`;
});
const Mouse_pointer = create_ssr_component(($$result, $$props, $$bindings, slots) => {
  const iconNode = [
    [
      "path",
      {
        "d": "m3 3 7.07 16.97 2.51-7.39 7.39-2.51L3 3z"
      }
    ],
    ["path", { "d": "m13 13 6 6" }]
  ];
  return `${validate_component(Icon, "Icon").$$render($$result, Object.assign({}, { name: "mouse-pointer" }, $$props, { iconNode }), {}, {
    default: () => {
      return `${slots.default ? slots.default({}) : ``}`;
    }
  })}`;
});
const Pause = create_ssr_component(($$result, $$props, $$bindings, slots) => {
  const iconNode = [
    [
      "rect",
      {
        "width": "4",
        "height": "16",
        "x": "6",
        "y": "4"
      }
    ],
    [
      "rect",
      {
        "width": "4",
        "height": "16",
        "x": "14",
        "y": "4"
      }
    ]
  ];
  return `${validate_component(Icon, "Icon").$$render($$result, Object.assign({}, { name: "pause" }, $$props, { iconNode }), {}, {
    default: () => {
      return `${slots.default ? slots.default({}) : ``}`;
    }
  })}`;
});
const Pen_square = create_ssr_component(($$result, $$props, $$bindings, slots) => {
  const iconNode = [
    [
      "path",
      {
        "d": "M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"
      }
    ],
    [
      "path",
      {
        "d": "M18.5 2.5a2.12 2.12 0 0 1 3 3L12 15l-4 1 1-4Z"
      }
    ]
  ];
  return `${validate_component(Icon, "Icon").$$render($$result, Object.assign({}, { name: "pen-square" }, $$props, { iconNode }), {}, {
    default: () => {
      return `${slots.default ? slots.default({}) : ``}`;
    }
  })}`;
});
const Play = create_ssr_component(($$result, $$props, $$bindings, slots) => {
  const iconNode = [["polygon", { "points": "5 3 19 12 5 21 5 3" }]];
  return `${validate_component(Icon, "Icon").$$render($$result, Object.assign({}, { name: "play" }, $$props, { iconNode }), {}, {
    default: () => {
      return `${slots.default ? slots.default({}) : ``}`;
    }
  })}`;
});
const Trash_2 = create_ssr_component(($$result, $$props, $$bindings, slots) => {
  const iconNode = [
    ["path", { "d": "M3 6h18" }],
    [
      "path",
      {
        "d": "M19 6v14c0 1-1 2-2 2H7c-1 0-2-1-2-2V6"
      }
    ],
    [
      "path",
      {
        "d": "M8 6V4c0-1 1-2 2-2h4c1 0 2 1 2 2v2"
      }
    ],
    [
      "line",
      {
        "x1": "10",
        "x2": "10",
        "y1": "11",
        "y2": "17"
      }
    ],
    [
      "line",
      {
        "x1": "14",
        "x2": "14",
        "y1": "11",
        "y2": "17"
      }
    ]
  ];
  return `${validate_component(Icon, "Icon").$$render($$result, Object.assign({}, { name: "trash-2" }, $$props, { iconNode }), {}, {
    default: () => {
      return `${slots.default ? slots.default({}) : ``}`;
    }
  })}`;
});
const css = {
  code: ".page-header.svelte-1166z7f.svelte-1166z7f{display:flex;justify-content:space-between;align-items:center;margin-bottom:2rem}.page-header.svelte-1166z7f h2.svelte-1166z7f{margin:0;color:#1e293b;font-size:1.875rem;font-weight:700}.form-row.svelte-1166z7f.svelte-1166z7f{display:grid;grid-template-columns:1fr 1fr;gap:1rem}.campaigns-list.svelte-1166z7f.svelte-1166z7f{display:flex;flex-direction:column;gap:1.5rem}.campaign-header.svelte-1166z7f.svelte-1166z7f{display:flex;justify-content:space-between;align-items:flex-start}.campaign-info.svelte-1166z7f.svelte-1166z7f{flex:1}.campaign-meta.svelte-1166z7f.svelte-1166z7f{margin:0.25rem 0 0 0;color:#64748b;font-size:0.875rem}.campaign-status.svelte-1166z7f.svelte-1166z7f{margin-left:1rem}.campaign-stats.svelte-1166z7f.svelte-1166z7f{display:flex;gap:2rem;padding:1rem;background:#f8fafc;border-radius:0.375rem}.stat.svelte-1166z7f.svelte-1166z7f{display:flex;align-items:center;gap:0.5rem;color:#64748b;font-size:0.875rem;font-weight:500}.campaign-actions.svelte-1166z7f.svelte-1166z7f{display:flex;gap:0.5rem;flex-wrap:wrap}.empty-state.svelte-1166z7f.svelte-1166z7f{text-align:center;padding:4rem 2rem;color:#64748b}.empty-state.svelte-1166z7f h3.svelte-1166z7f{margin:1rem 0 0.5rem 0;color:#374151;font-size:1.125rem;font-weight:600}.empty-state.svelte-1166z7f p.svelte-1166z7f{margin:0 0 2rem 0}@media(max-width: 768px){.page-header.svelte-1166z7f.svelte-1166z7f{flex-direction:column;align-items:flex-start;gap:1rem}.form-row.svelte-1166z7f.svelte-1166z7f{grid-template-columns:1fr}.campaign-header.svelte-1166z7f.svelte-1166z7f{flex-direction:column;align-items:flex-start;gap:1rem}.campaign-status.svelte-1166z7f.svelte-1166z7f{margin-left:0}.campaign-stats.svelte-1166z7f.svelte-1166z7f{flex-direction:column;gap:1rem}.campaign-actions.svelte-1166z7f.svelte-1166z7f{justify-content:flex-start}}",
  map: `{"version":3,"file":"+page.svelte","sources":["+page.svelte"],"sourcesContent":["<script>\\n\\timport { onMount } from 'svelte';\\n\\timport { Plus, Mail, Eye, MousePointer, Users, Calendar, Play, Pause, Edit, Trash2 } from 'lucide-svelte';\\n\\n\\t/** @type {Array<{id: number, subject: string, from_name: string, from_email: string, status: string, created_at: string, sent_at: string|null, scheduled_at: string|null}>} */\\n\\tlet campaigns = [];\\n\\tlet loading = false;\\n\\tlet showCreateCampaign = false;\\n\\n\\tonMount(async () => {\\n\\t\\tawait loadCampaigns();\\n\\t});\\n\\n\\tasync function loadCampaigns() {\\n\\t\\ttry {\\n\\t\\t\\tconst response = await fetch('/api/campaigns');\\n\\t\\t\\tif (response.ok) {\\n\\t\\t\\t\\tconst data = await response.json();\\n\\t\\t\\t\\tcampaigns = data.data || [];\\n\\t\\t\\t}\\n\\t\\t} catch (error) {\\n\\t\\t\\tconsole.error('Failed to load campaigns:', error);\\n\\t\\t}\\n\\t}\\n\\n\\t/** @param {Object} campaign */\\n\\tasync function createCampaign(campaign) {\\n\\t\\tloading = true;\\n\\t\\ttry {\\n\\t\\t\\tconst response = await fetch('/api/campaigns', {\\n\\t\\t\\t\\tmethod: 'POST',\\n\\t\\t\\t\\theaders: { 'Content-Type': 'application/json' },\\n\\t\\t\\t\\tbody: JSON.stringify(campaign)\\n\\t\\t\\t});\\n\\t\\t\\t\\n\\t\\t\\tif (response.ok) {\\n\\t\\t\\t\\tawait loadCampaigns();\\n\\t\\t\\t\\tshowCreateCampaign = false;\\n\\t\\t\\t} else {\\n\\t\\t\\t\\talert('Failed to create campaign');\\n\\t\\t\\t}\\n\\t\\t} catch (error) {\\n\\t\\t\\tconsole.error('Failed to create campaign:', error);\\n\\t\\t\\talert('Failed to create campaign');\\n\\t\\t} finally {\\n\\t\\t\\tloading = false;\\n\\t\\t}\\n\\t}\\n\\n\\t/** @param {number} campaignId */\\n\\tasync function testCampaign(campaignId) {\\n\\t\\tloading = true;\\n\\t\\ttry {\\n\\t\\t\\tconst response = await fetch(\`/api/campaigns/\${campaignId}/test\`, {\\n\\t\\t\\t\\tmethod: 'POST',\\n\\t\\t\\t\\theaders: { 'Content-Type': 'application/json' },\\n\\t\\t\\t\\tbody: JSON.stringify({ test_emails: ['test@example.com'] })\\n\\t\\t\\t});\\n\\t\\t\\t\\n\\t\\t\\tif (response.ok) {\\n\\t\\t\\t\\talert('Test email sent successfully');\\n\\t\\t\\t} else {\\n\\t\\t\\t\\talert('Failed to send test email');\\n\\t\\t\\t}\\n\\t\\t} catch (error) {\\n\\t\\t\\tconsole.error('Failed to send test email:', error);\\n\\t\\t\\talert('Failed to send test email');\\n\\t\\t} finally {\\n\\t\\t\\tloading = false;\\n\\t\\t}\\n\\t}\\n\\n\\t/** @param {number} campaignId @param {string} scheduledAt */\\n\\tasync function scheduleCampaign(campaignId, scheduledAt) {\\n\\t\\tloading = true;\\n\\t\\ttry {\\n\\t\\t\\tconst response = await fetch(\`/api/campaigns/\${campaignId}/schedule\`, {\\n\\t\\t\\t\\tmethod: 'POST',\\n\\t\\t\\t\\theaders: { 'Content-Type': 'application/json' },\\n\\t\\t\\t\\tbody: JSON.stringify({ scheduled_at: scheduledAt })\\n\\t\\t\\t});\\n\\t\\t\\t\\n\\t\\t\\tif (response.ok) {\\n\\t\\t\\t\\tawait loadCampaigns();\\n\\t\\t\\t\\talert('Campaign scheduled successfully');\\n\\t\\t\\t} else {\\n\\t\\t\\t\\talert('Failed to schedule campaign');\\n\\t\\t\\t}\\n\\t\\t} catch (error) {\\n\\t\\t\\tconsole.error('Failed to schedule campaign:', error);\\n\\t\\t\\talert('Failed to schedule campaign');\\n\\t\\t} finally {\\n\\t\\t\\tloading = false;\\n\\t\\t}\\n\\t}\\n\\n\\t/** @param {string} status */\\n\\tfunction getStatusColor(status) {\\n\\t\\tswitch (status) {\\n\\t\\t\\tcase 'draft':\\n\\t\\t\\t\\treturn 'gray';\\n\\t\\t\\tcase 'scheduled':\\n\\t\\t\\t\\treturn 'info';\\n\\t\\t\\tcase 'sending':\\n\\t\\t\\t\\treturn 'warning';\\n\\t\\t\\tcase 'sent':\\n\\t\\t\\t\\treturn 'success';\\n\\t\\t\\tcase 'failed':\\n\\t\\t\\t\\treturn 'error';\\n\\t\\t\\tdefault:\\n\\t\\t\\t\\treturn 'gray';\\n\\t\\t}\\n\\t}\\n\\n\\t/** @param {string} dateString */\\n\\tfunction formatDate(dateString) {\\n\\t\\treturn new Date(dateString).toLocaleDateString('en-US', {\\n\\t\\t\\tyear: 'numeric',\\n\\t\\t\\tmonth: 'short',\\n\\t\\t\\tday: 'numeric',\\n\\t\\t\\thour: '2-digit',\\n\\t\\t\\tminute: '2-digit'\\n\\t\\t});\\n\\t}\\n<\/script>\\n\\n<svelte:head>\\n\\t<title>Campaigns - Newsletter Platform</title>\\n</svelte:head>\\n\\n<div class=\\"container\\">\\n\\t<header class=\\"header\\">\\n\\t\\t<h1>Newsletter Platform</h1>\\n\\t\\t<nav class=\\"nav\\">\\n\\t\\t\\t<a href=\\"/\\" class=\\"nav-link\\">Dashboard</a>\\n\\t\\t\\t<a href=\\"/domains\\" class=\\"nav-link\\">Domains</a>\\n\\t\\t\\t<a href=\\"/lists\\" class=\\"nav-link\\">Lists</a>\\n\\t\\t\\t<a href=\\"/campaigns\\" class=\\"nav-link active\\">Campaigns</a>\\n\\t\\t\\t<a href=\\"/settings\\" class=\\"nav-link\\">Settings</a>\\n\\t\\t</nav>\\n\\t</header>\\n\\n\\t<main class=\\"main\\">\\n\\t\\t<div class=\\"page-header\\">\\n\\t\\t\\t<h2>Campaigns</h2>\\n\\t\\t\\t<button class=\\"btn btn-primary\\" on:click={() => showCreateCampaign = true}>\\n\\t\\t\\t\\t<Plus size={16} />\\n\\t\\t\\t\\tCreate Campaign\\n\\t\\t\\t</button>\\n\\t\\t</div>\\n\\n\\t\\t{#if showCreateCampaign}\\n\\t\\t\\t<div class=\\"card\\">\\n\\t\\t\\t\\t<h3>Create New Campaign</h3>\\n        <form on:submit|preventDefault={async (e) => {\\n          const form = e.target;\\n          if (!form) return;\\n          const formData = new FormData(/** @type {HTMLFormElement} */ (form));\\n          /** @param {string} name */\\n          const getValue = (name) => {\\n            const value = formData.get(name);\\n            return value && typeof value === 'string' ? value : '';\\n          };\\n          const campaign = {\\n            list_id: parseInt(getValue('list_id') || '0'),\\n            subject: getValue('subject'),\\n            from_name: getValue('from_name'),\\n            from_email: getValue('from_email'),\\n            reply_to: getValue('reply_to'),\\n            html: getValue('html'),\\n            text: getValue('text')\\n          };\\n\\t\\t\\t\\t\\tawait createCampaign(campaign);\\n\\t\\t\\t\\t}}>\\n\\t\\t\\t\\t\\t<div class=\\"form-row\\">\\n\\t\\t\\t\\t\\t\\t<div class=\\"form-group\\">\\n\\t\\t\\t\\t\\t\\t\\t<label for=\\"subject\\" class=\\"form-label\\">Subject</label>\\n\\t\\t\\t\\t\\t\\t\\t<input\\n\\t\\t\\t\\t\\t\\t\\t\\ttype=\\"text\\"\\n\\t\\t\\t\\t\\t\\t\\t\\tid=\\"subject\\"\\n\\t\\t\\t\\t\\t\\t\\t\\tname=\\"subject\\"\\n\\t\\t\\t\\t\\t\\t\\t\\tclass=\\"form-input\\"\\n\\t\\t\\t\\t\\t\\t\\t\\tplaceholder=\\"Your campaign subject\\"\\n\\t\\t\\t\\t\\t\\t\\t\\trequired\\n\\t\\t\\t\\t\\t\\t\\t/>\\n\\t\\t\\t\\t\\t\\t</div>\\n\\t\\t\\t\\t\\t\\t<div class=\\"form-group\\">\\n\\t\\t\\t\\t\\t\\t\\t<label for=\\"list_id\\" class=\\"form-label\\">List</label>\\n\\t\\t\\t\\t\\t\\t\\t<select id=\\"list_id\\" name=\\"list_id\\" class=\\"form-input form-select\\" required>\\n\\t\\t\\t\\t\\t\\t\\t\\t<option value=\\"1\\">Main List</option>\\n\\t\\t\\t\\t\\t\\t\\t</select>\\n\\t\\t\\t\\t\\t\\t</div>\\n\\t\\t\\t\\t\\t</div>\\n\\n\\t\\t\\t\\t\\t<div class=\\"form-row\\">\\n\\t\\t\\t\\t\\t\\t<div class=\\"form-group\\">\\n\\t\\t\\t\\t\\t\\t\\t<label for=\\"from_name\\" class=\\"form-label\\">From Name</label>\\n\\t\\t\\t\\t\\t\\t\\t<input\\n\\t\\t\\t\\t\\t\\t\\t\\ttype=\\"text\\"\\n\\t\\t\\t\\t\\t\\t\\t\\tid=\\"from_name\\"\\n\\t\\t\\t\\t\\t\\t\\t\\tname=\\"from_name\\"\\n\\t\\t\\t\\t\\t\\t\\t\\tclass=\\"form-input\\"\\n\\t\\t\\t\\t\\t\\t\\t\\tplaceholder=\\"Your Name\\"\\n\\t\\t\\t\\t\\t\\t\\t\\trequired\\n\\t\\t\\t\\t\\t\\t\\t/>\\n\\t\\t\\t\\t\\t\\t</div>\\n\\t\\t\\t\\t\\t\\t<div class=\\"form-group\\">\\n\\t\\t\\t\\t\\t\\t\\t<label for=\\"from_email\\" class=\\"form-label\\">From Email</label>\\n\\t\\t\\t\\t\\t\\t\\t<input\\n\\t\\t\\t\\t\\t\\t\\t\\ttype=\\"email\\"\\n\\t\\t\\t\\t\\t\\t\\t\\tid=\\"from_email\\"\\n\\t\\t\\t\\t\\t\\t\\t\\tname=\\"from_email\\"\\n\\t\\t\\t\\t\\t\\t\\t\\tclass=\\"form-input\\"\\n\\t\\t\\t\\t\\t\\t\\t\\tplaceholder=\\"news@example.com\\"\\n\\t\\t\\t\\t\\t\\t\\t\\trequired\\n\\t\\t\\t\\t\\t\\t\\t/>\\n\\t\\t\\t\\t\\t\\t</div>\\n\\t\\t\\t\\t\\t</div>\\n\\n\\t\\t\\t\\t\\t<div class=\\"form-group\\">\\n\\t\\t\\t\\t\\t\\t<label for=\\"reply_to\\" class=\\"form-label\\">Reply To (optional)</label>\\n\\t\\t\\t\\t\\t\\t<input\\n\\t\\t\\t\\t\\t\\t\\ttype=\\"email\\"\\n\\t\\t\\t\\t\\t\\t\\tid=\\"reply_to\\"\\n\\t\\t\\t\\t\\t\\t\\tname=\\"reply_to\\"\\n\\t\\t\\t\\t\\t\\t\\tclass=\\"form-input\\"\\n\\t\\t\\t\\t\\t\\t\\tplaceholder=\\"replies@example.com\\"\\n\\t\\t\\t\\t\\t\\t/>\\n\\t\\t\\t\\t\\t</div>\\n\\n\\t\\t\\t\\t\\t<div class=\\"form-group\\">\\n\\t\\t\\t\\t\\t\\t<label for=\\"html\\" class=\\"form-label\\">HTML Content</label>\\n\\t\\t\\t\\t\\t\\t<textarea\\n\\t\\t\\t\\t\\t\\t\\tid=\\"html\\"\\n\\t\\t\\t\\t\\t\\t\\tname=\\"html\\"\\n\\t\\t\\t\\t\\t\\t\\tclass=\\"form-input form-textarea\\"\\n\\t\\t\\t\\t\\t\\t\\tplaceholder=\\"<h1>Hello World</h1><p>Your email content here...</p>\\"\\n\\t\\t\\t\\t\\t\\t\\trows=\\"10\\"\\n\\t\\t\\t\\t\\t\\t\\trequired\\n\\t\\t\\t\\t\\t\\t></textarea>\\n\\t\\t\\t\\t\\t</div>\\n\\n\\t\\t\\t\\t\\t<div class=\\"form-group\\">\\n\\t\\t\\t\\t\\t\\t<label for=\\"text\\" class=\\"form-label\\">Text Content</label>\\n\\t\\t\\t\\t\\t\\t<textarea\\n\\t\\t\\t\\t\\t\\t\\tid=\\"text\\"\\n\\t\\t\\t\\t\\t\\t\\tname=\\"text\\"\\n\\t\\t\\t\\t\\t\\t\\tclass=\\"form-input form-textarea\\"\\n\\t\\t\\t\\t\\t\\t\\tplaceholder=\\"Hello World\\\\n\\\\nYour email content here...\\"\\n\\t\\t\\t\\t\\t\\t\\trows=\\"5\\"\\n\\t\\t\\t\\t\\t\\t\\trequired\\n\\t\\t\\t\\t\\t\\t></textarea>\\n\\t\\t\\t\\t\\t</div>\\n\\n\\t\\t\\t\\t\\t<div class=\\"card-footer\\">\\n\\t\\t\\t\\t\\t\\t<button type=\\"button\\" class=\\"btn btn-secondary\\" on:click={() => showCreateCampaign = false}>\\n\\t\\t\\t\\t\\t\\t\\tCancel\\n\\t\\t\\t\\t\\t\\t</button>\\n\\t\\t\\t\\t\\t\\t<button type=\\"submit\\" class=\\"btn btn-primary\\" disabled={loading}>\\n\\t\\t\\t\\t\\t\\t\\t{loading ? 'Creating...' : 'Create Campaign'}\\n\\t\\t\\t\\t\\t\\t</button>\\n\\t\\t\\t\\t\\t</div>\\n\\t\\t\\t\\t</form>\\n\\t\\t\\t</div>\\n\\t\\t{/if}\\n\\n\\t\\t<div class=\\"campaigns-list\\">\\n\\t\\t\\t{#each campaigns as campaign (campaign.id)}\\n\\t\\t\\t\\t<div class=\\"card\\">\\n\\t\\t\\t\\t\\t<div class=\\"card-header\\">\\n\\t\\t\\t\\t\\t\\t<div class=\\"campaign-header\\">\\n\\t\\t\\t\\t\\t\\t\\t<div class=\\"campaign-info\\">\\n\\t\\t\\t\\t\\t\\t\\t\\t<h3 class=\\"card-title\\">{campaign.subject}</h3>\\n\\t\\t\\t\\t\\t\\t\\t\\t<p class=\\"campaign-meta\\">\\n\\t\\t\\t\\t\\t\\t\\t\\t\\tFrom: {campaign.from_name} &lt;{campaign.from_email}&gt;\\n\\t\\t\\t\\t\\t\\t\\t\\t</p>\\n\\t\\t\\t\\t\\t\\t\\t\\t<p class=\\"campaign-meta\\">\\n\\t\\t\\t\\t\\t\\t\\t\\t\\tCreated: {formatDate(campaign.created_at)}\\n\\t\\t\\t\\t\\t\\t\\t\\t\\t{#if campaign.sent_at}\\n\\t\\t\\t\\t\\t\\t\\t\\t\\t\\t• Sent: {formatDate(campaign.sent_at)}\\n\\t\\t\\t\\t\\t\\t\\t\\t\\t{/if}\\n\\t\\t\\t\\t\\t\\t\\t\\t</p>\\n\\t\\t\\t\\t\\t\\t\\t</div>\\n\\t\\t\\t\\t\\t\\t\\t<div class=\\"campaign-status\\">\\n\\t\\t\\t\\t\\t\\t\\t\\t<span class=\\"badge badge-{getStatusColor(campaign.status)}\\">\\n\\t\\t\\t\\t\\t\\t\\t\\t\\t{campaign.status}\\n\\t\\t\\t\\t\\t\\t\\t\\t</span>\\n\\t\\t\\t\\t\\t\\t\\t</div>\\n\\t\\t\\t\\t\\t\\t</div>\\n\\t\\t\\t\\t\\t</div>\\n\\n\\t\\t\\t\\t\\t<div class=\\"card-body\\">\\n\\t\\t\\t\\t\\t\\t{#if campaign.status === 'sent'}\\n\\t\\t\\t\\t\\t\\t\\t<div class=\\"campaign-stats\\">\\n\\t\\t\\t\\t\\t\\t\\t\\t<div class=\\"stat\\">\\n\\t\\t\\t\\t\\t\\t\\t\\t\\t<Users size={16} />\\n\\t\\t\\t\\t\\t\\t\\t\\t\\t<span>1,250 sent</span>\\n\\t\\t\\t\\t\\t\\t\\t\\t</div>\\n\\t\\t\\t\\t\\t\\t\\t\\t<div class=\\"stat\\">\\n\\t\\t\\t\\t\\t\\t\\t\\t\\t<Eye size={16} />\\n\\t\\t\\t\\t\\t\\t\\t\\t\\t<span>24.5% open</span>\\n\\t\\t\\t\\t\\t\\t\\t\\t</div>\\n\\t\\t\\t\\t\\t\\t\\t\\t<div class=\\"stat\\">\\n\\t\\t\\t\\t\\t\\t\\t\\t\\t<MousePointer size={16} />\\n\\t\\t\\t\\t\\t\\t\\t\\t\\t<span>3.2% click</span>\\n\\t\\t\\t\\t\\t\\t\\t\\t</div>\\n\\t\\t\\t\\t\\t\\t\\t</div>\\n\\t\\t\\t\\t\\t\\t{/if}\\n\\t\\t\\t\\t\\t</div>\\n\\n\\t\\t\\t\\t\\t<div class=\\"card-footer\\">\\n\\t\\t\\t\\t\\t\\t<div class=\\"campaign-actions\\">\\n\\t\\t\\t\\t\\t\\t\\t{#if campaign.status === 'draft'}\\n\\t\\t\\t\\t\\t\\t\\t\\t<button \\n\\t\\t\\t\\t\\t\\t\\t\\t\\tclass=\\"btn btn-secondary btn-sm\\" \\n\\t\\t\\t\\t\\t\\t\\t\\t\\ton:click={() => testCampaign(campaign.id)}\\n\\t\\t\\t\\t\\t\\t\\t\\t\\tdisabled={loading}\\n\\t\\t\\t\\t\\t\\t\\t\\t>\\n\\t\\t\\t\\t\\t\\t\\t\\t\\t<Mail size={14} />\\n\\t\\t\\t\\t\\t\\t\\t\\t\\tTest Send\\n\\t\\t\\t\\t\\t\\t\\t\\t</button>\\n\\t\\t\\t\\t\\t\\t\\t\\t<button \\n\\t\\t\\t\\t\\t\\t\\t\\t\\tclass=\\"btn btn-primary btn-sm\\" \\n\\t\\t\\t\\t\\t\\t\\t\\t\\ton:click={() => scheduleCampaign(campaign.id, new Date().toISOString())}\\n\\t\\t\\t\\t\\t\\t\\t\\t\\tdisabled={loading}\\n\\t\\t\\t\\t\\t\\t\\t\\t>\\n\\t\\t\\t\\t\\t\\t\\t\\t\\t<Play size={14} />\\n\\t\\t\\t\\t\\t\\t\\t\\t\\tSend Now\\n\\t\\t\\t\\t\\t\\t\\t\\t</button>\\n\\t\\t\\t\\t\\t\\t\\t{/if}\\n\\n\\t\\t\\t\\t\\t\\t\\t{#if campaign.status === 'scheduled'}\\n\\t\\t\\t\\t\\t\\t\\t\\t<button \\n\\t\\t\\t\\t\\t\\t\\t\\t\\tclass=\\"btn btn-secondary btn-sm\\" \\n\\t\\t\\t\\t\\t\\t\\t\\t\\ton:click={() => scheduleCampaign(campaign.id, '')}\\n\\t\\t\\t\\t\\t\\t\\t\\t\\tdisabled={loading}\\n\\t\\t\\t\\t\\t\\t\\t\\t>\\n\\t\\t\\t\\t\\t\\t\\t\\t\\t<Pause size={14} />\\n\\t\\t\\t\\t\\t\\t\\t\\t\\tCancel\\n\\t\\t\\t\\t\\t\\t\\t\\t</button>\\n\\t\\t\\t\\t\\t\\t\\t{/if}\\n\\n\\t\\t\\t\\t\\t\\t\\t{#if campaign.status === 'sent'}\\n\\t\\t\\t\\t\\t\\t\\t\\t<button class=\\"btn btn-secondary btn-sm\\">\\n\\t\\t\\t\\t\\t\\t\\t\\t\\t<Eye size={14} />\\n\\t\\t\\t\\t\\t\\t\\t\\t\\tView Report\\n\\t\\t\\t\\t\\t\\t\\t\\t</button>\\n\\t\\t\\t\\t\\t\\t\\t{/if}\\n\\n\\t\\t\\t\\t\\t\\t\\t<button class=\\"btn btn-secondary btn-sm\\">\\n\\t\\t\\t\\t\\t\\t\\t\\t<Edit size={14} />\\n\\t\\t\\t\\t\\t\\t\\t\\tEdit\\n\\t\\t\\t\\t\\t\\t\\t</button>\\n\\t\\t\\t\\t\\t\\t\\t<button class=\\"btn btn-danger btn-sm\\">\\n\\t\\t\\t\\t\\t\\t\\t\\t<Trash2 size={14} />\\n\\t\\t\\t\\t\\t\\t\\t\\tDelete\\n\\t\\t\\t\\t\\t\\t\\t</button>\\n\\t\\t\\t\\t\\t\\t</div>\\n\\t\\t\\t\\t\\t</div>\\n\\t\\t\\t\\t</div>\\n\\t\\t\\t{/each}\\n\\n\\t\\t\\t{#if campaigns.length === 0}\\n\\t\\t\\t\\t<div class=\\"empty-state\\">\\n\\t\\t\\t\\t\\t<Mail size={48} />\\n\\t\\t\\t\\t\\t<h3>No campaigns yet</h3>\\n\\t\\t\\t\\t\\t<p>Create your first campaign to start sending newsletters.</p>\\n\\t\\t\\t\\t\\t<button class=\\"btn btn-primary\\" on:click={() => showCreateCampaign = true}>\\n\\t\\t\\t\\t\\t\\t<Plus size={16} />\\n\\t\\t\\t\\t\\t\\tCreate Campaign\\n\\t\\t\\t\\t\\t</button>\\n\\t\\t\\t\\t</div>\\n\\t\\t\\t{/if}\\n\\t\\t</div>\\n\\t</main>\\n</div>\\n\\n<style>\\n\\t.page-header {\\n\\t\\tdisplay: flex;\\n\\t\\tjustify-content: space-between;\\n\\t\\talign-items: center;\\n\\t\\tmargin-bottom: 2rem;\\n\\t}\\n\\n\\t.page-header h2 {\\n\\t\\tmargin: 0;\\n\\t\\tcolor: #1e293b;\\n\\t\\tfont-size: 1.875rem;\\n\\t\\tfont-weight: 700;\\n\\t}\\n\\n\\t.form-row {\\n\\t\\tdisplay: grid;\\n\\t\\tgrid-template-columns: 1fr 1fr;\\n\\t\\tgap: 1rem;\\n\\t}\\n\\n\\t.campaigns-list {\\n\\t\\tdisplay: flex;\\n\\t\\tflex-direction: column;\\n\\t\\tgap: 1.5rem;\\n\\t}\\n\\n\\t.campaign-header {\\n\\t\\tdisplay: flex;\\n\\t\\tjustify-content: space-between;\\n\\t\\talign-items: flex-start;\\n\\t}\\n\\n\\t.campaign-info {\\n\\t\\tflex: 1;\\n\\t}\\n\\n\\t.campaign-meta {\\n\\t\\tmargin: 0.25rem 0 0 0;\\n\\t\\tcolor: #64748b;\\n\\t\\tfont-size: 0.875rem;\\n\\t}\\n\\n\\t.campaign-status {\\n\\t\\tmargin-left: 1rem;\\n\\t}\\n\\n\\t.campaign-stats {\\n\\t\\tdisplay: flex;\\n\\t\\tgap: 2rem;\\n\\t\\tpadding: 1rem;\\n\\t\\tbackground: #f8fafc;\\n\\t\\tborder-radius: 0.375rem;\\n\\t}\\n\\n\\t.stat {\\n\\t\\tdisplay: flex;\\n\\t\\talign-items: center;\\n\\t\\tgap: 0.5rem;\\n\\t\\tcolor: #64748b;\\n\\t\\tfont-size: 0.875rem;\\n\\t\\tfont-weight: 500;\\n\\t}\\n\\n\\t.campaign-actions {\\n\\t\\tdisplay: flex;\\n\\t\\tgap: 0.5rem;\\n\\t\\tflex-wrap: wrap;\\n\\t}\\n\\n\\t.empty-state {\\n\\t\\ttext-align: center;\\n\\t\\tpadding: 4rem 2rem;\\n\\t\\tcolor: #64748b;\\n\\t}\\n\\n\\t.empty-state h3 {\\n\\t\\tmargin: 1rem 0 0.5rem 0;\\n\\t\\tcolor: #374151;\\n\\t\\tfont-size: 1.125rem;\\n\\t\\tfont-weight: 600;\\n\\t}\\n\\n\\t.empty-state p {\\n\\t\\tmargin: 0 0 2rem 0;\\n\\t}\\n\\n\\t@media (max-width: 768px) {\\n\\t\\t.page-header {\\n\\t\\t\\tflex-direction: column;\\n\\t\\t\\talign-items: flex-start;\\n\\t\\t\\tgap: 1rem;\\n\\t\\t}\\n\\n\\t\\t.form-row {\\n\\t\\t\\tgrid-template-columns: 1fr;\\n\\t\\t}\\n\\n\\t\\t.campaign-header {\\n\\t\\t\\tflex-direction: column;\\n\\t\\t\\talign-items: flex-start;\\n\\t\\t\\tgap: 1rem;\\n\\t\\t}\\n\\n\\t\\t.campaign-status {\\n\\t\\t\\tmargin-left: 0;\\n\\t\\t}\\n\\n\\t\\t.campaign-stats {\\n\\t\\t\\tflex-direction: column;\\n\\t\\t\\tgap: 1rem;\\n\\t\\t}\\n\\n\\t\\t.campaign-actions {\\n\\t\\t\\tjustify-content: flex-start;\\n\\t\\t}\\n\\t}\\n</style>\\n"],"names":[],"mappings":"AA0XC,0CAAa,CACZ,OAAO,CAAE,IAAI,CACb,eAAe,CAAE,aAAa,CAC9B,WAAW,CAAE,MAAM,CACnB,aAAa,CAAE,IAChB,CAEA,2BAAY,CAAC,iBAAG,CACf,MAAM,CAAE,CAAC,CACT,KAAK,CAAE,OAAO,CACd,SAAS,CAAE,QAAQ,CACnB,WAAW,CAAE,GACd,CAEA,uCAAU,CACT,OAAO,CAAE,IAAI,CACb,qBAAqB,CAAE,GAAG,CAAC,GAAG,CAC9B,GAAG,CAAE,IACN,CAEA,6CAAgB,CACf,OAAO,CAAE,IAAI,CACb,cAAc,CAAE,MAAM,CACtB,GAAG,CAAE,MACN,CAEA,8CAAiB,CAChB,OAAO,CAAE,IAAI,CACb,eAAe,CAAE,aAAa,CAC9B,WAAW,CAAE,UACd,CAEA,4CAAe,CACd,IAAI,CAAE,CACP,CAEA,4CAAe,CACd,MAAM,CAAE,OAAO,CAAC,CAAC,CAAC,CAAC,CAAC,CAAC,CACrB,KAAK,CAAE,OAAO,CACd,SAAS,CAAE,QACZ,CAEA,8CAAiB,CAChB,WAAW,CAAE,IACd,CAEA,6CAAgB,CACf,OAAO,CAAE,IAAI,CACb,GAAG,CAAE,IAAI,CACT,OAAO,CAAE,IAAI,CACb,UAAU,CAAE,OAAO,CACnB,aAAa,CAAE,QAChB,CAEA,mCAAM,CACL,OAAO,CAAE,IAAI,CACb,WAAW,CAAE,MAAM,CACnB,GAAG,CAAE,MAAM,CACX,KAAK,CAAE,OAAO,CACd,SAAS,CAAE,QAAQ,CACnB,WAAW,CAAE,GACd,CAEA,+CAAkB,CACjB,OAAO,CAAE,IAAI,CACb,GAAG,CAAE,MAAM,CACX,SAAS,CAAE,IACZ,CAEA,0CAAa,CACZ,UAAU,CAAE,MAAM,CAClB,OAAO,CAAE,IAAI,CAAC,IAAI,CAClB,KAAK,CAAE,OACR,CAEA,2BAAY,CAAC,iBAAG,CACf,MAAM,CAAE,IAAI,CAAC,CAAC,CAAC,MAAM,CAAC,CAAC,CACvB,KAAK,CAAE,OAAO,CACd,SAAS,CAAE,QAAQ,CACnB,WAAW,CAAE,GACd,CAEA,2BAAY,CAAC,gBAAE,CACd,MAAM,CAAE,CAAC,CAAC,CAAC,CAAC,IAAI,CAAC,CAClB,CAEA,MAAO,YAAY,KAAK,CAAE,CACzB,0CAAa,CACZ,cAAc,CAAE,MAAM,CACtB,WAAW,CAAE,UAAU,CACvB,GAAG,CAAE,IACN,CAEA,uCAAU,CACT,qBAAqB,CAAE,GACxB,CAEA,8CAAiB,CAChB,cAAc,CAAE,MAAM,CACtB,WAAW,CAAE,UAAU,CACvB,GAAG,CAAE,IACN,CAEA,8CAAiB,CAChB,WAAW,CAAE,CACd,CAEA,6CAAgB,CACf,cAAc,CAAE,MAAM,CACtB,GAAG,CAAE,IACN,CAEA,+CAAkB,CACjB,eAAe,CAAE,UAClB,CACD"}`
};
function getStatusColor(status) {
  switch (status) {
    case "draft":
      return "gray";
    case "scheduled":
      return "info";
    case "sending":
      return "warning";
    case "sent":
      return "success";
    case "failed":
      return "error";
    default:
      return "gray";
  }
}
function formatDate(dateString) {
  return new Date(dateString).toLocaleDateString("en-US", {
    year: "numeric",
    month: "short",
    day: "numeric",
    hour: "2-digit",
    minute: "2-digit"
  });
}
const Page = create_ssr_component(($$result, $$props, $$bindings, slots) => {
  let campaigns = [];
  $$result.css.add(css);
  return `${$$result.head += `<!-- HEAD_svelte-1yoe9jm_START -->${$$result.title = `<title>Campaigns - Newsletter Platform</title>`, ""}<!-- HEAD_svelte-1yoe9jm_END -->`, ""} <div class="container"><header class="header" data-svelte-h="svelte-1evoi29"><h1>Newsletter Platform</h1> <nav class="nav"><a href="/" class="nav-link">Dashboard</a> <a href="/domains" class="nav-link">Domains</a> <a href="/lists" class="nav-link">Lists</a> <a href="/campaigns" class="nav-link active">Campaigns</a> <a href="/settings" class="nav-link">Settings</a></nav></header> <main class="main"><div class="page-header svelte-1166z7f"><h2 class="svelte-1166z7f" data-svelte-h="svelte-d7ntm7">Campaigns</h2> <button class="btn btn-primary">${validate_component(Plus, "Plus").$$render($$result, { size: 16 }, {}, {})}
				Create Campaign</button></div> ${``} <div class="campaigns-list svelte-1166z7f">${each(campaigns, (campaign) => {
    return `<div class="card"><div class="card-header"><div class="campaign-header svelte-1166z7f"><div class="campaign-info svelte-1166z7f"><h3 class="card-title">${escape(campaign.subject)}</h3> <p class="campaign-meta svelte-1166z7f">From: ${escape(campaign.from_name)} &lt;${escape(campaign.from_email)}&gt;</p> <p class="campaign-meta svelte-1166z7f">Created: ${escape(formatDate(campaign.created_at))} ${campaign.sent_at ? `• Sent: ${escape(formatDate(campaign.sent_at))}` : ``} </p></div> <div class="campaign-status svelte-1166z7f"><span class="${"badge badge-" + escape(getStatusColor(campaign.status), true) + " svelte-1166z7f"}">${escape(campaign.status)} </span></div> </div></div> <div class="card-body">${campaign.status === "sent" ? `<div class="campaign-stats svelte-1166z7f"><div class="stat svelte-1166z7f">${validate_component(Users, "Users").$$render($$result, { size: 16 }, {}, {})} <span data-svelte-h="svelte-1hro0pg">1,250 sent</span></div> <div class="stat svelte-1166z7f">${validate_component(Eye, "Eye").$$render($$result, { size: 16 }, {}, {})} <span data-svelte-h="svelte-5zhha0">24.5% open</span></div> <div class="stat svelte-1166z7f">${validate_component(Mouse_pointer, "MousePointer").$$render($$result, { size: 16 }, {}, {})} <span data-svelte-h="svelte-1qitaxu">3.2% click</span></div> </div>` : ``}</div> <div class="card-footer"><div class="campaign-actions svelte-1166z7f">${campaign.status === "draft" ? `<button class="btn btn-secondary btn-sm" ${""}>${validate_component(Mail, "Mail").$$render($$result, { size: 14 }, {}, {})}
									Test Send</button> <button class="btn btn-primary btn-sm" ${""}>${validate_component(Play, "Play").$$render($$result, { size: 14 }, {}, {})}
									Send Now
								</button>` : ``} ${campaign.status === "scheduled" ? `<button class="btn btn-secondary btn-sm" ${""}>${validate_component(Pause, "Pause").$$render($$result, { size: 14 }, {}, {})}
									Cancel
								</button>` : ``} ${campaign.status === "sent" ? `<button class="btn btn-secondary btn-sm">${validate_component(Eye, "Eye").$$render($$result, { size: 14 }, {}, {})}
									View Report
								</button>` : ``} <button class="btn btn-secondary btn-sm">${validate_component(Pen_square, "Edit").$$render($$result, { size: 14 }, {}, {})}
								Edit</button> <button class="btn btn-danger btn-sm">${validate_component(Trash_2, "Trash2").$$render($$result, { size: 14 }, {}, {})}
								Delete</button> </div></div> </div>`;
  })} ${campaigns.length === 0 ? `<div class="empty-state svelte-1166z7f">${validate_component(Mail, "Mail").$$render($$result, { size: 48 }, {}, {})} <h3 class="svelte-1166z7f" data-svelte-h="svelte-1092a3e">No campaigns yet</h3> <p class="svelte-1166z7f" data-svelte-h="svelte-n0y20a">Create your first campaign to start sending newsletters.</p> <button class="btn btn-primary">${validate_component(Plus, "Plus").$$render($$result, { size: 16 }, {}, {})}
						Create Campaign</button></div>` : ``}</div></main> </div>`;
});
export {
  Page as default
};
