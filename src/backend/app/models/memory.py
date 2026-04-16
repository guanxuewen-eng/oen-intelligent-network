from sqlalchemy import Column, String, DateTime, Integer, Text
from sqlalchemy.sql import func
from app.database import Base


class Memory(Base):
    __tablename__ = "memories"

    id = Column(Integer, primary_key=True, autoincrement=True)
    agent_id = Column(String, nullable=False, index=True)
    content = Column(Text, nullable=False)
    memory_type = Column(String, nullable=True)  # short / long / episodic
    extra_meta = Column("extra_meta", Text, nullable=True)  # JSON string
    created_at = Column(DateTime(timezone=True), server_default=func.now())
    updated_at = Column(
        DateTime(timezone=True), server_default=func.now(), onupdate=func.now()
    )
