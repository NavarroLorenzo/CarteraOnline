import { createRoot } from "react-dom/client";
import { BrowserRouter } from "react-router-dom";
import App from "./app/App";
import { AuthProvider } from "./auth/AuthContext";
import { AppErrorBoundary } from "./components/app/AppErrorBoundary";
import "./styles.css";

createRoot(document.getElementById("root")!).render(
  <BrowserRouter>
    <AppErrorBoundary>
      <AuthProvider>
        <App />
      </AuthProvider>
    </AppErrorBoundary>
  </BrowserRouter>,
);
