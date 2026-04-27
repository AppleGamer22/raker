import { createRootRoute, Outlet } from "@tanstack/react-router";

import Header from "@/components/header";
import { Menu } from "@/components/menu";
import {
	Sidebar,
	SidebarGroup,
	SidebarMenu,
	SidebarMenuButton,
	SidebarMenuItem,
	SidebarProvider,
} from "@/components/ui/sidebar";
import { InstagramIcon } from "@/components/ui/svgs/instagram";
import { TikTokIcon } from "@/components/ui/svgs/tiktok";
import { VSCOIcon } from "@/components/ui/svgs/vsco";

export const Route = createRootRoute({
	component: Root,
	ssr: false,
});

function Root() {
	// TODO: https://ui.shadcn.com/blocks/sidebar#sidebar-16
	return (
		<>
			<Header />
			<Menu />
			<Outlet />
		</>
	);
}
