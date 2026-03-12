function renderFounders() {
  const container = document.getElementById("founder-grid");
  const list = Array.isArray(window.founders) ? window.founders : [];

  if (!container) return;

  if (list.length === 0) {
    container.innerHTML = "<p>No founder profiles yet.</p>";
    return;
  }

  container.innerHTML = list
    .map(
      (founder) => `
        <article class="founder-card">
          <h3 class="founder-name">${escapeHtml(founder.name || "Founder")}</h3>
          <p class="founder-role">${escapeHtml(founder.role || "")}</p>
          <p class="founder-bio">${escapeHtml(founder.bio || "")}</p>
        </article>
      `
    )
    .join("");
}

function escapeHtml(value) {
  return String(value)
    .replaceAll("&", "&amp;")
    .replaceAll("<", "&lt;")
    .replaceAll(">", "&gt;")
    .replaceAll('"', "&quot;")
    .replaceAll("'", "&#39;");
}

renderFounders();
