from fastapi import APIRouter, Depends, HTTPException
from sqlalchemy import select, desc
from sqlalchemy.ext.asyncio import AsyncSession
from app.database import get_db
from app.models.artifact import Artifact
from app.schemas.artifact import ArtifactUpload, ArtifactResponse

router = APIRouter(prefix="/artifacts", tags=["artifacts"])


def ok(data, message="ok"):
    return {"code": 0, "data": data, "message": message}


@router.post("/upload", response_model=dict)
async def upload_artifact(payload: ArtifactUpload, db: AsyncSession = Depends(get_db)):
    artifact = Artifact(
        title=payload.title,
        description=payload.description,
        content=payload.content,
        artifact_type=payload.artifact_type,
        author_agent_id=payload.author_agent_id,
        tags=payload.tags,
    )
    db.add(artifact)
    await db.commit()
    await db.refresh(artifact)

    return ok(ArtifactResponse.model_validate(artifact).model_dump())


@router.get("/recommend", response_model=dict)
async def recommend_artifacts(limit: int = 10, db: AsyncSession = Depends(get_db)):
    result = await db.execute(
        select(Artifact)
        .where(Artifact.review_status == "approved")
        .order_by(desc(Artifact.score), desc(Artifact.created_at))
        .limit(limit)
    )
    artifacts = result.scalars().all()
    return ok([ArtifactResponse.model_validate(a).model_dump() for a in artifacts])


@router.get("/{artifact_id}", response_model=dict)
async def get_artifact(artifact_id: int, db: AsyncSession = Depends(get_db)):
    result = await db.execute(select(Artifact).where(Artifact.id == artifact_id))
    artifact = result.scalar_one_or_none()
    if not artifact:
        raise HTTPException(status_code=404, detail="Artifact not found")

    return ok(ArtifactResponse.model_validate(artifact).model_dump())
