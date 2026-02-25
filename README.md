# infra-stream-sample
テンプレもとの変更反映方法
```bash
git remote add upstream https://github.com/nakamuraitsuki/infra-stream-temp.git
git fetch upstream
git merge upstream/main --allow-unrelated-histories --squash
```
