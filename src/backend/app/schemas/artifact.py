from pydantic import BaseModel
from typing import Optional
from datetime import datetime


class ArtifactUpload(BaseModel):
    title: str
    description: Optional[str] = None
    content: str
    artifact_type: str  # skill / plan / template
    author_agent_id: Optional[str] = None
    tags: Optional[str] = None


class ArtifactResponse(BaseModel):
    id: int
    title: str
    description: Optional[str]
    content: str
    artifact_type: str
    author_agent_id: Optional[str]
    review_status: str
    score: float
    tags: Optional[str]
    created_at: Optional[datetime]
    updated_at: Optional[datetime]

    model_config = {"from_attributes": True}
