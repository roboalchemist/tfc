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
# List workspaces
tfc ws list
tfc ws list --search "prod"

# Show workspace details
tfc ws show my-workspace

# List runs for a workspace
tfc run list --workspace my-workspace

# Show run details
tfc run show run-abc123

# View plan log
tfc plan log plan-abc123

# List variables
tfc var list --workspace ws-abc123

# List teams
tfc team list

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

| Command | Alias | Description |
|---------|-------|-------------|
| `workspace` | `ws` | Manage workspaces |
| `run` | | Manage runs |
| `plan` | | View plan details/logs |
| `apply` | | View apply details/logs |
| `state-version` | `sv` | Manage state versions |
| `var` | | Manage workspace variables |
| `varset` | `vs` | Manage variable sets |
| `org` | | View organizations |
| `team` | | Manage teams |
| `team-access` | `ta` | Manage team workspace access |
| `project` | `proj` | Manage projects |
| `policy` | `pol` | View policies |
| `policy-set` | `ps` | View policy sets |
| `policy-check` | `pc` | Manage policy checks |
| `run-task` | `rt` | Manage run tasks |
| `notification` | `notif` | Manage notifications |
| `agent-pool` | `ap` | View agent pools |
| `audit-trail` | `audit` | View audit events |
| `config-version` | `cv` | Manage config versions |

See [reference/commands.md](reference/commands.md) for full command details.
