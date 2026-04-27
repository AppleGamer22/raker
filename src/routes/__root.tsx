import { TanStackDevtools } from "@tanstack/react-devtools";
import { QueryClientProvider } from "@tanstack/react-query";
import { ReactQueryDevtoolsPanel } from "@tanstack/react-query-devtools";
import { createRootRoute, Outlet } from "@tanstack/react-router";
import { TanStackRouterDevtoolsPanel } from "@tanstack/react-router-devtools";

import { ThemeProvider } from "#/components/theme-provider";
import {
	Sidebar,
	SidebarGroup,
	SidebarMenu,
	SidebarMenuButton,
	SidebarMenuItem,
	SidebarProvider,
} from "#/components/ui/sidebar";
import { InstagramIcon } from "#/components/ui/svgs/instagram";
import { TikTokIcon } from "#/components/ui/svgs/tiktok";
import { VSCOIcon } from "#/components/ui/svgs/vsco";
import { ModeToggle } from "#/components/ui/theme-toggle";
import { queryClient, router } from "#/router";

export const Route = createRootRoute({
	component: Root,
});

function Root() {
	// TODO: https://ui.shadcn.com/blocks/sidebar#sidebar-16
	return (
		<>
			<ThemeProvider storageKey="raker-ui-theme">
				<QueryClientProvider client={queryClient}>
					<header className="sticky top-0 z-50 flex w-full items-center border-b bg-background">
						<div className="flex h-(--header-height) w-full items-center gap-2 px-4">
							<ModeToggle />
						</div>
					</header>
					<SidebarProvider className="flex flex-col">
						<Sidebar className="top-(--header-height) h-[calc(100svh-var(--header-height))]!">
							<SidebarGroup>
								<SidebarMenu>
									<SidebarMenuItem>
										<SidebarMenuButton>
											<InstagramIcon className="w-4 h-4" />
											Instagram Post
										</SidebarMenuButton>
									</SidebarMenuItem>
									<SidebarMenuItem>
										<SidebarMenuButton>
											<InstagramIcon className="w-4 h-4" />
											Instagram Highlight
										</SidebarMenuButton>
									</SidebarMenuItem>
									<SidebarMenuItem>
										<SidebarMenuButton>
											<InstagramIcon className="w-4 h-4" />
											Instagram Story
										</SidebarMenuButton>
									</SidebarMenuItem>
									<SidebarMenuItem>
										<SidebarMenuButton>
											<TikTokIcon className="w-4" />
											TikTok Post
										</SidebarMenuButton>
									</SidebarMenuItem>
									<SidebarMenuItem>
										<SidebarMenuButton>
											<VSCOIcon className="w-4" />
											VSCO Post
										</SidebarMenuButton>
									</SidebarMenuItem>
								</SidebarMenu>
							</SidebarGroup>
						</Sidebar>
					</SidebarProvider>
					<Outlet />
				</QueryClientProvider>
			</ThemeProvider>
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
						name: "Tanstack Query",
						render: <ReactQueryDevtoolsPanel client={queryClient} />,
					},
				]}
			/>
		</>
	);
}
