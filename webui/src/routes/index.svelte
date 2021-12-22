<script>
	import { onMount } from 'svelte';

	export let data = [];
	let pageId = 1;
	let pageSize = 10;

	async function fetchData() {
		const response = await fetch(
			`http://localhost:5000/movies?page_id=${pageId}&page_size=${pageSize}`
		);
		data = await response.json();
		console.log(data);
	}

	onMount(() => {
		fetchData();
	});
</script>

<div>
	{#each data as movie}
		<h4>{movie.film}</h4>
	{/each}
</div>
