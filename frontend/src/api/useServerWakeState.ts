import { useSyncExternalStore } from "react";
import { getServerWakeStateSnapshot, subscribeToServerWakeState } from "./client";

export function useServerWakeState() {
  return useSyncExternalStore(
    subscribeToServerWakeState,
    getServerWakeStateSnapshot,
    getServerWakeStateSnapshot,
  );
}
