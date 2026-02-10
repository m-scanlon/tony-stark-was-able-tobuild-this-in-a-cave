# Jarvis Memory & Project State Architecture

## Overview
This document defines how project state, memory, and versioning work inside the Jarvis system.  
The goal is to create a **living, versioned project model** that the AI can safely modify, audit, and roll back—without relying on a traditional relational database.

The architecture is built around two core components:

1. **Object Store** – the system of record for all project state and commits
2. **Vector Store** – a semantic index built from project data for fast retrieval

---

## Core Principles

### 1. Single Source of Truth
Project state must live in exactly one canonical location.

- All project data is stored as **versioned objects**
- The AI never edits state directly
- All changes occur through **explicit commits**

### 2. Append-Only Commit History
Project changes are stored as immutable commit objects.

Each commit:
- Contains a full snapshot of project state
- Includes metadata about the change
- Links to its parent commit
- Can be used for rollback

### 3. Vector Index is Derived Data
The vector database is **not** the source of truth.

- It indexes project state and documents
- It can be deleted and rebuilt at any time
- It exists only for semantic retrieval

---

## High-Level Architecture

```
+-------------------+
|      LLM Agent     |
+---------+---------+
          |
          v
+-------------------+
|  Commit Interface |
| (Change Requests) |
+---------+---------+
          |
          v
+-------------------+
|   Object Store     |
| (Project Commits)  |
+---------+---------+
          |
          v
+-------------------+
|   Vector Store     |
| (Semantic Index)   |
+-------------------+

---

## Object Store (System of Record)

The object store holds:

- Project state snapshots
- Commit history
- Attachments and documents

### Phase 1 (Local Mode)
Use the filesystem as the object store.

Example root: `.skyra/`

### Phase 2 (Distributed Mode)
Move the same structure to:

- AWS S3
- MinIO (self-hosted S3-compatible)

No structural changes required.

---

## Project Directory Structure

Example: `.skyra/projects/jarvis/`

```
jarvis/
├── HEAD.json
├── state.json
├── commits/
│   ├── 2026-02-09T21-10-33Z.json
│   └── 2026-02-09T21-25-04Z.json
└── attachments/
    ├── diagram.png
    └── notes.md
```

### Files Explained

#### `HEAD.json`
Pointer to the current commit.

```json
{
  "current_commit": "2026-02-09T21-25-04Z"
}
```

#### `state.json`
Materialized snapshot of the current state.

Derived from the commit referenced in HEAD.json.

---

## Commit Object Format

Each commit is immutable.

Example: `commits/2026-02-09T21-25-04Z.json`

```json
{
  "commit_id": "2026-02-09T21-10-33Z",
  "project": "jarvis",
  "parent": "2026-02-09T20-55-12Z",
  "actor": {
    "type": "ai",
    "model": "qwen2.5-coder:7b",
    "user": "mike"
  },
  "message": "Switch memory to vector index over project files",
  "timestamp": "2026-02-09T21-10-33Z",
  "changes": [
    {
      "op": "set",
      "path": "/memory/strategy",
      "value": "vector_index_over_files"
    }
  ],
  "snapshot": {
    "project": "jarvis",
    "goal": "Build a local-first AI assistant",
    "memory": {
      "strategy": "vector_index_over_files",
      "vector_db": "chroma"
    }
  }
}
```

---

## Implementation Steps

### Step 1 — Create the project structure
```
.skyra/projects/jarvis/
├── HEAD.json
├── state.json
└── commits/
```

### Step 2 — Create the genesis commit
The first snapshot of the project with initial state.

### Step 3 — Build a tiny commit tool
```
skyra commit
```
Which:
- Loads current state
- Shows proposed changes  
- Asks accept/reject
- Writes commit + updates HEAD

This is the core of the system - the transition from "AI experiments" to an actual personal operating system for projects.