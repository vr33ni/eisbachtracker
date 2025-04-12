# Eisbachtracker

A web app to track surfer activity and water conditions at the Eisbach wave in Munich 🌊🏄‍♂️

This repository contains both the frontend (Vue.js) and backend (Go) of the Eisbachtracker project.

---

## Structure

├── frontend/ → Vue 3 client (Vite) for the web interface  
├── go-server/ → Go backend API with PostgreSQL & Flyway  
└── README.md → This file

---

## About the project

- Users can view current water levels, temperature, and recent surfer counts.
- Surfer activity predictions based on recent data and water temperature.
- Data is collected and stored in a PostgreSQL database (Neon).
- The backend is hosted on Render, the frontend is deployed on GitHub Pages.

---

## More details

Check the respective READMEs for detailed setup and usage instructions:

- [Frontend README](./client/README.md) — Vue.js Client
- [Go Backend README](./go-server/README.md) — Go API + Database Migrations + Render Deploy

---

## License

MIT — Feel free to use, adapt, and surf responsibly 🏄‍♀️
