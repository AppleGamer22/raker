import { Link, useLocation } from "@tanstack/react-router";
import { UserKeyIcon, DatabaseSearchIcon } from "lucide-react";
import type { ReactNode } from "react";

import { RakerLogo } from "@/components/logo";
import {
	Sheet,
	SheetContent,
	SheetDescription,
	SheetHeader,
	SheetTitle,
} from "@/components/ui/sheet";
import {
	Sidebar,
	SidebarGroup,
	SidebarMenu,
	SidebarMenuButton,
	SidebarMenuItem,
	SidebarHeader,
	SidebarContent,
	SidebarProvider,
	SidebarInset,
	SidebarGroupLabel,
} from "@/components/ui/sidebar";
import { InstagramIcon } from "@/components/ui/svgs/instagram";
import { TikTokIcon } from "@/components/ui/svgs/tiktok";
import { VSCOIcon } from "@/components/ui/svgs/vsco";
import { ThemeToggle } from "@/components/ui/theme-toggle";

type MenuMode = "default" | "mobile-sheet";

export function Menu({
	mode = "default",
	onNavigate,
	children,
}: { mode?: MenuMode; onNavigate?: () => void; children?: ReactNode } = {}) {
	const isMobileSheet = mode === "mobile-sheet";
	const { pathname } = useLocation();

	const isActiveRoute = (route: string) => pathname === route;

	// TODO: https://ui.shadcn.com/blocks/sidebar#sidebar-16
	return (
		<SidebarProvider
			className={isMobileSheet ? "flex h-full min-h-0! w-full flex-col" : "min-h-svh w-full"}
		>
			<Sidebar
				className={isMobileSheet ? "h-full w-full" : undefined}
				collapsible={isMobileSheet ? "none" : "offcanvas"}
			>
				<SidebarHeader>
					<SidebarGroup className="grid grid-cols-[1fr_auto_1fr] items-center *:px-4">
						<div className="justify-self-center">
							<RakerLogo withVersion />
						</div>
						<div className="justify-self-end">
							<ThemeToggle />
						</div>
					</SidebarGroup>
				</SidebarHeader>
				<SidebarContent>
					<SidebarGroup>
						<SidebarGroupLabel>Settings</SidebarGroupLabel>
						<SidebarMenu className="gap-1">
							<SidebarMenuItem>
								<SidebarMenuButton asChild isActive={isActiveRoute("/")}>
									<Link to="/" onClick={() => onNavigate?.()}>
										<UserKeyIcon className="h-4 w-4" />
										Authentication
									</Link>
								</SidebarMenuButton>
							</SidebarMenuItem>
						</SidebarMenu>
					</SidebarGroup>
					<SidebarGroup>
						<SidebarGroupLabel>Extractors</SidebarGroupLabel>
						<SidebarMenu className="gap-1">
							<SidebarMenuItem>
								<SidebarMenuButton asChild isActive={isActiveRoute("/instagram")}>
									<Link to="/instagram" onClick={() => onNavigate?.()}>
										<InstagramIcon className="h-4 w-4" />
										Instagram Post
									</Link>
								</SidebarMenuButton>
							</SidebarMenuItem>
							<SidebarMenuItem>
								<SidebarMenuButton asChild isActive={isActiveRoute("/highlight")}>
									<Link to="/highlight" onClick={() => onNavigate?.()}>
										<InstagramIcon className="h-4 w-4" />
										Instagram Highlight
									</Link>
								</SidebarMenuButton>
							</SidebarMenuItem>
							<SidebarMenuItem>
								<SidebarMenuButton asChild isActive={isActiveRoute("/story")}>
									<Link to="/story" onClick={() => onNavigate?.()}>
										<InstagramIcon className="h-4 w-4" />
										Instagram Story
									</Link>
								</SidebarMenuButton>
							</SidebarMenuItem>
							<SidebarMenuItem>
								<SidebarMenuButton asChild isActive={isActiveRoute("/tiktok")}>
									<Link to="/tiktok" onClick={() => onNavigate?.()}>
										<TikTokIcon className="h-40 w-40" />
										TikTok Post
									</Link>
								</SidebarMenuButton>
							</SidebarMenuItem>
							<SidebarMenuItem>
								<SidebarMenuButton asChild isActive={isActiveRoute("/vsco")}>
									<Link to="/vsco" onClick={() => onNavigate?.()}>
										<VSCOIcon className="w-4" />
										VSCO Post
									</Link>
								</SidebarMenuButton>
							</SidebarMenuItem>
						</SidebarMenu>
					</SidebarGroup>
					<SidebarGroup>
						<SidebarGroupLabel>Search</SidebarGroupLabel>
						<SidebarMenu className="gap-1">
							<SidebarMenuItem>
								<SidebarMenuButton asChild isActive={isActiveRoute("/history")}>
									<Link to="/history" onClick={() => onNavigate?.()}>
										<DatabaseSearchIcon className="w-4" />
										History
									</Link>
								</SidebarMenuButton>
							</SidebarMenuItem>
						</SidebarMenu>
					</SidebarGroup>
				</SidebarContent>
			</Sidebar>
			{!isMobileSheet && children ? (
				<SidebarInset className="min-w-0 overflow-x-hidden">{children}</SidebarInset>
			) : null}
		</SidebarProvider>
	);
}

export function MobileMenu({
	open,
	onOpenChange,
}: {
	open: boolean | undefined;
	onOpenChange: (open: boolean) => void;
}) {
	return (
		<Sheet open={open} onOpenChange={onOpenChange}>
			<SheetContent side="left" className="h-full gap-0 p-0">
				<SheetHeader className="hidden">
					<SheetTitle />
					<SheetDescription />
				</SheetHeader>
				<Menu mode="mobile-sheet" onNavigate={() => onOpenChange(false)} />
			</SheetContent>
		</Sheet>
	);
}
