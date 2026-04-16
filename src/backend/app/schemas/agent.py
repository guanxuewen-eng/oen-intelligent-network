from pydantic import BaseModel
from typing import Optional
from datetime import datetime


class AgentRegister(BaseModel):
    agent_id: str
    model_type: str
    subscription_tier: str  # pro / plus / api / free
    resource_info: Optional[str] = None


class AgentHeartbeat(BaseModel):
    cpu_usage: Optional[float] = None
    memory_usage: Optional[float] = None
    status: Optional[str] = None
    qualification: Optional[str] = None  # 岗位资格评估结果


class AgentStatusResponse(BaseModel):
    id: int
    agent_id: str
    model_type: str
    subscription_tier: str
    resource_info: Optional[str]
    status: str
    cpu_usage: Optional[float]
    memory_usage: Optional[float]
    last_heartbeat: Optional[datetime]
    created_at: Optional[datetime]
    updated_at: Optional[datetime]

    model_config = {"from_attributes": True}


class AgentRegisterResponse(BaseModel):
    id: int
    agent_id: str
    model_type: str
    subscription_tier: str
    resource_info: Optional[str]
    status: str
    created_at: Optional[datetime]

    model_config = {"from_attributes": True}
