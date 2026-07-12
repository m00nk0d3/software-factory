"""Configuration management."""
from typing import Optional
from pydantic_settings import BaseSettings, SettingsConfigDict


class Settings(BaseSettings):
    """Application settings."""
    
    model_config = SettingsConfigDict(
        env_file=".env",
        env_file_encoding="utf-8",
        case_sensitive=False,
    )
    
    # Application
    app_name: str = "Software Factory"
    app_version: str = "0.1.0"
    debug: bool = False
    
    # API
    api_host: str = "0.0.0.0"
    api_port: int = 8000
    api_workers: int = 4
    
    # Database
    postgres_host: str = "localhost"
    postgres_port: int = 5432
    postgres_user: str = "factory"
    postgres_password: str = "factory"
    postgres_db: str = "factory"
    
    @property
    def database_url(self) -> str:
        """Get database URL."""
        return (
            f"postgresql+asyncpg://{self.postgres_user}:{self.postgres_password}"
            f"@{self.postgres_host}:{self.postgres_port}/{self.postgres_db}"
        )
    
    # Redis
    redis_host: str = "localhost"
    redis_port: int = 6379
    redis_db: int = 0
    redis_password: Optional[str] = None
    
    @property
    def redis_url(self) -> str:
        """Get Redis URL."""
        if self.redis_password:
            return f"redis://:{self.redis_password}@{self.redis_host}:{self.redis_port}/{self.redis_db}"
        return f"redis://{self.redis_host}:{self.redis_port}/{self.redis_db}"
    
    # LM Studio
    lm_studio_url: str = "http://localhost:1234/v1"
    lm_studio_timeout: float = 120.0
    
    # Sandcastle
    sandcastle_url: str = "http://localhost:3001"
    
    # Worktrees
    worktree_root: str = "~/worktrees"
    
    # Observability
    jaeger_host: str = "localhost"
    jaeger_port: int = 6831
    prometheus_port: int = 9091
    
    # Logging
    log_level: str = "INFO"
    log_format: str = "json"


settings = Settings()
