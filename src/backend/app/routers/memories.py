from fastapi import APIRouter, Depends, HTTPException
from sqlalchemy import select, desc
from sqlalchemy.ext.asyncio import AsyncSession
from app.database import get_db
from app.models.memory import Memory
from app.models.agent import Agent
from app.schemas.memory import MemoryStore, MemoryResponse

router = APIRouter(prefix="/memories", tags=["memories"])


def ok(data, message="ok"):
    return {"code": 0, "data": data, "message": message}


@router.post("/store", response_model=dict)
async def store_memory(payload: MemoryStore, db: AsyncSession = Depends(get_db)):
    # Validate agent exists
    result = await db.execute(select(Agent).where(Agent.agent_id == payload.agent_id))
    if not result.scalar_one_or_none():
        raise HTTPException(status_code=404, detail="Agent not found")

    memory = Memory(
        agent_id=payload.agent_id,
        content=payload.content,
        memory_type=payload.memory_type,
        extra_meta=payload.extra_meta,
    )
    db.add(memory)
    await db.commit()
    await db.refresh(memory)

    return ok(MemoryResponse.model_validate(memory).model_dump())


@router.get("/{agent_id}", response_model=dict)
async def get_memories(
    agent_id: str, limit: int = 50, db: AsyncSession = Depends(get_db)
):
    result = await db.execute(
        select(Memory)
        .where(Memory.agent_id == agent_id)
        .order_by(desc(Memory.created_at))
        .limit(limit)
    )
    memories = result.scalars().all()
    return ok([MemoryResponse.model_validate(m).model_dump() for m in memories])
