from sqlalchemy import Column, String, DateTime, Integer, Text, Float
from sqlalchemy.sql import func
from app.database import Base


class Artifact(Base):
    __tablename__ = "artifacts"

    id = Column(Integer, primary_key=True, autoincrement=True)
    title = Column(String, nullable=False)
    description = Column(Text, nullable=True)
    content = Column(Text, nullable=False)
    artifact_type = Column(String, nullable=False)  # skill / plan / template
    author_agent_id = Column(String, nullable=True)
    review_status = Column(String, default="pending")  # pending / approved / rejected
    score = Column(Float, default=0.0)  # recommendation score
    tags = Column(Text, nullable=True)  # JSON string
    created_at = Column(DateTime(timezone=True), server_default=func.now())
    updated_at = Column(
        DateTime(timezone=True), server_default=func.now(), onupdate=func.now()
    )
