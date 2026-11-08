# Product specification

This directory defines the product behavior and operator contracts. Group requirements by topic and split files as they grow. Give each rule one authoritative home and link to its ID from related requirements, issues, and tests. Keep implementation decisions under `design/`.

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
