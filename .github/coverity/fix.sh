echo "#fix me " | cat - .github/coverity_cov_build > temp && mv temp .github/coverity_cov_build
git config --global user.name "github-actions[bot]"
git config --global user.email "github-actions[bot]@users.noreply.github.com"
git checkout -b cov-build-fix
git add .github/coverity_cov_build
git commit -m "Add fix comment to .github/coverity_cov_build"
git push origin cov-build-fix
gh pr create --title "cov-build fix" --body "Automatically created PR to address build failure" --base main