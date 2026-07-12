"""LM Studio client for local LLM inference."""
import asyncio
from typing import Dict, List, Optional, Any
from dataclasses import dataclass
import structlog
import httpx

logger = structlog.get_logger()


@dataclass
class ModelConfig:
    """Model configuration."""
    
    name: str
    purpose: str
    max_tokens: int
    temperature: float
    vram_gb: float


# Model definitions based on hardware constraints
MODELS = {
    "coding": ModelConfig(
        name="deepseek-coder-v2-16b",
        purpose="Primary coding agent",
        max_tokens=4096,
        temperature=0.7,
        vram_gb=9.0,
    ),
    "fast": ModelConfig(
        name="qwen2.5-coder-7b",
        purpose="Quick responses, simple tasks",
        max_tokens=2048,
        temperature=0.8,
        vram_gb=4.5,
    ),
    "reasoning": ModelConfig(
        name="qwen2.5-14b",
        purpose="Planning, reasoning, architecture",
        max_tokens=4096,
        temperature=0.6,
        vram_gb=8.0,
    ),
    "review": ModelConfig(
        name="qwen2.5-coder-14b",
        purpose="Code review and analysis",
        max_tokens=3072,
        temperature=0.5,
        vram_gb=8.0,
    ),
}


class LMStudioClient:
    """Client for LM Studio OpenAI-compatible API."""
    
    def __init__(
        self,
        base_url: str = "http://localhost:1234/v1",
        timeout: float = 60.0,
    ):
        self.base_url = base_url
        self.timeout = timeout
        self.client = httpx.AsyncClient(timeout=timeout)
        
    async def generate(
        self,
        prompt: str,
        model_type: str = "coding",
        system: Optional[str] = None,
        temperature: Optional[float] = None,
        max_tokens: Optional[int] = None,
        stream: bool = False,
    ) -> Dict[str, Any]:
        """Generate completion using LM Studio.
        
        Args:
            prompt: User prompt
            model_type: Type of model to use (coding, fast, reasoning, review)
            system: System prompt
            temperature: Sampling temperature (overrides model default)
            max_tokens: Maximum tokens (overrides model default)
            stream: Whether to stream the response
            
        Returns:
            Dict with content and metadata
        """
        model = MODELS.get(model_type)
        if not model:
            raise ValueError(f"Unknown model type: {model_type}")
        
        # Build messages
        messages = []
        if system:
            messages.append({"role": "system", "content": system})
        messages.append({"role": "user", "content": prompt})
        
        # Prepare request
        request_data = {
            "model": model.name,
            "messages": messages,
            "temperature": temperature or model.temperature,
            "max_tokens": max_tokens or model.max_tokens,
            "stream": stream,
        }
        
        logger.info(
            "llm_request",
            model=model.name,
            prompt_length=len(prompt),
            model_type=model_type,
        )
        
        try:
            response = await self.client.post(
                f"{self.base_url}/chat/completions",
                json=request_data,
            )
            response.raise_for_status()
            
            result = response.json()
            
            # Extract completion
            content = result["choices"][0]["message"]["content"]
            usage = result.get("usage", {})
            
            logger.info(
                "llm_response",
                model=model.name,
                tokens=usage.get("total_tokens", 0),
                completion_tokens=usage.get("completion_tokens", 0),
            )
            
            return {
                "content": content,
                "tokens": usage.get("total_tokens", 0),
                "model": model.name,
                "finish_reason": result["choices"][0].get("finish_reason"),
            }
            
        except httpx.HTTPError as e:
            logger.error("llm_request_failed", error=str(e))
            raise
    
    async def list_models(self) -> List[Dict[str, Any]]:
        """List available models from LM Studio.
        
        Returns:
            List of model info dicts
        """
        try:
            response = await self.client.get(f"{self.base_url}/models")
            response.raise_for_status()
            result = response.json()
            return result.get("data", [])
        except httpx.HTTPError as e:
            logger.error("list_models_failed", error=str(e))
            return []
    
    async def check_health(self) -> bool:
        """Check if LM Studio is running and responsive.
        
        Returns:
            True if healthy, False otherwise
        """
        try:
            models = await self.list_models()
            return len(models) > 0
        except Exception:
            return False
    
    async def close(self) -> None:
        """Close the HTTP client."""
        await self.client.aclose()
