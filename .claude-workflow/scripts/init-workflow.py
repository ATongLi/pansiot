#!/usr/bin/env python3
"""
Claude Code Workflow Initialization Script

This script initializes a new project with the Claude Code workflow structure.
It creates all necessary directories, copies templates, and sets up initial configuration.
"""

import os
import sys
import json
import shutil
from pathlib import Path
from datetime import datetime

class WorkflowInitializer:
    def __init__(self, project_path="."):
        self.project_path = Path(project_path).resolve()
        self.workflow_dir = self.project_path / ".claude-workflow"
        self.date = datetime.now().strftime("%Y-%m-%d")

    def create_directories(self):
        """Create all required directories"""
        print("📁 Creating directory structure...")

        dirs = [
            "01-requirements/raw-requirements",
            "01-requirements/functional-requirements",
            "01-requirements/user-stories",
            "02-planning",
            "03-design/architecture/adr",
            "03-design/technical-solutions",
            "03-design/api-design",
            "03-design/database-design",
            "04-implementation/implementation-plans",
            "04-implementation/code-mapping",
            "04-implementation/implementation-logs",
            "05-verification/test-plans",
            "05-verification/test-cases",
            "05-verification/verification-reports",
            "06-documentation/technical-docs",
            "06-documentation/user-guides",
            "06-documentation/api-docs",
            "07-platforms/gateway/requirements",
            "07-platforms/gateway/design",
            "07-platforms/gateway/implementation",
            "07-platforms/hmi/requirements",
            "07-platforms/hmi/design",
            "07-platforms/hmi/implementation",
            "07-platforms/configuration/requirements",
            "07-platforms/configuration/design",
            "07-platforms/configuration/implementation",
            "07-platforms/cloud/requirements",
            "07-platforms/cloud/design",
            "07-platforms/cloud/implementation",
            "07-platforms/app/requirements",
            "07-platforms/app/design",
            "07-platforms/app/implementation",
            "07-platforms/edge-ai/requirements",
            "07-platforms/edge-ai/design",
            "07-platforms/edge-ai/implementation",
            "07-platforms/scada/requirements",
            "07-platforms/scada/design",
            "07-platforms/scada/implementation",
            "07-platforms/web-editor/requirements",
            "07-platforms/web-editor/design",
            "07-platforms/web-editor/implementation",
            "parallel-tasks",
            "templates",
            "state-history",
        ]

        for dir_path in dirs:
            full_path = self.workflow_dir / dir_path
            full_path.mkdir(parents=True, exist_ok=True)

        print("✅ Directory structure created")

    def copy_templates(self, template_source):
        """Copy template files from source"""
        print("📄 Copying template files...")

        template_dir = self.workflow_dir / "templates"
        templates = ["REQ-template.md", "FE-template.md", "SOL-template.md", "IMP-template.md"]

        for template in templates:
            source = Path(template_source) / ".claude-workflow" / "templates" / template
            if source.exists():
                shutil.copy2(source, template_dir / template)
                print(f"  ✅ Copied {template}")
            else:
                print(f"  ⚠️  Template not found: {template}")

        # Copy state templates
        state_templates = [
            "parallel-tasks/active-tasks-template.md",
            "parallel-tasks/task-dependencies-template.md",
            "dependency-backlog-template.md",
            "feature-to-code-map-template.md",
            "rt-matrix-template.md",
        ]

        for template in state_templates:
            source = Path(template_source) / ".claude-workflow" / template
            if source.exists():
                dest = self.workflow_dir / template.replace("-template", "")
                shutil.copy2(source, dest)
                print(f"  ✅ Copied {template}")

        print("✅ Template files copied")

    def create_initial_config(self, config_data):
        """Create initial configuration file"""
        print("⚙️  Creating configuration...")

        config_file = self.workflow_dir / "config.yml"

        default_config = {
            "project": {
                "name": config_data.get("project_name", "My Project"),
                "type": config_data.get("project_type", "application"),
                "version": "1.0.0",
            },
            "platforms": [
                {"id": "gateway", "name": "网关端", "enabled": config_data.get("gateway", False)},
                {"id": "hmi", "name": "HMI运行端", "enabled": config_data.get("hmi", False)},
                {"id": "configuration", "name": "组态端", "enabled": config_data.get("configuration", False)},
                {"id": "cloud", "name": "云平台端", "enabled": config_data.get("cloud", False)},
                {"id": "app", "name": "APP端", "enabled": config_data.get("app", False)},
                {"id": "edge-ai", "name": "边缘智能服务器", "enabled": config_data.get("edge-ai", False)},
                {"id": "scada", "name": "Scada软件", "enabled": config_data.get("scada", False)},
                {"id": "web-editor", "name": "Web可视化编辑器", "enabled": config_data.get("web-editor", False)},
            ],
            "workflow": {
                "enforce_order": True,
                "enable_parallel": True,
                "max_parallel_tasks": 5,
                "enable_dependency_ignore": True,
            },
            "documentation": {
                "auto_generate": True,
                "format": "markdown",
                "include_api_docs": True,
            },
            "traceability": {
                "auto_update": True,
                "mapping_granularity": "function",
            },
        }

        # Write YAML config
        yaml_content = """# Claude Code 工作流配置

project:
  name: "{name}"
  type: "{type}"
  version: "{version}"

# 支持的平台
platforms:
""".format(
            name=default_config["project"]["name"],
            type=default_config["project"]["type"],
            version=default_config["project"]["version"],
        )

        for platform in default_config["platforms"]:
            yaml_content += f"""  - id: {platform['id']}
    name: {platform['name']}
    enabled: {str(platform['enabled']).lower()}
"""

        yaml_content += """
# 工作流配置
workflow:
  enforce_order: true
  enable_parallel: true
  max_parallel_tasks: 5
  enable_dependency_ignore: true

# 文档配置
documentation:
  auto_generate: true
  format: markdown
  include_api_docs: true

# 追溯配置
traceability:
  auto_update: true
  mapping_granularity: function
"""

        with open(config_file, "w", encoding="utf-8") as f:
            f.write(yaml_content)

        print("✅ Configuration created")

    def create_initial_state(self):
        """Create initial state files"""
        print("📊 Creating initial state files...")

        # current-phase.md
        current_phase = self.workflow_dir / "current-phase.md"
        current_phase.write_text(
            f"""# 当前执行阶段

## 项目信息
- **项目名称**: {self.workflow_dir.parent.name}
- **当前阶段**: 阶段1 - 需求分析与规划
- **开始日期**: {self.date}
- **整体进度**: 0%

## 当前任务
- **任务ID**: 待定
- **任务名称**: 项目初始化
- **状态**: 进行中
- **当前步骤**: 配置项目
- **完成度**: 10%

## 任务堆栈
### 主任务链
1. ⏳ 项目初始化
   - ✅ 创建目录结构
   - ✅ 复制模板文件
   - ⏳ 配置项目参数

### 暂停的任务
无

### 待执行任务
- [待执行] 收集第一个需求

## 上下文信息
- **相关平台**: 待配置
- **相关文件**: .claude-workflow/config.yml
- **依赖项**: 无
- **阻塞项**: 无

## 快速恢复命令
"继续配置项目参数"
""",
            encoding="utf-8"
        )

        # active-tasks.md
        active_tasks = self.workflow_dir / "parallel-tasks" / "active-tasks.md"
        active_tasks.write_text(
            f"""# 活跃并行任务列表

更新时间: {self.date} 00:00

## 正在进行的功能

### 项目初始化
- **当前阶段**: 配置
- **完成度**: 10%
- **当前步骤**: 配置项目参数
- **状态**: 🟢 正常进行
- **依赖**: 无

## 并行统计
- **活跃功能数**: 0
- **分布阶段**: 配置(1)
- **预计完成日期**: 待定

## 阻塞警告
无
""",
            encoding="utf-8"
        )

        # task-dependencies.md
        task_deps = self.workflow_dir / "parallel-tasks" / "task-dependencies.md"
        task_deps.write_text(
            """# 任务依赖关系图

## 依赖关系可视化

```
无依赖关系
```

## 依赖类型定义

| 类型 | 符号 | 说明 | 处理方式 |
|------|------|------|---------|
| 强依赖 | ═══ | 必须等待依赖完成 | 串行开发 |
| 弱依赖 | ━━ | 可以并行，使用忽略机制 | 并行 + mock |
| 无依赖 | 无 | 完全独立 | 自由并行 |

## 依赖忽略记录
无

## 依赖关系表
| 当前功能 | 依赖功能 | 依赖类型 | 处理方式 | 状态 |
""",
            encoding="utf-8"
        )

        # dependency-backlog.md
        dep_backlog = self.workflow_dir / "dependency-backlog.md"
        dep_backlog.write_text(
            """# 依赖忽略跟踪表

## 待补齐依赖项

| ID | 当前模块 | 依赖模块 | 忽略内容 | Mock位置 | 补齐优先级 | 依赖状态 | 补齐期限 | 负责人 |

## 依赖补齐记录
无

## 统计信息

- **待补齐依赖数**: 0
- **P0 优先级**: 0
- **P1 优先级**: 0
- **P2 优先级**: 0
- **本月已补齐**: 0
""",
            encoding="utf-8"
        )

        # feature-to-code-map.md
        code_map = self.workflow_dir / "feature-to-code-map.md"
        code_map.write_text(
            """# 功能到代码映射表

## 映射规则
- 每个功能点对应具体的文件、类、函数
- 记录行号范围以便快速定位
- 包含依赖关系和影响范围

## 映射表

| 功能ID | 功能描述 | 平台 | 代码文件 | 类/组件 | 函数/方法 | 行号 | 依赖 | 状态 |

## 依赖关系图
无

## 影响分析查询
无

## 反向映射 (代码 → 功能)

| 代码文件 | 函数 | 实现功能 | 功能ID | 状态 |

## 平台分布

| 平台 | 功能数 | 已实现 | 进行中 | 待开始 |
|------|--------|--------|--------|--------|
| 网关 | 0 | 0 | 0 | 0 |
| HMI | 0 | 0 | 0 | 0 |
| 云平台 | 0 | 0 | 0 | 0 |
| APP | 0 | 0 | 0 | 0 |
| 边缘AI | 0 | 0 | 0 | 0 |
| Scada | 0 | 0 | 0 | 0 |
| Web编辑器 | 0 | 0 | 0 | 0 |

## 统计信息

- **总功能点数**: 0
- **已实现**: 0 (0%)
- **进行中**: 0 (0%)
- **待开始**: 0 (0%)
- **依赖忽略**: 0
""",
            encoding="utf-8"
        )

        # rt-matrix.md
        rt_matrix = self.workflow_dir / "rt-matrix.md"
        rt_matrix.write_text(
            f"""# 需求追溯矩阵 (Requirements Traceability Matrix)

## 追溯概览
- **总需求数**: 0
- **已实现**: 0 (0%)
- **已验证**: 0 (0%)
- **未开始**: 0 (0%)

## 前向追溯 (需求 → 交付物)

| 需求ID | 需求描述 | 功能需求 | 技术方案 | 实现计划 | 代码模块 | 测试用例 | 状态 |

## 后向追溯 (交付物 → 需求)

| 代码文件 | 功能描述 | 关联需求 | 测试覆盖 | 状态 |

## 影响分析矩阵

| 需求变更 | 影响功能 | 影响代码 | 影响测试 | 影响平台 | 风险等级 |

## 测试覆盖率
无

## 平台覆盖

| 平台 | 需求数 | 已实现 | 已验证 | 覆盖率 |
|------|--------|--------|--------|--------|
| 网关 | 0 | 0 | 0 | 0% |
| HMI | 0 | 0 | 0 | 0% |
| 云平台 | 0 | 0 | 0 | 0% |
| APP | 0 | 0 | 0 | 0% |
| 边缘AI | 0 | 0 | 0 | 0% |
| Scada | 0 | 0 | 0 | 0% |
| Web编辑器 | 0 | 0 | 0 | 0% |

## 缺口分析
无

## 追溯完整性检查
待初始化

## 变更历史

| 日期 | 变更类型 | 影响范围 | 变更人 | 审批人 |
| {self.date} | 项目初始化 | 全部 | System | - |

## 报告生成

- **最后更新**: {self.date}
- **更新人**: System
- **下次审查**: 待定
""",
            encoding="utf-8"
        )

        print("✅ Initial state files created")

    def create_readme(self):
        """Create README for workflow directory"""
        readme = self.workflow_dir / "README.md"

        readme.write_text(
            """# Claude Code 工作流管理

本目录使用 Claude Code 结构化研发工作流管理体系进行管理。

## 目录结构

```
.claude-workflow/
├── 01-requirements/        # 需求管理
├── 02-planning/           # 计划管理
├── 03-design/             # 设计管理
├── 04-implementation/      # 实现管理
├── 05-verification/       # 验证管理
├── 06-documentation/      # 文档管理
├── 07-platforms/          # 多平台管理
├── parallel-tasks/        # 并行任务管理
├── templates/             # 文档模板
├── state-history/         # 状态历史
├── config.yml             # 项目配置
├── current-phase.md       # 当前阶段
├── feature-to-code-map.md # 功能代码映射
├── dependency-backlog.md  # 依赖跟踪
└── rt-matrix.md           # 需求追溯矩阵
```

## 使用方式

### 开始新功能
```
"我们需要添加[新功能]"
→ workflow-orchestrator 自动激活
→ 跟随引导完成需求收集、方案设计、实现...
```

### 查看当前状态
```
"我上次做到哪了？"
→ workflow-coordinator 显示当前任务和进度
```

### 暂停和恢复
```
"暂停当前任务，处理其他任务"
→ 任务切换和上下文保存
"继续之前的任务"
→ 恢复到之前的任务
```

## 工作流程

1. **需求阶段** (requirement-manager)
   - 创建原始需求 (REQ-{N})
   - 分解功能需求 (FE-{N})

2. **设计阶段** (solution-designer)
   - 创建技术方案 (SOL-{N})
   - 设计架构和接口

3. **实现计划** (implementation-manager)
   - 创建实现计划 (IMP-{N})
   - 设计代码映射

4. **代码实现** (implementation-manager)
   - 按计划实现
   - 记录代码映射

5. **验证和文档** (verification-manager, documentation-generator)
   - 测试验证
   - 生成文档

## 更多信息

参考完整文档: `Claude-Code-Workflow-Design-v1.0.0.md`
""",
            encoding="utf-8"
        )

    def initialize(self, config_data=None, template_source=None):
        """Run complete initialization"""
        print("=" * 60)
        print("🚀 Claude Code Workflow Initialization")
        print("=" * 60)
        print()

        if config_data is None:
            config_data = {}

        self.create_directories()
        print()

        if template_source:
            self.copy_templates(template_source)
            print()
        else:
            print("⚠️  跳过模板复制（未指定模板源）")
            print()

        self.create_initial_config(config_data)
        print()

        self.create_initial_state()
        print()

        self.create_readme()
        print()

        print("=" * 60)
        print("✅ 初始化完成！")
        print("=" * 60)
        print()
        print("📝 下一步:")
        print("1. 审查配置文件: .claude-workflow/config.yml")
        print("2. 启动 Claude Code")
        print("3. 开始第一个需求: '我们需要添加[新功能]'")
        print()


def main():
    """Main entry point"""
    import argparse

    parser = argparse.ArgumentParser(
        description="Initialize Claude Code workflow structure"
    )
    parser.add_argument(
        "project_path",
        nargs="?",
        default=".",
        help="Project path (default: current directory)"
    )
    parser.add_argument(
        "--project-name",
        help="Project name"
    )
    parser.add_argument(
        "--project-type",
        default="application",
        help="Project type (default: application)"
    )
    parser.add_argument(
        "--platforms",
        help="Comma-separated list of platforms (gateway,hmi,cloud,app,edge-ai,scada,web-editor)"
    )
    parser.add_argument(
        "--template-source",
        help="Path to template source (claude-workflow-template directory)"
    )

    args = parser.parse_args()

    # Parse platforms
    config_data = {
        "project_name": args.project_name,
        "project_type": args.project_type,
    }

    if args.platforms:
        platforms = args.platforms.lower().split(",")
        for platform in platforms:
            if platform in ["gateway", "hmi", "configuration", "cloud", "app", "edge-ai", "scada", "web-editor"]:
                config_data[platform] = True

    # Initialize
    initializer = WorkflowInitializer(args.project_path)
    initializer.initialize(config_data, args.template_source)


if __name__ == "__main__":
    main()
