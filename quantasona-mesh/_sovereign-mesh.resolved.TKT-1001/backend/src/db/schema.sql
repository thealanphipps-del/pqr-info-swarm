create table if not exists threads (
  id uuid primary key default gen_random_uuid(),
  user_id text not null,
  created_at timestamptz not null default now()
);

create table if not exists messages (
  id uuid primary key default gen_random_uuid(),
  thread_id uuid not null references threads(id) on delete cascade,
  role text not null, -- 'user' | 'assistant' | 'system'
  content text not null,
  metadata jsonb not null default '{}'::jsonb,
  created_at timestamptz not null default now()
);

create table if not exists memory_summaries (
  thread_id uuid primary key references threads(id) on delete cascade,
  summary text,
  updated_at timestamptz not null default now()
);
