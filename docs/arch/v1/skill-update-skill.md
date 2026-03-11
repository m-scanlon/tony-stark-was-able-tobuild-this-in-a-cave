# Skill: Update Skill (Compatibility)

`update_skill` appears in older docs as the primitive used to roll forward a skill version.

Canonical v1 naming uses `provision_skill`:

- a new skill version is content-addressed
- it gets a new `skill_id`
- it is provisioned as a new executable entry
- old versions remain append-only

In other words, "update" is implemented as "provision a new version," not mutate in place.

## Mapping

- legacy term: `update_skill`
- canonical primitive: `provision_skill`

## Related

- `docs/arch/v1/skill.md`
- `docs/arch/v1/skill-lifecycle.md`
- `docs/arch/v1/crypto-protocol.md`
- `docs/arch/v1/skill-improvement.md`

