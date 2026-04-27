import { QueryClient } from "@tanstack/react-query";
import { createRouter } from "@tanstack/react-router";

import { routeTree } from "./routeTree.gen";

export const queryClient = new QueryClient();

export const router = createRouter({
	routeTree,
	context: {
		queryClient,
	},
	scrollRestoration: true,
	defaultPreload: "render",
	defaultPreloadStaleTime: 0,
});

export const getRouter = () => router;

declare module "@tanstack/react-router" {
	interface Register {
		router: ReturnType<typeof getRouter>;
	}
}
