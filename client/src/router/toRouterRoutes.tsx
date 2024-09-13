import { map } from "lodash/fp";
import { RouteObject } from "react-router-dom";

import { Route } from "./routes";
import { GlobalError } from "../components/error";

export const toRouterRoutes = (routes: Route[]): RouteObject[] =>
  map((route: Route) => ({
    ...route,
    id: route.path,
    element: <route.element />,
    errorElement: <GlobalError />,
  }))(routes);
