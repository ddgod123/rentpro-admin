# 代码同步脚本实现

## 功能概述

为了简化前后端代码的版本管理和远程仓库同步，实现了两个自动化同步脚本：
1. **完整版同步脚本** (`sync_repositories.sh`) - 功能全面的智能同步工具
2. **快速同步脚本** (`quick_sync.sh`) - 简化版本，适合日常快速使用

## 需求背景

### 开发痛点
- 前后端项目需要分别进行Git操作
- 手动同步容易遗漏项目或出错
- 提交信息不规范，缺乏统一性
- 需要重复执行相同的Git命令序列

### 解决方案
- 自动化Git操作流程
- 统一管理前后端项目同步
- 标准化提交信息格式
- 提供错误检查和安全保障

## 实现详情

### 1. 完整版同步脚本 (sync_repositories.sh)

#### 核心功能
```bash
#!/bin/bash
# 主要功能模块
- 环境检查 (网络、目录、Git仓库)
- 状态检测 (分支、提交、未推送更改)
- 智能提交 (交互式提交信息输入)
- 安全推送 (拉取合并 + 推送)
- 详细日志 (彩色输出、执行统计)
```

#### 技术特性
1. **彩色日志输出**
```bash
# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}
```

2. **智能路径检测**
```bash
# 自动获取脚本所在目录
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BACKEND_DIR="${SCRIPT_DIR}/rentpro-admin-main"
FRONTEND_DIR="${SCRIPT_DIR}/rent-foren"
```

3. **全面状态检查**
```bash
check_uncommitted_changes() {
    # 检查工作区更改
    if ! git diff --quiet; then
        return 1
    fi
    
    # 检查暂存区更改
    if ! git diff --cached --quiet; then
        return 1
    fi
    
    # 检查未跟踪文件
    local untracked=$(git ls-files --others --exclude-standard)
    if [ -n "$untracked" ]; then
        return 1
    fi
    
    return 0
}
```

4. **交互式提交**
```bash
commit_changes() {
    # 显示更改摘要
    git status --porcelain | head -10
    
    # 提示用户输入提交信息
    echo "请输入提交信息 (直接回车使用默认信息):"
    echo "默认: $default_message"
    read -r commit_message
    
    if [ -z "$commit_message" ]; then
        commit_message=$default_message
    fi
    
    git add .
    git commit -m "$commit_message"
}
```

5. **安全推送机制**
```bash
push_to_remote() {
    # 先拉取避免冲突
    if git pull origin "$branch" --rebase; then
        print_info "✓ 成功拉取远程更改"
    else
        print_error "拉取失败，可能有冲突"
        return 1
    fi
    
    # 推送到远程
    git push origin "$branch"
}
```

#### 命令行参数支持
```bash
# 显示帮助信息
./sync_repositories.sh --help

# 仅检查状态，不执行同步
./sync_repositories.sh --dry-run

# 正常执行
./sync_repositories.sh
```

### 2. 快速同步脚本 (quick_sync.sh)

#### 简化流程
```bash
#!/bin/bash
# 快速同步核心逻辑

# 生成时间戳提交信息
TIMESTAMP=$(date '+%Y-%m-%d %H:%M:%S')

# 同步后端
cd "$BACKEND_DIR"
git add .
git commit -m "feat: 后端功能更新 - $TIMESTAMP" || true
git pull origin main --rebase
git push origin main

# 同步前端
cd "$FRONTEND_DIR"
git add .
git commit -m "feat: 前端功能更新 - $TIMESTAMP" || true
git pull origin main --rebase
git push origin main
```

#### 设计原则
- **简洁高效**: 最少的交互，最快的执行
- **错误容忍**: 即使部分命令失败也继续执行
- **标准化**: 使用统一的提交信息格式

## 文件结构

```
/Users/mac/go/src/rentPro/houduan/
├── sync_repositories.sh     # 完整版同步脚本
├── quick_sync.sh            # 快速同步脚本
├── README_SYNC.md           # 使用说明文档
├── rentpro-admin-main/      # 后端项目目录
└── rent-foren/              # 前端项目目录
```

## 使用方式

### 1. 日常开发同步
```bash
# 快速同步，适合日常提交
./quick_sync.sh
```

