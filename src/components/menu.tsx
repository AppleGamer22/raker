import { UserKeyIcon, DatabaseSearchIcon } from "lucide-react";

import {
	Sidebar,
	SidebarGroup,
	SidebarMenu,
	SidebarMenuButton,
	SidebarMenuItem,
	SidebarHeader,
	SidebarContent,
	SidebarProvider,
	SidebarGroupLabel,
} from "@/components/ui/sidebar";
import { InstagramIcon } from "@/components/ui/svgs/instagram";
import { TikTokIcon } from "@/components/ui/svgs/tiktok";
import { VSCOIcon } from "@/components/ui/svgs/vsco";
import { ThemeToggle } from "@/components/ui/theme-toggle";

export function Menu() {
	// TODO: https://ui.shadcn.com/blocks/sidebar#sidebar-16
	return (
		<SidebarProvider className="flex min-h-0! flex-col">
			<Sidebar>
				<SidebarHeader>
					<SidebarGroup className="grid grid-cols-[1fr_auto_1fr] items-center">
						<div className="justify-self-center">
							<img alt="Raker Logo" src="/raker.svg" className="w-6" />
						</div>
						<div className="justify-self-end">
							<ThemeToggle />
						</div>
					</SidebarGroup>
				</SidebarHeader>
				<SidebarContent>
					<SidebarGroup>
						<SidebarGroupLabel>Settings</SidebarGroupLabel>
						<SidebarMenu>
							<SidebarMenuItem>
								<SidebarMenuButton>
									<UserKeyIcon className="h-4 w-4" />
									Authentication
								</SidebarMenuButton>
							</SidebarMenuItem>
						</SidebarMenu>
					</SidebarGroup>
					<SidebarGroup>
						<SidebarGroupLabel>Extractors</SidebarGroupLabel>
						<SidebarMenu>
							<SidebarMenuItem>
								<SidebarMenuButton>
									<InstagramIcon className="h-4 w-4" />
									Instagram Post
								</SidebarMenuButton>
							</SidebarMenuItem>
							<SidebarMenuItem>
								<SidebarMenuButton>
									<InstagramIcon className="h-4 w-4" />
									Instagram Highlight
								</SidebarMenuButton>
							</SidebarMenuItem>
							<SidebarMenuItem>
								<SidebarMenuButton>
									<InstagramIcon className="h-4 w-4" />
									Instagram Story
								</SidebarMenuButton>
							</SidebarMenuItem>
							<SidebarMenuItem>
								<SidebarMenuButton>
									<TikTokIcon className="h-40 w-40" />
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
					<SidebarGroup>
						<SidebarGroupLabel>Search</SidebarGroupLabel>
						<SidebarMenu>
							<SidebarMenuItem>
								<SidebarMenuButton>
									<DatabaseSearchIcon className="w-4" />
									History
								</SidebarMenuButton>
							</SidebarMenuItem>
						</SidebarMenu>
					</SidebarGroup>
				</SidebarContent>
			</Sidebar>
		</SidebarProvider>
	);
}
