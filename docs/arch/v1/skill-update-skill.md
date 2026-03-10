# Skill: Update Skill

Skill nodes are executable infrastructure. They cannot be mutated by general graph writes — not by reasoning, not by integrate, not by any other process. `update_skill` is the only path to modifying a skill node.

**Triggered by**: user approval of a skill change proposal
**Layer**: committed — skill nodes are first class
**Type**: system primitive skill — pre-provisioned in Redis at boot

---

## Why This Exists

A skill node is not just data. It is a contract — roadmap, boundary rules, state contract, validation criteria, Redis registration. Mutating it has execution consequences. The general write path (`write_node`, `write_edge`) is not aware of those consequences. `update_skill` is.

---

## What It Can Change

- Skill roadmap (tasks)
- Boundary rules
- Severity policy
- Replan budget
- Validation criteria
- Skill description (affects semantic discovery)

## What It Cannot Change

- Skill `id` — immutable
- Redis provisioning status — use `provision_skill` to register, deprovisioning does not exist

---

## The Process

```
proposed skill change arrives
  → update_skill validates the change against the skill contract
  → proposes change to user via propose_commit
  → user approves
  → skill node updated in committed layer
  → Redis entry updated to reflect new contract
```

---

## Rules

- reasoning and integrate are explicitly denied write access to skill nodes
- No other skill can mutate a skill node
- Every update requires user approval — skill nodes are committed layer, first class
- Every update is a commit — full audit trail

---

## Skill Contract

```
skill: update_skill
boundary_rules:
  write_node (skill type):  allow_always — this skill only
  write_node (other types): deny
  propose_commit:           allow_always
  redis_write:              allow_always (skill contract sync)
state_contract: committed (user approval required)
severity_policy:
  invalid contract change: halt and notify user
replan_budget: 0
```

---

## Related

- `docs/arch/v1/skill-lifecycle.md` — skill crystallization, provisioning
- `docs/arch/v1/skill-reasoning.md` — explicitly cannot write skill nodes
- `docs/arch/v1/skill-integrate.md` — explicitly cannot write skill nodes
- `docs/arch/v1/kernel.md` — Redis skill registry, trust boundary
