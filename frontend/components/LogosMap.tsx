"use client";

import clsx from "clsx";
import { useAppStore, selectBeings, selectSkyraWeight } from "@/lib/store";
import type { Being, TopologyNode } from "@/lib/protocol/types";
import { Logos } from "./Logos";

/**
 * TRANSITIONAL — exports `LogosMap` but renders the v.1 Universe-present view:
 * the topology tree + every being. Will be renamed to <UniversePresentView>
 * in the Phase 2 UI rework (task #22), and the React Flow visualization will
 * be reintroduced there with proper weighted-edge rendering.
 *
 * For now this is a readable, scrollable rendering of the universe state so
 * we can verify the wire shapes end-to-end.
 */

export function LogosMap() {
  const beings = useAppStore(selectBeings);
  const topology = useAppStore((s) => s.topology);
  const economics = useAppStore((s) => s.economics);
  const selectBeing = useAppStore((s) => s.selectBeing);
  const selectedBeing = useAppStore((s) => s.selectedBeing);
  const skyraWeight = useAppStore(selectSkyraWeight);

  if (beings.length === 0) {
    return (
      <div className="grid h-full place-items-center">
        <p className="font-mono text-sm text-ink-dim">waiting for universe snapshot…</p>
      </div>
    );
  }

  const focused = selectedBeing
    ? beings.find((b) => b.name === selectedBeing) ?? null
    : null;

  return (
    <div className="grid h-full grid-cols-[260px_minmax(0,1fr)_320px] divide-x divide-surface-edge">
      {/* Left rail — beings list. */}
      <aside className="min-h-0 overflow-y-auto p-3">
        <header className="mb-2 flex items-center justify-between">
          <span className="font-mono text-xs uppercase tracking-wider text-ink-dim">
            beings
          </span>
          <SkyraIndicator weight={skyraWeight} />
        </header>
        <div className="space-y-1.5">
          {beings.map((b) => (
            <Logos
              key={b.name}
              being={b}
              selected={b.name === selectedBeing}
              onClick={() => selectBeing(b.name === selectedBeing ? null : b.name)}
            />
          ))}
        </div>
      </aside>

      {/* Center — topology tree. */}
      <main className="min-h-0 overflow-y-auto p-4">
        <header className="mb-3">
          <h2 className="font-mono text-xs uppercase tracking-wider text-ink-dim">
            topology
          </h2>
          <p className="text-[10px] uppercase tracking-wider text-ink-faint">
            recursive weighted composition · placeholder rendering
          </p>
        </header>
        <TopologyTree node={topology} />
      </main>

      {/* Right rail — selected being detail. */}
      <aside className="min-h-0 overflow-y-auto p-3">
        {focused ? (
          <BeingDetail being={focused} />
        ) : (
          <div className="mt-12 text-center font-mono text-xs text-ink-dim">
            select a being
          </div>
        )}
      </aside>

      {/* Floating economics readout. */}
      <div className="pointer-events-none absolute bottom-3 left-1/2 -translate-x-1/2 font-mono text-[10px] uppercase tracking-wider text-ink-faint">
        {Object.entries(economics.fields)
          .map(([k, v]) => `${k}: ${v}`)
          .join(" · ")}
      </div>
    </div>
  );
}

function TopologyTree({ node, depth = 0 }: { node: TopologyNode; depth?: number }) {
  const isLeaf =
    (node.children?.length ?? 0) === 0 &&
    (node.relationships?.length ?? 0) === 0 &&
    (node.expressors?.length ?? 0) === 0;

  return (
    <div className={clsx(depth > 0 && "ml-4 border-l border-surface-edge pl-3")}>
      <div className="flex items-center gap-2 py-0.5 font-mono text-xs">
        <span className="text-ink">{node.id}</span>
        <span className="text-ink-dim">[{node.type}]</span>
        {typeof node.weight === "number" && (
          <span className="text-ink-faint">w={node.weight.toFixed(2)}</span>
        )}
        {isLeaf && <span className="text-ink-faint">·</span>}
      </div>
      {node.children?.map((c) => <TopologyTree key={c.id} node={c} depth={depth + 1} />)}
      {node.relationships?.map((c) => (
        <TopologyTree key={`r:${c.id}`} node={c} depth={depth + 1} />
      ))}
      {node.expressors?.map((c) => (
        <TopologyTree key={`e:${c.id}`} node={c} depth={depth + 1} />
      ))}
    </div>
  );
}

