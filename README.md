# tfc

A CLI for the [Terraform Cloud](https://app.terraform.io) (HCP Terraform) API. Manage workspaces, runs, plans, state, variables, teams, projects, and more from the command line.

## Installation

### Homebrew

```bash
brew install roboalchemist/tap/tfc
```

### Go

```bash
go install github.com/roboalchemist/tfc@latest
```

### From Source

```bash
git clone https://github.com/roboalchemist/tfc.git
cd tfc
make install
```

## Setup

```bash
export TFC_TOKEN="your-terraform-cloud-api-token"
export TFC_ORG="your-org"  # optional default organization
```

Generate a token at **User Settings > Tokens** in [Terraform Cloud](https://app.terraform.io/app/settings/tokens).

## Usage

```bash
# List workspaces
tfc ws list
tfc ws list --search "prod"

# Show workspace details
tfc ws show my-workspace

# List runs
tfc run list --workspace ws-abc123

# Show run details with plan/apply IDs
tfc run show run-abc123

# View plan log
tfc plan log plan-abc123

# List variables
tfc var list --workspace ws-abc123

# JSON output with jq filtering
tfc ws list --json --jq '.[].attributes.name'

# Field selection
tfc run show run-abc123 --json --fields id,attributes.status,plan_id
```

## Output Modes

| Flag | Description |
|------|-------------|
| *(default)* | Colored table |
| `--json` / `-j` | JSON |
| `--plaintext` | Tab-separated (for piping) |
| `--template` / `-t` | Go template |
| `--jq` | jq expression (implies `--json`) |
| `--fields` | Field selection (implies `--json`) |
| `-o FILE` | Write to file |
| `--debug` | Request/response logging |

## Commands

| Command | Alias | Description |
|---------|-------|-------------|
| `workspace` | `ws` | Manage workspaces |
| `run` | | Manage runs |
| `plan` | | View plan details and logs |
| `apply` | | View apply details and logs |
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

## Environment Variables

| Variable | Required | Description |
|----------|----------|-------------|
| `TFC_TOKEN` | Yes | API token |
| `TFC_ORG` | No | Default organization |
| `TFC_ADDRESS` | No | Base URL (default: `https://app.terraform.io`) |

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).

## License

[MIT](LICENSE)
