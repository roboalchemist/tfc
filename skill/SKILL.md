# tfc

CLI for Terraform Cloud (HCP Terraform) API. Manage workspaces, runs, plans, state, variables, teams, and more.

## Setup

```bash
# Required: API token
export TFC_TOKEN="your-terraform-cloud-token"

# Optional: default organization (avoids --org on every command)
export TFC_ORG="my-org"

# Optional: custom base URL (for Terraform Enterprise)
export TFC_ADDRESS="https://tfe.example.com"
```

## Quick Reference

```bash
# Workspaces
tfc ws list
tfc ws list --search "prod"
tfc ws show my-workspace

# Runs (full lifecycle)
tfc run list --workspace my-workspace
tfc run show run-abc123
tfc run create --workspace my-workspace --message "Deploy update"
tfc run apply run-abc123 --comment "Approved"
tfc run discard run-abc123 --comment "Not needed"
tfc run cancel run-abc123
tfc run cancel run-abc123 --force

# Plans & Applies
tfc plan show plan-abc123
tfc plan log plan-abc123
tfc apply show apply-abc123
tfc apply log apply-abc123

# Policy checks
tfc pc list --run run-abc123
tfc pc show polchk-abc123
tfc pc override polchk-abc123

# Variables, teams, projects, state versions
tfc var list --workspace ws-abc123
tfc var show var-abc123
tfc team list
tfc team show team-abc123
tfc proj list
tfc proj show prj-abc123
tfc sv list --workspace ws-abc123
tfc sv show sv-abc123
tfc org list
tfc org show my-org

# JSON output with jq
tfc ws list --json --jq '.[].name'
tfc run show run-abc123 --json --fields status,message
```

## Output Modes

| Flag | Description |
|------|-------------|
| (default) | Colored table output |
| `--json` / `-j` | JSON output |
| `--plaintext` | Tab-separated (for piping) |
| `--template` | Go template |
| `--jq` | jq expression (implies --json) |
| `--fields` | Field selection (implies --json) |
| `--no-color` | Disable colors |
| `-o FILE` | Write to file |

## Command Groups

| Command | Alias | Description | Status |
|---------|-------|-------------|--------|
| `workspace` | `ws` | Manage workspaces | list, show |
| `run` | | Manage runs | list, show, create, apply, discard, cancel |
| `plan` | | View plan details/logs | show, log |
| `apply` | | View apply details/logs | show, log |
| `state-version` | `sv` | Manage state versions | list, show |
| `var` | | Manage workspace variables | list, show |
| `varset` | `vs` | Manage variable sets | stub |
| `org` | | View organizations | list, show |
| `team` | | Manage teams | list, show |
| `team-access` | `ta` | Manage team workspace access | stub |
| `project` | `proj` | Manage projects | list, show |
| `policy` | `pol` | View policies | stub |
| `policy-set` | `ps` | View policy sets | stub |
| `policy-check` | `pc` | Manage policy checks | list, show, override |
| `run-task` | `rt` | Manage run tasks | stub |
| `notification` | `notif` | Manage notifications | stub |
| `agent-pool` | `ap` | View agent pools | stub |
| `audit-trail` | `audit` | View audit events | stub |
| `config-version` | `cv` | Manage config versions | stub |

**Status key**: Listed subcommands are fully implemented. "stub" = all subcommands return "not yet implemented".

See [reference/commands.md](reference/commands.md) for full command details.