function BeingDetail({ being }: { being: Being }) {
  return (
    <div className="space-y-3">
      <header>
        <div className="font-mono text-sm text-ink">{being.name}</div>
        <div className="text-[10px] uppercase tracking-wider text-ink-dim">
          {being.type} · weight {being.weight.toFixed(2)} · {being.status}
        </div>
        <p className="mt-1 text-xs italic text-ink-dim">{being.identity}</p>
        <p className="mt-1 text-xs text-ink-dim">{being.purpose}</p>
      </header>

      <Section title="relationships">
        {being.relationships.length === 0 ? (
          <Empty />
        ) : (
          being.relationships.map((r) => <WeightedRow key={r.target} target={r.target} weight={r.weight} usage={r.usage} />)
        )}
      </Section>

      <Section title="expressors">
        {being.expressors.length === 0 ? (
          <Empty />
        ) : (
          being.expressors.map((e) => <WeightedRow key={e.target} target={e.target} weight={e.weight} />)
        )}
      </Section>

      {being.memories.items && being.memories.items.length > 0 && (
        <Section title="memories">
          {being.memories.items.map((m) => (
            <div key={m.filename} className="text-xs text-ink-dim">
              <span className="font-mono">{m.filename}</span>
              <p className="mt-0.5 truncate text-ink-faint">{m.content}</p>
            </div>
          ))}
        </Section>
      )}

      {being.memories.skills && being.memories.skills.length > 0 && (
        <Section title="skills">
          {being.memories.skills.map((s) => (
            <div key={s.name} className="text-xs text-ink-dim">
              <span className="font-mono">{s.name}</span>
              <p className="mt-0.5 truncate text-ink-faint">{s.content}</p>
            </div>
          ))}
        </Section>
      )}
    </div>
  );
}

function Section({ title, children }: { title: string; children: React.ReactNode }) {
  return (
    <section>
      <h3 className="mb-1 font-mono text-[10px] uppercase tracking-wider text-ink-faint">
        {title}
      </h3>
      <div className="space-y-1">{children}</div>
    </section>
  );
}

function Empty() {
  return <p className="font-mono text-xs italic text-ink-faint">empty</p>;
}

function WeightedRow({ target, weight, usage }: { target: string; weight: number; usage?: number }) {
  return (
    <div className="flex items-center justify-between gap-2 font-mono text-xs">
      <span className="text-ink">{target}</span>
      <div className="flex items-center gap-2">
        {typeof usage === "number" && (
          <span className="text-ink-faint">×{usage}</span>
        )}
        <div className="h-1 w-16 overflow-hidden rounded bg-surface-edge">
          <div
            className="h-full bg-logos-being"
            style={{ width: `${Math.round(weight * 100)}%` }}
          />
        </div>
        <span className="text-ink-dim">{weight.toFixed(2)}</span>
      </div>
    </div>
  );
}

function SkyraIndicator({ weight }: { weight: number | null }) {
  if (weight == null) return null;
  // High weight = active = system shaping. Low = stable.
  const isShaping = weight > 0.7;
  return (
    <div
      className={clsx(
        "flex items-center gap-1.5 font-mono text-[10px] uppercase tracking-wider",
        isShaping ? "text-amber-300" : "text-emerald-300",
      )}
      title={`Skyra: ${weight.toFixed(2)} — ${isShaping ? "shaping" : "stable"}`}
    >
      <span
        className={clsx(
          "h-1.5 w-1.5 rounded-full",
          isShaping ? "bg-amber-300" : "bg-emerald-300",
        )}
      />
      skyra {isShaping ? "shaping" : "stable"}
    </div>
  );
}
