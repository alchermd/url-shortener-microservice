document.querySelector('form').addEventListener('submit', function(e) {
	e.preventDefault();

	const endpoint = this.getAttribute('action');
	const url = document.querySelector(`input[type="url"]`).value;
	
	fetch(endpoint, {
		method: "POST",
		headers: {
			"Content-Type": "application/json",
		},
		body: JSON.stringify({url}),
	})
	.then(res => res.json())
	.then(data => {
		if (data.error) {
			alert("Invalid URL");
			return;
		}

		document.querySelector("#short-url-container").style.display = "block";

		const shortUrl = `/api/shorturl/${data.short_url}`;
		const link = document.querySelector("#short-url");
		link.innerHTML = shortUrl;
		link.setAttribute('href', shortUrl)
	})
});