### 2. 重要功能发布
```bash
# 完整同步，可自定义提交信息
./sync_repositories.sh
```

### 3. 状态检查
```bash
# 仅查看状态，不执行同步
./sync_repositories.sh --dry-run
```

## 执行效果示例

### 完整版脚本输出
```
========================================
租房管理系统 - 代码同步脚本
时间: 2025-09-05 15:50:23
========================================
[INFO] 检查网络连接...
[SUCCESS] ✓ 网络连接正常
========================================
[INFO] 开始同步 后端项目...
[INFO] ✓ 后端项目 目录存在
[INFO] ✓ 后端项目 是Git仓库
[INFO] 后端项目 状态:
  分支: main
  最新提交: 61de46e
  提交信息: fix: 修复认证路由中的GORM模型错误
[INFO] ✓ 后端项目 工作区干净
[SUCCESS] ✓ 后端项目 成功推送到远程仓库
[SUCCESS] 🎉 所有项目同步完成!
[INFO] 总耗时: 15秒
```

### 快速脚本输出
```
开始同步前后端代码...
[1/2] 同步后端项目...
✓ 后端同步完成
[2/2] 同步前端项目...
✓ 前端同步完成
🎉 所有项目同步完成!
```

## 安全特性

### 1. 环境检查
- **网络连接检查**: 确保能连接GitHub
- **目录存在性检查**: 确保项目目录存在
- **Git仓库检查**: 确保是有效的Git仓库

### 2. 冲突预防
- **拉取后推送**: 先拉取远程更改再推送
- **Rebase策略**: 使用rebase避免不必要的合并提交
- **错误处理**: 遇到冲突时提供明确的错误信息

### 3. 数据保护
- **非破坏性操作**: 不会删除或覆盖本地更改
- **状态检查**: 提供详细的仓库状态信息
- **回滚能力**: 支持Git的标准回滚机制

## 提交信息规范

### 自动生成格式
```
feat: 后端功能更新 - 2025-09-05 15:50:23
feat: 前端功能更新 - 2025-09-05 15:50:23
```

### 自定义格式示例
```
feat: 完善楼盘管理功能
fix: 修复认证系统JWT配置问题
refactor: 优化前端组件结构
docs: 更新API文档
```

## 扩展功能

### 1. 支持的Git操作
- 自动添加所有更改 (`git add .`)
- 智能提交管理 (`git commit`)
- 安全推送策略 (`git pull --rebase` + `git push`)

### 2. 日志管理
- 彩色输出，提升可读性
- 分级日志 (INFO, SUCCESS, WARNING, ERROR)
- 执行时间统计

### 3. 错误处理
- 网络连接失败处理
- Git操作失败处理
- 用户输入验证

## 性能优化

### 1. 并行处理
虽然当前是顺序处理前后端项目，但脚本结构支持后续改为并行处理。

### 2. 增量检查
只有在检测到实际更改时才执行Git操作，避免不必要的网络请求。

### 3. 智能缓存
利用Git自身的缓存机制，减少重复的状态检查。

## 维护说明

### 路径配置
如需修改项目路径，编辑脚本开头的路径变量：
```bash
BACKEND_DIR="${SCRIPT_DIR}/your-backend-project"
FRONTEND_DIR="${SCRIPT_DIR}/your-frontend-project"
```

### 远程仓库配置
脚本假设远程仓库名为 `origin`，分支为 `main`。如需修改：
```bash
git push your-remote-name your-branch-name
```

### 日志级别调整
可以通过修改打印函数来调整日志输出级别。

## 相关文件

**脚本文件**:
- `sync_repositories.sh` - 完整版同步脚本
- `quick_sync.sh` - 快速同步脚本

**文档文件**:
- `README_SYNC.md` - 详细使用说明
- 本文档 - 技术实现说明

**项目目录**:
- `rentpro-admin-main/` - 后端项目
- `rent-foren/` - 前端项目

## 使用建议

1. **日常开发**: 使用快速脚本 (`quick_sync.sh`)
2. **重要更新**: 使用完整脚本并自定义提交信息
3. **状态检查**: 使用 `--dry-run` 参数查看当前状态
4. **故障排查**: 查看详细的错误日志和Git状态

这套同步脚本大大简化了前后端项目的版本管理流程，提高了开发效率并减少了操作错误。
