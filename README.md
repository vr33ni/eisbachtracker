# Eisbachtracker

A web app to track surfer activity and water conditions at the Eisbach wave in Munich 🌊🏄‍♂️

This repository contains a frontend (Vue.js), a backend (Go) and a machine learning model (Python).

---

## Structure

├── frontend/ → Vue 3 client (Vite) for the web interface  
├── go-server/ → Go backend API with PostgreSQL & Flyway  
├── ml-model/ → Flask API based on a linear regression model to create a prediction
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
- [Ml Model README](./ml-model/README.md) — Linear Regression Model + Flask API


---

## Future Ideas

- Add automatic canary deployments on pull requests
  - Deploy PRs to `/canary/{branch-name}` on GitHub Pages
  - Optional: Deploy backend canary to Render preview environment
  - Automatically comment preview URL on PR

## License

MIT — Feel free to use, adapt, and surf responsibly 🏄‍♀️
