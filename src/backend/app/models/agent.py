from sqlalchemy import Column, String, DateTime, Float, Integer, Text
from sqlalchemy.sql import func
from app.database import Base


class Agent(Base):
    __tablename__ = "agents"

    id = Column(Integer, primary_key=True, autoincrement=True)
    agent_id = Column(String, unique=True, nullable=False, index=True)
    model_type = Column(String, nullable=False)
    subscription_tier = Column(String, nullable=False)  # pro / plus / api / free
    resource_info = Column(Text, nullable=True)  # JSON string
    status = Column(String, default="online")  # online / offline / idle
    cpu_usage = Column(Float, default=0.0)
    memory_usage = Column(Float, default=0.0)
    last_heartbeat = Column(DateTime(timezone=True), server_default=func.now())
    created_at = Column(DateTime(timezone=True), server_default=func.now())
    updated_at = Column(
        DateTime(timezone=True), server_default=func.now(), onupdate=func.now()
    )
