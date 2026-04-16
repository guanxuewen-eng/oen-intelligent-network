from fastapi import APIRouter, Depends, HTTPException
from sqlalchemy import select, update
from sqlalchemy.ext.asyncio import AsyncSession
from sqlalchemy.sql import func
from app.database import get_db
from app.models.agent import Agent
from app.schemas.agent import (
    AgentRegister,
    AgentHeartbeat,
    AgentStatusResponse,
    AgentRegisterResponse,
)

router = APIRouter(prefix="/agents", tags=["agents"])


def ok(data, message="ok"):
    return {"code": 0, "data": data, "message": message}


@router.post("/register", response_model=dict)
async def register_agent(payload: AgentRegister, db: AsyncSession = Depends(get_db)):
    # Check if agent already exists
    result = await db.execute(select(Agent).where(Agent.agent_id == payload.agent_id))
    existing = result.scalar_one_or_none()
    if existing:
        raise HTTPException(status_code=409, detail="Agent already registered")

    agent = Agent(
        agent_id=payload.agent_id,
        model_type=payload.model_type,
        subscription_tier=payload.subscription_tier,
        resource_info=payload.resource_info,
    )
    db.add(agent)
    await db.commit()
    await db.refresh(agent)

    return ok(AgentRegisterResponse.model_validate(agent).model_dump())


@router.get("/{agent_id}/status", response_model=dict)
async def get_agent_status(agent_id: str, db: AsyncSession = Depends(get_db)):
    result = await db.execute(select(Agent).where(Agent.agent_id == agent_id))
    agent = result.scalar_one_or_none()
    if not agent:
        raise HTTPException(status_code=404, detail="Agent not found")

    return ok(AgentStatusResponse.model_validate(agent).model_dump())


@router.post("/{agent_id}/heartbeat", response_model=dict)
async def agent_heartbeat(
    agent_id: str, payload: AgentHeartbeat, db: AsyncSession = Depends(get_db)
):
    result = await db.execute(select(Agent).where(Agent.agent_id == agent_id))
    agent = result.scalar_one_or_none()
    if not agent:
        raise HTTPException(status_code=404, detail="Agent not found")

    update_data = {"last_heartbeat": func.now()}
    if payload.cpu_usage is not None:
        update_data["cpu_usage"] = payload.cpu_usage
    if payload.memory_usage is not None:
        update_data["memory_usage"] = payload.memory_usage
    if payload.status is not None:
        update_data["status"] = payload.status

    await db.execute(
        update(Agent).where(Agent.id == agent.id).values(**update_data)
    )
    await db.commit()
    await db.refresh(agent)

    return ok(AgentStatusResponse.model_validate(agent).model_dump())
