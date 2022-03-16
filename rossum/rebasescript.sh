#!/bin/bash
set -e

LATEST_UPSTREAM_RELEASE=$(curl -s https://api.github.com/repos/hashicorp/terraform-provider-aws/releases/latest | jq -r .tag_name)

(git remote | grep -q hashicorp) || git remote add hashicorp git@github.com:hashicorp/terraform-provider-aws.git

git fetch --tags hashicorp
git fetch --tags origin

git checkout main
git reset --hard "$LATEST_UPSTREAM_RELEASE"

for branch in fix-shared-aws_db_snapshot add-resource-db-snapshot-copy rossum; do
  git checkout "$branch"
  git rebase --onto "$LATEST_UPSTREAM_RELEASE" hashicorp/main
  git push -f --set-upstream origin "$branch"
  git checkout main
  git merge "$branch" --no-edit
done

cd "$(git rev-parse --show-toplevel)"
rm .github/workflows/*
cp rossum/release.yml .github/workflows/
cp rossum/README.md rossum/.goreleaser.yml .
git add .
git commit -m "Merge PRs and fix workflows"
git push -f --set-upstream origin main

git tag "$LATEST_UPSTREAM_RELEASE-1"
git push --tags

