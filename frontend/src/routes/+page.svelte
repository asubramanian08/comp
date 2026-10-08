<script lang="ts">
	import { onMount } from 'svelte';

	// The Go server. Override with VITE_API_URL (Vite reads it from the environment at build time).
	const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/graphql';

	let count = $state(0);
	let error = $state('');

	// Sends one GraphQL operation (e.g. "mutation { increment }") to the Go server and shows
	// the number that comes back in the reply's `field`.
	async function send(query: string, field: string) {
		try {
			const res = await fetch(API_URL, {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ query })
			});
			const json = await res.json();
			// GraphQL can reply HTTP 200 with an `errors` list, so check both.
			if (!res.ok || json.errors) throw new Error(json.errors?.[0]?.message ?? `HTTP ${res.status}`);
			// Replies can arrive out of order; the counter only goes up, so keep the larger number.
			count = Math.max(count, json.data[field]);
			error = '';
		} catch (e) {
			error = `Request failed: ${(e as Error).message}`;
		}
	}

	onMount(() => send('query { counter }', 'counter'));
</script>

<button onclick={() => send('mutation { increment }', 'increment')}>Increment</button>
<p>Count: {count}</p>
{#if error}<p style="color: red">{error}</p>{/if}
