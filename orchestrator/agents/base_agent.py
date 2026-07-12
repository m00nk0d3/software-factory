"""Base agent implementation."""
from abc import ABC, abstractmethod
from typing import Dict, Any, Optional
from dataclasses import dataclass
import structlog

from core.llm import LMStudioClient
from core.event_bus import EventBus, Event

logger = structlog.get_logger()


@dataclass
class AgentContext:
    """Agent execution context."""
    
    task_id: str
    worktree_path: str
    branch: str
    metadata: Dict[str, Any]


class BaseAgent(ABC):
    """Base class for all agents."""
    
    def __init__(
        self,
        agent_id: str,
        agent_type: str,
        llm_client: LMStudioClient,
        event_bus: EventBus,
        model_type: str = "coding",
    ):
        self.agent_id = agent_id
        self.agent_type = agent_type
        self.llm_client = llm_client
        self.event_bus = event_bus
        self.model_type = model_type
        self.logger = logger.bind(agent_id=agent_id, agent_type=agent_type)
    
    @abstractmethod
    async def execute(self, context: AgentContext) -> Dict[str, Any]:
        """Execute the agent's task.
        
        Args:
            context: Execution context
            
        Returns:
            Result dictionary
        """
        pass
    
    @abstractmethod
    def get_system_prompt(self) -> str:
        """Get the agent's system prompt.
        
        Returns:
            System prompt string
        """
        pass
    
    async def generate(
        self,
        prompt: str,
        system: Optional[str] = None,
        temperature: Optional[float] = None,
    ) -> str:
        """Generate completion using LLM.
        
        Args:
            prompt: User prompt
            system: System prompt (defaults to agent's system prompt)
            temperature: Sampling temperature
            
        Returns:
            Generated text
        """
        if system is None:
            system = self.get_system_prompt()
        
        result = await self.llm_client.generate(
            prompt=prompt,
            system=system,
            model_type=self.model_type,
            temperature=temperature,
        )
        
        return result["content"]
    
    async def emit_event(
        self,
        event_type: str,
        data: Dict[str, Any],
        trace_id: Optional[str] = None,
    ) -> None:
        """Emit an event to the event bus.
        
        Args:
            event_type: Type of event
            data: Event data
            trace_id: Trace ID for distributed tracing
        """
        import uuid
        from datetime import datetime
        
        event = Event(
            id=str(uuid.uuid4()),
            type=event_type,
            source=f"{self.agent_type}:{self.agent_id}",
            data=data,
            timestamp=datetime.utcnow(),
            trace_id=trace_id,
        )
        
        await self.event_bus.publish(event)
    
    async def run(self, context: AgentContext) -> Dict[str, Any]:
        """Run the agent with error handling and event emission.
        
        Args:
            context: Execution context
            
        Returns:
            Result dictionary
        """
        self.logger.info("agent_started", task_id=context.task_id)
        
        # Emit start event
        await self.emit_event(
            "agent.started",
            {
                "task_id": context.task_id,
                "agent_id": self.agent_id,
                "agent_type": self.agent_type,
            },
        )
        
        try:
            result = await self.execute(context)
            
            # Emit success event
            await self.emit_event(
                "agent.completed",
                {
                    "task_id": context.task_id,
                    "agent_id": self.agent_id,
                    "result": result,
                },
            )
            
            self.logger.info("agent_completed", task_id=context.task_id)
            
            return result
            
        except Exception as e:
            self.logger.error(
                "agent_failed",
                task_id=context.task_id,
                error=str(e),
            )
            
            # Emit failure event
            await self.emit_event(
                "agent.failed",
                {
                    "task_id": context.task_id,
                    "agent_id": self.agent_id,
                    "error": str(e),
                },
            )
            
            raise
