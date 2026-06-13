# Bus Reservation Frontend

This is the mobile-first frontend for the Bus Reservation Platform, built with React, TypeScript, and MUI.

## Setup

1.  **Install dependencies**:
    ```bash
    npm install
    ```

2.  **Run in development**:
    ```bash
    npm run dev
    ```
    The dev server includes a proxy to `http://localhost:8080` for API requests.

3.  **Build for production**:
    ```bash
    npm run build
    ```

## Key Features

- **Mobile-First Design**: Optimized for mobile WebView with bottom navigation and touch-friendly targets.
- **JWT Authentication**: Secure login/register flow with token-based authorization.
- **Real-time Inventory**: Atomically updated seat availability.
- **Order Management**: Asynchronous order processing with status tracking.
