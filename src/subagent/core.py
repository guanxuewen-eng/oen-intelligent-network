"""
OEN 子代理核心模块
职责：采集、评估、上传、拉取、汇报
"""
import json
import time
import logging
from dataclasses import dataclass, field, asdict
from pathlib import Path
from typing import Optional
from enum import Enum

logger = logging.getLogger("oen.subagent")


class SubscriptionTier(str, Enum):
    PRO = "pro"
    PLUS = "plus"
    API = "api"
    FREE = "free"


class QualificationStatus(str, Enum):
    ELIGIBLE = "eligible"
    INSUFFICIENT_RESOURCE = "insufficient_resource"
    MODEL_TOO_WEAK = "model_weak"
    NOT_APPLICABLE = "not_applicable"


@dataclass
class AgentProfile:
    """大模型档案"""
    model_type: str
    subscription_tier: str
    context_window: int
    token_limit: str
    tool_calls_limit: Optional[int] = None
    rate_limit: Optional[str] = None


@dataclass
class ResourceStatus:
    """资源状态"""
    cpu_usage: float
    memory_usage: float
    api_calls_remaining: Optional[int] = None
    context_used_pct: float = 0.0


@dataclass
class QualificationResult:
    """岗位资格评估结果"""
    expert_review: QualificationStatus
    website_ops: QualificationStatus
    suggestions: list = field(default_factory=list)


@dataclass
class WakeReport:
    """唤醒后汇报给用户大模型的内容"""
    agent_profile: dict
    resource_status: dict
    qualification: dict
    ops_status: dict
    recommendations: list
    timestamp: float = 0.0

    def __post_init__(self):
        if self.timestamp == 0.0:
            self.timestamp = time.time()


class OENSubAgent:
    """OEN 子代理核心"""

    def __init__(self, center_url: str, agent_id: str, memory_dir: Optional[str] = None):
        self.center_url = center_url.rstrip("/")
        self.agent_id = agent_id
        self.memory_dir = Path(memory_dir or f"./memory/{agent_id}")
        self.memory_dir.mkdir(parents=True, exist_ok=True)
        self.state_file = self.memory_dir / "state.json"
        self.evolution_log = self.memory_dir / "evolution-log.md"
        self.learnings_dir = self.memory_dir / "learnings"
        self.learnings_dir.mkdir(exist_ok=True)

        # 初始化学习文件
        (self.learnings_dir / "LEARNINGS.md").touch(exist_ok=True)
        (self.learnings_dir / "ERRORS.md").touch(exist_ok=True)
        (self.learnings_dir / "FEATURE_REQUESTS.md").touch(exist_ok=True)

    async def on_wake(self, agent_profile: dict, resource_status: dict) -> WakeReport:
        """唤醒后完整流程"""
        logger.info(f"SubAgent waking up for {self.agent_id}")

        # 1. 检查模型和资源
        profile = AgentProfile(
            model_type=agent_profile["model_type"],
            subscription_tier=agent_profile["subscription_tier"],
            context_window=agent_profile.get("context_window", 0),
            token_limit=agent_profile.get("token_limit", "unknown"),
        )
        resources = ResourceStatus(
            cpu_usage=resource_status.get("cpu_usage", 0),
            memory_usage=resource_status.get("memory_usage", 0),
        )

        # 2. 岗位资格评估
        qualification = self.evaluate_qualification(profile, resources)

        # 3. 采集本地经验并上传
        await self.collect_and_upload()

        # 4. 拉取推荐
        recommendations = await self.pull_recommendations()

        # 5. 更新记忆
        self.update_memory(profile, resources, qualification, recommendations)

        # 6. 检查运维状态
        ops_status = await self.check_ops_status()

        # 7. 生成汇报
        report = WakeReport(
            agent_profile=asdict(profile),
            resource_status=asdict(resources),
            qualification=asdict(qualification),
            ops_status=ops_status,
            recommendations=recommendations,
        )

        logger.info(f"SubAgent wake cycle complete for {self.agent_id}")
        return report

    def evaluate_qualification(self, profile: AgentProfile, resources: ResourceStatus) -> QualificationResult:
        """评估能否参与专家评审和网站运维"""
        # 专家评审资格：模型够聪明 + 资源够
        smart_models = ["claude-sonnet-4", "claude-opus", "gpt-4", "gpt-5", "qwen-max", "qwen-plus"]
        is_smart = profile.model_type.lower() in [m.lower() for m in smart_models]
        has_enough_resource = profile.subscription_tier in ["pro", "plus", "api"]
        has_context = profile.context_window >= 100000

        if is_smart and has_enough_resource and has_context:
            expert_status = QualificationStatus.ELIGIBLE
            expert_suggestions = []
        elif not is_smart:
            expert_status = QualificationStatus.MODEL_TOO_WEAK
            expert_suggestions = ["建议升级到更强大的模型（如 Claude Sonnet/Opus, GPT-4/5）"]
        elif not has_enough_resource:
            expert_status = QualificationStatus.INSUFFICIENT_RESOURCE
            expert_suggestions = ["建议升级订阅（Pro/Plus）或改用 API 付费模式"]
        else:
            expert_status = QualificationStatus.MODEL_TOO_WEAK
            expert_suggestions = ["模型上下文窗口不足，建议升级"]

        # 网站运维资格：资源充足 + 持续在线
        if has_enough_resource and resources.cpu_usage < 80 and resources.memory_usage < 80:
            ops_status = QualificationStatus.ELIGIBLE
            ops_suggestions = []
        else:
            ops_status = QualificationStatus.INSUFFICIENT_RESOURCE
            ops_suggestions = ["资源不足，无法承担运维监督工作"]

        return QualificationResult(
            expert_review=expert_status,
            website_ops=ops_status,
            suggestions=expert_suggestions + ops_suggestions,
        )

    async def collect_and_upload(self):
        """采集本地经验并上传到中心"""
        # 扫描 learnings 目录
        learnings = self._scan_learnings()
        if learnings:
            # TODO: POST /api/v1/artifacts/upload
            logger.info(f"Found {len(learnings)} learnings to upload")

    async def pull_recommendations(self) -> list:
        """从中心拉取推荐技能"""
        # TODO: GET /api/v1/artifacts/recommend
        logger.info(f"Pulling recommendations for {self.agent_id}")
        return []

    async def check_ops_status(self) -> dict:
        """检查网站运维状态"""
        # TODO: 检查中心双版本网站和 API 健康状态
        return {
            "public_site": "unknown",
            "admin_console": "unknown",
            "api_health": "unknown",
            "last_sync": None,
        }

    def update_memory(self, profile: AgentProfile, resources: ResourceStatus, qualification: QualificationResult, recommendations: list):
        """更新子代理记忆"""
        state = {
            "agent_id": self.agent_id,
            "last_wake": time.time(),
            "model_type": profile.model_type,
            "subscription_tier": profile.subscription_tier,
            "qualification": asdict(qualification),
            "wake_count": self._get_state().get("wake_count", 0) + 1,
        }
        self._save_state(state)

    def _scan_learnings(self) -> list:
        """扫描 learnings 目录"""
        learnings = []
        for f in self.learnings_dir.glob("*.md"):
            content = f.read_text().strip()
            if content and not content.startswith("# "):
                learnings.append({"file": f.name, "content": content})
        return learnings

    def _get_state(self) -> dict:
        if self.state_file.exists():
            return json.loads(self.state_file.read_text())
        return {}

    def _save_state(self, state: dict):
        self.state_file.write_text(json.dumps(state, indent=2, ensure_ascii=False))
