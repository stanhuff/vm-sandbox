# Product specification

This directory defines the product behavior and operator contracts. Group requirements by topic and split files as they grow. Give each rule one authoritative home and link to its ID from related requirements, issues, and tests.

Use the [glossary](glossary.md) for shared terms in requirements, stories, and tests.

## Requirements and design decisions

Keep the spec small. Give a behavior a requirement ID when changing it would warrant review as a product change. Leave other behavior open to implementation judgment. A story describes a user goal and tracks work toward it. The requirements it adds or changes continue to govern the product after the story closes.

The AI may draft requirements from the story and make routine design choices. Record a design decision under `design/` when future work should reuse it and the code alone leaves room for inconsistent choices. State the choice and why it was made, and link to any requirements it supports. Use a short architecture decision record when the rationale or tradeoffs matter. Design decisions can evolve as long as the implementation continues to satisfy the active requirements. Routine choices can live in code and tests.

## Early review

During refinement, flag a requirement as `review-early` when getting it wrong now would make a future change expensive. Add one sentence explaining what would need to be redesigned, rewritten, or migrated. Focus on the cost of changing the underlying decision, including compatibility and data migration where relevant. Leave routine requirements unmarked.

When a new requirement adds dependencies that make an existing decision expensive to reverse, flag the affected requirement in the same PR. The flag can change without replacing the requirement ID.

## Requirement IDs

Use a single repository-wide sequence, such as `REQ-0001`. The fixed `REQ-` marker makes references searchable. The number carries no category or file meaning, so requirements can move without being renumbered. IDs are unique across active requirements and the archive. Never reuse a published ID.

## Changes and clarifications

- A draft requirement can change freely until it first appears in a PR. From that point, keep its ID and meaning stable, including during review.
- Clarify an existing requirement in place only to resolve genuine ambiguity or correct wording without changing its meaning. A clarification must agree with existing implementation wherever it addresses the ambiguous point. If no implementation addresses that point, the clarification may choose an interpretation the original wording permits.
- Evaluate compatibility against the implementation before the PR's code changes. Supersede the requirement if the proposed change adds a new obligation, changes an explicit rule, or conflicts with existing behavior permitted by the original wording.
- If the implementation already violates an unambiguous requirement, fix the implementation and retain the ID.

## Supersession and withdrawal

Assign a new ID to replacement requirements. Mark the old ID as superseded and record the replacement ID or IDs. Mark a requirement withdrawn if it has no replacement.

Keep a retired-ID index in this directory so old issue and PR references can be resolved without reading archived text. Move withdrawn requirements into `archive/`. Move superseded requirements there once active references in the repo have been updated. Preserve the retired text and its ID.

## Cross-feature changes

During issue refinement, identify existing requirements affected by the proposed change. Update affected requirements, references, and tests together in the PR so the active spec remains consistent. Refer to requirements by ID so references remain usable when files move.
