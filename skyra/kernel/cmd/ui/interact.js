const conversationEl = document.getElementById("conversation");
const composerEl = document.getElementById("composer");
const messageEl = document.getElementById("message");

async function bootInteract() {
  const response = await fetch("/api/v1/interact");
  const snapshot = await response.json();
  renderMessages(snapshot.messages || []);
  connectInteractionStream();
}

function renderMessages(messages) {
  conversationEl.innerHTML = "";
  for (const message of messages) {
    appendMessage(message);
  }
  conversationEl.scrollTop = conversationEl.scrollHeight;
}

function appendMessage(message) {
  const article = document.createElement("article");
  article.className = `message ${message.role === "user" ? "message-user" : "message-assistant"}`;
  article.textContent = message.content || "";
  article.dataset.messageId = message.id || "";
  conversationEl.appendChild(article);
  conversationEl.scrollTop = conversationEl.scrollHeight;
}

function connectInteractionStream() {
  const stream = new EventSource("/api/v1/interact/stream");
  stream.addEventListener("interaction.message", (event) => {
    appendMessage(JSON.parse(event.data));
  });
}

composerEl.addEventListener("submit", async (event) => {
  event.preventDefault();
  const content = messageEl.value.trim();
  if (!content) return;

  messageEl.value = "";

    await fetch("/api/v1/stimuli", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        content,
        source: "human",
        type: "text",
      }),
    });
});

bootInteract().catch((error) => {
  conversationEl.innerHTML = `<article class="message message-assistant">Failed to load interaction surface: ${error.message}</article>`;
});
