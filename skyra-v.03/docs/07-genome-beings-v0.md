# Genome Beings v0

## Purpose

This document defines the genome layer of the current `v.03` model.

The genome layer exists so later runtime life is possible.

## Bootstrap Boundary

The creator writes `genome.skyra`.

The external launcher raises the host services first.

The creator, external launcher, and host services are pre-runtime beings.

Those host services make genome execution possible and host the beings that
come online when the genome executes.

The being factory being then creates the genome beings in dependency order and
routes each one through the kernel being for validation.

## The Being Factory

The being factory is a genome being hosted by the being factory service.

Its purpose is to:

- read singleton and template definitions from the genome
- instantiate beings and companion beings
- generate fresh verification keys for runtime and differentiated beings
- route created beings to the kernel

It does not validate.

It does not route arbitrary expression.

It creates beings.

## The Kernel

The kernel is a genome being hosted by the kernel service with special
authority.

Its purpose is to:

- validate nature
- validate registration-token entry into runtime registration
- verify signed envelopes on expression
- make beings live
- route expression

It does not create beings.

That work belongs to the being factory.

## The RDS Being

The RDS being is the non-cognitive shared storage being.

It stores:

- being records
- relationship records
- language records

It responds with acceptance or rejection.

## The Registry Beings

### Registry Writer

Writes being records and language records to the RDS being.

Those records include:

- identity
- purpose
- verification_key

### Theory Of Mind Being

Reads from the RDS being through their genome relationship.

Surfaces only the public slice:

- identity
- purpose

It does not expose verification keys.

## The Language Beings

### Language Writer

Records callable relationship language to language records and appends language
history through relationship records in the RDS being.

### Language Reader

Retrieves language records for a relationship.

## The Differentiator

The differentiator is a genome being responsible for revealing and reorganizing
overloaded identity.

It initiates differentiation through its own signed authority.

## The Genome Singletons

At minimum, the genome defines singleton beings for:

- being factory
- kernel
- RDS being
- registry writer
- language writer
- language reader
- differentiator
- main sensory being
- prefrontal being
- strategy being
- values being
- consequence being
- conflict being
- theory-of-mind being

## The Genome Templates

The genome also defines templates.

At minimum, the runtime needs templates for:

- personal hippocampus
- personal experience store being
- present being
- runtime birth
- differentiation birth
- purpose-bound boundary beings when required

Templates are not live beings by themselves.

They are creator-authored patterns the being factory instantiates later.

## Companion Beings

Every created being receives the companion beings required by its creation path.

Those companion beings are instantiated through genome templates or explicitly
seeded singleton definitions during genome bootstrap.

## Security

Genome beings receive creator-baked verification keys from `genome.skyra`.

Runtime and differentiated beings receive fresh verification keys from the
being factory on their creation path.

Every live genome being signs expression like any other live being.

## Short Framing

The genome layer is not mystical.

It is the creator-authored set of singleton beings, templates, and relationship
edges that the being factory and kernel use to make the runtime possible.
