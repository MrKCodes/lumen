import { createClient } from '@supabase/supabase-js';

const url = import.meta.env.VITE_SUPABASE_URL;
const anon = import.meta.env.VITE_SUPABASE_ANON_KEY;

if (!url || !anon) {
	// Fail loud in dev; the app is useless without these.
	console.error('Missing VITE_SUPABASE_URL / VITE_SUPABASE_ANON_KEY — copy .env.example to .env');
}

/**
 * Shared Supabase client. The anon key is public by design; Row-Level Security
 * in Postgres is what actually protects data.
 */
export const supabase = createClient(url ?? '', anon ?? '', {
	auth: {
		persistSession: true,
		autoRefreshToken: true,
		detectSessionInUrl: false
	}
});
