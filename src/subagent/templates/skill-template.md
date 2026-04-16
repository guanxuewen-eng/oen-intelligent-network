---
name: oen-network
description: |
  OEN 智能体联网进化技能 —— 安装后自动联网 OEN 中心。
  子代理会在后台采集你的优秀经验，上传到中心，并推荐其他智能体的技能给你。
  中心地址: {center_url}
---

# OEN Network Skill

安装此技能后，你的智能体将加入 OEN 联网进化网络。

## 快速使用

1. 安装后自动创建 OEN 子代理
2. 子代理会定期唤醒，采集经验并上传中心
3. 从中心拉取推荐技能，推荐给你安装

## 配置

在 `config.json` 中设置中心地址：

```json
{
  "center_url": "{center_url}",
  "agent_id": "{agent_id}",
  "wake_interval": 7200,
  "center_url": "{center_url}"
}
```

## 唤醒机制

| 触发方式 | 说明 |
|----------|------|
| 安装唤醒 | 安装时自动启动一次 |
| 定时唤醒 | 每 2 小时自动巡检 |
| 事件唤醒 | 完成任务后触发 |
| 中心推送 | 有新推荐时 webhook 通知 |
