# 🌊 EisbachTracker PWA

A Progressive Web App (PWA) that shows live water level and flow for the Eisbach River in Munich — using Vue 3, Vite, Tailwind, and GitHub Pages.

---

## 📦 Features

- ✅ Live data from Pegelalarm API
- ✅ Water level alert if below threshold
- ✅ Works offline (PWA)
- ✅ Installable on mobile & desktop
- ✅ GitHub Actions CI/CD

---

## 🚀 Getting Started

### Clone & install

    git clone https://github.com/your-username/eisbachtracker-pwa.git
    cd eisbachtracker-pwa
    pnpm install   # or npm install

### Add your `.env` file

Create a `.env` file at the root with:

    VITE_API_URL=https://api.pegelalarm.at/api/station/1.0/list?commonid=16515005-de&responseDetailLevel=high

### Run the app

    pnpm dev

Visit: http://localhost:5173

---

## 🛠 Build & Preview (to locally test the install button - only works on Chrome)

    pnpm build      # Builds to ./dist
    pnpm preview    # Locally preview the built PWA

---

## ⚙️ Deploying to GitHub Pages

Uses GitHub Actions + [`peaceiris/actions-gh-pages`](https://github.com/peaceiris/actions-gh-pages)

### Setup

Go to **Settings > Secrets > Actions** and add:

- `VITE_API_URL`: same API URL as above

Update `vite.config.ts`:

    base: '/eisbachtracker-pwa/'

Push to `main` or `master`.

The GitHub Actions workflow will:

- Inject your `.env` secret
- Build the site
- Deploy `dist/` to `gh-pages` branch

Then set GitHub Pages source to `gh-pages` branch.

---

## 📱 PWA

- Works offline
- Chrome/Edge: install button appears
- iOS Safari: “Add to Home Screen” manually

---

## 🔄 Reset install (for testing)

- Chrome: go to `chrome://apps`, remove app
- iOS: long press → remove from home screen
- Clear cache in browser dev tools

---

## 🧑‍💻 Author

Built by [@vr33ni](https://github.com/vr33ni)  
Data from [pegelalarm.at](https://api.pegelalarm.at)

---
