"""FastAPI application entry point."""
from contextlib import asynccontextmanager
from typing import AsyncIterator
import structlog
from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from redis.asyncio import Redis

from core.config import settings
from core.event_bus import EventBus
from core.llm import LMStudioClient

# Configure structured logging
structlog.configure(
    processors=[
        structlog.stdlib.filter_by_level,
        structlog.stdlib.add_logger_name,
        structlog.stdlib.add_log_level,
        structlog.stdlib.PositionalArgumentsFormatter(),
        structlog.processors.TimeStamper(fmt="iso"),
        structlog.processors.StackInfoRenderer(),
        structlog.processors.format_exc_info,
        structlog.processors.UnicodeDecoder(),
        structlog.processors.JSONRenderer() if settings.log_format == "json"
        else structlog.dev.ConsoleRenderer(),
    ],
    wrapper_class=structlog.stdlib.BoundLogger,
    context_class=dict,
    logger_factory=structlog.stdlib.LoggerFactory(),
    cache_logger_on_first_use=True,
)

logger = structlog.get_logger()


@asynccontextmanager
async def lifespan(app: FastAPI) -> AsyncIterator[None]:
    """Application lifespan manager."""
    logger.info("🏭 Starting Software Factory Orchestrator")
    
    # Initialize Redis
    redis = Redis.from_url(settings.redis_url, decode_responses=False)
    app.state.redis = redis
    
    # Initialize Event Bus
    event_bus = EventBus(redis)
    app.state.event_bus = event_bus
    
    # Initialize LLM Client
    llm_client = LMStudioClient(base_url=settings.lm_studio_url)
    app.state.llm_client = llm_client
    
    # Check LM Studio health
    is_healthy = await llm_client.check_health()
    if is_healthy:
        logger.info("✅ LM Studio is running")
    else:
        logger.warning("⚠️ LM Studio not detected")
    
    logger.info("✅ Orchestrator started successfully")
    
    yield
    
    # Shutdown
    logger.info("🛑 Shutting down Orchestrator")
    await llm_client.close()
    await redis.close()


# Create FastAPI app
app = FastAPI(
    title=settings.app_name,
    version=settings.app_version,
    lifespan=lifespan,
)

# CORS middleware
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)


@app.get("/health")
async def health_check():
    """Health check endpoint."""
    return {
        "status": "ok",
        "app": settings.app_name,
        "version": settings.app_version,
    }


@app.get("/")
async def root():
    """Root endpoint."""
    return {
        "name": settings.app_name,
        "version": settings.app_version,
        "docs": "/docs",
    }


if __name__ == "__main__":
    import uvicorn
    
    uvicorn.run(
        "main:app",
        host=settings.api_host,
        port=settings.api_port,
        reload=settings.debug,
        log_level=settings.log_level.lower(),
    )
