from fastapi import FastAPI
from pydantic import BaseModel, Field


app = FastAPI(title="skyra-listener", version="0.1.0")


class ListenerEvent(BaseModel):
    transcript: str = Field(..., min_length=1, max_length=2000)
    source: str = Field(default="voice", max_length=64)
    wake_word_detected: bool = True


class ListenerDecision(BaseModel):
    decision: str
    invoke_frontdoor_llm: bool
    reason: str


@app.get("/health")
def health() -> dict[str, str]:
    return {"status": "ok"}


@app.post("/listener/event", response_model=ListenerDecision)
def listener_event(event: ListenerEvent) -> ListenerDecision:
    text = event.transcript.strip().lower()

    if not event.wake_word_detected:
        return ListenerDecision(
            decision="ignore",
            invoke_frontdoor_llm=False,
            reason="wake word not detected",
        )

    if text in {"", "uh", "um", "hmm"}:
        return ListenerDecision(
            decision="ignore",
            invoke_frontdoor_llm=False,
            reason="low signal transcript",
        )

    if "cancel" in text or "never mind" in text:
        return ListenerDecision(
            decision="ignore",
            invoke_frontdoor_llm=False,
            reason="explicit cancel intent",
        )

    return ListenerDecision(
        decision="dispatch",
        invoke_frontdoor_llm=True,
        reason="wake word + usable transcript",
    )
