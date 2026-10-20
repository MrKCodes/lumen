import { supabase } from '$lib/supabase.js';

/**
 * Global auth/session state using Svelte 5 runes.
 * `session` and `profile` are reactive; import `auth` anywhere.
 */
class AuthState {
	session = $state(null);
	profile = $state(null); // { id, username, display_name, role }
	loading = $state(true);

	get user() {
		return this.session?.user ?? null;
	}
	get isAuthed() {
		return !!this.session;
	}
	get isAdmin() {
		return this.profile?.role === 'admin';
	}

	/** Call once on app start. Loads session + profile and subscribes to changes. */
	async init() {
		const { data } = await supabase.auth.getSession();
		this.session = data.session;
		await this.loadProfile();
		this.loading = false;

		supabase.auth.onAuthStateChange(async (_event, session) => {
			this.session = session;
			await this.loadProfile();
		});
	}

	async loadProfile() {
		if (!this.session) {
			this.profile = null;
			return;
		}
		const { data } = await supabase
			.from('profiles')
			.select('id, username, display_name, role')
			.eq('id', this.session.user.id)
			.single();
		this.profile = data ?? null;
	}

	/** Register with email + password. No email confirmation (configured in Supabase). */
	async register({ email, password, username, displayName }) {
		const { error } = await supabase.auth.signUp({
			email,
			password,
			options: { data: { username, display_name: displayName || username } }
		});
		if (error) throw error;
		// With confirmation off, signUp returns an active session immediately.
		await this.loadProfile();
	}

	async login({ email, password }) {
		const { error } = await supabase.auth.signInWithPassword({ email, password });
		if (error) throw error;
		await this.loadProfile();
	}

	async logout() {
		await supabase.auth.signOut();
		this.session = null;
		this.profile = null;
	}
}

export const auth = new AuthState();
