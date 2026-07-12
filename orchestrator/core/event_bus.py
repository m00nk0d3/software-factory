"""Event bus implementation using Redis Streams."""
import json
import uuid
from datetime import datetime
from typing import Any, Callable, Dict, List, Optional
from dataclasses import dataclass, asdict
import asyncio
import structlog
from redis.asyncio import Redis

logger = structlog.get_logger()


@dataclass
class Event:
    """Event data structure."""
    
    id: str
    type: str
    source: str
    data: Dict[str, Any]
    timestamp: datetime
    trace_id: Optional[str] = None
    
    def to_dict(self) -> Dict[str, Any]:
        """Convert event to dictionary."""
        d = asdict(self)
        d['timestamp'] = self.timestamp.isoformat()
        return d
    
    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> 'Event':
        """Create event from dictionary."""
        data['timestamp'] = datetime.fromisoformat(data['timestamp'])
        return cls(**data)


class EventBus:
    """Event bus using Redis Streams."""
    
    def __init__(self, redis: Redis, stream_name: str = "factory:events"):
        self.redis = redis
        self.stream_name = stream_name
        self.handlers: Dict[str, List[Callable]] = {}
        self._running = False
        self._consumer_tasks: List[asyncio.Task] = []
        
    async def publish(self, event: Event) -> str:
        """Publish event to the stream.
        
        Args:
            event: Event to publish
            
        Returns:
            Event ID from Redis
        """
        event_dict = event.to_dict()
        
        # Serialize data as JSON
        stream_data = {
            "id": event.id,
            "type": event.type,
            "source": event.source,
            "timestamp": event_dict["timestamp"],
            "data": json.dumps(event.data),
        }
        
        if event.trace_id:
            stream_data["trace_id"] = event.trace_id
        
        # Add to stream
        message_id = await self.redis.xadd(self.stream_name, stream_data)
        
        logger.info(
            "event_published",
            event_id=event.id,
            event_type=event.type,
            message_id=message_id.decode(),
        )
        
        return message_id.decode()
    
    def subscribe(self, event_type: str, handler: Callable) -> None:
        """Subscribe to events of a specific type.
        
        Args:
            event_type: Type of event to subscribe to
            handler: Async function to handle events
        """
        if event_type not in self.handlers:
            self.handlers[event_type] = []
        self.handlers[event_type].append(handler)
        
        logger.info("event_subscription_added", event_type=event_type)
    
    async def start_consuming(
        self,
        consumer_group: str = "factory",
        consumer_name: Optional[str] = None,
        batch_size: int = 10,
    ) -> None:
        """Start consuming events from the stream.
        
        Args:
            consumer_group: Redis consumer group name
            consumer_name: Consumer name (default: random UUID)
            batch_size: Number of messages to read per batch
        """
        if not consumer_name:
            consumer_name = f"consumer-{uuid.uuid4().hex[:8]}"
        
        # Create consumer group if it doesn't exist
        try:
            await self.redis.xgroup_create(
                self.stream_name,
                consumer_group,
                id="0",
                mkstream=True,
            )
            logger.info("consumer_group_created", group=consumer_group)
        except Exception as e:
            # Group already exists
            logger.debug("consumer_group_exists", group=consumer_group)
        
        self._running = True
        logger.info(
            "event_consumer_started",
            group=consumer_group,
            consumer=consumer_name,
        )
        
        while self._running:
            try:
                # Read messages
                messages = await self.redis.xreadgroup(
                    consumer_group,
                    consumer_name,
                    {self.stream_name: ">"},
                    count=batch_size,
                    block=1000,  # 1 second timeout
                )
                
                if not messages:
                    continue
                
                # Process messages
                for stream, message_list in messages:
                    for message_id, data in message_list:
                        await self._handle_message(
                            message_id.decode(),
                            data,
                            consumer_group,
                        )
                        
            except asyncio.CancelledError:
                logger.info("event_consumer_cancelled")
                break
            except Exception as e:
                logger.error("event_consumer_error", error=str(e))
                await asyncio.sleep(1)  # Back off on error
    
    async def _handle_message(
        self,
        message_id: str,
        data: Dict[bytes, bytes],
        consumer_group: str,
    ) -> None:
        """Handle a single message.
        
        Args:
            message_id: Redis message ID
            data: Message data
            consumer_group: Consumer group name
        """
        try:
            # Decode message
            decoded_data = {
                k.decode(): v.decode() for k, v in data.items()
            }
            
            # Parse event
            event = Event(
                id=decoded_data["id"],
                type=decoded_data["type"],
                source=decoded_data["source"],
                timestamp=datetime.fromisoformat(decoded_data["timestamp"]),
                data=json.loads(decoded_data["data"]),
                trace_id=decoded_data.get("trace_id"),
            )
            
            # Find handlers
            handlers = self.handlers.get(event.type, [])
            
            if handlers:
                # Execute handlers
                for handler in handlers:
                    try:
                        await handler(event)
                    except Exception as e:
                        logger.error(
                            "event_handler_error",
                            event_type=event.type,
                            error=str(e),
                        )
            
            # Acknowledge message
            await self.redis.xack(self.stream_name, consumer_group, message_id)
            
        except Exception as e:
            logger.error(
                "message_processing_error",
                message_id=message_id,
                error=str(e),
            )
    
    async def stop(self) -> None:
        """Stop consuming events."""
        self._running = False
        
        # Cancel all consumer tasks
        for task in self._consumer_tasks:
            task.cancel()
        
        await asyncio.gather(*self._consumer_tasks, return_exceptions=True)
        
        logger.info("event_bus_stopped")
