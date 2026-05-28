/**
 * Wire types for the v.1 runtime protocol.
 *
 * Source of truth: backend-spec-v1-runtime.pdf (Mike, 2026-05-13).
 *
 * The runtime owns all relational state. These types describe what flows over
 * the WebSocket — they are NOT a parallel domain model. Components derive views
 * from these messages; they do not maintain their own.
 *
 * Vocabulary: use the runtime names on the wire (Being, Reality, Thread,
 * Exchange, Topology, Relation). React components can be named however we want.
 *
 * Anything labeled INFERRED is a frontend guess on a shape Mike didn't fully
 * specify. Flag these in the next backend sync.
 */

// ---------------------------------------------------------------------------
// Beings
// ---------------------------------------------------------------------------

/** Genome type — distinguishes runtime origin of a being. */
export type BeingType = "llm" | "user" | "agent" | "process";

/** A weighted edge from a being to one of its Relationships (context source). */
export interface BeingRelationship {
  target: string;
  weight: number;
  /** INFERRED optional — Mike's example shows it for some entries, not all. */
  usage?: number;
}

/** A weighted edge from a being to one of its Expressors (execution capability). */
export interface BeingExpressor {
  target: string;
  weight: number;
}

/** A persisted memory item written by/for the being. */
export interface BeingMemoryItem {
  filename: string;
  content: string;
}

/** A named skill — durable competence the being can draw on. */
export interface BeingSkill {
  name: string;
  content: string;
}

export interface BeingMemories {
  items?: BeingMemoryItem[];
  skills?: BeingSkill[];
}

/**
 * A Being — an entity in the system. Could be LLM-backed (skyra),
 * user-backed (michael), agent-backed, or process-backed.
 */
export interface Being {
  name: string;
  type: BeingType;
  identity: string;
  purpose: string;
  /** INFERRED enum — Mike's example shows "active" but doesn't enumerate. */
  status: string;
  peers: string[];
  weight: number;
  relationships: BeingRelationship[];
  expressors: BeingExpressor[];
  memories: BeingMemories;
}

// ---------------------------------------------------------------------------
// Threads
// ---------------------------------------------------------------------------

export interface ThreadEdge {
  from: string;
  to: string;
}

/**
 * A Thread — append-only graph of beings participating in a routed loop.
 * Each thread has its own member list and edges; edges only ever add.
 */
export interface Thread {
  id: string;
  created_by: string;
  active: boolean;
  members: string[];
  edges: ThreadEdge[];
}

// ---------------------------------------------------------------------------
// Exchanges
// ---------------------------------------------------------------------------

export interface ExchangeEntry {
  index: number;
  from: string;
  content: string;
  ts: number;
}

/**
 * An Exchange — conversation history between two beings. Append-only.
 * Either party can append. No turns. No request/response.
 *
 * `key` is the canonical id, e.g. "michael:skyra". `parties` is the unordered
 * pair of being names.
 */
export interface Exchange {
  key: string;
  parties: [string, string];
  active: boolean;
  entries: ExchangeEntry[];
}

// ---------------------------------------------------------------------------
// Economics
// ---------------------------------------------------------------------------

/** Open-shape economics block. Mike: not yet enforced in the descent. */
export interface Economics {
  fields: Record<string, number>;
}

// ---------------------------------------------------------------------------
// Topology
// ---------------------------------------------------------------------------

/**
 * A node in the weighted recursive topology. Every node can contain
 * Relationships (context edges) and Expressors (execution edges). The root
 * node has children (top-level beings). Inside a being, Relationships and
 * Expressors are themselves Realities and can contain their own subtrees.
 *
 * INFERRED detail: Mike's example shows `children` only on the root and
 * `relationships`/`expressors` on inner nodes. We treat all three as
 * optional everywhere so the tree can grow uniformly.
 */
export interface TopologyNode {
  id: string;
  /** Reality kind — "Universe", "Self", "Relationship", "Memory", "Think", "Act", etc. */
  type: string;
  weight?: number;
  children?: TopologyNode[];
  relationships?: TopologyNode[];
  expressors?: TopologyNode[];
}

// ---------------------------------------------------------------------------
// Universe state (full snapshot)
// ---------------------------------------------------------------------------

export interface UniverseState {
  beings: Being[];
  threads: Thread[];
  exchanges: Exchange[];
  economics: Economics;
  topology: TopologyNode;
}

// ---------------------------------------------------------------------------
// Delta payloads
// ---------------------------------------------------------------------------

export interface EntryDelta {
  exchange: string;
  index: number;
  from: string;
  content: string;
}

/** INFERRED — Mike's section 5 names the type but doesn't show the payload. */
export interface EdgeDelta {
  thread_id: string;
  from: string;
  to: string;
}

/** INFERRED — "weight changed on a relationship or expressor." */
export interface WeightDelta {
  being: string;
  kind: "relationship" | "expressor";
  target: string;
  weight: number;
}

/** INFERRED — "new memory written." */
export interface MemoryDelta {
  being: string;
  /** Either an item or a skill — backend picks. */
  item?: BeingMemoryItem;
  skill?: BeingSkill;
}

/** INFERRED — "topology changed (new being adds a subtree)." */
export interface TopologyDelta {
  parent: string;
  subtree: TopologyNode;
}

export interface ErrorPayload {
  origin: string;
  message: string;
}

export interface AuthFailPayload {
  reason: string;
}

// ---------------------------------------------------------------------------
// Message envelope
// ---------------------------------------------------------------------------

/**
 * Every message — server or client — wears this envelope.
 * From Mike's section 6: { id, type, ts, payload }.
 */
interface Envelope {
  id: string;
  ts: number;
}

export type ServerMessage =
  | (Envelope & { type: "auth_ok"; payload: Record<string, never> })
  | (Envelope & { type: "auth_fail"; payload: AuthFailPayload })
  | (Envelope & { type: "universe"; payload: UniverseState })
  | (Envelope & { type: "entry"; payload: EntryDelta })
  | (Envelope & { type: "edge"; payload: EdgeDelta })
  | (Envelope & { type: "being"; payload: Being })
  | (Envelope & { type: "weight"; payload: WeightDelta })
  | (Envelope & { type: "memory"; payload: MemoryDelta })
  | (Envelope & { type: "topology"; payload: TopologyDelta })
  | (Envelope & { type: "error"; payload: ErrorPayload });

/** Client → server: auth handshake. First message after connect. */
export interface AuthPayload {
  token: string;
}

/** Client → server: user input. `target` optional; runtime routes by weight if omitted. */
export interface ImpulsePayload {
  origin: string;
  content: string;
  target?: string;
}

export type ClientMessage =
  | (Envelope & { type: "auth"; payload: AuthPayload })
  | (Envelope & { type: "impulse"; payload: ImpulsePayload });

export type AnyMessage = ServerMessage | ClientMessage;
