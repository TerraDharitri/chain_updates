import { createRoot } from "react-dom/client";
import App from "./App.tsx";
import "./index.css";
import { initApp } from "./lib";
import { config } from "./initConfig.ts";
console.log(config);
initApp(config).then(() => {
  (window as any).multiversxWallet = {};
  createRoot(document.getElementById("root")!).render(<App />);
});
