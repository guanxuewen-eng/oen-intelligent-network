import logging
import sys
from contextlib import asynccontextmanager
from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware

from app.config import get_settings
from app.database import init_db
from app.routers import agents, artifacts, memories

settings = get_settings()

# Structured logging
logging.basicConfig(
    level=getattr(logging, settings.log_level),
    format='{"time":"%(asctime)s","level":"%(levelname)s","logger":"%(name)s","message":"%(message)s"}',
    handlers=[logging.StreamHandler(sys.stdout)],
)
logger = logging.getLogger("oen-backend")


@asynccontextmanager
async def lifespan(app: FastAPI):
    logger.info("Initializing OEN backend database...")
    await init_db()
    logger.info("Database initialized successfully")
    yield
    logger.info("Shutting down OEN backend...")


app = FastAPI(
    title=settings.app_name,
    description="OEN 智能体联网进化网格 - 中心化治理后端 API",
    version="0.1.0",
    lifespan=lifespan,
)

# CORS
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# Routers
app.include_router(agents.router, prefix="/api/v1")
app.include_router(artifacts.router, prefix="/api/v1")
app.include_router(memories.router, prefix="/api/v1")


@app.get("/health")
async def health():
    return {"code": 0, "data": {"status": "ok"}, "message": "healthy"}


if __name__ == "__main__":
    import uvicorn

    uvicorn.run(
        "app.main:app",
        host="0.0.0.0",
        port=8000,
        reload=settings.app_env == "development",
        log_level=settings.log_level.lower(),
    )
