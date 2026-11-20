git apply vgenv-full.patch

当前将这套改动提交到独立分支（如 feat-vgenv）：
bash
git checkout -b feat-vgenv
git add .
git commit -m "feat: 支持 vgenv 视频模型声明式插件与全能参考"


后续当上游开源仓库有代码更新时，只需要两步：
bash
# 1. 切回 main 更新上游代码
git checkout main
git pull origin main
# 2. 切回您的功能分支变基（Rebase）
git checkout feat-vgenv
git rebase main