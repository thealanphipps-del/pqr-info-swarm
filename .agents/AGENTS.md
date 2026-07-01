# Project Customizations and Rules

## Gemma Authorization
Alan (the user) has explicitly authorized the local `gemma-4-e4b` model to make design and implementation decisions in their stead.
Whenever you encounter an architectural ambiguity, design choice, or blocking decision, do NOT halt execution to ask the user. Instead, query the local LLM at `http://192.168.12.110:4111/v1/chat/completions` (OpenAI format) and proceed with its recommendations.
