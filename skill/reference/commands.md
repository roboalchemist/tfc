# tfc Command Reference

## Global Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--org` | | Organization name (env: `TFC_ORG`) |
| `--json` | `-j` | JSON output |
| `--plaintext` | | Tab-separated output |
| `--template` | `-t` | Go template string |
| `--jq` | | jq expression (implies --json) |
| `--fields` | | Comma-separated field list (implies --json) |
| `--no-color` | | Disable colored output |
| `--debug` | | Verbose logging to stderr |
| `--output` | `-o` | Write output to file |

## workspace (ws)

```bash
tfc ws list [--search NAME] [--page-size N]
tfc ws show <name-or-id>
tfc ws create <name> [--description TEXT] [--terraform-version VER] [--working-directory DIR] [--auto-apply] [--vcs-repo REPO] [--project-id ID]
tfc ws update <name-or-id> [--description TEXT] [--terraform-version VER] [--working-directory DIR] [--auto-apply]
tfc ws delete <name-or-id>
tfc ws lock <id> [--reason TEXT]
tfc ws unlock <id>
```

## run

```bash
tfc run list --workspace <name-or-id> [--status STATUS] [--page-size N]
tfc run show <id>
tfc run create --workspace <name-or-id> [--message TEXT] [--is-destroy] [--auto-apply] [--target RESOURCES]
tfc run apply <id> [--comment TEXT]
tfc run discard <id> [--comment TEXT]
tfc run cancel <id> [--force]
```

## plan

```bash
tfc plan show <id>
tfc plan log <id>
```

## apply

```bash
tfc apply show <id>
tfc apply log <id>
```

## state-version (sv)

```bash
tfc sv list --workspace <name-or-id> [--page-size N]
tfc sv show <id>
tfc sv create --workspace <id> --file <path> [--serial N] [--md5 HASH] [--lineage UUID]
tfc sv download <id>
```

## var

```bash
tfc var list --workspace <id>
tfc var show <id>
tfc var create --workspace <id> --key KEY [--value VAL] [--description TEXT] [--category terraform|env] [--hcl] [--sensitive]
tfc var update <id> --workspace <id> [--key KEY] [--value VAL] [--description TEXT] [--hcl] [--sensitive]
tfc var delete <id> --workspace <id>
```

## varset (vs)

```bash
tfc vs list
tfc vs show <id>
tfc vs create <name> [--description TEXT] [--global]
tfc vs update <id> [--name NAME] [--description TEXT] [--global]
tfc vs delete <id>
tfc vs apply <id> --workspace <ids...>
tfc vs remove <id> --workspace <ids...>
```

## org

```bash
tfc org list
tfc org show [name]
```

## team

```bash
tfc team list
tfc team show <id>
tfc team create <name> [--visibility secret|organization] [--manage-workspaces] [--manage-modules] [--manage-providers] [--manage-policies]
tfc team update <id> [--name NAME] [--visibility secret|organization]
tfc team delete <id>
```

## team-access (ta)

```bash
tfc ta list --workspace <id>
tfc ta show <id>
tfc ta add --workspace <id> --team <id> [--access read|plan|write|admin|custom]
tfc ta update <id> [--access read|plan|write|admin|custom]
tfc ta remove <id>
```

## project (proj)

```bash
tfc proj list
tfc proj show <id>
tfc proj create <name> [--description TEXT]
tfc proj update <id> [--name NAME] [--description TEXT]
tfc proj delete <id>
```

## policy (pol)

```bash
tfc pol list
tfc pol show <id>
```

## policy-set (ps)

```bash
tfc ps list
tfc ps show <id>
```

## policy-check (pc)

```bash
tfc pc list --run <id>
tfc pc show <id>
tfc pc override <id>
```

## run-task (rt)

```bash
tfc rt list
tfc rt show <id>
tfc rt create <name> --url <URL> [--description TEXT] [--hmac-key KEY] [--enabled]
tfc rt update <id> [--name NAME] [--url URL] [--description TEXT] [--enabled]
tfc rt delete <id>
```

## notification (notif)

```bash
tfc notif list --workspace <id>
tfc notif show <id>
tfc notif create <name> --workspace <id> --destination-type <type> [--url URL] [--triggers EVENTS] [--enabled]
tfc notif update <id> [--name NAME] [--url URL] [--triggers EVENTS] [--enabled]
tfc notif delete <id>
```

## agent-pool (ap)

```bash
tfc ap list
tfc ap show <id>
```

## audit-trail (audit)

```bash
tfc audit list [--since TIMESTAMP] [--page-size N]
```

## config-version (cv)

```bash
tfc cv list --workspace <id> [--page-size N]
tfc cv show <id>
tfc cv create --workspace <id> [--auto-queue-runs] [--speculative]
tfc cv upload <id> <file>
```

## Environment Variables

| Variable | Required | Description |
|----------|----------|-------------|
| `TFC_TOKEN` | Yes | Terraform Cloud API token |
| `TFC_ORG` | No | Default organization name |
| `TFC_ADDRESS` | No | Base URL (default: `https://app.terraform.io`) |
