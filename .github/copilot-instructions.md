# Copilot Instructions

## Project Rules

This repository is a small typed request binding library for Fiber.

Keep all changes aligned with these rules:

- Do not turn `fiberkit` into a framework.
- Do not add DI, service resolution, auto-wiring, or reflection-based argument magic.
- Do not introduce global configuration unless explicitly requested.
- Do not add broad abstractions when a narrow helper is enough.
- Do not change the public API shape without a strong reason.
- Do not add new files or folders unless they are clearly justified.
- Do not touch sample or docs code unless the task requires it.
- Prefer explicit, small, readable code over clever code.
- Prefer local helpers over shared abstractions if the shared abstraction weakens clarity.
- Keep wrapper behavior predictable: parse, validate, return a simple error, call the handler.

## Implementation Rules

- Keep `Body`, `Query`, `Params`, `ParamsBody`, `Set`, and `Get` simple and direct.
- If you add behavior, keep it opt-in and local to the wrapper.
- Validation should remain based on `go-playground/validator/v10` unless the user asks otherwise.
- Error responses should stay predictable and easy to override only if a clear requirement appears.
- Preserve the typed handler and typed middleware use case.
- Avoid changing naming conventions unless a task explicitly asks for it.

## Testing Rules

- Any behavior change must be covered by tests.
- Do not remove existing tests unless they are obsolete because of an explicit change.
- Add focused tests for new behavior.
- Keep tests small and readable.
- Verify both success paths and failure paths when applicable.

## Editing Rules

- Make the smallest change that fully solves the request.
- Do not rewrite unrelated code.
- Do not reformat large areas of the codebase unless necessary.
- Preserve existing style in nearby code.
- If a task is ambiguous, choose the simplest interpretation that fits the current design.

## Commit Message Protocol

All commits MUST follow this format strictly:

`[EMOJI] [type]: short description`

The body is required. It should be concise, useful, and explain what changed and why.

Types:

- ✨ `feat`: New feature
- 🐛 `fix`: Bug fix
- ♻️ `refactor`: Code change that neither fixes a bug nor adds a feature
- 📝 `docs`: Documentation changes
- 🎨 `style`: Formatting, missing semi-colons, etc.
- 🔧 `chore`: Updating build tasks, package manager configs, etc.
- ✅ `test`: Adding missing tests or correcting existing tests

Example:

`✨ feat: add optional validation error formatting`
