# Agent notes

- [README.md](./README.md) covers building, the editing workflow for templates and generated files, adding components to the docs site and to JSON UI, and recording removed components in [MIGRATIONS.md](./MIGRATIONS.md).
- [DESIGN.md](./DESIGN.md) is authoritative for UI work: read it before producing or editing any UI.
- `examples/apps/` is a backend-free, HTMX-navigable prototype of the Access admin "Applications" area (served at `/examples/apps`). It mirrors `invopop/access` templates with stand-in models so the front end can be designed here and ported later. Read [examples/apps/HANDOFF.md](./examples/apps/HANDOFF.md) before changing it, and keep access's field and template names so the port stays mechanical.
