import { TransportProvider } from "@connectrpc/connect-query";
import { TanStackDevtools } from "@tanstack/react-devtools";
import { QueryClientProvider } from "@tanstack/react-query";
import { ReactQueryDevtoolsPanel } from "@tanstack/react-query-devtools";
import { RouterProvider } from "@tanstack/react-router";

import "./styles.css";
import { TanStackRouterDevtoolsPanel } from "@tanstack/react-router-devtools";
import { StrictMode } from "react";
import ReactDOM from "react-dom/client";

import { ThemeProvider } from "@/hooks/theme-provider";
import { UserProvider } from "@/hooks/user-provider";

import { queryClient, transport, router } from "./router";

const rootElement = document.getElementById("root");

if (!rootElement || !rootElement.innerHTML) {
	const root = ReactDOM.createRoot(rootElement!);
	root.render(
		<StrictMode>
			<ThemeProvider storageKey="raker-ui-theme">
				<TransportProvider transport={transport}>
					<QueryClientProvider client={queryClient}>
						<UserProvider>
							<RouterProvider router={router} />
						</UserProvider>
						<TanStackDevtools
							config={{
								position: "bottom-right",
							}}
							plugins={[
								{
									name: "TanStack Router",
									render: <TanStackRouterDevtoolsPanel router={router} />,
								},
								{
									name: "TanStack Query",
									render: <ReactQueryDevtoolsPanel client={queryClient} />,
								},
							]}
						/>
					</QueryClientProvider>
				</TransportProvider>
			</ThemeProvider>
		</StrictMode>,
	);
}
