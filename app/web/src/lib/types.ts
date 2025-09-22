// Type definitions for Newsletter Platform

export interface Campaign {
	id: number;
	subject: string;
	from_name: string;
	from_email: string;
	status: 'draft' | 'scheduled' | 'sending' | 'sent' | 'failed';
	created_at: string;
	sent_at: string | null;
	scheduled_at: string | null;
}

export interface Domain {
	id: number;
	domain: string;
	dkim_selector: string;
	verified_at: string | null;
	spf_record: string;
	dmarc_record: string;
	ptr_record: string;
}

export interface DomainStatus {
	domain: string;
	spf: { status: string; message: string };
	dkim: { status: string; message: string };
	dmarc: { status: string; message: string };
	ptr: { status: string; message: string };
	overall: string;
}

export interface Subscriber {
	id: number;
	email: string;
	status: 'active' | 'bounced' | 'complained' | 'unsubscribed';
	attributes: Record<string, any>;
	created_at: string;
	unsubscribed_at: string | null;
}

export interface List {
	id: number;
	name: string;
	description: string;
	created_at: string;
}

export interface Event {
	id: number;
	campaign_id: number;
	subscriber_id: number;
	type: string;
	meta: Record<string, any>;
	at: string;
}
