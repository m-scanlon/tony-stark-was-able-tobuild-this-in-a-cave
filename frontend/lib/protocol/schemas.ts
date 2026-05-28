/**
 * Zod schemas matching lib/protocol/types.ts (v.1 runtime).
 *
 * Every inbound WebSocket message is parsed through ServerMessageSchema before
 * it touches state. Anything that doesn't parse is logged as a drop and
 * thrown away — never silently passed through. This is how we catch backend
 * protocol drift the moment it happens.
 */

import { z } from "zod";

// ---------------------------------------------------------------------------
// Beings
// ---------------------------------------------------------------------------

export const BeingTypeSchema = z.enum(["llm", "user", "agent", "process"]);

export const BeingRelationshipSchema = z.object({
  target: z.string(),
  weight: z.number(),
  usage: z.number().optional(),
});

export const BeingExpressorSchema = z.object({
  target: z.string(),
  weight: z.number(),
});

export const BeingMemoryItemSchema = z.object({
  filename: z.string(),
  content: z.string(),
});

export const BeingSkillSchema = z.object({
  name: z.string(),
  content: z.string(),
});

export const BeingMemoriesSchema = z.object({
  items: z.array(BeingMemoryItemSchema).optional(),
  skills: z.array(BeingSkillSchema).optional(),
});

export const BeingSchema = z.object({
  name: z.string(),
  type: BeingTypeSchema,
  identity: z.string(),
  purpose: z.string(),
  status: z.string(),
  peers: z.array(z.string()),
  weight: z.number(),
  relationships: z.array(BeingRelationshipSchema),
  expressors: z.array(BeingExpressorSchema),
  memories: BeingMemoriesSchema,
});

// ---------------------------------------------------------------------------
// Threads
// ---------------------------------------------------------------------------

export const ThreadEdgeSchema = z.object({
  from: z.string(),
  to: z.string(),
});

export const ThreadSchema = z.object({
  id: z.string(),
  created_by: z.string(),
  active: z.boolean(),
  members: z.array(z.string()),
  edges: z.array(ThreadEdgeSchema),
});

// ---------------------------------------------------------------------------
// Exchanges
// ---------------------------------------------------------------------------

export const ExchangeEntrySchema = z.object({
  index: z.number(),
  from: z.string(),
  content: z.string(),
  ts: z.number(),
});

export const ExchangeSchema = z.object({
  key: z.string(),
  parties: z.tuple([z.string(), z.string()]),
  active: z.boolean(),
  entries: z.array(ExchangeEntrySchema),
});

// ---------------------------------------------------------------------------
// Economics
// ---------------------------------------------------------------------------

export const EconomicsSchema = z.object({
  fields: z.record(z.number()),
});

// ---------------------------------------------------------------------------
// Topology (recursive)
// ---------------------------------------------------------------------------

export interface TopologyNodeZ {
  id: string;
  type: string;
  weight?: number;
  children?: TopologyNodeZ[];
  relationships?: TopologyNodeZ[];
  expressors?: TopologyNodeZ[];
}

export const TopologyNodeSchema: z.ZodType<TopologyNodeZ> = z.lazy(() =>
  z.object({
    id: z.string(),
    type: z.string(),
    weight: z.number().optional(),
    children: z.array(TopologyNodeSchema).optional(),
    relationships: z.array(TopologyNodeSchema).optional(),
    expressors: z.array(TopologyNodeSchema).optional(),
  }),
);

// ---------------------------------------------------------------------------
// Universe state
// ---------------------------------------------------------------------------

export const UniverseStateSchema = z.object({
  beings: z.array(BeingSchema),
  threads: z.array(ThreadSchema),
  exchanges: z.array(ExchangeSchema),
  economics: EconomicsSchema,
  topology: TopologyNodeSchema,
});

// ---------------------------------------------------------------------------
// Delta payloads
// ---------------------------------------------------------------------------

export const EntryDeltaSchema = z.object({
  exchange: z.string(),
  index: z.number(),
  from: z.string(),
  content: z.string(),
});

export const EdgeDeltaSchema = z.object({
  thread_id: z.string(),
  from: z.string(),
  to: z.string(),
});

export const WeightDeltaSchema = z.object({
  being: z.string(),
  kind: z.enum(["relationship", "expressor"]),
  target: z.string(),
  weight: z.number(),
});

export const MemoryDeltaSchema = z.object({
  being: z.string(),
  item: BeingMemoryItemSchema.optional(),
  skill: BeingSkillSchema.optional(),
});

export const TopologyDeltaSchema = z.object({
  parent: z.string(),
  subtree: TopologyNodeSchema,
});

export const ErrorPayloadSchema = z.object({
  origin: z.string(),
  message: z.string(),
});

export const AuthFailPayloadSchema = z.object({
  reason: z.string(),
});

// ---------------------------------------------------------------------------
// Message envelope (server)
// ---------------------------------------------------------------------------

const envelope = z.object({
  id: z.string(),
  ts: z.number(),
});

export const ServerMessageSchema = z.discriminatedUnion("type", [
  envelope.extend({ type: z.literal("auth_ok"), payload: z.object({}) }),
  envelope.extend({ type: z.literal("auth_fail"), payload: AuthFailPayloadSchema }),
  envelope.extend({ type: z.literal("universe"), payload: UniverseStateSchema }),
  envelope.extend({ type: z.literal("entry"), payload: EntryDeltaSchema }),
  envelope.extend({ type: z.literal("edge"), payload: EdgeDeltaSchema }),
  envelope.extend({ type: z.literal("being"), payload: BeingSchema }),
  envelope.extend({ type: z.literal("weight"), payload: WeightDeltaSchema }),
  envelope.extend({ type: z.literal("memory"), payload: MemoryDeltaSchema }),
  envelope.extend({ type: z.literal("topology"), payload: TopologyDeltaSchema }),
  envelope.extend({ type: z.literal("error"), payload: ErrorPayloadSchema }),
]);

export type ServerMessageParsed = z.infer<typeof ServerMessageSchema>;
