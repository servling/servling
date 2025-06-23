🚀 Servling

    The container deployment platform that just works — Deploy, manage, and monitor your applications with zero complexity.

<div align="center">

</div>
✨ Why Choose Servling?

Tired of complex deployment pipelines? Servling transforms container deployment from a chore into a breeze. Whether you're a solo developer or managing enterprise applications, Servling adapts to your workflow.
🎯 Perfect For

    Self-hosters who want full control over their deployment infrastructure
    Developers deploying personal projects on their own servers
    Small teams managing applications on private infrastructure
    Anyone who values simplicity over complex orchestration platforms

🌟 Key Features
🖱️ One-Click Deployments

Deploy complex multi-service applications with a single click. No YAML wrestling required.
📊 Real-Time Monitoring

Watch your applications come alive with live status updates, logs, and health monitoring.
🎨 Modern Web Interface

Beautiful, responsive dashboard built with Nuxt.js — manage everything from anywhere.
🔒 Enterprise-Ready Security

Built-in authentication, authorization, and secure container isolation.
📋 Smart Templates

Pre-configured application templates get you from zero to deployed in minutes.
⚡ Self-Hosting Made Simple

Deploy on your own infrastructure without vendor lock-in or monthly fees.
🏗️ Architecture That Scales

┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   Web Interface │───▶│  Backend Server  │───▶│ Docker Engine   │
│    (Nuxt.js)    │    │     (Go)         │    │   (Containers)  │
└─────────────────┘    └──────────────────┘    └─────────────────┘
│
▼
┌─────────────────┐
│   PostgreSQL    │
│   (Database)    │
└─────────────────┘

Components:

    🎯 Backend Server: Robust Go API handling orchestration and management
    🌐 Web Interface: Intuitive Nuxt.js dashboard for all your deployment needs
    💾 PostgreSQL: Reliable data persistence for configurations and state

🚀 Quick Start
Prerequisites

bash

# Essential tools
Docker & Docker Compose
Go 1.21+
Node.js 18+
pnpm

⚡ Get Running in 60 Seconds

bash

# 1. Clone and enter
git clone https://github.com/servling/servling.git
cd servling

# 2. Spin up infrastructure
docker-compose -f compose.dev.yaml up -d

# 3. Install dependencies
pnpm install

# 4. Start backend (new terminal)
go run main.go

# 5. Start web interface (new terminal)
cd web && pnpm dev

🎉 That's it! Visit http://localhost:3000 and start deploying.
📁 Project Structure

servling/
├── 🖥️  apps/
│   ├── server/          # Go backend API
│   └── web/            # Nuxt.js frontend
├── 🗄️  ent/             # Database models & schema
├── 📦  packages/        # Shared utilities
└── 📋  schema/          # API specifications

💡 What Makes Servling Special
🔄 Real-Time Everything

Server-Sent Events keep your dashboard synchronized with container state changes instantly.
🎛️ Flexible Configuration

Environment variables, port mappings, labels, and volumes — all configurable through an intuitive UI.
🏢 Multi-Tenant Ready

Secure user authentication ensures teams can collaborate safely on shared infrastructure.
🔧 Developer Experience First

Built by developers, for developers. Every feature designed to reduce friction and increase productivity.
🛠️ Configuration

Servling can be configured using environment variables. Check the source code for all available configuration options.
🤝 Join the Community

We're building something amazing together! Here's how you can help:

    🌟 Star this repo if Servling helps with your self-hosting needs
    🐛 Report bugs to help us improve
    💡 Suggest features for future releases
    🔧 Submit PRs — all contributions welcome!

Contributing

    Fork the repository
    Create your feature branch (git checkout -b feature/amazing-feature)
    Commit your changes (git commit -m 'Add amazing feature')
    Push to the branch (git push origin feature/amazing-feature)
    Open a Pull Request

📄 License

Released under the MIT License — use it however you want!
<div align="center">

Ready to simplify your deployments?

⭐ Star on GitHub | 📚 View Source | 🐛 Report Issues

Made with ❤️ for developers who love self-hosting
</div>
