import express, { Request, Response } from 'express';
import cors from 'cors';
import pino from 'pino';
import { createWorktree, run } from '@ai-hero/sandcastle';
import { docker } from '@ai-hero/sandcastle/sandboxes/docker';
import type { BranchStrategy } from '@ai-hero/sandcastle';

const logger = pino({
  transport: {
    target: 'pino-pretty',
    options: {
      colorize: true,
      translateTime: 'HH:MM:ss',
      ignore: 'pid,hostname',
    },
  },
});

const app = express();
const PORT = process.env.PORT || 3001;

app.use(cors());
app.use(express.json());

// Health check
app.get('/health', (req: Request, res: Response) => {
  res.json({ status: 'ok', timestamp: new Date().toISOString() });
});

// Create worktree
app.post('/api/worktree/create', async (req: Request, res: Response) => {
  try {
    const { branch, taskId } = req.body;

    if (!branch || !taskId) {
      return res.status(400).json({
        error: 'Missing required fields: branch, taskId',
      });
    }

    logger.info({ branch, taskId }, 'Creating worktree');

    const branchStrategy: BranchStrategy = {
      type: 'branch',
      branch,
    };

    await using wt = await createWorktree({ branchStrategy });

    const result = {
      worktreePath: wt.worktreePath,
      branch: wt.branch,
      taskId,
      status: 'created',
    };

    logger.info(result, 'Worktree created successfully');

    res.json(result);
  } catch (error) {
    logger.error({ error }, 'Failed to create worktree');
    res.status(500).json({
      error: error instanceof Error ? error.message : 'Unknown error',
    });
  }
});

// List worktrees
app.get('/api/worktree/list', async (req: Request, res: Response) => {
  try {
    // This would need to track active worktrees
    // For now, return empty array
    res.json({ worktrees: [] });
  } catch (error) {
    logger.error({ error }, 'Failed to list worktrees');
    res.status(500).json({
      error: error instanceof Error ? error.message : 'Unknown error',
    });
  }
});

// Delete worktree
app.delete('/api/worktree/:taskId', async (req: Request, res: Response) => {
  try {
    const { taskId } = req.params;

    logger.info({ taskId }, 'Deleting worktree');

    // Sandcastle handles cleanup automatically with `using` syntax
    // This endpoint is for manual cleanup if needed

    res.json({ status: 'deleted', taskId });
  } catch (error) {
    logger.error({ error }, 'Failed to delete worktree');
    res.status(500).json({
      error: error instanceof Error ? error.message : 'Unknown error',
    });
  }
});

// Run agent in worktree
app.post('/api/agent/run', async (req: Request, res: Response) => {
  try {
    const { taskId, prompt, worktreePath, maxIterations = 5 } = req.body;

    if (!prompt || !worktreePath) {
      return res.status(400).json({
        error: 'Missing required fields: prompt, worktreePath',
      });
    }

    logger.info({ taskId, worktreePath }, 'Running agent');

    // This is a placeholder - actual agent implementation would come from Python
    // For now, just acknowledge the request
    res.json({
      status: 'started',
      taskId,
      message: 'Agent execution started (placeholder)',
    });
  } catch (error) {
    logger.error({ error }, 'Failed to run agent');
    res.status(500).json({
      error: error instanceof Error ? error.message : 'Unknown error',
    });
  }
});

app.listen(PORT, () => {
  logger.info(`🏭 Sandcastle Bridge API listening on port ${PORT}`);
  logger.info(`📋 Health check: http://localhost:${PORT}/health`);
});
