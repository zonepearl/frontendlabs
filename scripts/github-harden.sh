#!/usr/bin/env bash
# Applies this repository's GitHub protections through the REST API.
# Safe to re-run: rulesets are matched by name and updated in place.
#
#   GH_TOKEN=<token> scripts/github-harden.sh            # dry run: shows what would change
#   GH_TOKEN=<token> scripts/github-harden.sh --apply    # applies it
#   REQUIRE_SIGNED=1 GH_TOKEN=... scripts/github-harden.sh --apply   # also require signed commits
#
# Token: a fine-grained personal access token for zonepearl/frontendlabs only, with
# "Administration: read and write" (and read access to metadata). Create it, run
# this once, then delete the token. Never give agents a token with Administration.
set -euo pipefail

REPO="${REPO:-zonepearl/frontendlabs}"
API="https://api.github.com"
APPLY=0; [ "${1:-}" = "--apply" ] && APPLY=1
cd "$(dirname "$0")/.."

[ -n "${GH_TOKEN:-}" ] || { echo "GH_TOKEN is not set (see the header of this script)"; exit 2; }
command -v python3 >/dev/null || { echo "python3 is required"; exit 2; }

gh_api() {   # gh_api METHOD PATH [JSON-BODY] -> prints the response body; fails on HTTP >= 400
  local method=$1 path=$2 body=${3:-} out code
  out=$(mktemp)
  code=$(curl -sS -o "$out" -w '%{http_code}' -X "$method" \
    -H "Authorization: Bearer $GH_TOKEN" -H "Accept: application/vnd.github+json" \
    -H "X-GitHub-Api-Version: 2022-11-28" ${body:+-d "$body"} "$API$path")
  if [ "$code" -ge 400 ]; then echo "  ! $method $path -> HTTP $code: $(head -c 400 "$out")" >&2; rm -f "$out"; return 1; fi
  cat "$out"; rm -f "$out"
}

step() {    # step "description" METHOD PATH [JSON]: runs only with --apply
  local desc=$1; shift
  if [ "$APPLY" = 1 ]; then
    if gh_api "$@" >/dev/null; then echo "  ✓ $desc"; else echo "  ✗ $desc (see error above)"; FAILED=1; fi
  else
    echo "  · would: $desc"
  fi
}
FAILED=0

echo "Repository: $REPO   mode: $([ "$APPLY" = 1 ] && echo APPLY || echo 'dry run (add --apply)')"
repo_json=$(gh_api GET "/repos/$REPO") || { echo "Cannot read $REPO with this token (expired, wrong repo, or no access); nothing changed."; exit 3; }
admin=$(python3 -c 'import json,sys; print(json.loads(sys.argv[1]).get("permissions",{}).get("admin",False))' "$repo_json")
[ "$admin" = "True" ] || { echo "The token does not have admin rights on $REPO; nothing changed."; exit 3; }

echo "1. Rulesets (no bypass actors: they bind admins and agents too)"
existing=$(gh_api GET "/repos/$REPO/rulesets?per_page=100")
for file in .github/rulesets/*.json; do
  body=$(python3 - "$file" <<'PY'
import json, os, sys
r = json.load(open(sys.argv[1]))
if os.environ.get("REQUIRE_SIGNED") == "1" and r["target"] == "branch":
    r["rules"].append({"type": "required_signatures"})
print(json.dumps(r))
PY
)
  name=$(python3 -c 'import json,sys; print(json.loads(sys.argv[1])["name"])' "$body")
  id=$(python3 -c 'import json,sys; print(next((str(r["id"]) for r in json.loads(sys.argv[1]) if r["name"]==sys.argv[2]), ""))' "$existing" "$name")
  if [ -n "$id" ]; then step "update ruleset \"$name\"" PUT "/repos/$REPO/rulesets/$id" "$body"
  else step "create ruleset \"$name\"" POST "/repos/$REPO/rulesets" "$body"; fi
done

echo "2. Code security"
step "secret scanning + push protection, Dependabot security updates" PATCH "/repos/$REPO" \
  '{"security_and_analysis":{"secret_scanning":{"status":"enabled"},"secret_scanning_push_protection":{"status":"enabled"},"dependabot_security_updates":{"status":"enabled"}}}'
step "Dependabot vulnerability alerts" PUT "/repos/$REPO/vulnerability-alerts"
step "private vulnerability reporting (SECURITY.md)" PUT "/repos/$REPO/private-vulnerability-reporting"

echo "3. GitHub Actions"
step "allow only GitHub-owned actions (all workflows use actions/*)" PUT "/repos/$REPO/actions/permissions" \
  '{"enabled":true,"allowed_actions":"selected"}'
step "selected actions: GitHub-owned only" PUT "/repos/$REPO/actions/permissions/selected-actions" \
  '{"github_owned_allowed":true,"verified_allowed":false,"patterns_allowed":[]}'
step "default workflow token read-only; workflows cannot approve PRs" PUT "/repos/$REPO/actions/permissions/workflow" \
  '{"default_workflow_permissions":"read","can_approve_pull_request_reviews":false}'
step "fork pull requests need approval before workflows run" PUT "/repos/$REPO/actions/permissions/fork-pr-contributor-approval" \
  '{"approval_policy":"all_external_contributors"}'

echo "4. Merging and deployment"
step "squash/rebase merges only, delete merged branches, no auto-merge" PATCH "/repos/$REPO" \
  '{"allow_merge_commit":false,"allow_squash_merge":true,"allow_rebase_merge":true,"allow_auto_merge":false,"delete_branch_on_merge":true}'
step "github-pages environment: deploy from selected branches only" PUT "/repos/$REPO/environments/github-pages" \
  '{"deployment_branch_policy":{"protected_branches":false,"custom_branch_policies":true}}'
if [ "$APPLY" = 1 ]; then
  policies=$(gh_api GET "/repos/$REPO/environments/github-pages/deployment-branch-policies" || echo '{}')
  if python3 -c 'import json,sys; sys.exit(0 if any(p["name"]=="main" for p in json.loads(sys.argv[1]).get("branch_policies",[])) else 1)' "$policies"; then
    echo "  ✓ github-pages deploys only from main (already set)"
  else
    step "github-pages deploys only from main" POST "/repos/$REPO/environments/github-pages/deployment-branch-policies" '{"name":"main","type":"branch"}'
  fi
else
  echo "  · would: github-pages deploys only from main"
fi

if [ "$APPLY" = 1 ]; then
  echo "Active rulesets now:"
  gh_api GET "/repos/$REPO/rulesets" | python3 -c 'import json,sys; [print("  -", r["name"], "(", r["enforcement"], ")") for r in json.load(sys.stdin)]'
fi
[ "$FAILED" = 0 ] || { echo "Some steps failed (details above)."; exit 1; }
echo "Done."
