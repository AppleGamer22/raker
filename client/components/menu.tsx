import { Link, useLocation } from "@tanstack/react-router";
import { UserKeyIcon, DatabaseSearchIcon } from "lucide-react";
import type { ReactNode } from "react";

import { RakerLogo } from "@/components/logo";
import { ThemeToggle } from "@/components/theme-toggle";
import { Sheet, SheetContent, SheetDescription, SheetHeader, SheetTitle } from "@/components/ui/sheet";
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
import { SnapchatIcon } from "@/components/ui/svgs/snapchat";
import { TikTokIcon } from "@/components/ui/svgs/tiktok";
import { VSCOIcon } from "@/components/ui/svgs/vsco";
import { useUser } from "@/hooks/user-provider";
import { defaultPostTypes } from "@/lib/utils";

type MenuMode = "default" | "mobile-sheet";

export function Menu({
	mode = "default",
	onNavigate,
	children,
}: { mode?: MenuMode; onNavigate?: () => void; children?: ReactNode } = {}) {
	const isMobileSheet = mode === "mobile-sheet";
	const { pathname } = useLocation();
	const { username, categories, isCategoriesPending } = useUser();
	const isSignedIn = username !== null;

	const isActiveRoute = (route: string) => pathname === route;

	// TODO: https://ui.shadcn.com/blocks/sidebar#sidebar-16
	return (
		<SidebarProvider className={isMobileSheet ? "flex h-full min-h-0! w-full flex-col" : "min-h-svh w-full"}>
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
								<SidebarMenuButton
									isActive={isActiveRoute("/")}
									render={<Link to="/" onClick={() => onNavigate?.()} />}
								>
									<UserKeyIcon className="h-4 w-4" />
									Authentication
								</SidebarMenuButton>
							</SidebarMenuItem>
						</SidebarMenu>
					</SidebarGroup>
					<SidebarGroup>
						<SidebarGroupLabel>Extractors</SidebarGroupLabel>
						<SidebarMenu className="gap-1">
							<SidebarMenuItem>
								<SidebarMenuButton
									disabled={!isSignedIn}
									isActive={isActiveRoute("/instagram")}
									render={
										<Link
											disabled={!isSignedIn}
											to="/instagram"
											search={{ post: "", incognito: false }}
											onClick={() => onNavigate?.()}
										/>
									}
								>
									<InstagramIcon className="h-4 w-4" />
									Instagram Post
								</SidebarMenuButton>
							</SidebarMenuItem>
							<SidebarMenuItem>
								<SidebarMenuButton
									disabled={!isSignedIn}
									isActive={isActiveRoute("/highlight")}
									render={
										<Link
											disabled={!isSignedIn}
											to="/highlight"
											search={{ highlight: "" }}
											onClick={() => onNavigate?.()}
										/>
									}
								>
									<InstagramIcon className="h-4 w-4" />
									Instagram Highlight
								</SidebarMenuButton>
							</SidebarMenuItem>
							<SidebarMenuItem>
								<SidebarMenuButton
									disabled={!isSignedIn}
									isActive={isActiveRoute("/story")}
									render={
										<Link
											disabled={!isSignedIn}
											to="/story"
											search={{ owner: "" }}
											onClick={() => onNavigate?.()}
										/>
									}
								>
									<InstagramIcon className="h-4 w-4" />
									Instagram Story
								</SidebarMenuButton>
							</SidebarMenuItem>
							<SidebarMenuItem>
								<SidebarMenuButton
									disabled={!isSignedIn}
									isActive={isActiveRoute("/tiktok")}
									render={
										<Link
											disabled={!isSignedIn}
											to="/tiktok"
											search={{ owner: "", post: "", incognito: false }}
											onClick={() => onNavigate?.()}
										/>
									}
								>
									<TikTokIcon className="h-40 w-40" />
									TikTok Post
								</SidebarMenuButton>
							</SidebarMenuItem>
							<SidebarMenuItem>
								<SidebarMenuButton
									disabled={!isSignedIn}
									isActive={isActiveRoute("/snapchat")}
									render={
										<Link
											disabled={!isSignedIn}
											to="/snapchat"
											search={{ owner: "", highlight: "" }}
											onClick={() => onNavigate?.()}
										/>
									}
								>
									<SnapchatIcon className="h-40 w-40" />
									Snapchat Highlight
								</SidebarMenuButton>
							</SidebarMenuItem>
							<SidebarMenuItem>
								<SidebarMenuButton
									disabled={!isSignedIn}
									isActive={isActiveRoute("/vsco")}
									render={
										<Link
											disabled={!isSignedIn}
											to="/vsco"
											search={{ owner: "", post: "" }}
											onClick={() => onNavigate?.()}
										/>
									}
								>
									<VSCOIcon className="w-4" />
									VSCO Post
								</SidebarMenuButton>
							</SidebarMenuItem>
						</SidebarMenu>
					</SidebarGroup>
					<SidebarGroup>
						<SidebarGroupLabel>Search</SidebarGroupLabel>
						<SidebarMenu className="gap-1">
							<SidebarMenuItem>
								<SidebarMenuButton
									disabled={!isSignedIn}
									isActive={isActiveRoute("/history")}
									render={
										<Link
											disabled={!isSignedIn}
											to="/history"
											search={{
												categories: isCategoriesPending ? [] : categories,
												exclusive: false,
												owners: [],
												types: defaultPostTypes,
												page: 1n,
											}}
											onClick={() => onNavigate?.()}
										/>
									}
								>
									<DatabaseSearchIcon className="w-4" />
									History
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
