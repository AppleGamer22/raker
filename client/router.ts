import { createConnectTransport } from "@connectrpc/connect-web";
import { QueryClient } from "@tanstack/react-query";
import { createRouter } from "@tanstack/react-router";

import { routeTree } from "./routeTree.gen";

export const transport = createConnectTransport({ baseUrl: "/api" });

export const queryClient = new QueryClient({
	defaultOptions: {
		queries: {
			refetchOnWindowFocus: false,
		},
	},
});

export const router = createRouter({
	routeTree,
	context: {
		queryClient,
	},
	scrollRestoration: true,
	defaultPreload: "render",
	defaultPreloadStaleTime: 0,
});

declare module "@tanstack/react-router" {
	interface Register {
		router: typeof router;
	}
}
