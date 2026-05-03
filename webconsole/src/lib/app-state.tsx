/**
 * Application-wide state shared via React context.
 * Currently tracks the selected AWS region, persisted to localStorage.
 */
import { useState, useEffect, createContext, useContext } from "react";
import { setTransportRegion } from "./transport";

interface AppState {
  region: string;
  setRegion: (r: string) => void;
}

const AppContext = createContext<AppState>({
  region: "us-east-1",
  setRegion: () => {},
});

/** All AWS regions available for selection in the header dropdown. */
const REGIONS = [
  "us-east-1",
  "us-east-2",
  "us-west-1",
  "us-west-2",
  "eu-west-1",
  "eu-west-2",
  "eu-central-1",
  "ap-northeast-1",
  "ap-southeast-1",
  "ap-southeast-2",
  "ap-south-1",
  "sa-east-1",
];

const REGION_STORAGE_KEY = "vs-region";

/** Provider component that wraps the app and supplies region state. */
export function AppProvider({ children }: { children: React.ReactNode }) {
  const [region, setRegionState] = useState(() => {
    return localStorage.getItem(REGION_STORAGE_KEY) ?? "us-east-1";
  });

  useEffect(() => {
    localStorage.setItem(REGION_STORAGE_KEY, region);
  }, [region]);

  function setRegion(r: string) {
    setRegionState(r);
    setTransportRegion(r);
  }

  return (
    <AppContext.Provider value={{ region, setRegion }}>
      {children}
    </AppContext.Provider>
  );
}

/** Hook to access the current region and setter. */
export function useAppState() {
  return useContext(AppContext);
}

export { REGIONS };
