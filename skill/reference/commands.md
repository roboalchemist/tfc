# tfc Command Reference

Commands marked with (stub) return "not yet implemented".

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
tfc ws create <name> [...]                    # (stub)
tfc ws update <name-or-id> [...]              # (stub)
tfc ws delete <name-or-id>                    # (stub)
tfc ws lock <id> [--reason TEXT]              # (stub)
tfc ws unlock <id>                            # (stub)
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
tfc sv create --workspace <id> --file <path> [...]  # (stub)
tfc sv download <id>                                 # (stub)
```

## var

```bash
tfc var list --workspace <id>
tfc var show <id>
tfc var create --workspace <id> --key KEY [...]   # (stub)
tfc var update <id> --workspace <id> [...]        # (stub)
tfc var delete <id> --workspace <id>              # (stub)
```

## varset (vs) — all stubs

```bash
tfc vs list                                       # (stub)
tfc vs show <id>                                  # (stub)
tfc vs create <name> [...]                        # (stub)
tfc vs update <id> [...]                          # (stub)
tfc vs delete <id>                                # (stub)
tfc vs apply <id> --workspace <ids...>            # (stub)
tfc vs remove <id> --workspace <ids...>           # (stub)
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
tfc team create <name> [...]                      # (stub)
tfc team update <id> [...]                        # (stub)
tfc team delete <id>                              # (stub)
```

## team-access (ta) — all stubs

```bash
tfc ta list --workspace <id>                      # (stub)
tfc ta show <id>                                  # (stub)
tfc ta add --workspace <id> --team <id> [...]     # (stub)
tfc ta update <id> [...]                          # (stub)
tfc ta remove <id>                                # (stub)
```

## project (proj)

```bash
tfc proj list
tfc proj show <id>
tfc proj create <name> [--description TEXT]       # (stub)
tfc proj update <id> [...]                        # (stub)
tfc proj delete <id>                              # (stub)
```

## policy (pol) — all stubs

```bash
tfc pol list                                      # (stub)
tfc pol show <id>                                 # (stub)
```

## policy-set (ps) — all stubs

```bash
tfc ps list                                       # (stub)
tfc ps show <id>                                  # (stub)
```

## policy-check (pc)

```bash
tfc pc list --run <id>
tfc pc show <id>
tfc pc override <id>
```

## run-task (rt) — all stubs

```bash
tfc rt list                                       # (stub)
tfc rt show <id>                                  # (stub)
tfc rt create <name> --url <URL> [...]            # (stub)
tfc rt update <id> [...]                          # (stub)
tfc rt delete <id>                                # (stub)
```

## notification (notif) — all stubs

```bash
tfc notif list --workspace <id>                   # (stub)
tfc notif show <id>                               # (stub)
tfc notif create <name> --workspace <id> [...]    # (stub)
tfc notif update <id> [...]                       # (stub)
tfc notif delete <id>                             # (stub)
```

## agent-pool (ap) — all stubs

```bash
tfc ap list                                       # (stub)
tfc ap show <id>                                  # (stub)
```

## audit-trail (audit) — all stubs

```bash
tfc audit list [--since TIMESTAMP] [--page-size N]  # (stub)
```

## config-version (cv) — all stubs

```bash
tfc cv list --workspace <id> [--page-size N]      # (stub)
tfc cv show <id>                                  # (stub)
tfc cv create --workspace <id> [...]              # (stub)
tfc cv upload <id> <file>                         # (stub)
```

## Environment Variables

| Variable | Required | Description |
|----------|----------|-------------|
| `TFC_TOKEN` | Yes | Terraform Cloud API token |
| `TFC_ORG` | No | Default organization name |
| `TFC_ADDRESS` | No | Base URL (default: `https://app.terraform.io`) |
