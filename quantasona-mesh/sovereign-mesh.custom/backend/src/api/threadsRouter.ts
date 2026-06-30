import { Router } from 'express';
import { MessagesRepo } from '../db/repositories/threadsRepo';

export const threadsRouter = Router();
const messagesRepo = new MessagesRepo();

threadsRouter.get('/:id/messages', async (req, res) => {
  try {
    const messages = await messagesRepo.getThreadMessages(req.params.id);
    res.json({ messages });
  } catch (err: any) {
    res.status(500).json({ error: err.message });
  }
});
