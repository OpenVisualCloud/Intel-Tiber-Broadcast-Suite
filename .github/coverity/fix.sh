echo "#fix me " | cat - .github/coverity/cov_build > temp && mv temp .github/coverity/cov_build
git config --global user.name "github-actions[bot]"
git config --global user.email "github-actions[bot]@users.noreply.github.com"
git checkout -b cov-build-fix
git add .github/coverity/cov_build
git commit -m "Add fix comment to .github/coverity/cov_build"
git push origin cov-build-fix
if gh pr list --search "cov-build fix" --json title | grep -q '"title": "cov-build fix"'; then
  echo "PR with title 'cov-build fix' already exists. Skipping creation."
else
  gh pr create --title "cov-build fix" --body "Automatically created PR to address build failure" --base main
  gh pr comment --body "CC @zLukas"
fi
