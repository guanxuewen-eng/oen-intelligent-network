# ClawHub Self-Evolve 技能分析

## 工作机制

核心循环：4 步滴答巡航（Tick Handler）
唤醒/触发 → Step 1 扫描 → Step 2 记录 → Step 3 评估固化 → Step 4 启动新实验 → 休眠

| 步骤 | 做什么 | 关键文件 |
|------|--------|---------|
| Step 1: 状态同步 | 读 state.json，看有哪些实验在跑 | memory/evolve/state.json |
| Step 2: 记录观察 | 用 telemetry_hook 采集数据，写入 JSONL | memory/evolve/[cycle_id].jsonl |
| Step 3: 评估固化 | 到期对比基线，选赢家，写入配置文件 | AGENTS.md / TOOLS.md / evolution-log.md |
| Step 4: 启动新实验 | 找新瓶颈 → 搜索方案 → 设计实验 → 注册状态 | memory/evolve/candidates.md |

## 关键设计亮点

1. 防伪进化机制 —— 禁止纯文本排版，必须有物理变更
2. 并发状态机 —— 最多 10 个正交实验同时观察
3. 遥测钩子（Telemetry Hook）—— 标准化数据采集
4. 分级自主权 —— 绿色直接做、黄色做了通知、红色先问再做
5. 降噪防护 —— 无新数据不写日志
6. 熔断机制 —— 极端负向反应直接终止
7. 候选池 —— candidates.md 集中管理待进化方向

## OEN 借鉴方案

| ClawHub 机制 | OEN 对应设计 |
|-------------|-------------|
| 4 步滴答巡航 | 唤醒后：采集→上传→拉取→汇报→休眠 |
| telemetry_hook | 采集用户操作数据、技能执行结果 |
| state.json 状态机 | OEN 子代理独立记忆体系 |
| candidates.md 候选池 | 中心推荐引擎 + 本地技能需求池 |
| 防伪进化红线 | 必须有实际技能上传/安装 |
| 分级自主权 | 上传需确认、安装推荐需确认 |
| evolution-log.md | OEN 进化历程持久化 |

关键差异：ClawHub 是单机自我进化，OEN 是联网进化——采集后上传中心，中心汇聚全网经验再分发。
