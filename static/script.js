const linksTextarea = document.getElementById("links-textarea");
const divResults = document.querySelector(".results");

const checkLinksStatus = async (links) => {
	const resp = await fetch("/api/check-links", {
		method: "POST",
		headers: {
			"Content-Type": "application/json",
		},
		body: JSON.stringify({ urls: links }),
	});
	return await resp.json();
};

const checkLinks = async (btn) => {
	const links = linksTextarea.value.split("\n").filter((x) => x.trim().length > 0);
	if (links.length == 0) return;

	const btnText = btn.textContent;
	btn.disabled = true;
	btn.textContent = "Checando...";
	const data = await checkLinksStatus(links);
	btn.disabled = false;
	btn.textContent = btnText;

	const { urls, tookMs, ok } = data;
	if (!ok) return alert(`Houve um erro ao checar links. Erro: ${data.error}`);

	divResults.querySelector(".took-ms").textContent = `Levou ${Number(tookMs).toLocaleString("pt-BR")}ms`;

	divResults.querySelector(".data").innerHTML = urls
		.map((u) => {
			const { url, status, statusCode, isAlive } = u;

			return `<div class="url${isAlive ? "" : " not-ok"}">
			<a href="${url}" target="_blank" rel="noopener noreferrer">${url}</a>
			<span class="status-dot code-${String(statusCode).charAt(0)}"></span>
			<span class="status">${status}</span>
		</div>`;
		})
		.join("");
};
