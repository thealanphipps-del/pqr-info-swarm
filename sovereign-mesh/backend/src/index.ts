import express from 'express';
import cors from 'cors';
import { json } from 'body-parser';
import { chatRouter } from './api/chatRouter';
import { threadsRouter } from './api/threadsRouter';
import { loadConfig } from './config';

async function main() {
  const app = express();
  const config = loadConfig();

  app.use(cors());
  app.use(json());

  app.use('/api/chat', chatRouter);
  app.use('/api/threads', threadsRouter);

  // Health check endpoint
  app.get('/api/health', (req, res) => {
    res.json({ status: 'ONLINE', time: new Date() });
  });

  const port = config.PORT || 80;
  app.listen(port, () => {
    console.log(`=================================================`);
    console.log(`🚀 Sovereign Multi-Model Chat Backend is active`);
    console.log(`🔗 Gateway: http://localhost:${port}`);
    console.log(`=================================================`);
  });
}

main().catch((err) => {
  console.error('Fatal error starting Sovereign Gateway:', err);
  process.exit(1);
});
