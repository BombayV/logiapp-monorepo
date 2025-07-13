# Logi Web - SvelteKit Application

A modern SvelteKit application configured for Cloudflare Pages deployment using Bun.

## Development

Install dependencies and start the development server:

```bash
bun install
bun run dev
```

## Building

To create a production build optimized for Cloudflare Pages:

```bash
bun run build
```

## Cloudflare Deployment

This project is configured with the Cloudflare adapter for seamless deployment to Cloudflare Pages.

### Local Development with Cloudflare

Test your app locally with Cloudflare's runtime:

```bash
bun run build
bun run cf:dev
```

### Deploy to Cloudflare Pages

1. **First time setup**: Login to Wrangler
   ```bash
   bunx wrangler login
   ```

2. **Deploy your application**:
   ```bash
   bun run cf:deploy
   ```

3. **Monitor deployments**:
   ```bash
   bun run cf:tail
   ```

### Environment Variables

1. Copy the example environment file:
   ```bash
   cp .dev.vars.example .dev.vars
   ```

2. Add your environment variables to `.dev.vars` for local development

3. For production, set environment variables in the Cloudflare Pages dashboard or using Wrangler:
   ```bash
   bunx wrangler pages secret put VARIABLE_NAME
   ```

### Configuration

- **Adapter**: `@sveltejs/adapter-cloudflare` for Cloudflare Pages optimization
- **Runtime**: Cloudflare Workers with Node.js compatibility
- **Build Output**: `.svelte-kit/cloudflare/`

## Scripts

- `bun run dev` - Start development server
- `bun run build` - Build for production
- `bun run preview` - Preview production build locally
- `bun run cf:dev` - Start local Cloudflare development server
- `bun run cf:deploy` - Deploy to Cloudflare Pages
- `bun run cf:tail` - Monitor deployment logs
