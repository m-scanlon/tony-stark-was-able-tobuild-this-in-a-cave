const activeStepEl = document.getElementById("active-step");
const chainStateEl = document.getElementById("chain-state");
const perceptionEl = document.getElementById("perception");
const timelineEl = document.getElementById("timeline");

const thoughtState = {
  activeStep: null,
  activeChain: null,
  perception: null,
  suspended: [],
  steps: [],
};

async function bootThoughts() {
  const response = await fetch("/api/v1/thoughts");
  const snapshot = await response.json();
  thoughtState.activeStep = snapshot.active_step || null;
  thoughtState.activeChain = snapshot.active_chain || null;
  thoughtState.perception = snapshot.perception || null;
  thoughtState.suspended = snapshot.suspended_chains || [];
  thoughtState.steps = snapshot.steps || [];
  renderThoughts();
  connectThoughtStream();
}

function connectThoughtStream() {
  const stream = new EventSource("/api/v1/thoughts/stream");
  stream.addEventListener("thought.perception", (event) => {
    thoughtState.perception = JSON.parse(event.data);
    renderThoughts();
  });
  stream.addEventListener("thought.chain", (event) => {
    const chain = JSON.parse(event.data);
    thoughtState.activeChain = chain.status === "idle" ? null : chain;
    if (thoughtState.activeChain && thoughtState.activeChain.id) {
      thoughtState.suspended = thoughtState.suspended.filter((item) => item.id !== thoughtState.activeChain.id);
    }
    renderThoughts();
  });
  stream.addEventListener("thought.chain.suspended", (event) => {
    const chain = JSON.parse(event.data);
    thoughtState.suspended = [chain, ...thoughtState.suspended.filter((item) => item.id !== chain.id)];
    renderThoughts();
  });
  stream.addEventListener("thought.step.started", (event) => {
    thoughtState.activeStep = JSON.parse(event.data);
    renderThoughts();
  });
  stream.addEventListener("thought.step.delta", (event) => {
    thoughtState.activeStep = JSON.parse(event.data);
    renderThoughts();
  });
  stream.addEventListener("thought.step.completed", (event) => {
    const record = JSON.parse(event.data);
    thoughtState.steps.push(record);
    if (thoughtState.activeStep && thoughtState.activeStep.step_id === record.step_id) {
      thoughtState.activeStep = null;
    }
    renderThoughts();
  });
}

function renderThoughts() {
  renderActiveStep();
  renderChainState();
  renderPerception();
  renderTimeline();
}

function renderActiveStep() {
  const step = thoughtState.activeStep;
  if (!step) {
    activeStepEl.className = "active-step empty";
    activeStepEl.textContent = "No active step.";
    return;
  }

  activeStepEl.className = "active-step";
  activeStepEl.textContent = [
    `#${step.step_index} ${step.frame}`,
    `status: ${step.status}`,
    step.primitive_choice ? `choice: ${step.primitive_choice}` : "",
    "",
    step.raw_output || "",
  ]
    .filter(Boolean)
    .join("\n");
}

function renderChainState() {
  chainStateEl.textContent = JSON.stringify(
    {
      active_chain: thoughtState.activeChain,
      suspended_chains: thoughtState.suspended,
    },
    null,
    2,
  );
}

function renderPerception() {
  perceptionEl.textContent = JSON.stringify(thoughtState.perception || {}, null, 2);
}

function renderTimeline() {
  timelineEl.innerHTML = "";
  for (const step of thoughtState.steps) {
    const card = document.createElement("article");
    card.className = "step-card";

    const meta = document.createElement("div");
    meta.className = "step-meta";
    meta.innerHTML = `
      <span class="pill">#${escapeHtml(step.step_index)}</span>
      <span class="pill">${escapeHtml(step.frame || "")}</span>
      <span class="pill">${escapeHtml(step.status || "completed")}</span>
      <span class="pill">${escapeHtml(step.primitive_choice || "complete")}</span>
    `;

    const output = document.createElement("pre");
    output.className = "step-output";
    output.textContent = step.raw_output || "";

    const perception = document.createElement("pre");
    perception.className = "step-perception";
    perception.textContent = JSON.stringify(step.perception_snapshot || {}, null, 2);

    card.append(meta, output, perception);
    timelineEl.appendChild(card);
  }
}

function escapeHtml(value) {
  return String(value)
    .replaceAll("&", "&amp;")
    .replaceAll("<", "&lt;")
    .replaceAll(">", "&gt;")
    .replaceAll('"', "&quot;")
    .replaceAll("'", "&#39;");
}

bootThoughts().catch((error) => {
  timelineEl.innerHTML = `<article class="step-card">Failed to load thought surface: ${escapeHtml(error.message)}</article>`;
});
