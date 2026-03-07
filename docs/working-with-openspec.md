# Working With OpenSpec

OpenSpec is the default implementation workflow in this repository.

## Standard Flow

1. Explore or frame the idea.
2. Create a change with `/opsx:propose "idea"`.
3. Implement the change with `/opsx:apply`.
4. Archive the finished change with `/opsx:archive`.

## Where Things Live

- Active changes: `openspec/changes/<change-name>/`
- Archived stable specs: `openspec/specs/`
- Project-level OpenSpec rules: `openspec/config.yaml`

## FastSpec Relationship

Use OpenSpec for work planning and change control.
Use FastSpec YAML for the durable structured knowledge that survives beyond a single change.
