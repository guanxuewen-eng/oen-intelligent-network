from pydantic import BaseModel
from typing import Optional
from datetime import datetime


class MemoryStore(BaseModel):
    agent_id: str
    content: str
    memory_type: Optional[str] = None
    extra_meta: Optional[str] = None


class MemoryResponse(BaseModel):
    id: int
    agent_id: str
    content: str
    memory_type: Optional[str]
    extra_meta: Optional[str]
    created_at: Optional[datetime]
    updated_at: Optional[datetime]

    model_config = {"from_attributes": True}
