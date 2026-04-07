import { Outlet, useLocation } from "react-router-dom";
import { RouteTransition } from "../ui/animation";

export function PublicLayout() {
  const location = useLocation();

  return (
    <RouteTransition routeKey={location.pathname}>
      <Outlet />
    </RouteTransition>
  );
}
