/** Format minutes as "2h 15m" / "45m". */
export function fmtMinutes(min) {
	min = Math.round(min || 0);
	const h = Math.floor(min / 60);
	const m = min % 60;
	if (h && m) return `${h}h ${m}m`;
	if (h) return `${h}h`;
	return `${m}m`;
}

export function fmtHours(min) {
	return (Math.round(((min || 0) / 60) * 10) / 10).toString();
}

/** ISO date (YYYY-MM-DD) for today, local time. */
export function todayISO() {
	const d = new Date();
	const off = d.getTimezoneOffset();
	return new Date(d.getTime() - off * 60000).toISOString().slice(0, 10);
}

/** Short human date: "Oct 15". */
export function fmtDate(iso) {
	const d = new Date(iso + 'T00:00:00');
	return d.toLocaleDateString(undefined, { month: 'short', day: 'numeric' });
}

/** Group entries into a map keyed by YYYY-MM-DD → total minutes. */
export function minutesByDay(entries) {
	const map = {};
	for (const e of entries) map[e.entry_date] = (map[e.entry_date] || 0) + (e.minutes || 0);
	return map;
}

/**
 * Build a GitHub-style contribution grid for the last `weeks` weeks.
 * Returns array of {date, minutes, level 0-4}.
 */
export function heatmap(entries, weeks = 18) {
	const byDay = minutesByDay(entries);
	const days = weeks * 7;
	const cells = [];
	const today = new Date();
	for (let i = days - 1; i >= 0; i--) {
		const d = new Date(today);
		d.setDate(today.getDate() - i);
		const iso = d.toISOString().slice(0, 10);
		const min = byDay[iso] || 0;
		const level = min === 0 ? 0 : min < 30 ? 1 : min < 60 ? 2 : min < 120 ? 3 : 4;
		cells.push({ date: iso, minutes: min, level });
	}
	return cells;
}

/** Current consecutive-day streak ending today (days with >0 minutes). */
export function currentStreak(entries) {
	const byDay = minutesByDay(entries);
	let streak = 0;
	const d = new Date();
	for (;;) {
		const iso = d.toISOString().slice(0, 10);
		if (byDay[iso] > 0) {
			streak++;
			d.setDate(d.getDate() - 1);
		} else break;
	}
	return streak;
}

export const PALETTE = ['#c6ff4a', '#6c5cff', '#ff4d8d', '#35e6ff', '#3ddc97', '#ffcf5c', '#ff8f5c'];
