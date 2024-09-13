import { createBrowserRouter, RouteObject } from 'react-router-dom';

export const createRouter = (basename: string, routes: RouteObject[]) =>
  createBrowserRouter(routes, { basename });
