/**
 * Application-wide state shared via React context.
 * Currently tracks the selected AWS region, persisted to localStorage.
 */
import { useState, useEffect, createContext, useContext } from "react";
import { setTransportRegion } from "./transport";

interface AppState {
  region: string;
  setRegion: (r: string) => void;
  sidebarCollapsed: boolean;
  setSidebarCollapsed: (v: boolean) => void;
}

const AppContext = createContext<AppState>({
  region: "us-east-1",
  setRegion: () => {},
  sidebarCollapsed: false,
  setSidebarCollapsed: () => {},
});

/** All AWS regions available for selection in the header dropdown. */
const REGION_LABELS: Record<string, string> = {
  "us-east-1": "US East (N. Virginia)",
  "us-east-2": "US East (Ohio)",
  "us-west-1": "US West (N. California)",
  "us-west-2": "US West (Oregon)",
  "eu-west-1": "EU (Ireland)",
  "eu-west-2": "EU (London)",
  "eu-central-1": "EU (Frankfurt)",
  "ap-northeast-1": "Asia Pacific (Tokyo)",
  "ap-southeast-1": "Asia Pacific (Singapore)",
  "ap-southeast-2": "Asia Pacific (Sydney)",
  "ap-south-1": "Asia Pacific (Mumbai)",
  "sa-east-1": "South America (São Paulo)",
};

const REGIONS = Object.keys(REGION_LABELS);

const REGION_STORAGE_KEY = "vs-region";

/** Provider component that wraps the app and supplies region state. */
export function AppProvider({ children }: { children: React.ReactNode }) {
  const [region, setRegionState] = useState(() => {
    return localStorage.getItem(REGION_STORAGE_KEY) ?? "us-east-1";
  });
  const [sidebarCollapsed, setSidebarCollapsed] = useState(false);

  useEffect(() => {
    localStorage.setItem(REGION_STORAGE_KEY, region);
  }, [region]);

  function setRegion(r: string) {
    setRegionState(r);
    setTransportRegion(r);
  }

  return (
    <AppContext.Provider value={{ region, setRegion, sidebarCollapsed, setSidebarCollapsed }}>
      {children}
    </AppContext.Provider>
  );
}

/** Hook to access the current region and setter. */
export function useAppState() {
  return useContext(AppContext);
}

export { REGIONS, REGION_LABELS };
